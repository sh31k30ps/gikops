package project

import (
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
)

type Project struct {
	Name       string
	Components []ProjectComponent
	Clusters   []cluster.Cluster
	Level      int
	Origin     string
}

func NewConfig() *Project {
	return &Project{
		Components: []ProjectComponent{},
		Clusters:   []cluster.Cluster{},
		Level:      0,
		Origin:     "",
	}
}

func (p *Project) GetCluster(name string) cluster.Cluster {
	for _, cluster := range p.Clusters {
		if cluster.Name() == name {
			return cluster
		}
	}
	return nil
}

func (p *Project) GetClustersNames() []string {
	names := []string{}
	for _, cluster := range p.Clusters {
		names = append(names, cluster.Name())
	}
	return names
}

func (p *Project) GetComponent(name string) *ProjectComponent {
	for _, component := range p.Components {
		if component.Name == name {
			return &component
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
