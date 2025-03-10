package component

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/sh31k30ps/gikopsctl/pkg"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/directories"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra.Command for component management
func NewCommand(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "component",
		Short: "Manage Kubernetes components",
	}

	cmd.PersistentFlags().BoolP("all", "a", false, "Apply all components")
	cmd.PersistentFlags().StringP("folder", "f", "", "Folder to apply")
	cmd.RegisterFlagCompletionFunc("folder", flagCompletionFolder)

	cmd.AddCommand(
		newInitCmd(logger),
		newApplyCmd(logger),
		newCheckCmd(logger),
	)

	return cmd
}

func flagCompletionFolder(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectCfg, _ := pkg.GetCurrentProject()
	if projectCfg == nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	return directories.GetComponentsRoots(projectCfg), cobra.ShellCompDirectiveDefault
}

func validArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	all, _ := cmd.Flags().GetBool("all")
	if all {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	projectCfg, _ := pkg.GetCurrentProject()
	if projectCfg == nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	result := directories.GetRootsComponents(projectCfg)

	folder, _ := cmd.Flags().GetString("folder")

	if folder != "" {
		result = directories.GetRootComponents(projectCfg, folder)
	}

	_, unusedIngredients := lo.Difference(args, result)
	return unusedIngredients, cobra.ShellCompDirectiveNoFileComp
}

func getComponents(folder string, all bool, args []string) ([]string, error) {
	projectCfg, _ := pkg.GetCurrentProject()
	if projectCfg == nil {
		return nil, fmt.Errorf("project file not found")
	}

	if folder == "" {
		if all {
			return directories.GetRootsComponents(projectCfg), nil
		}
		return getComponentsFromGlobal(args)
	}

	roots := directories.GetComponentsRoots(projectCfg)
	if len(roots) == 0 {
		return nil, fmt.Errorf("no components found")
	}
	if !slices.Contains(roots, folder) {
		return nil, fmt.Errorf("folder %s not found", folder)
	}

	if all {
		return getComponentsFromFolder(folder, directories.GetRootComponents(projectCfg, folder))
	}

	return getComponentsFromFolder(folder, args)
}

func getComponentsFromGlobal(args []string) ([]string, error) {
	projectCfg, _ := pkg.GetCurrentProject()
	if projectCfg == nil {
		return nil, fmt.Errorf("project filenot found")
	}
	components := directories.GetRootsComponents(projectCfg)

	for _, arg := range args {
		if !slices.Contains(components, arg) {
			return nil, fmt.Errorf("component %s not exists", arg)
		}
	}

	return args, nil
}

func getComponentsFromFolder(folder string, args []string) ([]string, error) {
	projectCfg, _ := pkg.GetCurrentProject()
	if projectCfg == nil {
		return nil, fmt.Errorf("project filenot found")
	}
	roots := directories.GetComponentsRoots(projectCfg)
	if len(roots) == 0 {
		return nil, fmt.Errorf("no components found")
	}
	if !slices.Contains(roots, folder) {
		return nil, fmt.Errorf("folder %s not found", folder)
	}

	components := directories.GetRootComponents(projectCfg, folder)

	for id, arg := range args {
		if !slices.Contains(components, args[id]) {
			return nil, fmt.Errorf("component %s not exists", args[id])
		}
		args[id] = strings.Join([]string{folder, arg}, "/")
	}

	return args, nil
}
