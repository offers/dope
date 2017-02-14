package out

import (
	"fmt"

	"github.com/fatih/color"
)

func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func Info(a ...interface{}) (n int, err error) {
	return colorPrintln(color.FgCyan, a...)
}

func Infof(format string, a ...interface{}) (n int, err error) {
	return colorPrintf(format, color.FgCyan, a...)
}

func Success(a ...interface{}) (n int, err error) {
	return colorPrintln(color.FgGreen, a...)
}

func Successf(format string, a ...interface{}) (n int, err error) {
	return colorPrintf(format, color.FgGreen, a...)
}

func Notice(a ...interface{}) (n int, err error) {
	return colorPrintln(color.FgYellow, a...)
}

func Error(a ...interface{}) (n int, err error) {
	return colorPrintln(color.FgRed, a...)
}

func colorPrintln(c color.Attribute, a ...interface{}) (n int, err error) {
	color.Set(c)
	defer color.Unset()
	return fmt.Println(a...)
}

func colorPrintf(format string, c color.Attribute, a ...interface{}) (n int, err error) {
	color.Set(c)
	defer color.Unset()
	return fmt.Printf(format, a...)
}
