package component

import (
	"fmt"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/directories"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

func newCheckCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check a component configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			folder, _ := cmd.Flags().GetString("folder")

			projectCfg, _ := pkg.GetCurrentProject()
			if projectCfg == nil {
				return fmt.Errorf("project file not found")
			}
			status := cli.StatusForLogger(logger)
			components := []string{}
			if all {
				components = directories.GetRootsComponents(projectCfg)
				if folder != "" {
					components = directories.GetRootComponents(projectCfg, folder)
				}
			} else {
				if folder != "" {
					for id, cpmt := range args {
						args[id] = filepath.Join(folder, cpmt)
					}
				}
				components = args
			}
			errors := map[string]error{}

			for _, component := range components {
				status.Start(fmt.Sprintf("Checking component %s", component))
				if _, err := pkg.GetComponent(component); err != nil {
					errors[component] = err
					status.End(false)
				}
				status.End(true)
			}

			if len(errors) > 0 {
				for component, err := range errors {
					logger.Error(fmt.Sprintf("Component %s have en error: %s", component, err))
				}
				return fmt.Errorf("errors found while checking components: %v", errors)
			}

			return nil
		},
	}

	return cmd
}
