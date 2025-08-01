package component

import (
	"github.com/sh31k30ps/gikops/pkg/component"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/services"
	"github.com/spf13/cobra"
)

func newDeleteCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a component",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			folder, _ := cmd.Flags().GetString("folder")
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
			force, _ := cmd.Flags().GetBool("force")
			mgr := component.NewManager(logger)
			return mgr.DeleteComponents(components, env, component.ApplyMode(mode), force)
		},
		ValidArgsFunction: validArgsFunction,
	}

	cmd.Flags().StringP("env", "e", "local", "Environment to target")
	cmd.Flags().StringP("mode", "m", component.ApplyModeAll.String(), "Mode to apply (all, crds, manifests)")
	cmd.Flags().BoolP("force", "", false, "Force delete component")
	cmd.RegisterFlagCompletionFunc("mode", flagCompletionMode)
	cmd.RegisterFlagCompletionFunc("env", flagCompletionEnv)

	return cmd
}
