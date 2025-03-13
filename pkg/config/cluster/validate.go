package cluster

import "errors"

func ValidateAll(cl []Cluster) []error {
	errs := []error{}
	visitedClusters := make(map[string]bool)
	for _, c := range cl {
		errs = append(errs, Validate(c)...)
		if visitedClusters[c.Name()] {
			errs = append(errs, errors.New("cluster name must be unique"))
		}
		visitedClusters[c.Name()] = true
	}
	return errs
}

func Validate(c Cluster) []error {
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
