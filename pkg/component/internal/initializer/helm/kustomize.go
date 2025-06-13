package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/internal/kustomize"
)

func createKustomizeFile() error {
	kustomizeFile := filepath.Join("base", "kustomization.yaml")
	var err error
	if _, err = os.Stat(kustomizeFile); os.IsNotExist(err) {
		createKustomizeFileContent()
	}
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to create kustomize file: %w", err)
	}
	return nil
}

func createKustomizeFileContent() error {
	peviousDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(peviousDir)
	if err := os.Chdir("base"); err != nil {
		return fmt.Errorf("failed to change to component directory %s: %w", "base", err)
	}
	var baseFiles []string
	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			baseFiles = append(baseFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read base directory: %w", err)
	}
	if err := kustomize.Init(); err != nil {
		return fmt.Errorf("failed to init kustomize: %w", err)
	}
	if err := kustomize.AddResources(baseFiles); err != nil {
		return fmt.Errorf("failed to add resources: %w", err)
	}
	return nil
}
