package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

type Installer interface {
	Install(c project.ProjectCluster) error
	Uninstall(c project.ProjectCluster) error
}

func GetInstaller(logger log.Logger, clusterName string) (Installer, project.ProjectCluster, error) {
	config, err := services.GetCurrentProject()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current project: %w", err)
	}

	cluster := config.GetCluster(clusterName)
	if cluster == nil {
		return nil, nil, fmt.Errorf("cluster %s not found", clusterName)
	}
	switch cluster.(type) {
	case *project.KindCluster:
		return NewKindInstaller(logger), cluster, nil
	case *project.BasicCluster:
		return NewNoopInstaller(logger), cluster, nil
	}
	return nil, nil, fmt.Errorf("unsupported cluster type: %T", cluster)
}
