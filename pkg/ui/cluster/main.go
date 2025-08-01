package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/services"
	"github.com/sh31k30ps/gikops/pkg/ui"
	"github.com/sh31k30ps/gikops/pkg/ui/cluster/clustertype/clusterbasic"
	"github.com/sh31k30ps/gikops/pkg/ui/cluster/clustertype/clusterkind"
	"github.com/sh31k30ps/gikops/pkg/ui/standard"
)

const (
	questionAddCluster = "Do you want to add another cluster?"
	questionSelectType = "Select cluster type"
)

type UIClusterRequester struct {
	results *UIClusterResults
}

func NewRequester() *UIClusterRequester {
	return &UIClusterRequester{
		results: &UIClusterResults{},
	}
}

func (ui *UIClusterRequester) Request() (ui.UIRequestResult, error) {
	stop := false
	for !stop {
		_, err := ui.requestUnkwnownClusters()
		if err != nil {
			return nil, err
		}
		if resp, err := standard.PromptYesNo(questionAddCluster); err != nil || !resp {
			stop = true
		}
	}

	return ui.results.Clusters, nil
}

func (ui *UIClusterRequester) Config() (config.ConfigObject, error) {
	cfg, err := services.GetCurrentProject()
	if err != nil {
		return nil, err
	}
	if ui.results == nil {
		return cfg, nil
	}
	cfg.Clusters = append(cfg.Clusters, ui.results.Clusters...)
	return cfg, nil
}

func (ui *UIClusterRequester) RequestSpecificCluster(cType cluster.ClusterType) (cluster.Cluster, error) {
	var requester UIClusterInternalRequester
	if ui.results.Clusters == nil {
		ui.results.Clusters = []cluster.Cluster{}
	}
	switch cType {
	case cluster.ClusterTypeKind:
		requester = clusterkind.NewRequester()
	case cluster.ClusterTypeBasic:
		requester = clusterbasic.NewRequester()
	default:
		return nil, fmt.Errorf("invalid cluster type: %s", cType)
	}
	cluster, err := requester.Request()
	if err != nil {
		return nil, err
	}
	ui.results.Clusters = append(ui.results.Clusters, cluster)
	return cluster, nil
}

func (ui *UIClusterRequester) requestUnkwnownClusters() (cluster.Cluster, error) {
	cType, err := ui.selectTypeCluster()
	if err != nil {
		return nil, err
	}
	c, err := ui.RequestSpecificCluster(cType)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ui *UIClusterRequester) selectTypeCluster() (cluster.ClusterType, error) {
	clusterType, err := standard.PromptChoice(questionSelectType, cluster.ClusterTypesLabels)
	if err != nil {
		return "", fmt.Errorf("error running provider selection: %w", err)
	}
	return cluster.ClusterType(clusterType), nil
}
