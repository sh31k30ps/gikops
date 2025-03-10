package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/directories"
)

var (
	currentProject *project.Project
	components     map[string]*component.Component
)

const (
	maxDepth = 2
)

func GetComponent(name string) (*component.Component, error) {
	if components == nil {
		components = make(map[string]*component.Component)
	}
	if components[name] != nil {
		return components[name], nil
	}
	pCfg, err := GetCurrentProject()
	if err != nil {
		return nil, err
	}
	fileName, err := directories.GetComponentConfigFile(pCfg, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get component config file: %w", err)
	}

	cCfg, errs := manager.LoadComponent(fileName)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to load component: %v", errs)
	}
	components[name] = cCfg
	return cCfg, nil
}

func GetLoadedComponents() map[string]*component.Component {
	return components
}

func GetCurrentProject() (*project.Project, error) {
	fileName := project.GetProjectFileName()
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
