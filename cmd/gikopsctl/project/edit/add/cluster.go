package add

import (
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func newClusterCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Add a cluster to the project",
	}
	return cluster.CmdWithAddCluster(cmd, logger)

}
