package component

import (
	"github.com/sh31k30ps/gikops/pkg/cli"
	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/component"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/services"
	"github.com/sh31k30ps/gikops/pkg/ui"
	"github.com/sh31k30ps/gikops/pkg/ui/helm"
	"github.com/sh31k30ps/gikops/pkg/ui/kustomize"
	"github.com/sh31k30ps/gikops/pkg/ui/standard"
)

const (
	questionAddKustomize = "Do you want to use kustomize?"
	questionAddHelm      = "Do you want to use helm?"
)

type UIComponentRequester struct {
	results *UIComponentResults
	logger  log.Logger
	status  *cli.Status
}

func NewRequester(logger log.Logger) *UIComponentRequester {
	return &UIComponentRequester{
		logger:  logger,
		status:  cli.StatusForLogger(logger),
		results: &UIComponentResults{},
	}
}

func (ui *UIComponentRequester) Request(componentName string) (ui.UIRequestResult, error) {
	ui.results.Name = componentName
	if componentName == "" {
		res, err := standard.PromptName("", "component")
		if err != nil {
			return nil, err
		}
		ui.results.Name = res
	}

	res, err := standard.Prompt("", "Specific namespace of the component")
	if err != nil {
		return nil, err
	}
	if res != "" {
		ui.results.Namespace = res
	}

	if hasHelm, err := standard.PromptYesNo(questionAddHelm); err != nil || hasHelm {
		requester := helm.NewRequester()
		h, err := requester.Request()
		if err != nil {
			return nil, err
		}
		ui.results.Helm = h.(*helm.UIHelmResults)
	}

	if hasKustomize, err := standard.PromptYesNo(questionAddKustomize); err != nil || hasKustomize {
		requester := kustomize.NewRequester()
		k, err := requester.Request()
		if err != nil {
			return nil, err
		}
		ui.results.Kustomize = k.(*kustomize.UIKustomizeResults)
	}

	return ui.results, nil
}

func (ui *UIComponentRequester) Config() (config.ConfigObject, error) {
	var (
		err error
		cfg *component.Component
	)

	if cfg, err = services.GetComponent(ui.results.Name); err != nil {
		cfg = &component.Component{}
		component.SetComponentDefaults(cfg)
	}

	cfg.Name = ui.results.Name
	cfg.Namespace = ui.results.Namespace

	if ui.results.Helm != nil {
		cfg.Helm = &component.HelmConfig{
			Chart: ui.results.Helm.Chart,
			Repo:  ui.results.Helm.Repo,
			URL:   ui.results.Helm.RepoURL,
		}
		if ui.results.Helm.CRDs != nil {
			cfg.Helm.CRDsChart = ui.results.Helm.CRDs
		}
	}

	if ui.results.Kustomize != nil {
		cfg.Kustomize = &component.KustomizeConfig{
			URLs: ui.results.Kustomize.URLs,
		}
	}

	return cfg, nil
}
