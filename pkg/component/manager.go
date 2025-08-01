package component

import (
	"github.com/sh31k30ps/gikops/pkg/cli"
	"github.com/sh31k30ps/gikops/pkg/component/internal/initializer/helm"
	"github.com/sh31k30ps/gikops/pkg/component/internal/initializer/kustomize"
	"github.com/sh31k30ps/gikops/pkg/log"
)

type Manager struct {
	logger      log.Logger
	status      *cli.Status
	initializer []Initializer
}

func NewManager(logger log.Logger) *Manager {
	status := cli.StatusForLogger(logger)
	return &Manager{
		logger: logger,
		status: status,
		initializer: []Initializer{
			helm.NewInitializer(logger, status),
			kustomize.NewInitializer(logger, status),
		},
	}
}
