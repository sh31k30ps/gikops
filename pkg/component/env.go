package component

import (
	"fmt"
	"slices"

	"github.com/sh31k30ps/gikops/pkg/services"
)

func checkEnvironment(env string) error {
	if env == "" {
		return fmt.Errorf("environment is required")
	}

	pCfg, err := services.GetCurrentProject()
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	if len(pCfg.Clusters) == 0 {
		return fmt.Errorf("project %s does not have any clusters", pCfg.Name)
	}

	for _, cluster := range pCfg.Clusters {
		if cluster.Name() == env {
			return nil
		}
	}

	return fmt.Errorf("environment %s is not available for project %s", env, pCfg.Name)
}

func checkComponentEnvironment(component, env string) error {
	if err := checkEnvironment(env); err != nil {
		return fmt.Errorf("failed to check environment: %w", err)
	}

	cfg, err := services.GetComponent(component)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}

	if len(cfg.Clusters) == 0 {
		return nil
	}

	if !slices.Contains(cfg.Clusters, env) {
		return fmt.Errorf("environment %s is not available for component %s", env, component)
	}
	return nil
}
