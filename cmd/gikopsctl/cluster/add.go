package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg"
	"github.com/sh31k30ps/gikops/pkg/cluster"
	cfgcluster "github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func newAddCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a cluster to the project",
	}
	return CmdWithAddCluster(cmd, logger)
}

func CmdWithAddCluster(cmd *cobra.Command, logger log.Logger) *cobra.Command {
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		icmd, err := pkg.GetCommand(pkg.CommandCluster, logger)
		if err != nil {
			return err
		}
		if icmd, ok := icmd.(*cluster.Command); ok {
			if len(args) == 0 {
				if err := icmd.Create(); err != nil {
					return err
				}
				logger.V(0).Info("Cluster created successfully")
				return nil
			}
			if err := icmd.CreateSpecific(cfgcluster.ClusterType(args[0])); err != nil {
				return err
			}
			logger.V(0).Info("Cluster created successfully")
			return nil
		}
		return fmt.Errorf("invalid command")
	}
	cmd.ValidArgsFunction = func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return cfgcluster.ClusterTypesLabels, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return cmd
}
