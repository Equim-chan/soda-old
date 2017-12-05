package main

import (
	"errors"

	"ekyu.moe/soda/codec"
	"ekyu.moe/soda/packager"
)

func decrypt() error {
	// Prompt input method
	informln("For the encrypted text:")
	read, err := promptInputReader()
	if err != nil {
		return err
	}

	// Prompt output method
	informln("For the plain text:")
	write, err := promptOutputWriter()
	if err != nil {
		return err
	}

	// Read payload
	payloadStr, err := read()
	if err != nil {
		return err
	}

	// Decode payload
	payload := codec.DetectCodecAndDecode(string(payloadStr))

	// Validate length (4 crc32 + 24 nonce)
	if len(payload) <= 28 {
		return errors.New("wrong payload size")
	}

	// Check and detach crc32
	encrypted, ok := packager.DetachCrc32(payload)
	if !ok {
		return errors.New("crc32 checksum failed")
	}

	// Open it
	packet, err := session.Open(encrypted)
	if err != nil {
		return err
	}
	defer packet.Destroy()

	// Unpack packet
	plain, err := packager.Unpack(packet)
	if err != nil {
		return err
	}
	defer plain.Destroy()

	return write(plain.Buffer())
}
