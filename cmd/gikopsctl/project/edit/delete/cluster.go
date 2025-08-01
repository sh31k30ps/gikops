package delete

import (
	"github.com/sh31k30ps/gikops/cmd/gikopsctl/cluster"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func newClusterCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Remove cluster from the project",
	}
	return cluster.CmdWithDeleteCluster(cmd, logger)
}
