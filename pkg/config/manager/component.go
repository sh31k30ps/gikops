package manager

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/api/config/v1alpha1"
	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func LoadComponent(file string) (*component.Component, []error) {
	file = getComponentFile(file)
	c, err := GetConfigManager().Load(file)
	if err != nil {
		return nil, []error{fmt.Errorf("error loading component: %w", err)}
	}
	if c, ok := c.(*component.Component); ok {
		return c, nil
	}
	return nil, []error{fmt.Errorf("invalid component configuration")}
}

func SaveComponent(path string, cpnt *component.Component) error {
	if envPath := os.Getenv(component.ComponentFileEnvVar); envPath != "" {
		path = envPath
	}
	if path == "" {
		path = component.DefaultComponentFile
	}
	return GetConfigManager().Save(
		path,
		cpnt,
		v1alpha1.Version,
		v1alpha1.ComponentKind,
	)
}

func ComponentFileExists(file string) bool {
	file = getComponentFile(file)
	jsonPath := file + "." + string(config.ConfigExtensionJSON)
	yamlPath := file + "." + string(config.ConfigExtensionYAML)

	if _, err := os.Stat(jsonPath); err == nil {
		return true
	}
	if _, err := os.Stat(yamlPath); err == nil {
		return true
	}

	return false
}

func getComponentFile(file string) string {
	if file != "" {
		return file
	}
	if envPath := os.Getenv(component.ComponentFileEnvVar); envPath != "" {
		return envPath
	}
	return component.DefaultComponentFile
}

func GetComponentFileName(cfg *project.Project) string {
	comptName := os.Getenv(component.ComponentFileEnvVar)
	if comptName == "" {
		comptName = component.ComponentFileName
	}
	return comptName
}
