package cluster

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage Local Kubernetes clusters",
	}

	cmd.AddCommand(
		newInstallCommand(logger),
		newUninstallCommand(logger),
	)

	cmd.PersistentFlags().StringP("config", "c", "", "Path to the project configuration file")
	return cmd
}
