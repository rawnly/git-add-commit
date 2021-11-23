package term

import (
	"os"
	"os/exec"
	"strings"
)

// RunCommand Execute a command ignoring output
func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Clear clears terminal
func Clear() error {
	return RunCommand("clear")
}

// OpenEditor Edit commit message via VIM
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