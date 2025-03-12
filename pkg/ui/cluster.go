package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func (ui *UIProject) requestClusters() error {
	clusters := []project.ProjectCluster{}
	stop := false
	for !stop {
		cluster, err := ui.requestCluster()
		if err != nil {
			return err
		}
		clusters = append(clusters, cluster)
		if resp, err := ui.promptYesNo("Do you want to add another cluster?"); err != nil || !resp {
			stop = true
		}
	}
	ui.Results.Clusters = clusters
	return nil
}

func (ui *UIProject) requestCluster() (project.ProjectCluster, error) {
	clusterType, err := ui.selectTypeCluster()
	if err != nil {
		return nil, err
	}
	var (
		cluster project.ProjectCluster
		cErr    error
	)
	switch clusterType {
	case project.ClusterTypeKind:
		cluster, cErr = ui.RequestKindCluster()
		if cErr != nil {
			return nil, cErr
		}
	case project.ClusterTypeBasic:
		cluster, cErr = ui.RequestBasicCluster()
		if cErr != nil {
			return nil, cErr
		}
	default:
		return nil, fmt.Errorf("invalid cluster type: %s", clusterType)
	}
	return cluster, nil
}

func (ui *UIProject) selectTypeCluster() (project.ClusterType, error) {
	p := tea.NewProgram(newChoiceModel("Select cluster type", project.ClusterTypesLabels))
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running provider selection: %w", err)
	}
	choiceModel, ok := m.(ChoiceModel)
	if !ok {
		return "", fmt.Errorf("could not convert provider model")
	}
	if choiceModel.quitting {
		return "", fmt.Errorf("provider selection cancelled")
	}

	return project.ClusterType(choiceModel.choices[choiceModel.selected]), nil
}

func defaultKindCluster() *project.KindCluster {
	cluster := project.NewKindCluster()
	cluster.SetName("local")
	cluster.SetConfig(&project.KindConfig{
		Provider:   project.KindConfigProviderDocker,
		ConfigFile: "kind.yaml",
		OverridesFolder: []string{
			"overrides",
		},
	})
	return cluster
}
