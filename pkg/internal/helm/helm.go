package helm

import (
	"fmt"
	"os/exec"

	"github.com/sh31k30ps/gikopsctl/pkg/tools"
)

func getCmdArgs() (string, []string, error) {
	tool, err := tools.GetToolResolver().GetTool("helm")
	if err != nil {
		return "", nil, fmt.Errorf("helm is not installed or accessible: %w", err)
	}
	return tool.ResolvedName, tool.GetCmdArgs(), nil
}

func GetDefaults(chart string, version string) ([]byte, error) {
	c, args, err := getCmdArgs()
	if err != nil {
		return nil, fmt.Errorf("failed to get helm command: %w", err)
	}
	args = append(args, "show", "values", chart)
	if version != "" {
		args = append(args, "--version", version)
	}
	cmd := exec.Command(c, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get helm values: %w", err)
	}
	return output, nil
}

func GetTemplates(name, chart, values, version string, options ...string) ([]byte, error) {
	c, args, err := getCmdArgs()
	if err != nil {
		return nil, fmt.Errorf("failed to get helm command: %w", err)
	}
	args = append(args, "template", name, chart)
	if values != "" {
		args = append(args, "--values", values)
	}
	if version != "" {
		args = append(args, "--version", version)
	}
	args = append(args, options...)
	cmd := exec.Command(c, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get helm templates: %w", err)
	}
	return output, nil
}
