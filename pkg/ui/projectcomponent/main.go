package projectcomponent

import (
	"fmt"
	"path/filepath"

	"github.com/sh31k30ps/gikops/assets"
	"github.com/sh31k30ps/gikops/pkg/config/project"
	"github.com/sh31k30ps/gikops/pkg/ui"
	"github.com/sh31k30ps/gikops/pkg/ui/standard"
)

const (
	questionAddComponentsFolders = "Do you want to add components folders?"
	questionAddMoreComponents    = "Do you want to add more components folders?"
)

type UIComponentRequester struct {
	results *UIComponentResults
}

func NewRequester() *UIComponentRequester {
	return &UIComponentRequester{
		results: &UIComponentResults{
			Components: []project.ProjectComponent{},
		},
	}
}
func (ui *UIComponentRequester) Request() (ui.UIRequestResult, error) {
	core, err := ui.promptCoreComponents()
	if err != nil {
		return nil, err
	}
	ui.results.Components = append(ui.results.Components, *core)
	if reps, err := standard.PromptYesNo(questionAddComponentsFolders); err != nil || !reps {
		return ui.results.Components, nil
	}
	if err := ui.requestComponentsFolders(); err != nil {
		return nil, err
	}
	return ui.results.Components, nil
}

func (ui *UIComponentRequester) promptCoreComponents() (*project.ProjectComponent, error) {
	core := &project.ProjectComponent{
		Name:    "core",
		Require: []string{},
	}
	componentsDirs, compDirErr := assets.GetSubdirectories("components")
	if compDirErr != nil {
		return nil, fmt.Errorf("error getting components: %w", compDirErr)
	}
	for _, compDir := range componentsDirs {
		components, compErr := assets.GetSubdirectories(compDir)
		if compErr != nil {
			return nil, fmt.Errorf("error getting %s components: %w", compDir, compErr)
		}
		// Clean component paths to get only the last part
		cleanComponents := make([]string, len(components))
		for i, comp := range components {
			cleanComponents[i] = filepath.Base(comp)
		}
		selected, err := ui.promptSubComponents(filepath.Base(compDir), cleanComponents)
		if err != nil {
			return nil, fmt.Errorf("error running component %s selection: %w", compDir, err)
		}
		for id, comp := range selected {
			selected[id] = filepath.Join(filepath.Base(compDir), comp)
		}
		core.Require = append(core.Require, selected...)
	}

	return core, nil
}

func (ui *UIComponentRequester) promptSubComponents(name string, components []string) ([]string, error) {
	selected, err := PromptComponents(name, components)
	if err != nil {
		return nil, fmt.Errorf("error running component %s selection: %w", name, err)
	}
	return selected, nil
}

func (ui *UIComponentRequester) requestComponentsFolders() error {
	stop := false
	for !stop {
		folder, err := standard.PromptName("", "components folder")
		if err != nil {
			return err
		}
		ui.results.Components = append(ui.results.Components, project.ProjectComponent{
			Name: folder,
		})
		if resp, err := standard.PromptYesNo(questionAddMoreComponents); err != nil || !resp {
			stop = true
		}
	}
	return nil
}
