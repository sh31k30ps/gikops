package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/helm"
)

func getHelmDefaultFile(cfg *component.HelmChart, prefix string) string {
	defaultsFile := "values.yaml"
	if prefix != "" {
		defaultsFile = prefix + "-" + defaultsFile
	}

	return filepath.Join("default", defaultsFile)
}

func processHelmDefaults(cfg *component.HelmChart, prefix string) error {
	if cfg.Chart == "" {
		return nil
	}
	defaultsPath := getHelmDefaultFile(cfg, prefix)
	if _, err := os.Stat(defaultsPath); err == nil {
		return nil
	}
	output, err := helm.GetDefaults(cfg.Chart, cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to get helm values: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(defaultsPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(defaultsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write values file: %w", err)
	}
	return nil
}
