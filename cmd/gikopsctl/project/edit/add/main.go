package add

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add parameters to the project",
	}

	cmd.AddCommand(newClusterCommand(logger))
	cmd.AddCommand(newComponentsCommand(logger))

	return cmd
}
