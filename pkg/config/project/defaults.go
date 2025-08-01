package project

import (
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
)

func SetProjectDefaults(p *Project) {
	if p == nil {
		p = NewConfig()
	}
	if len(p.Clusters) == 0 {
		kindCluster := cluster.NewKindCluster()
		kindCluster.SetName("local")
		kindCluster.SetConfig(cluster.NewKindConfig())
		p.Clusters = []cluster.Cluster{kindCluster}
	}

	for _, cl := range p.Clusters {
		switch c := cl.(type) {
		case *cluster.KindCluster:
			cluster.SetKindClusterDefaults(c)
		}
	}
}
