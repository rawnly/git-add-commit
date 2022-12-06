package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/eiannone/keyboard"
	"github.com/mgutz/ansi"
	"github.com/rawnly/git-add-commit/git"
	"github.com/rawnly/git-add-commit/term"
)

var version = "development"

var cli struct {
	Version bool   `help:"Print version" short:"v"`
	Commit  string `arg:"" help:"Your commit message" name:"commit message" optional:""`
	Path    string `arg:"." help:"Path to files" optional:""`
	Remote  string `help:"Specify remote" name:"remote" default:"origin"`
}

func printVersion() {
	fmt.Println(fmt.Sprintf("Version %s", version))
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	_ = kong.Parse(&cli)

	if cli.Version {
		printVersion()
		return
	}

	commitMessage := strings.TrimRight(strings.TrimSpace(cli.Commit), "\n")
	pathspec := strings.TrimRight(strings.TrimSpace(cli.Path), "\n")

	if len(commitMessage) == 0 {
		newCommit, err := term.OpenEditor(commitMessage)
		handleCommandError(err)
		commitMessage = strings.TrimRight(strings.TrimSpace(newCommit), "\n")
	}

	if len(commitMessage) == 0 {
		printError("Please provide a valid commit message.")
		os.Exit(1)
	}

	status, err := git.Status(pathspec)
	handleCommandError(err)

	if len(status) == 0 {
		printWarn("Clean working tree. Nothing to commit.")
		return
	}

	execute(status, commitMessage, pathspec)
}

func showUsage() {
	fmt.Printf("Press %s to continue or %s to abort.\n", boldText("[ENTER]"), boldText("[ESC]"))
	fmt.Printf("Press %s to continue and %s.\n", boldText("[p]"), boldText("PUSH"))
	fmt.Println()
	fmt.Println(dimText("Press [d] to run diff"))
	fmt.Printf(
		"%s %s\n",
		dimText("Press [e] to edit the"),
		ansi.Color("commit message", "yellow+d"),
	)
	fmt.Println(dimText("Press [q] to quit"))
}

func showStatus(status []string, commitMessage string) {
	fmt.Println()
	fmt.Printf(
		"Committing the following files with: [ %s ]\n",
		ansi.Color(commitMessage, "yellow+hbu"),
	)
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
	os.Exit(1)
}

func prompt(status []string, commitMessage string) {
	handleCommandError(term.Clear())
	showStatus(status, commitMessage)
	showUsage()
}

func execute(status []string, commitMessage string, filesPath string) {
	keys := map[string]rune{
		"Q": 113, "E": 101,
		"D": 100, "P": 112,
	}

	prompt(status, commitMessage)

	char, key, err := keyboard.GetSingleKey()

	if err != nil {
		panic(err.Error())
	}

	switch key {
	case keyboard.KeyEnter: // ENTER
		handleCommandError(git.Add(filesPath))
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
		case keys["Q"]:
			fmt.Println()
			printError("Operation Aborted.")
			os.Exit(0)
		case keys["E"]:
			newCommit, err := term.OpenEditor(commitMessage)
			handleCommandError(err)
			_ = term.Clear()
			execute(status, strings.TrimRight(strings.TrimSpace(newCommit), "\n"), filesPath)
			break
		case keys["D"]:
			handleCommandError(git.Diff())
			_ = term.Clear()
			execute(status, commitMessage, filesPath)
			break
		case keys["P"]:
			handleCommandError(git.AddAll())
			fmt.Println()
			handleCommandError(git.Commit(commitMessage))

			err := git.Push(cli.Remote, git.CurrentBranch())
			handleCommandError(err)

			os.Exit(0)
		default:
			_ = term.Clear()
			execute(status, commitMessage, filesPath)
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
