package gikopsctl

import (
	"slices"

	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/check"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/completion"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/component"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/setup"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/version"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/tools"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewRootCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gikopsctl",
		Short: "A tool for managing Kubernetes components",
		Long: `gikopsctl is a command-line tool that helps manage Kubernetes components
and environments, replacing traditional Makefile and bash script approaches with
a more robust Go-based solution.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			exception := []string{"version", "completion", "check"}
			// Skip tool verification for version and completion commands
			if slices.Contains(exception, cmd.Name()) || (cmd.Parent() != nil && slices.Contains(exception, cmd.Parent().Name())) {
				return nil
			}

			// Update logger verbosity if verbose flag is set
			if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
				if l, ok := logger.(*cli.Logger); ok {
					l.SetVerbosity(1)
				}
			}

			return tools.VerifyTools()
		},
	}

	// Add common flags that will be inherited by all commands
	addCommonFlags(cmd)

	cmd.AddCommand(
		cluster.NewCommand(logger),
		setup.NewCommand(logger),
		component.NewCommand(logger),
		version.NewCommand(logger),
		completion.NewCommand(logger),
		check.NewCommand(logger),
	)

	return cmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
}
