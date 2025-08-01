package project

import (
	"fmt"
	"os"
	"slices"

	"github.com/sh31k30ps/gikops/assets"
	"github.com/sh31k30ps/gikops/pkg/cli"
	"github.com/sh31k30ps/gikops/pkg/cluster"
	"github.com/sh31k30ps/gikops/pkg/component"
	"github.com/sh31k30ps/gikops/pkg/config/manager"
	cfgproject "github.com/sh31k30ps/gikops/pkg/config/project"
	"github.com/sh31k30ps/gikops/pkg/internal/git"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/services"
	uiproject "github.com/sh31k30ps/gikops/pkg/ui/project"
)

type (
	Command struct {
		logger log.Logger
		status *cli.Status
		ui     *uiproject.UIProjectRequester
	}
	EditProjectMode string
	AddProjectMode  string
)

const (
	EditName     EditProjectMode = "name"
	AddComponant AddProjectMode  = "component"
)

var (
	AvailableEditModes = []EditProjectMode{
		EditName,
	}
	AvailableAddModes = []AddProjectMode{
		AddComponant,
	}
)

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

func (c *Command) Add(mode string, args ...string) error {
	if !slices.Contains(AvailableAddModes, AddProjectMode(mode)) {
		return fmt.Errorf("invalid mode %s", mode)
	}
	switch AddProjectMode(mode) {
	case AddComponant:
		return c.addComponent(args...)
	default:
		return fmt.Errorf("invalid mode %s", mode)
	}
}

func (c *Command) Edit(mode string, args ...interface{}) error {
	if !slices.Contains(AvailableEditModes, EditProjectMode(mode)) {
		return fmt.Errorf("invalid mode %s", mode)
	}
	cfg, err := services.GetCurrentProject()
	if err != nil {
		return err
	}
	switch EditProjectMode(mode) {
	case EditName:
		cfg.Name = args[0].(string)
	default:
		return fmt.Errorf("invalid mode %s", mode)
	}
	if err := manager.SaveProject(services.GetConfigFile(), cfg); err != nil {
		return err
	}
	c.logger.V(0).Info("Project file edited successfully")
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
	c.logger.V(0).Info("Installation complete")
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

func (c *Command) addComponent(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("component name is required")
	}
	name := args[0]

	c.status.Start(fmt.Sprintf("Adding component %s", name))

	cfg, err := services.GetCurrentProject()
	if err != nil {
		c.status.End(false)
		return err
	}

	cfg.Components = append(cfg.Components, cfgproject.ProjectComponent{
		Name: name,
	})

	c.logger.V(1).Info(fmt.Sprintf("Saving project file %s", services.GetConfigFile()))
	if err := manager.SaveProject(services.GetConfigFile(), cfg); err != nil {
		c.status.End(false)
		return err
	}

	c.logger.V(1).Info(fmt.Sprintf("Creating component directory %s", name))
	if err := os.MkdirAll(name, 0755); err != nil {
		c.status.End(false)
		return err
	}
	c.status.End(true)

	c.logger.V(0).Info(fmt.Sprintf("componenst directory %s Created successfully", name))
	return nil
}
