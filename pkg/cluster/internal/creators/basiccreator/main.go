package basiccreator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/log"
)

type Creator struct {
	logger log.Logger
}

func NewCreator(logger log.Logger) *Creator {
	return &Creator{logger: logger}
}

func (c *Creator) Create(cfg config.ConfigObject) error {
	cfgCluster, ok := cfg.(*cluster.BasicCluster)
	if !ok {
		return fmt.Errorf("config is not a BasicCluster")
	}
	c.logger.V(2).Info(fmt.Sprintf("Creating basic cluster %s", cfgCluster.Name()))
	folder := filepath.Join("clusters", cfgCluster.Name())
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	c.logger.V(2).Info(fmt.Sprintf("Creating basic cluster %s folder", cfgCluster.Name()))
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	return nil
}
