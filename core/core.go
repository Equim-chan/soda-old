// package core is only for internal use and not thread safe
package core // import "ekyu.moe/soda/core"

import (
	"crypto/rand"
	"errors"
	"fmt"
	"unsafe"

	"github.com/awnumar/memguard"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/salsa20/salsa"
)

var (
	zeros [16]byte
)

type Session interface {
	PublicKey() *[32]byte
	Compute(*[32]byte) error
	Seal(*memguard.LockedBuffer) ([]byte, error)
	Open([]byte) (*memguard.LockedBuffer, error)
}

type session struct {
	pub         *[32]byte
	pri         *memguard.LockedBuffer
	shared      *memguard.LockedBuffer
	sharedArray *[32]byte

	haveShared bool
}

// NewSession create a session and generates a key pair with their memory locked.
func NewSession() (Session, error) {
	s := &session{
		pub:        new([32]byte),
		haveShared: false,
	}

	pri, err := memguard.NewImmutableRandom(32)
	if err != nil {
		return nil, err
	}

	priArray := (*[32]byte)(unsafe.Pointer(&pri.Buffer()[0]))

	curve25519.ScalarBaseMult(s.pub, priArray)

	s.pri = pri
	return s, nil
}

func (s *session) PublicKey() *[32]byte {
	return s.pub
}

// Compute computes the shared secret. On success, the private key of the
// session will be destroyed.
func (s *session) Compute(pub *[32]byte) error {
	if s.haveShared {
		return errors.New("compute shared secret: already have shared secret")
	}

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
		return fmt.Errorf("compute shared secret: %s", err)
	}

	s.pri.Destroy()
	s.pri = nil
	s.haveShared = true

	return nil
}

// Seal encrypts the plain text. On success, the plain text will be
// destroyed and return the nonce+encrypted
func (s *session) Seal(plain *memguard.LockedBuffer) ([]byte, error) {
	defer plain.Destroy()

	nonce := make([]byte, 24)

	// Generate nonce
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	nonceArray := (*[24]byte)(unsafe.Pointer(&nonce[0]))

	// Seal
	payload := secretbox.Seal(nonce, plain.Buffer(), nonceArray, s.sharedArray)
	return payload, nil
}

func (s *session) Open(encrypted []byte) (*memguard.LockedBuffer, error) {
	if len(encrypted) < 24+secretbox.Overhead {
		return nil, errors.New("open: wrong size")
	}

	nonce := (*[24]byte)(unsafe.Pointer(&encrypted[0]))

	raw, ok := secretbox.Open(nil, encrypted[24:], nonce, s.sharedArray)
	if !ok {
		return nil, errors.New("open: authentication failed")
	}

	plain, err := memguard.NewImmutableFromBytes(raw)
	if err != nil {
		return nil, fmt.Errorf("open: %s", err)
	}

	return plain, nil
}
