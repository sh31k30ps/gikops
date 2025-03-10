package component

import (
	"fmt"
	"slices"

	"github.com/sh31k30ps/gikopsctl/pkg"
)

func checkEnvironment(env string) error {
	if env == "" {
		return fmt.Errorf("environment is required")
	}

	pCfg, err := pkg.GetCurrentProject()
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	if len(pCfg.Environments) == 0 {
		return fmt.Errorf("project %s does not have any environments", pCfg.Name)
	}

	if !slices.Contains(pCfg.Environments, env) {
		return fmt.Errorf("environment %s is not available for project %s", env, pCfg.Name)
	}

	return nil
}

func checkComponentEnvironment(component, env string) error {
	if err := checkEnvironment(env); err != nil {
		return fmt.Errorf("failed to check environment: %w", err)
	}

	cfg, err := pkg.GetComponent(component)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}

	if len(cfg.EnvironmentAvailability) == 0 {
		return nil
	}

	if !slices.Contains(cfg.EnvironmentAvailability, env) {
		return fmt.Errorf("environment %s is not available for component %s", env, component)
	}
	return nil
}
