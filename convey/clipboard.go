package convey

import (
	"github.com/atotto/clipboard"
)

func ClipboardWrite(text []byte) error {
	return clipboard.WriteAll(string(text))
}

func ClipboardRead() ([]byte, error) {
	str, err := clipboard.ReadAll()
	return []byte(str), err
}
