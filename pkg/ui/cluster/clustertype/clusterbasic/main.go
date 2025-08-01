package clusterbasic

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/ui/standard"
)

type Requester struct {
}

func NewRequester() *Requester {
	return &Requester{}
}

func (ui *Requester) Request() (cluster.Cluster, error) {
	cName, cNameErr := standard.PromptName("", "cluster")
	if cNameErr != nil {
		return nil, fmt.Errorf("cluster name is required")
	}
	if cName == "" {
		return nil, fmt.Errorf("cluster name is required")
	}

	c := cluster.NewBasicCluster()
	c.SetName(cName)
	return c, nil
}
