package convey

import (
	"strings"

	"github.com/mattn/go-colorable"
)

var (
	stdout = colorable.NewColorableStdout()

	warningBanner = []byte("\x1b[0;31m" + strings.Repeat("+", 60) + "\x1b[0m\n")
	dimBegin      = []byte("\x1b[0;90m")
	dimEnd        = []byte("\x1b[0m\n")
	newLine       = []byte("\n")
)

// Terminal is write-only.
func TerminalWrite(text []byte) error {
	if _, err := stdout.Write(warningBanner); err != nil {
		return err
	}
	if _, err := stdout.Write(dimBegin); err != nil {
		return err
	}

	if _, err := stdout.Write(text); err != nil {
		return err
	}

	if _, err := stdout.Write(dimEnd); err != nil {
		return err
	}
	if _, err := stdout.Write(warningBanner); err != nil {
		return err
	}

	if _, err := stdout.Write(newLine); err != nil {
		return err
	}

	return nil
}
