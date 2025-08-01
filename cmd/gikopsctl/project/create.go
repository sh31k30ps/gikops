package project

import (
	"github.com/sh31k30ps/gikops/pkg"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func newCreateCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a project configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			icmd, err := pkg.GetCommand(pkg.CommandProject, logger)
			if err != nil {
				return err
			}
			return icmd.Create()
		},
	}

	return cmd
}
