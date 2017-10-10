package main

import (
	"strings"

	"ekyu.moe/soda/i18n"

	"ekyu.moe/base91"
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

func init() {
	surveyCore.SelectFocusIcon = ">"
}

func promptLocale() (i18n.Locale, error) {
	question := &survey.Select{
		Message: "",
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

func promptOutput(isEncrypt bool) (int, error) {
	msg := i18n.PROMPT_OUTPUT_PLAIN
	if isEncrypt {
		msg = i18n.PROMPT_OUTPUT_ENCRYPTED
	}
	question := &survey.Select{
		Message: msg,
		Options: []string{i18n.PROMPT_OUTPUT_EDITOR, i18n.PROMPT_OUTPUT_TERMINAL},
		Help:    i18n.PROMPT_OUTPUT_HELP,
	}

	printer := ""
	if err := survey.AskOne(question, &printer, nil); err != nil {
		return -1, err
	}

	if printer == i18n.PROMPT_OUTPUT_TERMINAL {
		return OUTPUT_TERMINAL, nil
	}
	return OUTPUT_EDITOR, nil
}
