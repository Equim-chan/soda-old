// package core is only for internal use and not thread safe
package core // import "ekyu.moe/soda/core"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"

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
//     seq := uint64(1)
//     // little endian
//     nonce := nonceSeed[:8] XOR seq + nonceSeed[8:]
//     seq++
//     ...
// else
//     nonceSeed := SHAKE128(BobPub + AlicePub)[:24]
//     seq := uint64(1)
//     // big endian
//     nonce := nonceSeed[:16] + nonceSeed[16:] XOR seq
//     seq++
//     ...
//
// An example when Alice and Bob both have a seq number of 233,333,333,333
// (decimal):
//                    0           4           8            16          20
//     nonce seed   : 50 09 7e a0 48 ef 43 db a5 24 ... 31 34 aa 7e 91 d3 0c 6e 40
//     Alice's seq  : 55 1d c0 53 36 00 00 00
//     Alice's nonce: 05 14 be f3 7e ef 43 db a5 24 ... 31 34 aa 7e 91 d3 0c 6e 40
//     Bob's seq    :                                      00 00 00 36 53 c0 1d 55
//     Bob's nonce  : 50 09 7e a0 48 ef 43 db a5 24 ... 31 34 aa 7e a7 80 cc 73 15
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

	// log.Println("isAlice:", s.isAlice)

	// Here we go
	sha3.ShakeSum128(s.nonceSeed[:], seeder)
	s.seq = 1

	// log.Printf("nonceSeed: %x\n", *s.nonceSeed)

	return nil
}

func computeNounce(nonceSeed *[24]byte, seq uint64, isAlice bool) *[24]byte {
	buf := make([]byte, 8)
	nonce := new([24]byte)

	if isAlice {
		binary.LittleEndian.PutUint64(buf, seq)
		// log.Printf("buf: %x\n", buf)
		bytesutil.XorBytes(nonce[:8], nonceSeed[:8], buf)
		copy(nonce[8:], nonceSeed[8:])
	} else {
		binary.BigEndian.PutUint64(buf, seq)
		// log.Printf("buf: %x\n", buf)
		copy(nonce[:16], nonceSeed[:16])
		bytesutil.XorBytes(nonce[16:], nonceSeed[16:], buf)
	}

	return nonce
}

// The header's format is a bit like UTF-8.
//
// The followings are all in big endian.
//     L      H
//     xxxxxxx1                                                                // [1, 127]
//     xxxxxx10 xxxxxxxx                                                       // [128, 16383]
//     xxxxx100 xxxxxxxx xxxxxxxx xxxxxxxx                                     // [16384, 1 << 29 - 1]
//     xxxx1000 xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx // [1 << 29, 1 << 60]
//
// TODO(Equim): Use LEB128 instead. The current one is full of bugs.
func (s *Session) generateSeqHeader() []byte {
	switch {
	case s.seq <= 1<<7-1:
		head := make([]byte, 1)
		head[0] = uint8(s.seq&0x7f) | 0x80
		return head

	case s.seq <= 1<<14-1:
		head := make([]byte, 2)
		head[1] = uint8(s.seq & 0xff)
		head[0] = uint8((s.seq>>8)&0x3f) | 0x40
		return head

	case s.seq <= 1<<29-1:
		head := make([]byte, 4)
		head[3] = uint8(s.seq & 0xff)
		head[2] = uint8((s.seq >> 8) & 0xff)
		head[1] = uint8((s.seq >> 16) & 0xff)
		head[0] = uint8((s.seq>>24)&0x1f) | 0x20
		return head

	case s.seq <= 1<<60-1:
		head := make([]byte, 8)
		head[7] = uint8(s.seq & 0xff)
		head[6] = uint8((s.seq >> 8) & 0xff)
		head[5] = uint8((s.seq >> 16) & 0xff)
		head[4] = uint8((s.seq >> 24) & 0xff)
		head[3] = uint8((s.seq >> 32) & 0xff)
		head[2] = uint8((s.seq >> 40) & 0xff)
		head[1] = uint8((s.seq >> 48) & 0xff)
		head[0] = uint8((s.seq>>56)&0x0f) | 0x10
		return head
	}

	// Assert
	panic("seq out of range (1 << 60)")

	return nil
}

// Parse the seq header from an encrypted text. It returns the seq number and
// the length of it in the encrypted text. If it doesn't contain any valid info
// (e.g. e is 0 bytes long), 0 and -1 are returned.
func parseSeqHeader(e []byte) (uint64, int8) {
	// log.Printf("e head 4: %b %b %b %b\n", e[0], e[1], e[2], e[3])

	if len(e) < 1 {
		return 0, -1
	}

	switch {
	case e[0]&0x80 == 0x80:
		seq := uint64(e[0] & 0x7f)
		return seq, 1

	case e[0]&0x40 == 0x40:
		seq := uint64(e[1]) + (uint64(e[0])<<8)&0x3f
		return seq, 2

	case e[0]&0x20 == 0x20:
		seq := uint64(e[3]) +
			uint64(e[2])<<8 +
			uint64(e[1])<<16 +
			(uint64(e[0])<<24)&0x1f
		return seq, 4

	case e[0]&0x10 == 0x10:
		seq := uint64(e[7]) +
			uint64(e[6])<<8 +
			uint64(e[5])<<16 +
			uint64(e[4])<<24 +
			uint64(e[3])<<32 +
			uint64(e[2])<<40 +
			uint64(e[1])<<48 +
			(uint64(e[0])<<56)&0x0f
		return seq, 8
	}

	return 0, -1
}

// Seal encrypts the plain text.
// On success, the plain text will be destroyed and return the header+ciphertext.
func (s *Session) Seal(plain *memguard.LockedBuffer) ([]byte, error) {
	if s.seq == 0 {
		return nil, errors.New("seal: no shared key")
	}
	// TODO: > MAX_SEQ

	defer plain.Destroy()

	// Generate nonce
	nonce := computeNounce(s.nonceSeed, s.seq, s.isAlice)
	// log.Printf("nonce: %x\n", *nonce)

	// Generate header
	header := s.generateSeqHeader()

	s.seq++

	// Seal
	payload := secretbox.Seal(header, plain.Buffer(), nonce, s.sharedArray)

	return payload, nil
}

func (s *Session) Open(encrypted []byte) (*memguard.LockedBuffer, error) {
	if s.seq == 0 {
		return nil, errors.New("open: no shared key")
	}

	// Strip seq header
	seq, n := parseSeqHeader(encrypted)
	if n < 0 {
		return nil, errors.New("open: wrong size of encrypted text")
	}

	// log.Printf("seq: %v\n", seq)

	// Compute nonce
	nonce := computeNounce(s.nonceSeed, seq, !s.isAlice)
	// log.Printf("nonce: %x\n", *nonce)

	raw, ok := secretbox.Open(nil, encrypted[n:], nonce, s.sharedArray)
	if !ok {
		return nil, errors.New("open: authentication failed")
	}

	plain, err := memguard.NewImmutableFromBytes(raw)
	if err != nil {
		return nil, errors.New("open: " + err.Error())
	}

	return plain, nil
}
