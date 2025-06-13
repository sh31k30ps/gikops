package component

import (
	"github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/spf13/cobra"
)

func newApplyCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [component-name]",
		Short: "Apply a component to the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			folder, _ := cmd.Flags().GetString("folder")
			only, _ := cmd.Flags().GetBool("only")
			if d := services.GetComponentFolderFromDepth(); folder == "" && d != "" {
				folder = d
			}
			if c := services.GetComponentFromDepth(); c != "" && len(args) == 0 {
				args = []string{c}
			}
			components, err := getComponents(folder, all, args)
			if err != nil {
				return err
			}
			env, _ := cmd.Flags().GetString("env")
			mode, _ := cmd.Flags().GetString("mode")
			build, _ := cmd.Flags().GetBool("only-build")
			mgr := component.NewManager(logger)
			return mgr.ApplyComponents(components, env, component.ApplyMode(mode), only, build)
		},
		ValidArgsFunction: validArgsFunction,
	}

	cmd.Flags().StringP("env", "e", "local", "Environment to target")
	cmd.Flags().StringP("mode", "m", component.ApplyModeAll.String(), "Mode to apply (all, crds, manifests)")
	cmd.Flags().BoolP("only", "o", false, "Only initialize the component in the current folder")
	cmd.Flags().BoolP("only-build", "b", false, "Only build step")
	cmd.RegisterFlagCompletionFunc("mode", flagCompletionMode)
	cmd.RegisterFlagCompletionFunc("env", flagCompletionEnv)

	return cmd
}

func flagCompletionMode(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		component.ApplyModeAll.String(),
		component.ApplyModeCRDs.String(),
		component.ApplyModeManifests.String(),
	}, cobra.ShellCompDirectiveDefault
}

func flagCompletionEnv(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string
	if projectCfg, err := services.GetCurrentProject(); err == nil {
		names = projectCfg.GetClustersNames()
	}
	return names, cobra.ShellCompDirectiveDefault
}
