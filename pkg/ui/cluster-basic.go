package ui

import "github.com/sh31k30ps/gikopsctl/pkg/config/project"

func (ui *UIProject) requestBasicCluster() (*project.BasicCluster, error) {
	return project.NewBasicCluster(), nil
}
