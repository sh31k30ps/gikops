package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/cluster"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

// newInstallCommand returns a new cobra.Command for environment installation
func newInstallCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install  [cluster-name]",
		Short: "Install the Kubernetes environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Args = cobra.ExactArgs(1)
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
