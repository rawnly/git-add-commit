package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/mgutz/ansi"
	"github.com/rawnly/git-add-commit/git"
	"github.com/rawnly/git-add-commit/term"
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
		newCommit, err := term.OpenEditor(commitMessage)
		handleCommandError(err)
		commitMessage = strings.TrimRight(strings.TrimSpace(newCommit), "\n")
	}

	if len(commitMessage) == 0 {
		printError("Please provide a valid commit message.")
		os.Exit(1)
	}

	status, err := git.Status()
	handleCommandError(err)


	if len(status) == 0 {
		printWarn("Clean working tree. Nothing to commit.")
		return
	}

	execute(status, commitMessage)
}

func showUsage() {
	fmt.Printf("Press %s to continue or %s to abort.\n", boldText("[ENTER]"), boldText("[ESC]"))
	fmt.Printf("Press %s to continue and %s.\n", boldText("[p]"), boldText("PUSH"))
	fmt.Println()
	fmt.Println(dimText("Press [d] to run diff"))
	fmt.Printf("%s %s\n", dimText("Press [e] to edit the"), ansi.Color("commit message", "yellow+d"))
	fmt.Println(dimText("Press [q] to quit"))
}

func showStatus(status []string, commitMessage string) {
	fmt.Println()
	fmt.Printf("Committing the following files with: [ %s ]\n", ansi.Color(commitMessage, "yellow+hbu"))
	fmt.Println(dimText("-----------------"))
	for _, s := range status {
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
	fmt.Println(dimText("-----------------"))
	fmt.Println()
}

func handleCommandError(err error) {
	if err == nil {
		return
	}

	printError(fmt.Sprintf("An error is occurred: %s", err.Error()))
}

func prompt(status []string, commitMessage string) {
	handleCommandError(term.Clear())
	showStatus(status, commitMessage)
	showUsage()
}

func execute(status []string, commitMessage string) {
	const QChar = 113
	const EChar = 101
	const DChar = 100
	const PChar = 112

	prompt(status, commitMessage)

	char, key, err := keyboard.GetSingleKey()

	if err != nil {
		panic(err.Error())
	}

	switch key {
	case keyboard.KeyEnter: // ENTER
		handleCommandError(git.AddAll())
		fmt.Println()
		handleCommandError(git.Commit(commitMessage))
		os.Exit(0)
	case keyboard.KeyEsc:
		fmt.Println()
		printError("Operation Aborted.")
		os.Exit(0)
	case keyboard.KeyCtrlC:
		fmt.Println()
		printError("Operation Aborted.")
		os.Exit(0)
	default:
		switch char {
		case QChar:
			fmt.Println()
			printError("Operation Aborted.")
			os.Exit(0)
		case EChar:
			newCommit, err := term.OpenEditor(commitMessage)
			handleCommandError(err)

			commitMessage = newCommit
		case DChar:
			handleCommandError(git.Diff())
			break
		case PChar:
			handleCommandError(git.AddAll())
			fmt.Println()
			handleCommandError(git.Commit(commitMessage))

			output, err := git.Push(git.CurrentBranch())
			handleCommandError(err)
			fmt.Println(output)
			os.Exit(0)
		default:
			fmt.Println(char)
			execute(status, commitMessage)
			break
		}
		break
	}
}



func printWarn(text string) {
	content := ansi.Color(text, "black:yellow+h")
	log.SetPrefix(ansi.Color("WARNING: ", "black:yellow+h"))
	log.Println(content)
}

func printError(text string) {
	log.SetPrefix(ansi.Color("ERROR: ", "white:red"))
	log.Fatal(ansi.Color(text, "white:red"))
}

func boldText(text string) string {
	return ansi.Color(text, "default+b")
}

func dimText(text string) string {
	return ansi.Color(text, "default+d")
}
