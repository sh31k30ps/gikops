package check

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/tools"
	"github.com/spf13/cobra"
)

const (
	success = " \x1b[32m✓\x1b[0m"
	failure = " \x1b[31m✗\x1b[0m"
	warning = " \x1b[33m⚠\x1b[0m"
)

func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check the required tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools := tools.GetTools()
			logger.V(0).Info("Dependencies:")
			for _, tool := range tools {
				status := success
				message := "up to date"

				if !tool.IsInstalled {
					message = "not installed"
					status = warning
					if tool.IsMandatory {
						status = failure
					}
					logger.V(0).Info(fmt.Sprintf("  %s %s: %s", status, tool.Name, message))
					continue
				}

				if !tool.IsUpToDate {
					status = warning
					if tool.IsMandatory {
						status = failure
					}
					message = "not up to date"
				}

				logger.V(0).Info(fmt.Sprintf("  %s %s: %s", status, tool.Name, message))

				if tool.IsInstalled {
					logger.V(0).Info(fmt.Sprintf("    	Minimal version: %s", tool.MinVersion))
					logger.V(0).Info(fmt.Sprintf("    	Current version: %s", tool.CurrentVersion))
				}
			}
			return nil
		},
	}

	return cmd
}
