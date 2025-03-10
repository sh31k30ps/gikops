package directories

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func GetComponentConfigFile(projectConfig *project.Project, name string) (string, error) {
	yamlCfg := filepath.Join(name, manager.GetComponentFileName(projectConfig)+"."+string(config.ConfigExtensionYAML))
	jsonCfg := filepath.Join(name, manager.GetComponentFileName(projectConfig)+"."+string(config.ConfigExtensionJSON))
	if _, err := os.Stat(yamlCfg); err == nil {
		return yamlCfg, nil
	}
	if _, err := os.Stat(jsonCfg); err == nil {
		return jsonCfg, nil
	}
	return "", fmt.Errorf("component config not found")
}
