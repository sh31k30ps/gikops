package manager

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikops/api/config/v1alpha1"
	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/project"
)

func LoadProject(file string) (*project.Project, []error) {
	file = getProjectFile(file)
	p, err := GetConfigManager().Load(file)
	if err != nil {
		return nil, []error{fmt.Errorf("error loading project: %w", err)}
	}
	if p, ok := p.(*project.Project); ok {
		return p, nil
	}
	return nil, []error{fmt.Errorf("invalid project configuration")}
}

func SaveProject(path string, p *project.Project) error {
	if envPath := os.Getenv(project.ProjectFileEnvVar); envPath != "" {
		path = envPath
	}
	if path == "" {
		path = project.DefaultProjectFile
	}
	return GetConfigManager().Save(
		path,
		p,
		v1alpha1.Version,
		v1alpha1.ProjectKind,
	)
}

func ProjectFileExists(file string) bool {
	if file != "" {
		if _, err := os.Stat(file); err == nil {
			return true
		}
		return false
	}
	file = getProjectFile(file)
	jsonPath := project.ProjectFileName + "." + string(config.ConfigExtensionJSON)

	if _, err := os.Stat(file); err == nil {
		return true
	}

	if _, err := os.Stat(jsonPath); err == nil {
		return true
	}

	return false
}

func getProjectFile(file string) string {
	if file != "" {
		return file
	}
	if envPath := os.Getenv(project.ProjectFileEnvVar); envPath != "" {
		return envPath
	}
	return project.DefaultProjectFile
}
