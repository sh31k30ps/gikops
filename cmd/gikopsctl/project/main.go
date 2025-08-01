package project

import (
	"github.com/sh31k30ps/gikops/cmd/gikopsctl/project/edit"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	cmd.AddCommand(
		newCreateCmd(logger),
		newInstallCmd(logger),
		edit.NewCommand(logger),
	)

	return cmd
}
