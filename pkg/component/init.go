package component

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/internal/dependencies"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

func (m *Manager) InitComponents(components []string, only bool) error {
	m.logger.V(0).Info("Initializing components")
	if !only {
		m.status.Start("Checking dependencies")
		dg := dependencies.NewDependencyGraph()
		var errs []error
		components, errs = dg.Resolve(components, services.GetComponent)
		if len(errs) > 0 {
			m.status.End(false)
			return fmt.Errorf("dependencies check failed: %v", errs)
		}
		m.status.End(true)
		m.logger.V(0).Info(fmt.Sprintf("Components to initialize: %v", components))
	}
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
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(currentDir)
	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("failed to change to component directory %s: %w", name, err)
	}
	for _, initializer := range m.initializer {
		if err := initializer.Init(name); err != nil {
			return fmt.Errorf("failed to initialize component %s: %w", name, err)
		}
	}
	return nil
}
