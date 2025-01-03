package logging

import (
	"fmt"

	"github.com/TwiN/go-color"
)

var Verbosity int = 0

func Log(msg string, verbosity int) {
	if verbosity >= Verbosity {
		fmt.Println(msg)
	}
}

func Info(msg string) {
	fmt.Println(msg)
}

func Debug(msg string) {
	if Verbosity > 0 {
		fmt.Println(msg)
	}
}

func Debug2(msg string) {
	if Verbosity > 1 {
		fmt.Println(msg)
	}
}

func Debug3(msg string) {
	if Verbosity > 2 {
		fmt.Println(msg)
	}
}

func Debug4(msg string) {
	if Verbosity > 3 {
		fmt.Println(msg)
	}
}

func Error(msg string) {
	fmt.Println(color.InRed(msg))
}
