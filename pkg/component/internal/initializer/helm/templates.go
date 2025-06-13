package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/component/internal/initializer/common"
	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/helm"
)

func processHelmTemplate(name string, cfg *component.HelmChart, prefix string, keepTmp bool, options ...string) error {
	if cfg.Chart == "" {
		return nil
	}
	if err := os.MkdirAll("base", 0755); err != nil {
		return err
	}

	output, err := helm.GetTemplates(name, cfg.Chart, getHelmDefaultFile(cfg, prefix), cfg.Version, options...)
	if err != nil {
		return fmt.Errorf("failed to get helm templates: %w", err)
	}
	if prefix != "" {
		file := filepath.Join("base", prefix+".yaml")
		if err := os.WriteFile(file, output, 0644); err != nil {
			return fmt.Errorf("failed to write values file: %w", err)
		}
		return nil
	}

	tmpFile := filepath.Join("base", "tmp.yaml")
	if err := os.WriteFile(tmpFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write values file: %w", err)
	}

	if !keepTmp {
		defer os.Remove(tmpFile)
		if err := common.CutTemplate(tmpFile); err != nil {
			return fmt.Errorf("failed to process template: %w", err)
		}
	}

	return nil
}
