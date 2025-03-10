package completion

import (
	"errors"

	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/completion/bash"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/completion/fish"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/completion/powershell"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl/completion/zsh"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

const longDescription = `
Output shell completion code for the specified shell (bash, zsh, fish, or powershell).
This depends on the bash-completion binary. Example installation instructions:

Bash:
  $ source <(gikopsctl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gikopsctl completion bash > /etc/bash_completion.d/gikopsctl
  # macOS:
  $ gikopsctl completion bash > $(brew --prefix)/etc/bash_completion.d/gikopsctl

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  # Option 1:
  $ gikopsctl completion zsh > "${fpath[1]}/_gikopsctl"
  # Option 2: If using oh-my-zsh
  $ mkdir -p $ZSH/completions/
  $ gikopsctl completion zsh > $ZSH/completions/_gikopsctl

Fish:
  $ gikopsctl completion fish > ~/.config/fish/completions/gikopsctl.fish

PowerShell:
  PS> gikopsctl completion powershell > gikopsctl.ps1
  PS> ./gikopsctl.ps1

  # To load completions for every new session, add the output to your powershell profile:
  PS> gikopsctl completion powershell >> $PROFILE

Note for zsh users: zsh completions are only supported in versions of zsh >= 5.2
`

// NewCommand returns a new cobra.Command for shell completion
func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Output shell completion code for the specified shell (bash, zsh, fish, or powershell)",
		Long:  longDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return errors.New("subcommand is required")
		},
	}

	cmd.AddCommand(bash.NewCommand(logger))
	cmd.AddCommand(zsh.NewCommand(logger))
	cmd.AddCommand(fish.NewCommand(logger))
	cmd.AddCommand(powershell.NewCommand(logger))

	return cmd
}
