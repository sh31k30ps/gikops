package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/installer"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// newUninstallCommand returns a new cobra.Command for environment uninstallation
func newUninstallCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the Kubernetes environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile, _ := cmd.Flags().GetString("config")
			config, errs := manager.LoadProject(configFile)
			if len(errs) > 0 {
				return fmt.Errorf("failed to load project: %v", errs)
			}
			inst := installer.NewInstaller(config, logger)

			return inst.Uninstall()
		},
	}

	return cmd
}
