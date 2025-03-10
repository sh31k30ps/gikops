package component

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/dependencies"
)

func (m *Manager) InitComponents(components []string) error {
	m.logger.V(0).Info("Initializing components")
	m.status.Start("Checking dependencies")
	dg := dependencies.NewDependencyGraph()
	components, errs := dg.Resolve(components, pkg.GetComponent)
	if len(errs) > 0 {
		m.status.End(false)
		return fmt.Errorf("dependencies check failed: %v", errs)
	}
	m.status.End(true)

	m.logger.V(0).Info(fmt.Sprintf("Components to initialize: %v", components))
	for _, component := range components {
		if err := m.InitComponent(component); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) InitComponent(componentName string) error {
	m.status.Start(fmt.Sprintf("Initializing component %s", componentName))
	if err := m.initSingleComponent(componentName); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to initialize component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) initSingleComponent(name string) error {
	cCfg, err := pkg.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(currentDir)
	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("failed to change to component directory %s: %w", name, err)
	}

	if cCfg.Helm != nil {
		if err := processHookInit(cCfg.Helm.Before); err != nil {
			return fmt.Errorf("failed to process before hooks: %w", err)
		}

		if err := setupHelmRepo(name, cCfg); err != nil {
			return err
		}

		if err := processHookInit(cCfg.Helm.After); err != nil {
			return fmt.Errorf("failed to process after hooks: %w", err)
		}
	}

	return nil
}
