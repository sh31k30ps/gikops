package add

import (
	"github.com/sh31k30ps/gikops/pkg"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func newComponentsCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "component",
		Short: "Add component to the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			icmd, err := pkg.GetCommand(pkg.CommandProject, logger)
			if err != nil {
				return err
			}
			return icmd.Add("component", args...)
		},
	}

	return cmd
}
