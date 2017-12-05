package main

import (
	"errors"

	"github.com/awnumar/memguard"

	"ekyu.moe/soda/packager"
)

func encrypt() error {
	// Prompt input method
	informln("For the plain text:")
	read, err := promptInputReader()
	if err != nil {
		return err
	}

	// Prompt output codec
	informln("For the encrypted text:")
	encode, err := promptOutputCodec()
	if err != nil {
		return err
	}

	// Prompt output method
	write, err := promptOutputWriter()
	if err != nil {
		return err
	}

	// Read plain text
	raw, err := read()
	if err != nil {
		return err
	}

	// Validate length
	if len(raw) == 0 {
		return errors.New("plain text cannot be empty")
	}

	plain, err := memguard.NewImmutableFromBytes(raw)
	if err != nil {
		return err
	}
	defer plain.Destroy()

	// Pack it
	// It will try to compress the plain text
	// and the packet will be destroyed after packing
	packet, err := packager.Pack(plain)
	if err != nil {
		return err
	}
	defer packet.Destroy()

	// Seal
	// The packet will be destroyed after sealing
	encrypted, err := session.Seal(packet)
	if err != nil {
		return err
	}

	// Attach crc32
	payload := packager.AttachCrc32(encrypted)

	// Encode the packet
	payloadStr := encode(payload)

	// Output the payload
	return write([]byte(payloadStr))
}
