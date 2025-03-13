package cluster

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
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
		newAddCommand(logger),
		newDeleteCommand(logger),
	)

	cmd.PersistentFlags().StringP("config", "c", "", "Path to the project configuration file")

	return cmd
}

func validArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectCfg, _ := services.GetCurrentProject()
	if projectCfg == nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	clusters := []string{}
	for _, cluster := range projectCfg.Clusters {
		clusters = append(clusters, cluster.Name())
	}

	return clusters, cobra.ShellCompDirectiveNoFileComp
}
