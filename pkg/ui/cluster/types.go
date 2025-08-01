package cluster

import (
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
)

type UIClusterResults struct {
	Clusters []cluster.Cluster
}

type UIClusterInternalRequester interface {
	Request() (cluster.Cluster, error)
}
