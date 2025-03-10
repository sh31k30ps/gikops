package project

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config"
)

type KindConfigProvider string

const (
	ProjectFileName       = "gikops"
	ProjectFileExtension  = string(config.ConfigExtensionYAML)
	ProjectFileEnvVar     = "GIKOPS_PROJECT_FILE"
	DefaultProjectFile    = ProjectFileName + "." + ProjectFileExtension
	DefaultKindConfigFile = "kind.yaml"
)

const (
	KindConfigProviderDocker  KindConfigProvider = "docker"
	KindConfigProviderPodman  KindConfigProvider = "podman"
	KindConfigProviderNerdctl KindConfigProvider = "nerdctl"
)

var KindConfigProvidersLabels = []string{
	string(KindConfigProviderDocker),
	string(KindConfigProviderPodman),
	string(KindConfigProviderNerdctl),
}

var KindConfigProviders = []KindConfigProvider{
	KindConfigProviderDocker,
	KindConfigProviderPodman,
	KindConfigProviderNerdctl,
}
