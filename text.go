package main

import (
	"fmt"
	"github.com/mgutz/ansi"
)

func PrintWarn(text string) {
	fmt.Println(ansi.Color(text, "black:yellow+h"))
}

func PrintError(text string) {
	fmt.Println(ansi.Color(text, "white:red"))
}

func BoldText(text string) string {
	return ansi.Color(text, "default+b")
}

func DimText(text string) string {
	return ansi.Color(text, "default+d")
}