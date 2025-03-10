package directories

import (
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func GetComponentsRoots(projectConfig *project.Project) []string {
	existingRoots := []string{}
	for _, root := range projectConfig.Components.Folders {
		if _, err := os.Stat(root); err == nil {
			existingRoots = append(existingRoots, root)
		}
	}
	return existingRoots
}

func GetRootComponents(projectConfig *project.Project, root string) []string {
	components := []string{}
	if _, err := os.Stat(root); err == nil {
		entries, err := os.ReadDir(root)
		if err != nil {
			return components
		}
		for _, entry := range entries {
			if entry.IsDir() && IsComponentDir(root, entry.Name(), projectConfig.Components.FileName) {
				components = append(components, entry.Name())
			}
		}
	}
	return components
}

func GetRootsComponents(projectConfig *project.Project) []string {
	roots := GetComponentsRoots(projectConfig)
	components := []string{}
	for _, root := range roots {
		rootComponents := GetRootComponents(projectConfig, root)
		// if len(roots) > 1 {
		for id, component := range rootComponents {
			rootComponents[id] = filepath.Join(root, component)
		}
		// }
		components = append(components, rootComponents...)
	}
	return components
}

func IsComponentDir(root, name, fileName string) bool {
	if _, err := os.Stat(filepath.Join(root, name, fileName+"."+string(config.ConfigExtensionYAML))); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(root, name, fileName+"."+string(config.ConfigExtensionJSON))); err == nil {
		return true
	}
	return false
}
