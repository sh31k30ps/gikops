package project

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/assets"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	cpnt "github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

type ProjectCreator struct {
	logger log.Logger
	status *cli.Status
	config *project.Project
}

func NewProjectCreator(logger log.Logger) *ProjectCreator {
	return &ProjectCreator{
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}

func (pc *ProjectCreator) Create(
	cfg *project.Project,
) error {
	project.SetProjectDefaults(cfg)
	if errs := project.Validate(*cfg); len(errs) > 0 {
		return fmt.Errorf("invalid project config: %v", errs)
	}
	pc.config = cfg
	pc.status.Start("Creating project file")
	if err := pc.createProjectFile(); err != nil {
		pc.status.End(false)
		return err
	}
	pc.status.End(true)
	pc.status.Start("Copying basic files")
	if err := pc.copyBasicFiles(); err != nil {
		pc.status.End(false)
		return err
	}
	pc.status.End(true)
	pc.status.Start("Creating clusters")
	for _, cluster := range pc.config.Clusters {
		if err := pc.createCluster(cluster); err != nil {
			pc.status.End(false)
			return err
		}
	}
	pc.status.End(true)

	pc.status.Start("Installing components")
	for _, component := range pc.config.Components {
		if err := cpnt.InstallComponents(component); err != nil {
			pc.status.End(false)
			return err
		}
	}
	pc.status.End(true)
	return nil
}

func (pc *ProjectCreator) createProjectFile() error {
	pc.logger.V(1).Info("Creating project file")
	if err := manager.SaveProject(services.GetConfigFile(), pc.config); err != nil {
		return err
	}
	return nil
}

func (pc *ProjectCreator) copyBasicFiles() error {
	pc.logger.V(1).Info("Copying .gitignore")
	gitignore, err := assets.GetGitignore()
	if err != nil {
		return err
	}
	if err := os.WriteFile(".gitignore", gitignore, 0644); err != nil {
		return err
	}
	return nil
}

func (pc *ProjectCreator) createCluster(cluster project.ProjectCluster) error {
	pc.logger.V(1).Info(fmt.Sprintf("Creating cluster %s", cluster.Name()))
	switch c := cluster.(type) {
	case *project.KindCluster:
		return pc.createKindCluster(c)
	case *project.BasicCluster:
		return nil
	default:
		return fmt.Errorf("unsupported cluster type: %T", c)
	}
}
