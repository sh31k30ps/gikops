package cluster

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type NoopInstaller struct {
	logger log.Logger
}

func NewNoopInstaller(logger log.Logger) *NoopInstaller {
	return &NoopInstaller{
		logger: logger,
	}
}

func (i *NoopInstaller) Install(c project.ProjectCluster) error {
	i.logger.V(0).Info("Nothing to install for this cluster")
	return nil
}

func (i *NoopInstaller) Uninstall(c project.ProjectCluster) error {
	i.logger.V(0).Info("Nothing to uninstall for this cluster")
	return nil
}
