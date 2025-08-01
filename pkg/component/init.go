package component

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikops/pkg/config/component"
	"github.com/sh31k30ps/gikops/pkg/internal/dependencies"
	"github.com/sh31k30ps/gikops/pkg/internal/kustomize"
	"github.com/sh31k30ps/gikops/pkg/services"
)

func (m *Manager) InitComponents(components []string, only bool, keepTmp bool) error {
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
		if err := m.InitComponent(component, keepTmp); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) InitComponent(componentName string, keepTmp bool) error {
	m.status.Start(fmt.Sprintf("Initializing component %s", componentName))
	if err := m.initSingleComponent(componentName, keepTmp); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to initialize component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) initSingleComponent(name string, keepTmp bool) error {
	cfg, err := services.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component %s: %w", name, err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(currentDir)
	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("failed to change to component directory %s: %w", name, err)
	}
	for _, initializer := range m.initializer {
		if err := initializer.Init(name, keepTmp); err != nil {
			return fmt.Errorf("failed to initialize component %s: %w", name, err)
		}
	}
	if err := initClusters(cfg); err != nil {
		return fmt.Errorf("failed to initialize clusters for component %s: %w", name, err)
	}

	return nil
}

func initClusters(cfg *component.Component) error {
	if cfg == nil {
		return fmt.Errorf("missing configuration ")
	}
	clusters := cfg.Clusters
	if len(clusters) == 0 {
		project, err := services.GetCurrentProject()
		if err != nil {
			return fmt.Errorf("failed to get current project: %w", err)
		}
		clusters = []string{}
		for _, cluster := range project.Clusters {
			clusters = append(clusters, cluster.Name())
		}
	}

	for _, cluster := range clusters {
		if err := createComponentClusterKustomize(cluster); err != nil {
			return fmt.Errorf("failed to create kustomize for cluster %s: %w", cluster, err)
		}
	}
	return nil
}

func createComponentClusterKustomize(clusterName string) error {
	if _, err := os.Stat(filepath.Join(clusterName, "kustomization.yaml")); os.IsNotExist(err) {
		if err := os.MkdirAll(clusterName, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", clusterName, err)
		}
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		defer os.Chdir(currentDir)
		if err := os.Chdir(clusterName); err != nil {
			return fmt.Errorf("failed to change to component directory %s: %w", clusterName, err)
		}
		if err := kustomize.Init(); err != nil {
			return fmt.Errorf("failed to initialize kustomize: %w", err)
		}
		if err := kustomize.AddResource("../base"); err != nil {
			return fmt.Errorf("failed to add resource: %w", err)
		}
	}
	return nil
}
