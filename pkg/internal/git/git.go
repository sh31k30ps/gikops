package git

import (
	"fmt"
	"os/exec"

	"github.com/sh31k30ps/gikopsctl/pkg/tools"
)

func Clone(url string, path string) error {
	c, args, err := getCmdArgs()
	if err != nil {
		return fmt.Errorf("failed to get git command: %w", err)
	}

	cmd := exec.Command(c, append(args, "clone", url, path)...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

func Init(path string) error {
	c, args, err := getCmdArgs()
	if err != nil {
		return fmt.Errorf("failed to get git command: %w", err)
	}

	cmd := exec.Command(c, append(args, "init", path)...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	return nil
}

func getCmdArgs() (string, []string, error) {
	tool, err := tools.GetToolResolver().GetTool("git")
	if err != nil {
		return "", nil, fmt.Errorf("git is not installed or accessible: %w", err)
	}
	return tool.ResolvedName, tool.GetCmdArgs(), nil
}
