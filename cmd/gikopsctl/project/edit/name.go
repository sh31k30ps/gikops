package edit

import (
	"errors"

	"github.com/sh31k30ps/gikops/pkg"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

func newEditNameCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "name",
		Short: "Set the name of the project",
		Long:  `Set the name of the project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expected exactly one argument")
			}
			icmd, err := pkg.GetCommand(pkg.CommandProject, logger)
			if err != nil {
				return err
			}
			return icmd.Edit("name", args[0])
		},
	}
	return cmd
}
