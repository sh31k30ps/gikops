package kubectl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

func ChangeContext(context string) error {
	config, err := services.GetCurrentProject()
	if err != nil {
		return fmt.Errorf("failed to get current project: %w", err)
	}

	cluster := config.GetCluster(context)
	if cluster == nil {
		return fmt.Errorf("cluster %s not found", context)
	}
	cmd := exec.Command("kubectl", "config", "use-context", cluster.GetContext())
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to change context: %w", err)
	}
	return nil
}

func ChangeNamespace(namespace string) error {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to change namespace: %w", err)
	}
	return nil
}

func CreateNamespace(namespace string) error {
	cmd := exec.Command("kubectl", "create", "namespace", namespace)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}
	return nil
}

func NamespaceExists(namespace string) (bool, error) {
	cmd := exec.Command("kubectl", "get", "namespace", namespace)
	if output, err := cmd.CombinedOutput(); err != nil {
		if strings.Contains(string(output), "not found") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if namespace exists : %s : %w", string(output), err)
	}
	return true, nil
}

func CreateCRDs(crds string) error {
	cmd := exec.Command("kubectl", "create", "-f", crds)
	if output, err := cmd.CombinedOutput(); err != nil {
		if strings.Contains(string(output), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to create CRDs %s: %s: %w", crds, string(output), err)
	}
	return nil
}

func Apply(file string) error {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	if output, err := cmd.CombinedOutput(); err != nil {
		lines := strings.Split(string(output), "\n")
		return fmt.Errorf("failed to apply: %s : %w", &lines, err)
	}
	return nil
}

func WaittingForResourcesBeReady(resources []string) error {
	args := []string{"rollout", "status", "--timeout=2m"}
	args = append(args, resources...)
	cmd := exec.Command("kubectl", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to wait for resource to be ready: %s: %w", string(output), err)
	}

	return nil
}
