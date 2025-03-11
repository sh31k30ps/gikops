package remove

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove parameters from the project",
	}

	return cmd
}
