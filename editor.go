package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	BOM    = []byte{0xef, 0xbb, 0xbf}
	editor = "vim"
)

func init() {
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}
}

// 为了保持记事本的兼容性，这个函数可以将 LF 换成 CRLF
func forceCrlf(s string) string {
	lfText := strings.Replace(s, "\r\n", "\n", -1)
	return strings.Replace(lfText, "\n", "\r\n", -1)
}

// 为了保持终端的兼容性，这个函数可以将 CRLF 换成 LF
func forceLf(s string) string {
	crlfText := strings.Replace(s, "\n", "\r\n", -1)
	return strings.Replace(crlfText, "\r\n", "\n", -1)
}

func toggleEditorWithText(text string) error {
	f, err := ioutil.TempFile("", "soda")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	body := append(BOM, []byte(text)...)
	if _, err := f.Write(body); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
