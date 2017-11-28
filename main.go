package main // import "ekyu.moe/soda"

//go:generate goversioninfo -icon=icon.ico

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"

	"ekyu.moe/base91"
	"ekyu.moe/util/cli"
	"golang.org/x/crypto/nacl/box"

	"ekyu.moe/soda/i18n"
)

const (
	CMD_ENC = iota
	CMD_DEC
	CMD_RAND
	CMD_CLS
	CMD_EXIT
	OUTPUT_TERMINAL
	OUTPUT_EDITOR
)

var (
	shared        = new([32]byte)
	nonce         = new([24]byte)
	seq    uint64 = 1
)

func main() {
	l, err := promptLocale()
	if err != nil {
		gracefulFatal(err)
	}
	i18n.SetLocale(l)

	myPub, myPri, err := box.GenerateKey(rand.Reader)
	if err != nil {
		gracefulFatal(err)
	}

	myPubWithCrc32 := make([]byte, 36)

	c := crc32.Checksum(myPub[:], crc32.IEEETable)
	binary.BigEndian.PutUint32(myPubWithCrc32[:4], c)
	copy(myPubWithCrc32[4:], myPub[:])

	fmt.Printf("%s\n\n    ", i18n.YOUR_PUB)
	colorDim.Println(base91.EncodeToString(myPubWithCrc32))
	fmt.Println()

	hisPub, err := promptKey()
	if err != nil {
		gracefulFatal(err)
	}

	box.Precompute(shared, hisPub, myPri)
	colorGreenBold.Printf("=============== %s ===============\n", i18n.SESSION_BEGIN)

	for {
		quit, err := mainLoop()
		if err != nil {
			gracefulError(err)
		}
		if quit {
			break
		}
	}
}

func mainLoop() (bool, error) {
	// 清空 nonce，以防 compromised
	for i := 0; i < 24; i++ {
		nonce[i] = 0x00
	}

	colorMagentaBold.Printf("[#%d]\n", seq)

	cmd, err := promptCmd()
	if err != nil {
		return false, err
	}

	switch cmd {
	case CMD_ENC:
		plain, err := promptPlain()
		if err != nil {
			return false, err
		}

		payload, err := pack(plain)
		if err != nil {
			return false, err
		}

		output, err := promptOutput(true)
		if err != nil {
			return false, err
		}
		if output == OUTPUT_TERMINAL {
			fmt.Println(i18n.ENCRYPTED_BELOW)
			colorRed.Println(strings.Repeat("+", 60))
			colorDim.Println(payload)
			colorRed.Println(strings.Repeat("+", 60))
		} else {
			toggleEditorWithText(payload)
		}

	case CMD_DEC:
		encrypted, err := promptEncrypted()
		if err != nil {
			return false, err
		}

		// 校验已在 promptEncrypted 完成，拿到的 encrypted 是不含 crc32 的
		plain, err := unpack(encrypted)
		if err != nil {
			return false, err
		}

		output, err := promptOutput(false)
		if err != nil {
			return false, err
		}
		if output == OUTPUT_TERMINAL {
			fmt.Println(i18n.PLAIN_BELOW)
			colorRed.Println(strings.Repeat("+", 60))
			colorDim.Println(forceLf(plain))
			colorRed.Println(strings.Repeat("+", 60))
		} else {
			toggleEditorWithText(forceCrlf(plain))
		}

	case CMD_CLS:
		if err = cli.ClearScreen(); err != nil {
			return false, err
		}

	case CMD_RAND:
		r, err := uuidv4()
		if err != nil {
			return false, err
		}
		colorRed.Println(strings.Repeat("+", 60))
		colorDim.Println(r)
		colorRed.Println(strings.Repeat("+", 60))

	case CMD_EXIT:
		return true, nil
	}

	seq++
	fmt.Println()

	return false, nil
}
