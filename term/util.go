package term

import (
	"os"
	"os/exec"
	"strings"
)

// RunCommand Execute a command ignoring output
func RunCommand(command string, args ...string) (string, error) {
	stdout, err := exec.Command(command, args...).CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

// Clear clears terminal
func Clear() error {
	cmd := exec.Command("clear")

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunOSCommand(command string, args ...string) error {
	cmd := exec.Command(command,  args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// OpenEditor Edit commit message via VIM
func OpenEditor(content string) (string, error) {
	const FileName = ".commit"

	if e := writeFile(FileName, content); e != nil {
		return "", e
	}

	if err := RunOSCommand("vim", FileName); err != nil {
		return content, err
	}

	data, err := os.ReadFile(FileName)

	if err != nil {
		return content, err
	}

	if err = os.Remove(FileName); err != nil {
		println(err.Error())
		return "", err
	}

	return string(data), nil
}

// RemoveEmptyStrings removes empty strings from an array of strings
func RemoveEmptyStrings(arr []string) []string {
	var items []string

	for _, item := range arr {
		trimmed := strings.TrimRight(strings.TrimSpace(item), "\n")

		if len(trimmed) > 0 {
			items = append(items, trimmed)
		}
	}

	return items
}

// writeFile Writes a file truncating it if exists.
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