package config

import (
	"path/filepath"
	"strings"
)

type ConfigExtension string
type ConfigType string

const (
	ConfigExtensionJSON ConfigExtension = "json"
	ConfigExtensionYAML ConfigExtension = "yaml"
)

const (
	ConfigTypeProject   ConfigType = "project"
	ConfigTypeComponent ConfigType = "component"
)

var ConfigExtensions = []ConfigExtension{
	ConfigExtensionJSON,
	ConfigExtensionYAML,
}

var ConfigExtensionsLabels = []string{
	string(ConfigExtensionJSON),
	string(ConfigExtensionYAML),
}

var ConfigTypes = []ConfigType{
	ConfigTypeProject,
	ConfigTypeComponent,
}

var ConfigTypesLabels = []string{
	string(ConfigTypeProject),
	string(ConfigTypeComponent),
}

// GetConfigExtension returns the type of configuration file based on its extension
// Returns "json" for .json files, "yaml" for .yaml or .yml files, and "" for unknown
func GetConfigExtension(filename string) ConfigExtension {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return ConfigExtensionJSON
	case ".yaml", ".yml":
		return ConfigExtensionYAML
	default:
		return ""
	}
}
