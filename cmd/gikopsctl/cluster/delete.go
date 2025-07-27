package cluster

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func newDeleteCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [cluster-name]",
		Short: "Delete a cluster from the project",
	}
	return CmdWithDeleteCluster(cmd, logger)
}

func CmdWithDeleteCluster(cmd *cobra.Command, logger log.Logger) *cobra.Command {
	cmd.Args = cobra.ExactArgs(1)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("cluster name is required")
		}
		icmd, err := pkg.GetCommand(pkg.CommandCluster, logger)
		if err != nil {
			return err
		}
		if icmd, ok := icmd.(*cluster.Command); ok {
			if err := icmd.Delete(args[0]); err != nil {
				return err
			}
			logger.V(0).Info(fmt.Sprintf("Cluster %s deleted successfully", args[0]))
			return nil
		}
		return fmt.Errorf("invalid command")
	}
	cmd.ValidArgsFunction = validArgsFunction
	return cmd
}
