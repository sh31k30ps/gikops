package helm

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

type Initializer struct {
	logger log.Logger
	status *cli.Status
}

func NewInitializer(logger log.Logger, status *cli.Status) *Initializer {
	return &Initializer{
		logger: logger,
		status: status,
	}
}

func (i *Initializer) Init(name string) error {
	cCfg, err := services.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	if cCfg.Helm != nil {
		if err := processHookInit(cCfg.Helm.Before); err != nil {
			return fmt.Errorf("failed to process before hooks: %w", err)
		}

		if err := setupHelmRepo(name, cCfg); err != nil {
			return err
		}

		if err := processHookInit(cCfg.Helm.After); err != nil {
			return fmt.Errorf("failed to process after hooks: %w", err)
		}
	}
	return nil
}
