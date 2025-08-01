package delete

import (
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete parameters from the project",
	}

	cmd.AddCommand(newClusterCommand(logger))

	return cmd
}
