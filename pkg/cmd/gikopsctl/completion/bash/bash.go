package bash

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for bash completion
func NewCommand(logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Output shell completion code for bash",
		Long: `Output shell completion code for bash.
To load completions:

# for bash users
$ source <(gikopsctl completion bash)

# To load completions for each session, execute once:
## Linux:
$ gikopsctl completion bash > /etc/bash_completion.d/gikopsctl
## macOS:
$ gikopsctl completion bash > $(brew --prefix)/etc/bash_completion.d/gikopsctl`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		},
	}
}
