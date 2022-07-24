package git

import (
	"github.com/rawnly/git-add-commit/term"
	"github.com/rawnly/gitgud/git"
	"github.com/rawnly/gitgud/run"
	"strings"
)

// Commit `git commit -n -a -m [message]`
func Commit(message string) error {
	cmd := run.NewGitBuilder("commit").
		BoolFlag("-n").
		BoolFlag("-a").
		StringFlag("-m", message).
		Build()

	return cmd.Run()
}

// Push `git push origin ${branch}`
func Push(origin string, branch string) error {
	return term.RunOSCommand("git", "push", origin, branch)
}

// CurrentBranch Get current branch
func CurrentBranch() string {
	var currentBranch string

	b, err := term.RunCommand("git", "branch")

	if err != nil {
		return currentBranch
	}

	output := strings.ReplaceAll(b, " ", "")

	for _, b := range strings.Split(output, "\n") {
		if strings.Contains(b, "*") {
			currentBranch = b
		}
	}

	return strings.ReplaceAll(currentBranch, "*", "")
}

// Status `git status -s`
func Status() ([]string, error) {
	colorUiConfig, err := git.Config("color.ui").Output()

	if err != nil {
		if err := git.SetConfig("color.ui", "auto").Run(); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := git.SetConfig("color.ui", "always").Run(); err != nil {
		return nil, err
	}

	b, err := git.Status(&git.StatusOptions{Short: true}).Output()

	if err != nil {
		return nil, err
	}

	if err := git.SetConfig("color.ui", strings.TrimSpace(strings.Trim(string(colorUiConfig), "\n"))).Run(); err != nil {
		return nil, err
	}

	output := term.RemoveEmptyStrings(strings.Split(string(b), "\n"))

	return output, nil
}

// Diff `git diff`
func Diff() error {
	return run.Git("diff").RunInTerminal()
}

// AddAll `git add -A .`
func AddAll() error {
	return run.Git("add", "-A", ".").RunInTerminal()
}
