package git

import (
	"github.com/rawnly/git-add-commit/term"
	"os/exec"
	"strings"
)

// SetConfig `git config <key>`
func SetConfig(key string, value string) error {
	_, err := term.RunCommand("git", "config", key, value)

	return err
}

// Config `git config <key> <value>`
func Config(key string) (string, error) {
	out, err := exec.Command("git", "config", key).Output()

	if err != nil {
		return "", err
	}

	return string(out), err
}

// Commit `git commit -n -a -m [message]`
func Commit(message string) error {
	_, err := term.RunCommand("git", "commit", "-n", "-a", "-m", message)
	return err
}

// Status `git status -s`
func Status() ([]string, error) {
	colorUiConfig, err := Config("color.ui")

	if err != nil {
		return nil, err
	}

	if err := SetConfig("color.ui", "always"); err != nil {
		return nil, err
	}

	b, err := exec.Command("git", "status", "-s").Output()

	if err != nil {
		return nil, err
	}

	if err := SetConfig("color.ui", strings.TrimSpace(strings.Trim(colorUiConfig, "\n"))); err != nil {
		return nil, err
	}

	output := term.RemoveEmptyStrings(strings.Split(string(b), "\n"))

	return output, nil
}

// Diff `git diff`
func Diff() error {
	return term.RunOSCommand("git", "diff")
}

// AddAll `git add -A .`
func AddAll() error {
	return term.RunOSCommand("git", "add", "-A", ".")
}
