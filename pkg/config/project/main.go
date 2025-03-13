package project

import (
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
)

type Project struct {
	Name       string
	Components []ProjectComponent
	Clusters   []cluster.Cluster
}

func NewConfig() *Project {
	return &Project{
		Components: []ProjectComponent{},
		Clusters:   []cluster.Cluster{},
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

func GetProjectFileName() string {
	projectFileName := os.Getenv(ProjectFileEnvVar)
	if projectFileName == "" {
		return DefaultProjectFile
	}
	return projectFileName
}
