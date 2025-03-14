package project

import (
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/sh31k30ps/gikopsctl/pkg/ui"
	uicluster "github.com/sh31k30ps/gikopsctl/pkg/ui/cluster"
	component "github.com/sh31k30ps/gikopsctl/pkg/ui/projectcomponent"
	"github.com/sh31k30ps/gikopsctl/pkg/ui/standard"
)

const (
	questionCreateBasicCluster = "Do you want to create only standard kind cluster?"
	questionAddComponents      = "Do you want to add components?"
)

type UIProjectRequester struct {
	logger  log.Logger
	results *UIProjectResults
}

func NewRequester(logger log.Logger) *UIProjectRequester {
	return &UIProjectRequester{
		logger:  logger,
		results: &UIProjectResults{},
	}
}

func (ui *UIProjectRequester) Config() config.ConfigObject {
	cfg, err := services.GetCurrentProject()
	if err != nil {
		cfg = project.NewConfig()
	}

	cfg.Name = ui.results.Name
	cfg.Components = ui.results.Components
	cfg.Clusters = ui.results.Clusters

	return cfg
}

func (ui *UIProjectRequester) Request(args ...interface{}) (ui.UIRequestResult, error) {
	if err := ui.promptProjectName(); err != nil {
		return nil, err
	}
	basic, bErr := ui.promptCreateBasicCluster()
	if bErr != nil {
		return nil, bErr
	}
	ui.results.Clusters = []cluster.Cluster{}
	if basic {
		cl := cluster.DefaultKindCluster()
		cluster.SetKindClusterDefaults(cl)
		ui.results.Clusters = append(ui.results.Clusters, cl)
	} else {
		clusterRequester := uicluster.NewRequester()
		clusters, err := clusterRequester.Request()
		if err != nil {
			return nil, err
		}
		ui.results.Clusters = clusters.([]cluster.Cluster)
	}

	if reps, err := standard.PromptYesNo(questionAddComponents); err != nil || !reps {
		return ui.results, nil
	}
	componentRequester := component.NewRequester()
	components, err := componentRequester.Request()
	if err != nil {
		return nil, err
	}
	ui.results.Components = components.([]project.ProjectComponent)

	return ui.results, nil
}

func (ui *UIProjectRequester) promptProjectName() error {
	if dir, err := filepath.Abs("."); err == nil {
		ui.results.Name = filepath.Base(dir)
	}
	res, err := standard.PromptName(ui.results.Name, "project")
	if err != nil {
		return err
	}
	ui.results.Name = res
	return nil
}

func (ui *UIProjectRequester) promptCreateBasicCluster() (bool, error) {
	return standard.PromptYesNo(questionCreateBasicCluster)
}
