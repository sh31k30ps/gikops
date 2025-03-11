package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// newInstallCommand returns a new cobra.Command for environment installation
func newInstallCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the Kubernetes environment and components",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("cluster name is required")
			}
			inst, cluster, err := cluster.GetInstaller(logger, args[0])
			if err != nil {
				return err
			}
			return inst.Install(cluster)
		},
		ValidArgsFunction: validArgsFunction,
	}

	return cmd
}
