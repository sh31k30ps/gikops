package tools

import (
	"fmt"
	"os/exec"
)

type Tool struct {
	Name          string
	MinVersion    string
	VersionArgs   []string
	VersionGet    func(string) string
	CmdArgs       []string
	IsMandatory   bool
	Alternatives  []string
	IsAlternative bool
}

type ToolResolver struct {
	registry map[string]*Tool
	cache    map[string]*ResolvedTool
	executor func(name string, args ...string) *exec.Cmd
}

type ResolvedTool struct {
	Tool         *Tool
	Version      string
	IsInstalled  bool
	IsUpToDate   bool
	ResolvedName string
}

var toolResolver *ToolResolver

func GetToolResolver() *ToolResolver {
	if toolResolver == nil {
		toolResolver = NewToolResolver(nil)
	}
	return toolResolver
}

func NewToolResolver(executor func(name string, args ...string) *exec.Cmd) *ToolResolver {
	if executor == nil {
		executor = exec.Command
	}
	return &ToolResolver{
		registry: ToolRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: executor,
	}
}

func (r *ToolResolver) GetTool(name string) (*ResolvedTool, error) {
	if resolved, exists := r.cache[name]; exists {
		return resolved, nil
	}

	tool, exists := r.registry[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found in registry", name)
	}

	resolved := r.resolveTool(tool)
	if resolved.IsInstalled && resolved.IsUpToDate {
		r.cache[name] = resolved
		return resolved, nil
	}

	for _, altName := range tool.Alternatives {
		alt, err := r.GetTool(altName)
		if err != nil {
			continue
		}

		if alt.IsInstalled && alt.IsUpToDate {
			alt.ResolvedName = name
			r.cache[name] = alt
			return alt, nil
		}
	}

	if tool.IsMandatory {
		return nil, fmt.Errorf("tool %s (or its alternatives) is not available or up to date", name)
	}

	r.cache[name] = resolved
	return resolved, nil
}

func (r *ToolResolver) resolveTool(tool *Tool) *ResolvedTool {
	cmd := r.executor(tool.Name, tool.VersionArgs...)
	output, err := cmd.Output()

	resolved := &ResolvedTool{
		Tool:         tool,
		ResolvedName: tool.Name,
	}

	if err != nil {
		resolved.IsInstalled = false
		resolved.IsUpToDate = false
		return resolved
	}

	version := tool.VersionGet(string(output))
	resolved.Version = version
	resolved.IsInstalled = true
	resolved.IsUpToDate = compareVersions(version, tool.MinVersion) >= 0

	return resolved
}

func (rt *ResolvedTool) GetCmdArgs() []string {
	args := []string{}
	args = append(args, rt.Tool.CmdArgs...)
	return args
}
