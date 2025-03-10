package version

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/internal/version"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for version information
func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gikopsctl version %s\n", version.Version)
			fmt.Printf("  Git commit: %s\n", version.GitCommit)
			fmt.Printf("  Built: %s\n", version.BuildTime)
			fmt.Printf("  Go version: %s\n", version.GoVersion)
		},
	}

	return cmd
}
