package noopinstaller

import (
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type Installer struct {
	logger log.Logger
}

func NewInstaller(logger log.Logger) *Installer {
	return &Installer{
		logger: logger,
	}
}

func (i *Installer) Install(c cluster.Cluster) error {
	i.logger.V(0).Info("Nothing to install for this cluster")
	return nil
}

func (i *Installer) Uninstall(c cluster.Cluster) error {
	i.logger.V(0).Info("Nothing to uninstall for this cluster")
	return nil
}
