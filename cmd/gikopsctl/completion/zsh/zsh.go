package zsh

import (
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for zsh completion
func NewCommand(logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Output shell completion code for zsh",
		Long: `Output shell completion code for zsh.
To load completions:

# If shell completion is not already enabled in your environment,
# you will need to enable it. You can execute the following once:
$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
## Option 1:
$ gikopsctl completion zsh > "${fpath[1]}/_gikopsctl"

## Option 2: If using oh-my-zsh
$ mkdir -p $ZSH/completions/
$ gikopsctl completion zsh > $ZSH/completions/_gikopsctl

# You will need to start a new shell for this setup to take effect.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		},
	}
}
