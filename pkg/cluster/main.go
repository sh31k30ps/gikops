package cluster

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	cfgcluster "github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	uicluster "github.com/sh31k30ps/gikopsctl/pkg/ui/cluster"
)

type Command struct {
	logger log.Logger
	status *cli.Status
	ui     *uicluster.UIClusterRequester
}

func NewCommand(logger log.Logger) *Command {
	return &Command{
		logger: logger,
		status: cli.StatusForLogger(logger),
		ui:     uicluster.NewRequester(),
	}
}

func (c *Command) Create() error {
	if _, err := c.ui.Request(); err != nil {
		return err
	}
	c.logger.V(0).Info("Apply clusters configuration")
	cfg, err := c.ui.Config()
	if err != nil {
		return err
	}
	if err := manager.SaveProject(services.GetConfigFile(), cfg.(*project.Project)); err != nil {
		return err
	}
	if _, err := services.ReloadCurrentProject(); err != nil {
		return err
	}
	c.logger.V(1).Info("Project file saved")
	if err := c.GenerateClusters(); err != nil {

		return err
	}
	return nil
}

func (c *Command) CreateSpecific(cType cfgcluster.ClusterType) error {
	cl, err := c.ui.RequestSpecificCluster(cType)
	if err != nil {
		return err
	}
	cfg, err := c.ui.Config()
	if err != nil {
		return err
	}
	if err := manager.SaveProject(services.GetConfigFile(), cfg.(*project.Project)); err != nil {
		return err
	}
	c.logger.V(1).Info("Project file saved")

	if err := c.GenerateCluster(cl); err != nil {
		return err
	}
	return nil
}

func (c *Command) Edit() error {
	return nil
}

func (c *Command) Delete(id interface{}) error {
	if id == nil {
		return fmt.Errorf("id is required")
	}
	if id, ok := id.(string); ok {
		if id == "" {
			return fmt.Errorf("id is required")
		}

		cfg, err := services.GetCurrentProject()
		if err != nil {
			return err
		}
		cls := []cfgcluster.Cluster{}
		for _, cl := range cfg.Clusters {
			if cl.Name() != id {
				cls = append(cls, cl)
			}
		}
		cfg.Clusters = cls
		if err := manager.SaveProject(services.GetConfigFile(), cfg); err != nil {
			return err
		}
		c.logger.V(1).Info("Project file saved")

		if err := os.RemoveAll(filepath.Join("clusters", id)); err != nil {
			return err
		}
		c.logger.V(1).Info("Cluster directory deleted")
		return nil
	}
	return fmt.Errorf("id is required")
}

func (c *Command) Add() error {
	return nil
}

func (c *Command) GenerateClusters() error {
	cfg, err := services.GetCurrentProject()
	if err != nil {
		return err
	}
	c.status.Start("Generating clusters")
	for _, cl := range cfg.Clusters {
		c.logger.V(1).Info(fmt.Sprintf("Generating cluster %s", cl.Name()))
		if err := c.GenerateCluster(cl); err != nil {
			c.status.End(false)
			return err
		}
	}
	c.status.End(true)
	return nil
}

func (c *Command) GenerateCluster(cluster cfgcluster.Cluster) error {
	creator := GetCreatorFromConfig(cluster, c.logger)
	if err := creator.Create(cluster); err != nil {
		return err
	}
	return nil
}
