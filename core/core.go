// Package core is only for internal use and not thread safe.
package core // import "ekyu.moe/soda/core"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"

	"ekyu.moe/leb128"
	"ekyu.moe/util/bytesutil"
	"github.com/awnumar/memguard"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/salsa20/salsa"
	"golang.org/x/crypto/sha3"
)

var (
	zeros [16]byte
)

// Session holds the key pair, the shared secret, shared nonce seed and seq
// number. Sensitive data is locked using memguard.
type Session struct {
	pub *[32]byte
	pri *memguard.LockedBuffer

	shared      *memguard.LockedBuffer
	sharedArray *[32]byte

	nonceSeed *[24]byte
	seq       uint64
	isAlice   bool
}

// NewSession create a session and generates a key pair with their memory
// locked.
func NewSession() (*Session, error) {
	// Init private key
	pri, err := memguard.NewImmutableRandom(32)
	if err != nil {
		return nil, err
	}

	priArray := (*[32]byte)(unsafe.Pointer(&pri.Buffer()[0]))
	pub := new([32]byte)

	// Calculate public key
	curve25519.ScalarBaseMult(pub, priArray)

	return &Session{
		pub: pub,
		pri: pri,
		seq: 0,
	}, nil
}

func (s *Session) PublicKey() *[32]byte {
	return s.pub
}

func (s *Session) Seq() uint64 {
	return s.seq
}

// Compute computes the shared secret and shared nonce seed. On success, the
// private key of the session will be destroyed.
//
// This is how the shared nonce seed is computed on both sides: In big endian,
// compare Alice's and Bob's public keys. In the case that Alice's public key is
// greater, then
//     nonceSeed := SHAKE128(AlicePub + BobPub)[:24]
// else
//     nonceSeed := SHAKE128(BobPub + AlicePub)[:24]
func (s *Session) Compute(pub *[32]byte) error {
	if s.seq != 0 {
		return errors.New("compute: already have shared secret")
	}

	// We can't let them the same because the calculation of shared nonce seed
	// depends on the difference. It would be a serious BIG FAIL if two public
	// keys are the same, even though the ScalarMult doesn't really care about
	// it.
	if bytes.Compare(pub[:], s.pub[:]) == 0 {
		return errors.New("compute: two public keys are the same")
	}

	// Compute shared secret
	shared, err := memguard.NewMutable(32)
	if err != nil {
		return err
	}

	s.shared = shared
	s.sharedArray = (*[32]byte)(unsafe.Pointer(&s.shared.Buffer()[0]))
	priArray := (*[32]byte)(unsafe.Pointer(&s.pri.Buffer()[0]))

	curve25519.ScalarMult(s.sharedArray, priArray, pub)
	salsa.HSalsa20(s.sharedArray, &zeros, s.sharedArray, &salsa.Sigma)

	if err := s.shared.MakeImmutable(); err != nil {
		s.shared.Destroy()
		s.shared = nil
		s.sharedArray = nil
		return errors.New("compute: " + err.Error())
	}

	// Destroy private key
	s.pri.Destroy()
	s.pri = nil

	// Compute shared nonce seed.
	seeder := make([]byte, 64)
	s.nonceSeed = new([24]byte)

	// Decide who is Alice.
L1:
	for i, v := range s.pub {
		switch {
		case v > pub[i]:
			s.isAlice = true
			copy(seeder, s.pub[:])
			copy(seeder[32:], pub[:])
			break L1

		case v < pub[i]:
			s.isAlice = false
			copy(seeder, pub[:])
			copy(seeder[32:], s.pub[:])
			break L1
		}
	}

	// Here we go
	sha3.ShakeSum128(s.nonceSeed[:], seeder)
	s.seq = 1

	return nil
}

// Seal encrypts and authenticates the plain text.
// On success, the plain text will be destroyed and return the header+ciphertext.
func (s *Session) Seal(plain *memguard.LockedBuffer) ([]byte, error) {
	if s.seq == 0 {
		return nil, errors.New("seal: no shared key")
	}

	defer plain.Destroy()

	// Generate nonce
	nonce := computeNounce(s.nonceSeed, s.seq, s.isAlice)

	// Encode header
	header := leb128.AppendUleb128(nil, s.seq)
	s.seq++

	// Seal
	payload := secretbox.Seal(header, plain.Buffer(), nonce, s.sharedArray)

	return payload, nil
}

// Open authenticates and decrypts a message.
func (s *Session) Open(payload []byte) (*memguard.LockedBuffer, error) {
	if s.seq == 0 {
		return nil, errors.New("open: no shared key")
	}

	// Strip seq header
	seq, n := leb128.DecodeUleb128(payload)
	if n == 0 || int(n) >= len(payload) {
		return nil, errors.New("open: bad seq header")
	}

	// Compute nonce
	nonce := computeNounce(s.nonceSeed, seq, !s.isAlice)

	raw, ok := secretbox.Open(nil, payload[n:], nonce, s.sharedArray)
	if !ok {
		return nil, errors.New("open: authentication failed")
	}

	plain, err := memguard.NewImmutableFromBytes(raw)
	if err != nil {
		return nil, errors.New("open: " + err.Error())
	}

	return plain, nil
}

// For Alice:
//     // little endian
//     nonce := nonceSeed[:8] XOR seq + nonceSeed[8:]
// For Bob:
//     // big endian
//     nonce := nonceSeed[:16] + nonceSeed[16:] XOR seq
//
// An example when Alice and Bob both have a seq number of 233,333,333,333
// (decimal):
//                    0           4           8            16          20
//     nonce seed   : 50 09 7e a0 48 ef 43 db a5 24 ... 31 34 aa 7e 91 d3 0c 6e 40
//     Alice's seq  : 55 1d c0 53 36 00 00 00
//     Alice's nonce: 05 14 be f3 7e ef 43 db a5 24 ... 31 34 aa 7e 91 d3 0c 6e 40
//     Bob's seq    :                                      00 00 00 36 53 c0 1d 55
//     Bob's nonce  : 50 09 7e a0 48 ef 43 db a5 24 ... 31 34 aa 7e a7 80 cc 73 15
//
// The caller must increase seq after the call.
func computeNounce(nonceSeed *[24]byte, seq uint64, isAlice bool) *[24]byte {
	buf := make([]byte, 8)
	nonce := new([24]byte)

	if isAlice {
		binary.LittleEndian.PutUint64(buf, seq)
		bytesutil.XorBytes(nonce[:8], nonceSeed[:8], buf)
		copy(nonce[8:], nonceSeed[8:])
	} else {
		binary.BigEndian.PutUint64(buf, seq)
		bytesutil.XorBytes(nonce[16:], nonceSeed[16:], buf)
		copy(nonce[:16], nonceSeed[:16])
	}

	return nonce
}
