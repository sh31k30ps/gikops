package project

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/sh31k30ps/gikopsctl/pkg/ui"
)

type SetupConfig struct {
	Name       string
	Components map[string][]string
	Clusters   []project.ProjectCluster
}

func NewSetupConfig() *SetupConfig {
	return &SetupConfig{
		Components: make(map[string][]string),
	}
}

type Setup struct {
	logger log.Logger
	status *cli.Status
}

func NewSetup(logger log.Logger) *Setup {
	return &Setup{
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}

func (s *Setup) Run() error {
	if manager.ProjectFileExists(services.GetConfigFile()) {
		return fmt.Errorf("project file already exists")
	}

	ui := ui.NewUIProject(s.logger)
	if err := ui.Request(); err != nil {
		return err
	}

	s.logger.V(0).Info("Start project creation")
	pc := NewProjectCreator(s.logger)
	if err := pc.Create(ui.Results.GetConfig()); err != nil {
		return err
	}

	return nil
}
