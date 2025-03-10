package component

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type PodsGenerators struct {
	Deployments  []string
	Daemonsets   []string
	Statefulsets []string
}

func (pg *PodsGenerators) GetDeployments() []string {
	resources := []string{}
	for _, deployment := range pg.Deployments {
		resources = append(resources, fmt.Sprintf("deployment/%s", deployment))
	}
	return resources
}

func (pg *PodsGenerators) GetDaemonsets() []string {
	resources := []string{}
	for _, daemonset := range pg.Daemonsets {
		resources = append(resources, fmt.Sprintf("daemonset/%s", daemonset))
	}
	return resources
}

func (pg *PodsGenerators) GetStatefulsets() []string {
	resources := []string{}
	for _, statefulset := range pg.Statefulsets {
		resources = append(resources, fmt.Sprintf("statefulset/%s", statefulset))
	}
	return resources
}

func getPodsGeneratorsFromFile(file string) (*PodsGenerators, error) {
	// Read the computed manifest file
	manifestBytes, err := os.ReadFile("computed.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read computed manifest: %w", err)
	}

	// Split the manifest into individual resources
	resources := strings.Split(string(manifestBytes), "\n---\n")

	// Track names of workload resources
	pg := &PodsGenerators{
		Deployments:  []string{},
		Daemonsets:   []string{},
		Statefulsets: []string{},
	}

	// Parse each resource
	for _, resource := range resources {
		var obj map[string]interface{}
		if err := yaml.Unmarshal([]byte(resource), &obj); err != nil {
			continue
		}

		// Get kind and name
		kind, ok := obj["kind"].(string)
		if !ok {
			continue
		}

		metadata, ok := obj["metadata"].(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := metadata["name"].(string)
		if !ok {
			continue
		}

		// Add name to appropriate slice based on kind
		switch kind {
		case "Deployment":
			pg.Deployments = append(pg.Deployments, name)
		case "DaemonSet":
			pg.Daemonsets = append(pg.Daemonsets, name)
		case "StatefulSet":
			pg.Statefulsets = append(pg.Statefulsets, name)
		}
	}

	return pg, nil
}
