package component

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"

	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
)

type ApplyMode string

const (
	ApplyModeAll       ApplyMode = "all"
	ApplyModeCRDs      ApplyMode = "crds"
	ApplyModeManifests ApplyMode = "manifests"
)

func (m ApplyMode) String() string {
	return string(m)
}

func (m ApplyMode) IsValid() bool {
	return m == ApplyModeAll || m == ApplyModeCRDs || m == ApplyModeManifests
}

type Manager struct {
	logger log.Logger
	status *cli.Status
}

func NewManager(logger log.Logger) *Manager {
	return &Manager{
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}
