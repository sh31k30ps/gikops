package project

import (
	"errors"
)

func Validate(p Project) []error {
	var errs []error
	if p.Name == "" {
		errs = append(errs, errors.New("project name is required"))
	}
	if len(p.Clusters) == 0 {
		errs = append(errs, errors.New("a project cluster must be defined"))
	}
	visitedClusters := make(map[string]bool)
	for _, cluster := range p.Clusters {
		errs = append(errs, ValidateCluster(cluster)...)
		if visitedClusters[cluster.Name()] {
			errs = append(errs, errors.New("cluster name must be unique"))
		}
		visitedClusters[cluster.Name()] = true
	}
	visitedComponents := make(map[string]bool)
	for _, component := range p.Components {
		errs = append(errs, ValidateProjectComponent(component)...)
		if visitedComponents[component.Name] {
			errs = append(errs, errors.New("component name must be unique"))
		}
		visitedComponents[component.Name] = true
	}

	return errs
}

func ValidateProjectComponent(c ProjectComponent) []error {
	var errs []error
	if c.Name == "" {
		errs = append(errs, errors.New("a project component name must be defined"))
	}
	return errs
}

func ValidateCluster(c ProjectCluster) []error {
	var errs []error
	if c.Name() == "" {
		errs = append(errs, errors.New("kind cluster name is required"))
	}
	switch c := c.(type) {
	case *KindCluster:
		errs = append(errs, ValidateKindCluster(c)...)
	case *BasicCluster:
		break
	default:
		errs = append(errs, errors.New("unsupported cluster type"))
	}
	return errs
}

func ValidateKindCluster(c *KindCluster) []error {
	var errs []error
	cfg := c.Config().(*KindConfig)
	if cfg == nil {
		errs = append(errs, errors.New("kind config is required"))
		return errs
	}
	errs = append(errs, ValidateKindConfig(cfg)...)
	return errs
}

func ValidateKindConfig(c *KindConfig) []error {
	var errs []error
	if c.ConfigFile == "" {
		errs = append(errs, errors.New("kind config file is required"))
	}
	if c.Provider != KindConfigProviderDocker && c.Provider != KindConfigProviderPodman && c.Provider != KindConfigProviderNerdctl {
		errs = append(errs, errors.New("kind config provider must be docker, podman or nerdctl"))
	}
	return errs
}
