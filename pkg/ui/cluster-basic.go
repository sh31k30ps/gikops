package ui

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func (ui *UIProject) RequestBasicCluster() (*project.BasicCluster, error) {
	cName, cNameErr := ui.selectName("", "cluster")
	if cNameErr != nil {
		return nil, fmt.Errorf("cluster name is required")
	}
	if cName == "" {
		return nil, fmt.Errorf("cluster name is required")
	}

	cluster := project.NewBasicCluster()
	cluster.SetName(cName)
	return cluster, nil
}
