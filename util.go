package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GitDiff() error {
	return RunCommand("git", "diff")
}

func GitAddAll() error {
	return RunCommand("git", "add", "-A", ".")
}

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func GitCommit(message string) error {
	return RunCommand("git", "commit", "-n", "-a", "-m", message)
}

func GitStatus() ([]string, error) {
	colorui, err := GitConfig("color.ui")

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

	out := string(b)

	output := strings.Split(out, "\n")

	if err := GitConfigSet("color.ui", strings.TrimSpace(strings.Trim(colorui, "\n"))); err != nil {
		return nil, err
	}

	fmt.Println(colorui)

	return output, nil
}

func Clear() error {
	return RunCommand("clear")
}

func writeFile(filename string, content string) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		return err
	}

	return nil
}


func OpenEditor(content string) (string, error) {
	const FileName = ".commit"

	if e := writeFile(FileName, content); e != nil {
		return "", e
	}

	err := RunCommand("vim",  FileName)

	if err != nil {
		return content, err
	}

	data, err := os.ReadFile(FileName)

	if err != nil {
		return content, err
	}

	if err = os.Remove(FileName); err != nil {
		return "", err
	}

	return string(data), nil
}

func GitConfigSet(key string, value string) error {
	return RunCommand("git", "config", key, value)
}

func GitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", key).Output()

	if err != nil {
		return "", err
	}

	return string(out), err
}