package component

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikops/pkg/directories"
	"github.com/sh31k30ps/gikops/pkg/internal/dependencies"
	"github.com/sh31k30ps/gikops/pkg/internal/kubectl"
	"github.com/sh31k30ps/gikops/pkg/services"
)

func (m *Manager) DeleteComponents(components []string, env string, mode ApplyMode, force bool) error {
	m.logger.V(0).Info("Deleting component")
	// TODO : Reverse resolve dependencies

	dg := dependencies.NewDependencyGraph()
	cfg, err := services.GetCurrentProject()
	if err != nil {
		return fmt.Errorf("failed to get project config: %w", err)
	}
	existingComponents := directories.GetRootsComponents(cfg)
	for _, component := range components {
		components, err := dg.ReverseResolve(component, existingComponents, services.GetComponent)
		if err != nil {
			return fmt.Errorf("failed to reverse resolve dependencies: %v", err)
		}
		if len(components) > 0 && !force {
			return fmt.Errorf("component %s has dependencies that must be deleted first: %v", component, components)
		}
		if len(components) > 0 && force {
			m.logger.Warnf("Deleting component %s with dependencies: %v", component, components)
		}
	}

	m.status.Start("Checking environment")
	if err := checkEnvironment(env); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to check component environment: %w", err)
	}
	// TODO: check if the components are available in the environment
	m.status.End(true)

	if mode == ApplyModeAll || mode == ApplyModeManifests {
		m.logger.V(0).Info("Deleting manifests")
		for _, component := range components {
			if err := m.DeleteComponentManifests(component, env); err != nil {
				return err
			}
		}
	}

	if mode == ApplyModeCRDs || mode == ApplyModeAll {
		m.logger.V(0).Info("Deleting CRDs")
		for _, component := range components {
			if err := m.DeleteComponentCRDs(component, env); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) DeleteComponentManifests(componentName, env string) error {
	m.status.Start(fmt.Sprintf("Deleting manifests for component %s", componentName))
	if err := m.deleteSingleManifests(componentName, env); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to delete manifests for component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) DeleteComponentCRDs(componentName, env string) error {
	m.status.Start(fmt.Sprintf("Deleting CRDs for component %s", componentName))
	if err := m.deleteSingleCRDs(componentName, env); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to delete CRDs for component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) deleteSingleManifests(componentName, env string) error {
	computedFile := filepath.Join(componentName, env, "computed.yaml")
	if _, err := os.Stat(computedFile); os.IsNotExist(err) {
		return nil
	}
	if err := kubectl.ChangeContext(env); err != nil {
		return fmt.Errorf("failed to change context: %w", err)
	}
	if err := kubectl.Delete(computedFile); err != nil {
		return fmt.Errorf("failed to delete manifests: %w", err)
	}
	return nil
}

func (m *Manager) deleteSingleCRDs(componentName, env string) error {
	cfg, err := services.GetComponent(componentName)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	if cfg.Files.SkipCRDs {
		return nil
	}
	if err := kubectl.ChangeContext(env); err != nil {
		return fmt.Errorf("failed to change context: %w", err)
	}
	if err := kubectl.Delete(filepath.Join(componentName, "base", cfg.Files.CRDs)); err != nil {
		return fmt.Errorf("failed to delete manifests: %w", err)
	}
	return nil
}
