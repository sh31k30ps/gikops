package project

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	cfgproject "github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/sh31k30ps/gikopsctl/pkg/ui"
	uiproject "github.com/sh31k30ps/gikopsctl/pkg/ui/project"
)

type Command struct {
	logger log.Logger
	status *cli.Status
	ui     ui.UIRequester
}

func NewCommand(logger log.Logger) *Command {
	return &Command{
		logger: logger,
		status: cli.StatusForLogger(logger),
		ui:     uiproject.NewRequester(logger),
	}
}

func (c *Command) Create() error {
	if manager.ProjectFileExists(services.GetConfigFile()) {
		return fmt.Errorf("project file already exists")
	}

	if _, err := c.ui.Request(); err != nil {
		return err
	}

	c.logger.V(0).Info("Create project")
	if ui, ok := c.ui.(ui.UIRequesterConfigAware); ok {
		cfg := ui.Config().(*cfgproject.Project)
		if err := manager.SaveProject(services.GetConfigFile(), cfg); err != nil {
			return err
		}

		if _, err := services.ReloadCurrentProject(); err != nil {
			return err
		}
		c.logger.V(1).Info("Project file saved")

		if err := c.Generate(); err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf("ui is not config aware")
}

func (c *Command) Add() error {
	return nil
}

func (c *Command) Edit() error {
	return nil
}

func (c *Command) Delete(id interface{}) error {
	return nil
}

func (c *Command) Generate() error {
	c.status.Start("Copying basic files")
	if err := c.copyBasicFiles(); err != nil {
		c.status.End(false)
		return err
	}
	c.status.End(true)
	cluster.NewCommand(c.logger).GenerateClusters()

	c.status.Start("Installing components")
	cfg, err := services.GetCurrentProject()
	if err != nil {
		c.status.End(false)
		return err
	}
	for _, cpnt := range cfg.Components {
		if err := component.InstallComponents(cpnt); err != nil {
			c.status.End(false)
			return err
		}
	}
	c.status.End(true)

	return nil
}

func (c *Command) copyBasicFiles() error {
	c.logger.V(1).Info("Copying .gitignore")
	gitignore, err := assets.GetGitignore()
	if err != nil {
		return err
	}
	if err := os.WriteFile(".gitignore", gitignore, 0644); err != nil {
		return err
	}
	return nil
}
