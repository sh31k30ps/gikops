package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func (ui *UIProject) requestKindCluster() (*project.KindCluster, error) {
	provider, err := ui.selectKindProvider()
	if err != nil {
		return nil, err
	}
	cluster := &project.KindCluster{}
	cluster.SetConfig(project.KindConfig{
		Provider: provider,
	})
	return cluster, nil
}

func (ui *UIProject) selectKindProvider() (project.KindConfigProvider, error) {
	p := tea.NewProgram(newChoiceModel("Select container runtime provider", project.KindConfigProvidersLabels))
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

	return project.KindConfigProvider(choiceModel.choices[choiceModel.selected]), nil
}
