package project

import (
	"errors"

	"github.com/sh31k30ps/gikops/pkg/config/cluster"
)

func Validate(p Project) []error {
	var errs []error
	if p.Name == "" {
		errs = append(errs, errors.New("project name is required"))
	}
	if len(p.Clusters) == 0 {
		errs = append(errs, errors.New("a project cluster must be defined"))
	}

	errs = append(errs, cluster.ValidateAll(p.Clusters)...)

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
