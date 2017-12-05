package main // import "ekyu.moe/soda"

//go:generate goversioninfo -icon=icon.ico

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
	session core.Session
	seq     uint64 = 1

	Version   = "(dev)"
	BuildDate = "(unknown)"
	GitHash   = "(unknwon)"
)

func main() {
	defer memguard.DestroyAll()

	// Make sure we are in a tty
	if !terminal.IsTerminal(int(os.Stdout.Fd())) || !terminal.IsTerminal(int(os.Stdin.Fd())) {
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
	// Print seq number
	printSeq()

	// Prompt command
	cmd, err := promptCmd()
	if err != nil {
		return false, err
	}

	switch cmd {
	case CMD_ENC:
		{
			// Prompt input method
			informln("For the plain text:")
			read, err := promptInputReader()
			if err != nil {
				return false, err
			}

			// Prompt output codec
			informln("For the encrypted text:")
			encode, err := promptOutputCodec()
			if err != nil {
				return false, err
			}

			// Prompt output method
			write, err := promptOutputWriter()
			if err != nil {
				return false, err
			}

			// Read plain text
			raw, err := read()
			if err != nil {
				return false, err
			}

			// Validate length
			if len(raw) == 0 {
				return false, errors.New("plain text cannot be empty")
			}

			plain, err := memguard.NewImmutableFromBytes(raw)
			if err != nil {
				return false, errors.New("crc32 checksum failed")
			}

			// Pack it
			// It will try to compress the plain text
			// and the packet will be destroyed after packing
			packet, _ := packager.Pack(plain)

			// Seal
			// The packet will be destroyed after sealing
			encrypted, err := session.Seal(packet)
			if err != nil {
				return false, err
			}

			// Attach crc32
			payload := packager.AttachCrc32(encrypted)

			// Encode the packet
			payloadStr := encode(payload)

			// Output the payload
			write([]byte(payloadStr))

			// if output == OUTPUT_TERMINAL {
			// 	fmt.Println(i18n.ENCRYPTED_BELOW)
			// 	colorRed.Println(strings.Repeat("+", 60))
			// 	colorDim.Println(payload)
			// 	colorRed.Println(strings.Repeat("+", 60))
			// } else {
			// 	toggleEditorWithText(payload)
			// }
		}

	case CMD_DEC:
		{
			// Prompt input method
			informln("For the encrypted text:")
			read, err := promptInputReader()
			if err != nil {
				return false, err
			}

			// Prompt output method
			informln("For the plain text:")
			write, err := promptOutputWriter()
			if err != nil {
				return false, err
			}

			// Read payload
			payloadStr, err := read()
			if err != nil {
				return false, err
			}

			// Decode payload
			payload := codec.DetectCodecAndDecode(string(payloadStr))

			// Validate length (4 crc32 + 24 nonce)
			if len(payload) <= 28 {
				return false, errors.New("wrong payload size")
			}

			// Check and detach crc32
			encrypted, ok := packager.DetachCrc32(payload)
			if !ok {

			}

			// Open it
			packet, err := session.Open(encrypted)
			if err != nil {
				return false, err
			}

			// Unpack packet
			plain, err := packager.Unpack(packet)
			if err != nil {
				return false, err
			}

			// Write
			write(plain.Buffer())

			// Destroy the plain text
			plain.Destroy()
		}

	// if output == OUTPUT_TERMINAL {
	// 	fmt.Println(i18n.PLAIN_BELOW)
	// 	colorRed.Println(strings.Repeat("+", 60))
	// 	colorDim.Println(forceLf(plain))
	// 	colorRed.Println(strings.Repeat("+", 60))
	// } else {
	// 	toggleEditorWithText(forceCrlf(plain))
	// }

	case CMD_RAND:
		{
			// Prompt output method
			write, err := promptOutputWriter()
			if err != nil {
				return false, err
			}

			r, err := uuidv4()
			if err != nil {
				return false, err
			}

			write([]byte(r))

			// colorRed.Println(strings.Repeat("+", 60))
			// colorDim.Println(r)
			// colorRed.Println(strings.Repeat("+", 60))
		}

	case CMD_CLS:
		return false, cli.ClearScreen()

	case CMD_EXIT:
		return true, nil
	}

	seq++
	fmt.Println()

	return false, nil
}
