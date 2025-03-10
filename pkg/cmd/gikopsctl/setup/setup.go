package setup

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/setup"
	"github.com/spf13/cobra"
)

func NewCommand(logger log.Logger) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			set := setup.NewSetup(logger)
			return set.Run()
		},
	}
	return cmd
}
