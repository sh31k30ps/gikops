package clusterkind

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/ui/standard"
)

const (
	questionSelectProvider = "Select container runtime provider"
)

type Requester struct{}

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
	provider, err := ui.selectKindProvider()
	if err != nil {
		return nil, err
	}
	c := &cluster.KindCluster{}
	c.SetConfig(&cluster.KindConfig{
		Provider: provider,
	})
	c.SetName(cName)
	cluster.SetKindClusterDefaults(c)
	return c, nil
}

func (ui *Requester) selectKindProvider() (cluster.KindConfigProvider, error) {
	provider, err := standard.PromptChoice(questionSelectProvider, cluster.KindConfigProvidersLabels)
	if err != nil {
		return "", fmt.Errorf("error running provider selection: %w", err)
	}
	return cluster.KindConfigProvider(provider), nil
}
