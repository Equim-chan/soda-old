package main

import (
	"encoding/hex"

	"github.com/awnumar/memguard"
)

func uuidv4() error {
	// Prompt output method
	write, err := promptOutputWriter()
	if err != nil {
		return err
	}

	uuid, err := memguard.NewMutableRandom(16)
	if err != nil {
		return err
	}
	defer uuid.Destroy()

	src := uuid.Buffer()
	// Per 4.4, set bits for version and `clock_seq_hi_and_reserved`
	src[6] = (src[6] & 0x0f) | 0x40
	src[8] = (src[8] & 0x3f) | 0x80
	if err := uuid.MakeImmutable(); err != nil {
		return err
	}

	ascii, err := memguard.NewMutable(36) // 16 * 2 + 4
	if err != nil {
		return err
	}
	defer ascii.Destroy()

	// no extra memory copy, no string
	dst := ascii.Buffer()
	dst[8] = '-'
	dst[13] = '-'
	dst[18] = '-'
	dst[23] = '-'

	hex.Encode(dst[:8], src[:4])
	hex.Encode(dst[9:13], src[4:6])
	hex.Encode(dst[14:18], src[6:8])
	hex.Encode(dst[19:23], src[8:10])
	hex.Encode(dst[24:], src[10:])

	uuid.Destroy()
	if err := ascii.MakeImmutable(); err != nil {
		return err
	}

	return write(dst)
}
