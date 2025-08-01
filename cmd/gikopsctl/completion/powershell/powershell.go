package powershell

import (
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for powershell completion
func NewCommand(logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "powershell",
		Short: "Output shell completion code for powershell",
		Long: `Output shell completion code for powershell.
To load completions:

# Create the completion file:
PS> gikopsctl completion powershell > gikopsctl.ps1

# Source the completion file:
PS> ./gikopsctl.ps1

# To load completions for every new session, add the output to your powershell profile:
PS> gikopsctl completion powershell >> $PROFILE`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
		},
	}
}
