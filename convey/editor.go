package convey

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

var (
	bom    = []byte{0xef, 0xbb, 0xbf}
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

func EditorWrite(text []byte) error {
	f, err := ioutil.TempFile("", "soda")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(bom); err != nil {
		return err
	}
	if _, err := f.Write(text); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// TODO: check
	return cmd.Run()
}

func EditorRead() ([]byte, error) {
	f, err := ioutil.TempFile("", "soda")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(bom); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Reopen the file
	text, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(text, bom) {
		text = text[3:]
	}

	return text, nil
}
