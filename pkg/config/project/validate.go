package project

import (
	"errors"
)

func ValidateProject(p Project) []error {
	var errs []error
	if p.Name == "" {
		errs = append(errs, errors.New("project name is required"))
	}
	if p.Components == nil {
		errs = append(errs, errors.New("project components are required"))
	} else {
		errs = append(errs, ValidateProjectComponents(*p.Components)...)
	}
	if p.LocalCluster == nil {
		errs = append(errs, errors.New("project local cluster is required"))
	} else {
		errs = append(errs, ValidateKindCluster(*p.LocalCluster)...)
	}
	return errs
}

func ValidateProjectComponents(c ProjectComponents) []error {
	var errs []error
	if len(c.Folders) == 0 {
		errs = append(errs, errors.New("project components folders are required"))
	}
	if c.FileName == "" {
		errs = append(errs, errors.New("project components file name is required"))
	}
	return errs
}

func ValidateKindCluster(c KindCluster) []error {
	var errs []error
	if c.KindConfig == nil {
		errs = append(errs, errors.New("kind config is required"))
	} else {
		errs = append(errs, ValidateKindConfig(*c.KindConfig)...)
	}
	return errs
}

func ValidateKindConfig(c KindConfig) []error {
	var errs []error
	if c.ConfigFile == "" {
		errs = append(errs, errors.New("kind config file is required"))
	}
	if c.Provider != KindConfigProviderDocker && c.Provider != KindConfigProviderPodman && c.Provider != KindConfigProviderNerdctl {
		errs = append(errs, errors.New("kind config provider must be docker, podman or nerdctl"))
	}
	return errs
}
