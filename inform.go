package main

import (
	"fmt"

	"github.com/mattn/go-colorable"
	surveyTerm "gopkg.in/AlecAivazis/survey.v1/terminal"

	"ekyu.moe/soda/i18n"
)

var (
	stdout = colorable.NewColorableStdout()

	green = []byte{0x1b, '[', '3', '2', 'm'}
	dim   = []byte{0x1b, '[', '9', '0', 'm'}
	reset = []byte{0x1b, '[', '0', 'm'}
)

// Print in green
func informf(format string, a ...interface{}) {
	stdout.Write(green)
	defer stdout.Write(reset)

	fmt.Fprintf(stdout, format, a...)
}

// Print in green
func informln(a ...interface{}) {
	stdout.Write(green)
	defer stdout.Write(reset)

	fmt.Fprintln(stdout, a...)
}

// Print in dim
func hintf(format string, a ...interface{}) {
	stdout.Write(dim)
	defer stdout.Write(reset)

	fmt.Fprintf(stdout, format, a...)
}

// Special case
func printID() {
	fmt.Fprintf(stdout, "\x1b[1;35m[#%d]\x1b[0m\n", id)
}

func perror(err error) {
	if err == surveyTerm.InterruptErr {
		fmt.Frintln(stdout)
	} else {
		fmt.Fprintf(stdout, "\n  \x1b[1;31m%s\x1b[0m\n    \x1b[1;31m%s\x1b[0m\n\n", i18n.EXCEPTION_OCCURRED, err)
	}
}
