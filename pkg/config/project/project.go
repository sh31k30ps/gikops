package project

import "os"

type Project struct {
	Name         string
	Components   *ProjectComponents
	LocalCluster *KindCluster
	Environments []string
}

func NewProject() *Project {
	return &Project{
		Components:   NewProjectComponents(),
		LocalCluster: NewKindCluster(),
		Environments: []string{},
	}
}

func (p *Project) GetLocalClusterName() string {
	if p.LocalCluster == nil || p.LocalCluster.Name == "" {
		return p.Name
	}
	return p.LocalCluster.Name
}

type ProjectComponents struct {
	Folders  []string
	FileName string
}

func NewProjectComponents() *ProjectComponents {
	return &ProjectComponents{
		Folders:  []string{},
		FileName: "",
	}
}

type KindCluster struct {
	Name       string
	KindConfig *KindConfig
}

func NewKindCluster() *KindCluster {
	return &KindCluster{
		Name:       "",
		KindConfig: NewKindConfig(),
	}
}

type KindConfig struct {
	ConfigFile      string
	OverridesFolder []string
	Provider        KindConfigProvider
}

func NewKindConfig() *KindConfig {
	return &KindConfig{
		OverridesFolder: []string{},
		Provider:        "",
		ConfigFile:      "",
	}
}

func GetProjectFileName() string {
	projectFileName := os.Getenv(ProjectFileEnvVar)
	if projectFileName == "" {
		return DefaultProjectFile
	}
	return projectFileName
}
