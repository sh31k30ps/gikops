package common

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/sh31k30ps/gikops/pkg/config/component"
)

func CleanBaseDir(cfg *component.Component) error {
	// Clean up base directory except kustomization.yaml and configured files
	baseDir := filepath.Join("base")

	files, err := os.ReadDir(baseDir)
	if !os.IsNotExist(err) && err != nil {
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return fmt.Errorf("failed to create base directory: %w", err)
		}
		return nil
	}

	// Keep track of files that should be preserved
	preserveFiles := []string{
		"kustomization.yaml",
	}

	if cfg.Files != nil && len(cfg.Files.Keep) > 0 {
		preserveFiles = append(preserveFiles, cfg.Files.Keep...)
	}

	// Remove files that aren't in the preserve list
	for _, f := range files {
		if !slices.Contains(preserveFiles, f.Name()) {
			filePath := filepath.Join(baseDir, f.Name())
			if f.IsDir() {
				if err := os.RemoveAll(filePath); err != nil {
					return fmt.Errorf("failed to remove directory %s: %w", f.Name(), err)
				}
			} else {
				if err := os.Remove(filePath); err != nil {
					return fmt.Errorf("failed to remove file %s: %w", f.Name(), err)
				}
			}
		}
	}
	return nil
}
