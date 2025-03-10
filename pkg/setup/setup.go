package setup

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type Setup struct {
	logger         log.Logger
	status         *cli.Status
	UI             *UI
	ProjectCreator *ProjectCreator
}

func NewSetup(logger log.Logger) *Setup {
	return &Setup{
		logger:         logger,
		status:         cli.StatusForLogger(logger),
		UI:             NewUI(logger),
		ProjectCreator: NewProjectCreator(logger),
	}
}

func (s *Setup) Run() error {
	if manager.ProjectFileExists("") {
		return fmt.Errorf("project file already exists")
	}

	if err := s.UI.Request(); err != nil {
		return err
	}

	s.logger.V(0).Info("\n\nStart project creation")
	if err := s.ProjectCreator.Create(
		s.UI.createConfig(),
		s.UI.cfg,
	); err != nil {
		return err
	}

	return nil
}
