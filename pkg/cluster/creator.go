package cluster

import (
	"github.com/sh31k30ps/gikops/pkg/cluster/internal/creators/basiccreator"
	"github.com/sh31k30ps/gikops/pkg/cluster/internal/creators/kindcreator"
	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/log"
)

type Creator interface {
	Create(cfg config.ConfigObject) error
}

func GetCreator(clusterType cluster.ClusterType, logger log.Logger) Creator {
	switch clusterType {
	case cluster.ClusterTypeKind:
		return kindcreator.NewCreator(logger)
	case cluster.ClusterTypeBasic:
		return basiccreator.NewCreator(logger)
	default:
		return nil
	}
}

func GetCreatorFromConfig(c cluster.Cluster, logger log.Logger) Creator {
	switch c.(type) {
	case *cluster.KindCluster:
		return kindcreator.NewCreator(logger)
	case *cluster.BasicCluster:
		return basiccreator.NewCreator(logger)
	default:
		return nil
	}
}
