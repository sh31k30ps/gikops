package component

import "github.com/sh31k30ps/gikopsctl/pkg/config"

const (
	ComponentFileName      = "gikcpnt"
	ComponentFileExtension = string(config.ConfigExtensionYAML)
	ComponentFileEnvVar    = "GIKOPS_COMPONENT_FILE"
	DefaultComponentFile   = ComponentFileName + "." + ComponentFileExtension
	DefaultCRDsFileName    = "crds.yaml"
)
