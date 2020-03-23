package logger

import (
	"github.com/0xb10c/memo/config"
)

const prefix = "\033["
const resetAll = prefix + "0m"
const resetFontColor = prefix + "39m"

func colorIsEnabled() bool {
	return config.GetBool("log.colorizeOutput")
}

// Red colorizes the input red
func Red(input string) string {
	if colorIsEnabled() {
		return prefix + "31m" + input + resetFontColor
	}
	return input
}

// Blue colorizes the input red
func Blue(input string) string {
	if colorIsEnabled() {
		return prefix + "36m" + input + resetFontColor
	}
	return input
}

// Yellow colorizes the input red
func Yellow(input string) string {
	if colorIsEnabled() {
		return prefix + "93m" + input + resetFontColor
	}
	return input
}

// Dim colorizes the input red
func Dim(input string) string {
	if colorIsEnabled() {
		return prefix + "2m" + input + prefix + "22m"
	}
	return input
}
