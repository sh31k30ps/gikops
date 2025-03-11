package component

import (
	"github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
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
			components, err := getComponents(folder, all, args)
			if err != nil {
				return err
			}
			mgr := component.NewManager(logger)
			return mgr.InitComponents(components, only)
		},
		ValidArgsFunction: validArgsFunction,
	}

	cmd.Flags().BoolP("only", "o", false, "Only initialize the component in the current folder")

	return cmd
}
