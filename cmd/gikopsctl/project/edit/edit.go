package edit

import (
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/project/edit/add"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/project/edit/delete"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a project configuration",
	}

	cmd.AddCommand(
		add.NewCommand(logger),
		delete.NewCommand(logger),
	)

	return cmd
}
