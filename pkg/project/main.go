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
	"github.com/sh31k30ps/gikopsctl/pkg/internal/git"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	uiproject "github.com/sh31k30ps/gikopsctl/pkg/ui/project"
)

type Command struct {
	logger log.Logger
	status *cli.Status
	ui     *uiproject.UIProjectRequester
}

func NewCommand(logger log.Logger) *Command {
	return &Command{
		logger: logger,
		status: cli.StatusForLogger(logger),
		ui:     uiproject.NewRequester(logger),
	}
}

func (c *Command) Create(args ...interface{}) error {
	if manager.ProjectFileExists(services.GetConfigFile()) {
		return fmt.Errorf("project file already exists")
	}

	if _, err := c.ui.Request(); err != nil {
		return err
	}

	c.logger.V(0).Info("Create project")

	cfg := c.ui.Config().(*cfgproject.Project)
	if err := manager.SaveProject(services.GetConfigFile(), cfg); err != nil {
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

func (c *Command) Add() error {
	return nil
}

func (c *Command) Edit() error {
	return nil
}

func (c *Command) Delete(id interface{}) error {
	return nil
}

func (c *Command) Install() error {
	c.status.Start("Copying basic files")
	if err := c.copyBasicFiles(); err != nil {
		c.status.End(false)
		return err
	}
	c.status.End(true)

	c.status.Start("Installing components")
	if err := component.NewCommand(c.logger).Install(); err != nil {
		c.status.End(false)
		return err
	}
	c.status.End(true)

	c.status.Start("Installing Clusters")
	if err := cluster.NewCommand(c.logger).Install(); err != nil {
		c.status.End(false)
		return err
	}
	c.status.End(true)

	c.status.Start("Initializing git repository")
	if err := git.Init("."); err != nil {
		c.status.End(false)
		return err
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
