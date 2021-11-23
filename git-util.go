package main

import (
	"os/exec"
	"strings"
)

// GitConfigSet `git config <key>`
func GitConfigSet(key string, value string) error {
	return RunCommand("git", "config", key, value)
}

// GitConfig `git config <key> <value>`
func GitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", key).Output()

	if err != nil {
		return "", err
	}

	return string(out), err
}

// GitCommit `git commit -n -a -m [message]`
func GitCommit(message string) error {
	return RunCommand("git", "commit", "-n", "-a", "-m", message)
}

// GitStatus `git status -s`
func GitStatus() ([]string, error) {
	colorUiConfig, err := GitConfig("color.ui")

	if err != nil {
		return nil, err
	}

	if err := GitConfigSet("color.ui", "always"); err != nil {
		return nil, err
	}

	b, err := exec.Command("git", "status", "-s").Output()

	if err != nil {
		return nil, err
	}

	if err := GitConfigSet("color.ui", strings.TrimSpace(strings.Trim(colorUiConfig, "\n"))); err != nil {
		return nil, err
	}

	output := RemoveEmptyStrings(strings.Split(string(b), "\n"))

	return output, nil
}

// GitDiff `git diff`
func GitDiff() error {
	return RunCommand("git", "diff")
}

// GitAddAll `git add -A .`
func GitAddAll() error {
	return RunCommand("git", "add", "-A", ".")
}
