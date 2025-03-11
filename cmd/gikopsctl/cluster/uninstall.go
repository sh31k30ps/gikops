package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// newUninstallCommand returns a new cobra.Command for environment uninstallation
func newUninstallCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the Kubernetes environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("cluster name is required")
			}
			inst, cluster, err := cluster.GetInstaller(logger, args[0])
			if err != nil {
				return err
			}
			return inst.Uninstall(cluster)
		},
		ValidArgsFunction: validArgsFunction,
	}

	return cmd
}
