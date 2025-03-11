package project

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/project"
	"github.com/spf13/cobra"
)

func newCreateCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a project configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			set := project.NewSetup(logger)
			return set.Run()
		},
	}

	return cmd
}
