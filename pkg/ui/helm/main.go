package helm

import (
	"fmt"

	cfgcomponent "github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/ui"
	"github.com/sh31k30ps/gikopsctl/pkg/ui/standard"
)

type UIHelmRequester struct {
	results *UIHelmResults
}

type UIHelmResults struct {
	Repo    string
	RepoURL string
	Chart   *cfgcomponent.HelmChart
	CRDs    *cfgcomponent.HelmChart
}

func NewRequester() *UIHelmRequester {
	return &UIHelmRequester{
		results: &UIHelmResults{
			Chart: cfgcomponent.NewHelmChart(),
		},
	}
}

func (ui *UIHelmRequester) Request() (ui.UIRequestResult, error) {
	repo, err := standard.Prompt("", "Helm repository name?")
	if err != nil {
		return nil, err
	}
	if repo == "" {
		return nil, fmt.Errorf("repository name cannot be empty")
	}
	ui.results.Repo = repo

	repoURL, err := standard.Prompt("", "Helm repository URL?")
	if err != nil {
		return nil, err
	}
	if repoURL == "" {
		return nil, fmt.Errorf("repository URL cannot be empty")
	}
	ui.results.RepoURL = repoURL

	chart, err := requestChart()
	if err != nil {
		return nil, err
	}
	ui.results.Chart = chart

	withCrds, err := standard.PromptYesNo("Different CRDs chart?")
	if err != nil {
		return nil, err
	}
	if withCrds {
		chart, err := requestChart()
		if err != nil {
			return nil, err
		}
		ui.results.CRDs = chart
	}

	return ui.results, nil
}

func requestChart() (*cfgcomponent.HelmChart, error) {
	chart := cfgcomponent.NewHelmChart()
	chartName, err := standard.Prompt("", "Helm chart name?")
	if err != nil {
		return nil, err
	}
	if chartName == "" {
		return nil, fmt.Errorf("chart name cannot be empty")
	}
	chart.Chart = chartName

	chartURL, err := standard.Prompt("", "Specifique version?")
	if err != nil {
		return nil, err
	}
	if chartURL != "" {
		chart.Version = chartURL
	}

	return chart, nil
}
