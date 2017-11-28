package main

import (
	"fmt"
	"os"

	surveyTerm "gopkg.in/AlecAivazis/survey.v1/terminal"

	"ekyu.moe/soda/i18n"
)

func gracefulFatal(err error) {
	if err == surveyTerm.InterruptErr {
		// 直接退
		os.Exit(2)
	}
	colorRed.Printf("\n%s：\n%s\n", i18n.EXCEPTION_OCCURRED, err)
	gracefulError(err)
	colorDim.Println(i18n.PRESS_ENTER_TO_EXIT)
	fmt.Scanln()
	os.Exit(2)
}

func gracefulError(err error) {
	if err == surveyTerm.InterruptErr {
		// 不算错误，不处理
		fmt.Println()
		return
	}
	colorRed.Printf("\n%s：\n%s\n", i18n.EXCEPTION_OCCURRED, err)
}
