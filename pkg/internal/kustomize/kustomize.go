package kustomize

import (
	"fmt"
	"os/exec"
)

func Build(name string) error {
	cmd := exec.Command("kustomize", "build", "-o", "computed.yaml")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build: %s : %w", string(output), err)
	}
	return nil
}
