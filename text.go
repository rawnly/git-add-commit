package main

import (
	"github.com/mgutz/ansi"
	"log"
)

func PrintWarn(text string) {
	content := ansi.Color(text, "black:yellow+h")
	log.SetPrefix(ansi.Color("WARNING: ", "black:yellow+h"))
	log.Println(content)
}

func PrintError(text string) {
	log.SetPrefix(ansi.Color("ERROR: ", "white:red"))
	log.Fatal(ansi.Color(text, "white:red"))
}

func BoldText(text string) string {
	return ansi.Color(text, "default+b")
}

func DimText(text string) string {
	return ansi.Color(text, "default+d")
}