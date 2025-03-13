package cluster

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
)

type UIClusterResults struct {
	Clusters []cluster.Cluster
}

type UIClusterInternalRequester interface {
	Request() (cluster.Cluster, error)
}
