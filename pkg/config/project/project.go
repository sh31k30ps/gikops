package project

import "os"

type Project struct {
	Name       string
	Components []ProjectComponent
	Clusters   []ProjectCluster
}

func NewProject() *Project {
	return &Project{
		Components: []ProjectComponent{},
		Clusters:   []ProjectCluster{},
	}
}

func (p *Project) GetCluster(name string) ProjectCluster {
	for _, cluster := range p.Clusters {
		if cluster.Name() == name {
			return cluster
		}
	}
	return nil
}

func GetProjectFileName() string {
	projectFileName := os.Getenv(ProjectFileEnvVar)
	if projectFileName == "" {
		return DefaultProjectFile
	}
	return projectFileName
}
