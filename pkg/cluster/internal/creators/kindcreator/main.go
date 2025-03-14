package kindcreator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type Creator struct {
	logger log.Logger
}

func NewCreator(logger log.Logger) *Creator {
	return &Creator{logger: logger}
}

func (c *Creator) Create(cfg config.ConfigObject) error {
	cfgCluster, ok := cfg.(*cluster.KindCluster)
	if !ok {
		return fmt.Errorf("config is not a KindCluster")
	}
	c.logger.V(2).Info(fmt.Sprintf("Cleaning up kind cluster %s folder", cfgCluster.Name()))
	folder := filepath.Join("clusters", cfgCluster.Name())
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	c.logger.V(2).Info(fmt.Sprintf("Creating kind cluster %s folder", cfgCluster.Name()))
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	c.logger.V(2).Info(fmt.Sprintf("Creating kind cluster %s config file", cfgCluster.Name()))
	kind, err := assets.GetKindConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(folder, "kind.yaml"), kind, 0644); err != nil {
		return err
	}
	c.logger.V(2).Info("Copying overrides")
	if err := c.copyKindOverrides(folder, cfgCluster); err != nil {
		return err
	}
	return nil
}

func (c *Creator) copyKindOverrides(base string, cfg *cluster.KindCluster) error {
	for _, folder := range cfg.Config().(*cluster.KindConfig).OverridesFolder {
		if err := c.copyKindOverride(base, folder); err != nil {
			return err
		}
	}
	return nil
}

func (c *Creator) copyKindOverride(base string, folder string) error {
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	files, err := assets.GetFilesFromSubdirectory(folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		c.logger.V(2).Info(fmt.Sprintf("Copying file '%s'", file))
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
