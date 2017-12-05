package main

import (
	"fmt"
	"strings"

	"ekyu.moe/base91"
	"github.com/atotto/clipboard"
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"

	"ekyu.moe/soda/codec"
	"ekyu.moe/soda/convey"
	"ekyu.moe/soda/i18n"
)

const (
	CMD_ENC = iota
	CMD_DEC
	CMD_RAND
	CMD_CLS
	CMD_EXIT
)

func init() {
	surveyCore.SelectFocusIcon = ">"
	surveyCore.HelpIcon = ""
}

func promptLocale() (i18n.Locale, error) {
	question := &survey.Select{
		Message: fmt.Sprintf("soda %s at %s build %s", Version, GitHash, BuildDate),
		Options: []string{"English", "日本語", "中文 (繁體)", "中文 (简体)"},
	}

	l := ""
	if err := survey.AskOne(question, &l, nil); err != nil {
		return i18n.EN_US, err
	}

	switch l {
	case "日本語":
		return i18n.JA, nil
	case "中文 (繁體)":
		return i18n.ZH_TW, nil
	case "中文 (简体)":
		return i18n.ZH_CN, nil
	case "English":
		fallthrough
	default:
		return i18n.EN_US, nil
	}
}

func promptKey() (*[32]byte, error) {
	question := &survey.Input{
		Message: i18n.INPUT_PUB,
		Help:    i18n.INPUT_PUB_HELP,
	}

	encodedKey := ""
	if err := survey.AskOne(question, &encodedKey, pubValidator); err != nil {
		return nil, err
	}

	var key [32]byte
	// strip the first 4 bytes of crc32
	copy(key[:], base91.DecodeString(encodedKey)[4:])

	return &key, nil
}

func promptCmd() (int, error) {
	question := &survey.Select{
		Message: i18n.PROMPT_CMD,
		Options: []string{
			i18n.PROMPT_CMD_ENC,
			i18n.PROMPT_CMD_DEC,
			i18n.PROMPT_CMD_CLS,
			i18n.PROMPT_CMD_RAND,
			i18n.PROMPT_CMD_EXIT,
		},
		Help: i18n.PROMPT_CMD_HELP,
	}

	action := ""
	if err := survey.AskOne(question, &action, nil); err != nil {
		return -1, err
	}

	switch action {
	case i18n.PROMPT_CMD_ENC:
		return CMD_ENC, nil
	case i18n.PROMPT_CMD_DEC:
		return CMD_DEC, nil
	case i18n.PROMPT_CMD_CLS:
		return CMD_CLS, nil
	case i18n.PROMPT_CMD_RAND:
		return CMD_RAND, nil
	case i18n.PROMPT_CMD_EXIT:
		fallthrough
	default:
		return CMD_EXIT, nil
	}
}

func promptPlain() (string, error) {
	question := &survey.Editor{
		Message: i18n.PROMPT_PLAIN,
	}

	plain := ""
	if err := survey.AskOne(question, &plain, plainValidator); err != nil {
		return "", err
	}

	return strings.TrimSpace(plain), nil
}

func promptEncrypted() ([]byte, error) {
	question := &survey.Editor{
		Message: i18n.PROMPT_ENCRYPTED,
	}

	encrypted := ""
	if err := survey.AskOne(question, &encrypted, encryptedValidator); err != nil {
		return nil, err
	}

	// strip the first 4 bytes of crc32
	return base91.DecodeString(encrypted)[4:], nil
}

func promptOutputWriter() (convey.WriteFunc, error) {
	question := &survey.Select{
		Message: "Please select your output method",
		Options: []string{"Terminal", "Editor", "Clipboard"},
		// Options: []string{i18n.PROMPT_OUTPUT_EDITOR, i18n.PROMPT_OUTPUT_TERMINAL, "Clipboard"},
		Help: "TODO", // i18n.PROMPT_OUTPUT_HELP,
	}

	writer := ""
	if err := survey.AskOne(question, &writer, nil); err != nil {
		return nil, err
	}

	switch writer {
	case "Terminal":
		return convey.TerminalWrite, nil
	case "Clipboard":
		if clipboard.Unsupported {
			fmt.Println("Sorry but clipboard is not supported on your platform, fallback to editor")
			return convey.EditorWrite, nil
		}
		return convey.ClipboardWrite, nil
	case "Editor":
		fallthrough
	default:
		return convey.EditorWrite, nil
	}
}

func promptOutputCodec() (codec.EncodeFunc, error) {
	question := &survey.Select{
		Message: "Please select your output codec",
		Options: []string{"ASCII", "Emoji", "EmojiTag"},
		Help:    "TODO  (like >OwJh>}A) (like 👾🍧🙆🍬🙇🌱) (like :pizza::sushi::beer:)",
	}

	encode := ""
	if err := survey.AskOne(question, &encode, nil); err != nil {
		return nil, err
	}

	switch encode {
	case "Emoji":
		return codec.EmojiEncode, nil
	case "EmojiTag":
		return codec.EmojiTagEncode, nil
	case "ASCII":
		fallthrough
	default:
		return codec.Base91Encode, nil
	}
}

func promptInputReader() (convey.ReadFunc, error) {
	question := &survey.Select{
		Message: "Please select your input method",
		Options: []string{"Editor", "Clipboard"},
		// Options: []string{i18n.PROMPT_OUTPUT_EDITOR, i18n.PROMPT_OUTPUT_TERMINAL, i18n.PROMPT_OUTPUT_CLIPBOARD},
		Help: "TODO", // i18n.PROMPT_OUTPUT_HELP,
	}

	reader := ""
	if err := survey.AskOne(question, &reader, nil); err != nil {
		return nil, err
	}

	switch reader {
	case "Clipboard":
		if clipboard.Unsupported {
			fmt.Println("Sorry but clipboard is not supported on your platform, fallback to editor")
			return convey.EditorRead, nil
		}
		return convey.ClipboardRead, nil
	case "Editor":
		fallthrough
	default:
		return convey.EditorRead, nil
	}
}
