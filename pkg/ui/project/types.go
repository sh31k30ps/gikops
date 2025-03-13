package project

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

type UIProjectResults struct {
	Name       string
	Components []project.ProjectComponent
	Clusters   []cluster.Cluster
}
