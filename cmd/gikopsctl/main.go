package main

import (
	"os"
	"slices"

	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/check"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/cluster"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/completion"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/component"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/project"
	"github.com/sh31k30ps/gikopsctl/cmd/gikopsctl/version"
	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/sh31k30ps/gikopsctl/pkg/tools"
	"github.com/spf13/cobra"
)

const (
	longDescription = `GikOps is a tool that allows managing various Kubernetes clusters as well
as the different tools and applications deployed on them.
It enforces the GitOps principle by keeping all configurations in a versioned project.

It is also possible to manage dependencies between tools on a global scale.
The goal of this tool is to enable managing local Kubernetes clusters (development),
pre-production, and production clusters within the same configuration project to
ensure consistency across different clusters while maintaining differences.

It is important to have all the most similar environments possible to avoid
compatibility issues.
`
)

func NewRootCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "gikopsctl",
		Short:        "A tool for managing Kubernetes components",
		Long:         longDescription,
		SilenceUsage: true,
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

			if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
				services.SetConfigFile(configFile)
			}

			return tools.VerifyTools()
		},
	}

	// Add common flags that will be inherited by all commands
	addCommonFlags(cmd)

	cmd.AddCommand(
		cluster.NewCommand(logger),
		project.NewCommand(logger),
		component.NewCommand(logger),
		version.NewCommand(logger),
		completion.NewCommand(logger),
		check.NewCommand(logger),
	)

	return cmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	cmd.PersistentFlags().StringP("config", "c", "", "Project configuration file")
}

func main() {
	rootCmd := NewRootCmd(services.NewLogger("general"))
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
