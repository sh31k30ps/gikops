package kustomize

import (
	"fmt"
	"os/exec"

	"github.com/sh31k30ps/gikopsctl/pkg/tools"
)

func getCmdArgs() (string, []string, error) {
	tool, err := tools.GetToolResolver().GetTool("kustomize")
	if err != nil {
		return "", nil, fmt.Errorf("kustomize is not installed or accessible: %w", err)
	}
	return tool.ResolvedName, tool.GetCmdArgs(), nil
}

func Build(name string) error {
	c, args, err := getCmdArgs()
	if err != nil {
		return fmt.Errorf("failed to get kustomize command: %w", err)
	}

	cmd := exec.Command(c, append(args, "build", "-o", "computed.yaml")...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build: %s : %w", string(output), err)
	}
	return nil
}
