package main

import (
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/mattn/go-colorable"
	surveyTerm "gopkg.in/AlecAivazis/survey.v1/terminal"

	"ekyu.moe/soda/i18n"
)

var (
	stdout = colorable.NewColorableStdout()
)

// Print in green
func informf(format string, a ...interface{}) {
	fmt.Fprint(stdout, "\x1b[32m")
	fmt.Fprintf(stdout, format, a...)
	fmt.Fprint(stdout, "\x1b[0m")
}

// Print in green
func informln(a ...interface{}) {
	fmt.Fprint(stdout, "\x1b[32m")
	fmt.Fprintln(stdout, a...)
	fmt.Fprint(stdout, "\x1b[0m")
}

func printId() {
	fmt.Fprintf(stdout, "\x1b[1;35m[#%d]\x1b[0m\n", id)
}

func fatal(err error) {
	memguard.DestroyAll()

	if err == surveyTerm.InterruptErr {
		// 直接退
		memguard.SafeExit(2)
	}

	fmt.Fprintf(stdout, "\n  \x1b[1;31m%s\x1b[0m\n    \x1b[1;31m%s\x1b[0m\n\n", i18n.EXCEPTION_OCCURRED, err)
	fmt.Fprintln(stdout, "\x1b[90m"+i18n.PRESS_ENTER_TO_EXIT+"\x1b[0m")
	fmt.Scanln()

	memguard.SafeExit(2)
}

func perror(err error) {
	if err == surveyTerm.InterruptErr {
		// 不算错误，不处理
		fmt.Println()
		return
	}
	fmt.Fprintf(stdout, "\n  \x1b[1;31m%s\x1b[0m\n    \x1b[1;31m%s\x1b[0m\n\n", i18n.EXCEPTION_OCCURRED, err)
}
