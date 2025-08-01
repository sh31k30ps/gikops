package project

import (
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/config/project"
)

type UIProjectResults struct {
	Name       string
	Components []project.ProjectComponent
	Clusters   []cluster.Cluster
}
