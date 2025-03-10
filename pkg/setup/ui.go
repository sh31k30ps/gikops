package setup

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type UI struct {
	logger log.Logger
	cfg    *SetupConfig
}

func NewUI(logger log.Logger) *UI {
	return &UI{
		logger: logger,
		cfg:    NewSetupConfig(),
	}
}

func (ui *UI) Request() error {
	if err := ui.selectName(); err != nil {
		return err
	}

	if err := ui.selectProvider(); err != nil {
		return err
	}

	if err := ui.selectComponents(); err != nil {
		return err
	}

	return nil
}

func (ui *UI) createConfig() *project.Project {
	cfg := project.NewProject()
	cfg.Name = ui.cfg.Name
	cfg.LocalCluster.KindConfig.Provider = ui.cfg.Provider

	return cfg
}

func (ui *UI) selectName() error {
	dir, err := filepath.Abs(".")
	if err == nil {
		ui.cfg.Name = filepath.Base(dir)
	}
	p := tea.NewProgram(newNameModel(ui.cfg.Name))
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running name selection: %w", err)
	}
	nameModel, ok := m.(nameModel)
	if !ok {
		return fmt.Errorf("could not convert name model")
	}
	if nameModel.textInput.Value() != "" {
		ui.cfg.Name = nameModel.textInput.Value()
	}
	return nil
}

func (ui *UI) selectProvider() error {
	p := tea.NewProgram(newProviderModel(project.KindConfigProvidersLabels))
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running provider selection: %w", err)
	}
	providerModel, ok := m.(providerModel)
	if !ok {
		return fmt.Errorf("could not convert provider model")
	}
	if providerModel.quitting {
		return fmt.Errorf("provider selection cancelled")
	}

	ui.cfg.Provider = project.KindConfigProvider(
		project.KindConfigProvidersLabels[providerModel.selected],
	)
	return nil
}

func (ui *UI) selectComponents() error {
	componentsDirs, compDirErr := assets.GetSubdirectories("components")
	if compDirErr != nil {
		return fmt.Errorf("error getting components: %w", compDirErr)
	}
	for _, compDir := range componentsDirs {
		components, compErr := assets.GetSubdirectories(compDir)
		if compErr != nil {
			return fmt.Errorf("error getting %s components: %w", compDir, compErr)
		}
		// Clean component paths to get only the last part
		cleanComponents := make([]string, len(components))
		for i, comp := range components {
			cleanComponents[i] = filepath.Base(comp)
		}
		selected, err := ui.selectSubComponents(filepath.Base(compDir), cleanComponents)
		if err != nil {
			return fmt.Errorf("error running component %s selection: %w", compDir, err)
		}
		ui.cfg.Components[compDir] = selected
	}
	return nil
}

func (ui *UI) selectSubComponents(name string, components []string) ([]string, error) {
	p := tea.NewProgram(newComponentModel(name, components))
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running component %s selection: %w", name, err)
	}

	componentModel, ok := m.(componentModel)
	if !ok {
		return nil, fmt.Errorf("could not convert component model")
	}

	if componentModel.quitting {
		return nil, fmt.Errorf("component selection cancelled")
	}

	// Convert selected map to slice
	var selectedComponents []string
	for i, comp := range components {
		if componentModel.selected[i] {
			selectedComponents = append(selectedComponents, comp)
		}
	}

	return selectedComponents, nil
}
