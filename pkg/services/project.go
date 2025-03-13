package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

var (
	currentProject *project.Project
	configFile     string
)

const (
	maxDepth = 2
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
	if err := os.Chdir(filepath.Dir(fileName)); err != nil {
		return nil, fmt.Errorf("failed to change directory: %w", err)
	}
	return currentProject.(*project.Project), nil
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
