package ui

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type UIProjectResults struct {
	Name              string
	Components        project.ProjectComponent
	Clusters          []project.ProjectCluster
	ComponentsFolders []project.ProjectComponent
}

func (r *UIProjectResults) GetConfig() *project.Project {
	return &project.Project{
		Name:       r.Name,
		Components: r.ComponentsFolders,
		Clusters:   r.Clusters,
	}
}

type UIProject struct {
	logger  log.Logger
	Results *UIProjectResults
}

func NewUIProject(logger log.Logger) *UIProject {
	return &UIProject{
		logger:  logger,
		Results: &UIProjectResults{},
	}
}

func (ui *UIProject) Request() error {
	if err := ui.selectProjectName(); err != nil {
		return err
	}
	// TODO: add cluster
	basic, bErr := ui.promptYesNo("Do you want to create only standard kind cluster?")
	if bErr != nil {
		return fmt.Errorf("project creation cancelled")
	}
	ui.Results.Clusters = []project.ProjectCluster{}
	if basic {
		ui.Results.Clusters = append(ui.Results.Clusters, defaultKindCluster())
	} else {
		if err := ui.requestClusters(); err != nil {
			return err
		}
	}

	if reps, err := ui.promptYesNo("Do you want to add core components?"); err != nil || !reps {
		return nil
	}
	ui.Results.ComponentsFolders = []project.ProjectComponent{}

	ui.Results.Components = project.ProjectComponent{
		Name:    "core",
		Require: []string{},
	}
	if err := ui.selectComponents(); err != nil {
		return err
	}
	ui.Results.ComponentsFolders = append(ui.Results.ComponentsFolders, ui.Results.Components)
	if reps, err := ui.promptYesNo("Do you want to add components folders?"); err != nil || !reps {
		return nil
	}
	if err := ui.requestComponentsFolders(); err != nil {
		return err
	}

	return nil
}

func (ui *UIProject) selectProjectName() error {
	if dir, err := filepath.Abs("."); err == nil {
		ui.Results.Name = filepath.Base(dir)
	}
	res, err := ui.selectName(ui.Results.Name, "project")
	if err != nil {
		return err
	}
	ui.Results.Name = res
	return nil
}

func (ui *UIProject) selectComponents() error {
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
		ui.Results.Components.Require = append(ui.Results.Components.Require, selected...)
	}
	return nil
}

func (ui *UIProject) selectSubComponents(name string, components []string) ([]string, error) {
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
			selectedComponents = append(selectedComponents, filepath.Join(name, comp))
		}
	}

	return selectedComponents, nil
}

func (ui *UIProject) requestComponentsFolders() error {
	stop := false
	for !stop {
		folder, err := ui.selectName("", "components folder")
		if err != nil {
			return err
		}
		ui.Results.ComponentsFolders = append(ui.Results.ComponentsFolders, project.ProjectComponent{
			Name: folder,
		})
		if resp, err := ui.promptYesNo("Do you want to add another comment folder?"); err != nil || !resp {
			stop = true
		}
	}
	return nil
}
