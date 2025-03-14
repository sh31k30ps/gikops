package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cluster/internal/installers/kindinstaller"
	"github.com/sh31k30ps/gikopsctl/pkg/cluster/internal/installers/noopinstaller"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

type Installer interface {
	Install(c cluster.Cluster) error
	Uninstall(c cluster.Cluster) error
}

func GetInstaller(logger log.Logger, clusterName string) (Installer, cluster.Cluster, error) {
	config, err := services.GetCurrentProject()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current project: %w", err)
	}

	c := config.GetCluster(clusterName)
	if c == nil {
		return nil, nil, fmt.Errorf("cluster %s not found", clusterName)
	}
	switch c.(type) {
	case *cluster.KindCluster:
		return kindinstaller.NewInstaller(logger), c, nil
	case *cluster.BasicCluster:
		return noopinstaller.NewInstaller(logger), c, nil
	}
	return nil, nil, fmt.Errorf("unsupported cluster type: %T", c)
}
