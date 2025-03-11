package fish

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for fish completion
func NewCommand(logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "fish",
		Short: "Output shell completion code for fish",
		Long: `Output shell completion code for fish.
To load completions:

# To load completions for each session, execute once:
$ gikopsctl completion fish > ~/.config/fish/completions/gikopsctl.fish

# Alternatively, you can pipe the completion directly:
$ gikopsctl completion fish | source`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		},
	}
}
