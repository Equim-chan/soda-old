package main

import (
	"fmt"

	"github.com/awnumar/memguard"
	surveyTerm "gopkg.in/AlecAivazis/survey.v1/terminal"

	"ekyu.moe/soda/i18n"
)

func gracefulFatal(err error) {
	memguard.DestroyAll()

	if err == surveyTerm.InterruptErr {
		// 直接退
		memguard.SafeExit(2)
	}

	colorRed.Printf("\n%s：\n%s\n", i18n.EXCEPTION_OCCURRED, err)
	gracefulError(err)
	colorDim.Println(i18n.PRESS_ENTER_TO_EXIT)
	fmt.Scanln()

	memguard.SafeExit(2)
}

func gracefulError(err error) {
	if err == surveyTerm.InterruptErr {
		// 不算错误，不处理
		fmt.Println()
		return
	}
	colorRed.Printf("\n%s：\n%s\n", i18n.EXCEPTION_OCCURRED, err)
}
