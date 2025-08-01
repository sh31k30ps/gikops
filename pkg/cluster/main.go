package cluster

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikops/pkg/cli"
	"github.com/sh31k30ps/gikops/pkg/component"
	cfgcluster "github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/config/manager"
	"github.com/sh31k30ps/gikops/pkg/config/project"
	"github.com/sh31k30ps/gikops/pkg/directories"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/services"
	uicluster "github.com/sh31k30ps/gikops/pkg/ui/cluster"
)

type Command struct {
	logger log.Logger
	status *cli.Status
	ui     *uicluster.UIClusterRequester
	cmpt   *component.Command
}

func NewCommand(logger log.Logger) *Command {
	return &Command{
		logger: logger,
		status: cli.StatusForLogger(logger),
		ui:     uicluster.NewRequester(),
		cmpt:   component.NewCommand(logger),
	}
}

func (c *Command) Create(args ...interface{}) error {
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
	if err := c.Install(); err != nil {
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

func (c *Command) Edit(mode string, args ...interface{}) error {
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

		sCl := cfg.GetCluster(id)
		if sCl == nil {
			return fmt.Errorf("cluster %s not found", id)
		}
		cls := []cfgcluster.Cluster{}
		for _, cl := range cfg.Clusters {
			if cl.Name() != sCl.Name() {
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
		if err := c.cmpt.CleanComponentsCluster(sCl); err != nil {
			return err
		}
		c.logger.V(1).Info("Cluster directory deleted")
		return nil
	}
	return fmt.Errorf("id is required")
}

func (c *Command) Add(mode string, args ...string) error {
	return nil
}

func (c *Command) Install() error {
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
	project, err := services.GetCurrentProject()
	if err != nil {
		return err
	}
	components := directories.GetRootsComponents(project)
	for _, cmpt := range components {
		if err := c.cmpt.AddCluster(cmpt, cluster); err != nil {
			return err
		}
	}
	return nil
}
