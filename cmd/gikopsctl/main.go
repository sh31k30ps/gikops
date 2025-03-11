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
	longDescription = `gikopsctl is a command-line tool that helps manage Kubernetes components
and environments, replacing traditional Makefile and bash script approaches with
a more robust Go-based solution.

The project is organized into a main project folder and a components folder.
The main project folder contains the project configuration and the components folder contains the components.

The project configuration is stored in the gikops.yaml file.
The components are stored in the components folder.

The components are used to manage the Kubernetes components.
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
