package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func (pc *ProjectCreator) createKindCluster(c *project.KindCluster) error {
	pc.logger.V(1).Info(fmt.Sprintf("Cleaning up kind cluster %s folder", c.Name()))
	folder := filepath.Join("clusters", c.Name())
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	pc.logger.V(1).Info(fmt.Sprintf("Creating kind cluster %s folder", c.Name()))
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	pc.logger.V(1).Info(fmt.Sprintf("Creating kind cluster %s config file", c.Name()))
	kind, err := assets.GetKindConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(folder, "kind.yaml"), kind, 0644); err != nil {
		return err
	}
	pc.logger.V(1).Info("Copying overrides")
	if err := pc.copyKindOverrides(folder, c); err != nil {
		return err
	}
	return nil
}

func (pc *ProjectCreator) copyKindOverrides(base string, c *project.KindCluster) error {
	for _, folder := range c.Config().(*project.KindConfig).OverridesFolder {
		if err := pc.copyKindOverride(base, folder); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProjectCreator) copyKindOverride(base string, folder string) error {
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	files, err := assets.GetFilesFromSubdirectory(folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		pc.logger.V(2).Info(fmt.Sprintf("Copying file '%s'", file))
		content, err := assets.GetFile(file)
		if err != nil {
			return err
		}
		destFile := filepath.Join(base, folder, filepath.Base(file))
		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(destFile, content, 0644); err != nil {
			return err
		}
	}
	return nil
}
