package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikops/pkg/config/manager"
	"github.com/sh31k30ps/gikops/pkg/config/project"
)

var (
	currentProject *project.Project
	configFile     string
)

const (
	maxDepth = 4
)

func GetCurrentProject() (*project.Project, error) {
	fileName := configFile
	if fileName == "" {
		fileName = project.GetProjectFileName()
	}
	var err error
	if currentProject == nil {
		currentProject, err = getProjectFile(fileName, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get current project file: %w", err)
		}
	}
	return currentProject, nil
}

func getProjectFile(fileName string, depth int) (*project.Project, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max depth reached")
	}
	currentProject, err := manager.GetConfigManager().Load(fileName)
	if err != nil {
		return getProjectFile(filepath.Join("..", fileName), depth+1)
	}

	if cfg, ok := currentProject.(*project.Project); ok {
		cfg.Level = depth
		if d, err := os.Getwd(); err == nil {
			cfg.Origin = d
		}
		if err := os.Chdir(filepath.Dir(fileName)); err != nil {
			return nil, fmt.Errorf("failed to change directory: %w", err)
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("failed to cast to project.Project")
}

func ReloadCurrentProject() (*project.Project, error) {
	currentProject = nil
	return GetCurrentProject()
}

func SetConfigFile(fileName string) {
	configFile = fileName
}

func GetConfigFile() string {
	return configFile
}
