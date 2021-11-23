package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/mgutz/ansi"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	args := os.Args[1:]

	commitMessage := ""

	if len(args) > 0 {
		commitMessage = strings.TrimRight(strings.TrimSpace(args[0]), "\n")
	}

	if len(args) == 0 {
		newCommit, err := OpenEditor(commitMessage)
		handleCommandError(err)
		commitMessage = strings.TrimRight(strings.TrimSpace(newCommit), "\n")
	}

	if len(commitMessage) == 0 {
		PrintError("Please provide a valid commit message.")
		os.Exit(1)
	}

	status, err := GitStatus()
	handleCommandError(err)


	if len(status) == 0 {
		PrintWarn("Clean working tree. Nothing to commit.")
		return
	}

	execute(status, commitMessage)
}

func showUsage() {
	fmt.Printf("Press %s to continue or %s to abort.\n", BoldText("[ENTER]"), BoldText("[ESC]"))
	fmt.Println(DimText("Press [d] to run diff"))
	fmt.Printf("%s %s\n", DimText("Press [e] to edit the"), ansi.Color("commit message", "yellow+d"))
	fmt.Println(DimText("Press [q] to quit"))
}

func showStatus(status []string, commitMessage string) {
	fmt.Println()
	fmt.Printf("Committing the following files with: [ %s ]\n", ansi.Color(commitMessage, "yellow+hbu"))
	fmt.Println(DimText("-----------------"))
	for _, s := range status {
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
	fmt.Println(DimText("-----------------"))
	fmt.Println()
}

func handleCommandError(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
	PrintError("An error is occurred.")
	//os.Exit(1)
}

func prompt(status []string, commitMessage string) {
	handleCommandError(Clear())
	showStatus(status, commitMessage)
	showUsage()
}

func execute(status []string, commitMessage string) {
	const QChar = 113
	const EChar = 101
	const DChar = 100

	prompt(status, commitMessage)

	char, key, err := keyboard.GetSingleKey()

	if err != nil {
		panic(err.Error())
	}

	switch key {
	case keyboard.KeyEnter: // ENTER
		handleCommandError(GitAddAll())
		fmt.Println()
		handleCommandError(GitCommit(commitMessage))
		os.Exit(0)
	case keyboard.KeyEsc:
		fmt.Println()
		PrintError("Operation Aborted.")
		os.Exit(0)
	case keyboard.KeyCtrlC:
		fmt.Println()
		PrintError("Operation Aborted.")
		os.Exit(0)
	default:
		switch char {
		case QChar:
			fmt.Println()
			PrintError("Operation Aborted.")
			os.Exit(0)
		case EChar:
			newCommit, err := OpenEditor(commitMessage)
			handleCommandError(err)

			commitMessage = newCommit
		case DChar:
			handleCommandError(GitDiff())
			break
		default:
			execute(status, commitMessage)
			break
		}
		break
	}
}
