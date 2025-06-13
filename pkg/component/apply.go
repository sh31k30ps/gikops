package component

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/dependencies"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/kubectl"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/kustomize"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

func (m *Manager) ApplyComponents(components []string, env string, mode ApplyMode, only, onlyBuild bool) error {
	m.logger.V(0).Info("Applying components")
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
	}

	m.status.Start("Checking environment")
	if err := checkEnvironment(env); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to check component environment: %w", err)
	}
	// TODO: check if the components are available in the environment
	m.status.End(true)

	if mode == ApplyModeCRDs || mode == ApplyModeAll {
		m.logger.V(0).Info(fmt.Sprintf("CRDs to apply: %v", components))
		for _, component := range components {
			if err := m.ApplyComponentCRDs(component, env); err != nil {
				return err
			}
		}
	}

	if mode == ApplyModeAll || mode == ApplyModeManifests {
		m.logger.V(0).Info(fmt.Sprintf("Components to apply: %v", components))
		for _, component := range components {
			if err := m.ApplyComponent(component, env, onlyBuild); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) ApplyComponentCRDs(componentName, env string) error {
	m.status.Start(fmt.Sprintf("Applying CRDs for component %s", componentName))
	if err := m.applySingleCRDs(componentName, env); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to apply CRDs for component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) ApplyComponent(componentName, env string, onlyBuild bool) error {
	m.status.Start(fmt.Sprintf("Applying component %s", componentName))
	if err := checkComponentEnvironment(componentName, env); err != nil {
		return fmt.Errorf("failed to check component environment: %w", err)
	}
	if err := m.applySingleComponent(componentName, env, onlyBuild); err != nil {
		m.status.End(false)
		return fmt.Errorf("failed to apply component %s: %w", componentName, err)
	}
	m.status.End(true)
	return nil
}

func (m *Manager) applySingleCRDs(name string, env string) error {
	cfg, err := services.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	if cfg.Files.SkipCRDs {
		return nil
	}
	if err := kubectl.ChangeContext(env); err != nil {
		return fmt.Errorf("failed to change context: %w", err)
	}
	if err := kubectl.CreateCRDs(filepath.Join(name, "base", cfg.Files.CRDs)); err != nil {
		return err
	}

	return nil
}

func (m *Manager) applySingleComponent(name, env string, onlyBuild bool) error {
	if err := kubectl.ChangeContext(env); err != nil {
		return fmt.Errorf("failed to change context: %w", err)
	}

	if _, err := os.Stat(filepath.Join(name, env, "kustomization.yaml")); os.IsNotExist(err) {
		return fmt.Errorf("kustomize.yaml not found in %s/%s: %w", name, env, err)
	}

	cfg, err := services.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}

	namespace := component.GetComponentPrefix(name)
	if cfg.Namespace != "" {
		namespace = cfg.Namespace
	}
	exists, errNs := kubectl.NamespaceExists(namespace)
	if errNs != nil {
		return fmt.Errorf("failed to check if namespace exists: %w", errNs)
	}
	if !exists {
		if err := kubectl.CreateNamespace(namespace); err != nil {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	}
	if err := kubectl.ChangeNamespace(namespace); err != nil {
		return fmt.Errorf("failed to set namespace: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	os.Chdir(fmt.Sprintf("%s/%s", name, env))
	defer os.Chdir(currentDir)

	if cfg.Exec != nil && len(cfg.Exec.Before) > 0 {
		for _, command := range cfg.Exec.Before {
			m.logger.V(1).Info(fmt.Sprintf("Running before script: %s", command))
			if err := m.logger.V(1).CmdOutput(exec.Command("sh", "-c", command)); err != nil {
				return fmt.Errorf("failed to run before script: %w", err)
			}
		}
	}

	m.logger.V(1).Info("Building")
	if err := kustomize.Build(name); err != nil {
		return fmt.Errorf("failed to build: %w", err)
	}

	if !onlyBuild {
		m.logger.V(1).Info("Applying")
		if err := kubectl.Apply("computed.yaml"); err != nil {
			return fmt.Errorf("failed to apply: %w", err)
		}

		if err := m.waitResources(); err != nil {
			return fmt.Errorf("failed to wait for resources: %w", err)
		}
	}

	if cfg.Exec != nil && len(cfg.Exec.After) > 0 {
		for _, command := range cfg.Exec.After {
			m.logger.V(1).Info(fmt.Sprintf("Running after script: %s", command))
			if err := m.logger.V(1).CmdOutput(exec.Command("sh", "-c", command)); err != nil {
				return fmt.Errorf("failed to run after script: %w", err)
			}
		}
	}

	return nil
}

func (m *Manager) waitResources() error {
	pg, err := getPodsGeneratorsFromFile("")
	if err != nil {
		return fmt.Errorf("failed to get pods generators: %w", err)
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(pg.GetDeployments())+len(pg.GetDaemonsets())+len(pg.GetStatefulsets()))
	waitResources := func(resources []string) {
		m.logger.V(1).Info(fmt.Sprintf("Waiting for resources %s to be ready", resources))
		defer wg.Done()
		if err := kubectl.WaittingForResourcesBeReady(resources); err != nil {
			errChan <- fmt.Errorf("error waiting for resource %s: %w", resources, err)
		}
		m.logger.V(1).Info(fmt.Sprintf("Resources %s are ready", resources))
	}

	if len(pg.GetDeployments()) > 0 {
		wg.Add(1)
		go waitResources(pg.GetDeployments())
	}
	if len(pg.GetDaemonsets()) > 0 {
		wg.Add(1)
		go waitResources(pg.GetDaemonsets())
	}
	if len(pg.GetStatefulsets()) > 0 {
		wg.Add(1)
		go waitResources(pg.GetStatefulsets())
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return err
	}
	return nil
}
