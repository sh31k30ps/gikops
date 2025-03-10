package setup

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type ProjectCreator struct {
	logger log.Logger
	status *cli.Status
	config *project.Project
	setup  *SetupConfig
}

func NewProjectCreator(logger log.Logger) *ProjectCreator {
	return &ProjectCreator{
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}

func (pc *ProjectCreator) Create(
	projectConfig *project.Project,
	setup *SetupConfig,
) error {
	project.SetProjectDefaults(projectConfig)
	if errs := project.ValidateProject(*projectConfig); len(errs) > 0 {
		return fmt.Errorf("invalid project config: %v", errs)
	}
	pc.config = projectConfig
	pc.setup = setup
	pc.status.Start("Creating project file")
	if err := pc.createProjectFile(); err != nil {
		pc.status.End(false)
		return err
	}
	if err := pc.copyBasicFiles(); err != nil {
		pc.status.End(false)
		return err
	}
	if err := pc.copySelectedComponentsTypes(); err != nil {
		pc.status.End(false)
		return err
	}
	if err := pc.copyOverrides(); err != nil {
		pc.status.End(false)
		return err
	}
	pc.status.End(true)

	return nil
}

func (pc *ProjectCreator) createProjectFile() error {
	pc.logger.V(1).Info("Creating project file")
	if err := manager.SaveProject("", pc.config); err != nil {
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
	pc.logger.V(1).Info(fmt.Sprintf("Copying kind.yaml for %s", pc.config.LocalCluster.KindConfig.ConfigFile))
	kind, err := assets.GetKindConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile("kind.yaml", kind, 0644); err != nil {
		return err
	}
	return nil
}

func (pc *ProjectCreator) copySelectedComponentsTypes() error {
	pc.logger.V(1).Info("Copying selected components")

	// Remove core directory if it exists
	if err := os.RemoveAll("core"); err != nil {
		return err
	}

	for componentType, componentsList := range pc.setup.Components {
		if err := pc.copyComponents(componentType, componentsList); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProjectCreator) copyComponents(typeName string, list []string) error {
	pc.logger.V(2).Info(fmt.Sprintf("Copying components for '%s'", typeName))
	for _, component := range list {
		if err := pc.copyComponent(typeName, component); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProjectCreator) copyComponent(typeName string, componentName string) error {
	pc.logger.V(2).Info(fmt.Sprintf("Copying component '%s'", componentName))
	dir := filepath.Join(typeName, componentName)
	files, err := assets.GetFilesFromSubdirectory(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		pc.logger.V(3).Info(fmt.Sprintf("Copying file '%s'", file))
		content, err := assets.GetFile(file)
		if path.Base(file) == "component.yaml" {
			file = pc.config.Components.FileName + ".yaml"
		}
		destFile := filepath.Join(
			"core",
			componentName,
			strings.TrimPrefix(file, filepath.Join(typeName, componentName)),
		)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(destFile, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProjectCreator) copyOverrides() error {
	pc.logger.V(1).Info("Copying overrides")
	for _, folder := range pc.config.LocalCluster.KindConfig.OverridesFolder {
		if err := pc.copyOverride(folder); err != nil {
			return err
		}
	}
	return nil
}

func (pc *ProjectCreator) copyOverride(folder string) error {
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
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(file, content, 0644); err != nil {
			return err
		}
	}
	return nil
}
