package project

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config"
)

const (
	ProjectFileName      = "gikops"
	ProjectFileExtension = string(config.ConfigExtensionYAML)
	ProjectFileEnvVar    = "GIKOPS_PROJECT_FILE"
	DefaultProjectFile   = ProjectFileName + "." + ProjectFileExtension
)
