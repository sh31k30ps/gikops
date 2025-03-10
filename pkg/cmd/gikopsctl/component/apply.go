package component

import (
	"github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func newApplyCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [component-name]",
		Short: "Apply a component to the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			folder, _ := cmd.Flags().GetString("folder")

			components, err := getComponents(folder, all, args)
			if err != nil {
				return err
			}
			env, _ := cmd.Flags().GetString("env")
			mode, _ := cmd.Flags().GetString("mode")
			mgr := component.NewManager(logger)
			return mgr.ApplyComponents(components, env, component.ApplyMode(mode))
		},
		ValidArgsFunction: validArgsFunction,
	}

	cmd.Flags().StringP("env", "e", "local", "Environment to target")
	cmd.Flags().StringP("mode", "m", component.ApplyModeAll.String(), "Mode to apply (all, crds, manifests)")
	cmd.RegisterFlagCompletionFunc("mode", flagCompletionMode)

	return cmd
}

func flagCompletionMode(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		component.ApplyModeAll.String(),
		component.ApplyModeCRDs.String(),
		component.ApplyModeManifests.String(),
	}, cobra.ShellCompDirectiveDefault
}
