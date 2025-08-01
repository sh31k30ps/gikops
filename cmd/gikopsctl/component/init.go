package component

import (
	"github.com/sh31k30ps/gikops/pkg/component"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/services"
	"github.com/spf13/cobra"
)

func newInitCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [component-name]",
		Short: "Initialize a component",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			folder, _ := cmd.Flags().GetString("folder")
			only, _ := cmd.Flags().GetBool("only")
			keepTmp, _ := cmd.Flags().GetBool("keep-tmp")
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
			mgr := component.NewManager(logger)
			return mgr.InitComponents(components, only, keepTmp)
		},
		ValidArgsFunction: validArgsFunction,
	}

	cmd.Flags().BoolP("only", "o", false, "Only initialize the component in the current folder")
	cmd.Flags().BoolP("keep-tmp", "k", false, "Keep the temporary Helm chart file")

	return cmd
}
