package component

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/spf13/cobra"
)

func newCreateCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a component",
		RunE: func(cmd *cobra.Command, args []string) error {
			folder, _ := cmd.Flags().GetString("folder")
			if d := services.GetComponentFolderFromDepth(); d != "" {
				folder = d
			}
			if folder == "" {
				return fmt.Errorf("component folder is required")
			}
			icmd, err := pkg.GetCommand(pkg.CommandComponent, logger)
			if err != nil {
				return err
			}
			return icmd.Create(folder)
		},
	}

	return cmd
}
