package cluster

import (
	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	cfgcluster "github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
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
				return icmd.Create()
			}
			return icmd.CreateSpecific(cfgcluster.ClusterType(args[0]))
		}
		return icmd.Create()
	}
	cmd.ValidArgsFunction = func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return cfgcluster.ClusterTypesLabels, cobra.ShellCompDirectiveNoFileComp
	}
	return cmd
}
