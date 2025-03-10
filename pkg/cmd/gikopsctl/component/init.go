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

			components, err := getComponents(folder, all, args)
			if err != nil {
				return err
			}
			mgr := component.NewManager(logger)
			return mgr.InitComponents(components)
		},
		ValidArgsFunction: validArgsFunction,
	}

	return cmd
}
