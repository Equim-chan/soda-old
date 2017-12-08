package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"ekyu.moe/util/cli"
	"github.com/awnumar/memguard"
	"golang.org/x/crypto/ssh/terminal"

	"ekyu.moe/soda/codec"
	"ekyu.moe/soda/core"
	"ekyu.moe/soda/i18n"
	"ekyu.moe/soda/packager"
)

var (
	session *core.Session
	id      uint64 = 1
)

func main() {
	defer memguard.DestroyAll()

	// Make sure we are in a tty
	if !terminal.IsTerminal(int(os.Stdout.Fd())) ||
		!terminal.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "soda: soda only works in a tty.")
		memguard.SafeExit(1)
	}

	// Prompt locale
	l, err := promptLocale()
	if err != nil {
		fatal(err)
	}
	i18n.SetLocale(l)

	informln("\nYour key pair is to be generated.")

	// Prompt output codec
	informln("For your own public key:")
	encode, err := promptOutputCodec()
	if err != nil {
		fatal(err)
	}

	// Prompt output method
	write, err := promptOutputWriter()
	if err != nil {
		fatal(err)
	}

	// Generate session (key pair)
	session, err = core.NewSession()
	if err != nil {
		fatal(err)
	}

	// Append crc32 to the head
	packet := packager.AttachCrc32(session.PublicKey()[:])

	// Encode public key
	myPubStr := encode(packet)

	// Output public key
	if err := write([]byte(myPubStr)); err != nil {
		fatal(err)
	}

	for {
		// Prompt input method
		informln("\nFor your partner's public key:")
		read, err := promptInputReader()
		if err != nil {
			// this one is fatal
			fatal(err)
		}

		// Read partner's public key
		hisPubStr, err := read()
		if err != nil {
			perror(err)
			continue
		}

		// Decode public key
		packet := codec.DetectCodecAndDecode(string(hisPubStr))

		// Validate length
		if len(packet) != 36 {
			perror(errors.New("wrong public key size"))
			continue
		}

		// Check crc32
		hisPub, ok := packager.DetachCrc32(packet)
		if !ok {
			perror(errors.New("crc32 checksum failed"))
			continue
		}

		// Compute shared secret
		hisPubArray := (*[32]byte)(unsafe.Pointer(&hisPub[0]))
		if err := session.Compute(hisPubArray); err != nil {
			perror(err)
			continue
		}

		break
	}

	// Session begins
	informf("\n\x1b[1m================= %s =================\x1b[0m\n", i18n.SESSION_BEGIN)

	for {
		quit, err := mainLoop()
		if err != nil {
			perror(err)
		}
		if quit {
			break
		}
	}

	memguard.SafeExit(0)
}

func mainLoop() (bool, error) {
	// Print the ID (how many times mainLoop has been called without error)
	printId()

	// Prompt command
	cmd, err := promptCmd()
	if err != nil {
		return false, err
	}

	switch cmd {
	case CMD_ENC:
		err = encrypt()

	case CMD_DEC:
		err = decrypt()

	case CMD_RAND:
		err = uuidv4()

	case CMD_CLS:
		err = cli.ClearScreen()

	case CMD_EXIT:
		return true, nil
	}

	if err == nil {
		id++
		fmt.Println()
	}

	return false, err
}
