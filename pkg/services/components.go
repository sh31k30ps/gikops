package services

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/config/component"
	"github.com/sh31k30ps/gikops/pkg/config/manager"
	"github.com/sh31k30ps/gikops/pkg/directories"
)

var (
	components map[string]*component.Component
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
