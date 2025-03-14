package tools

import (
	"fmt"
	"slices"
)

// VerifyTools checks if all required tools are installed and meet version requirements
func VerifyTools() error {
	for id, tool := range ToolRegistry {
		if !tool.IsMandatory {
			continue
		}

		_, err := GetToolResolver().GetTool(id)
		if err != nil {
			return fmt.Errorf("%s n'est pas installé ou n'est pas accessible: %w", tool.Name, err)
		}
	}
	return nil
}

type ToolInfo struct {
	Id             string
	Name           string
	Version        string
	MinVersion     string
	IsMandatory    bool
	IsInstalled    bool
	IsUpToDate     bool
	UseAlternative bool
	ResolvedName   string
	Alternatives   []string
}

func ListTools() []ToolInfo {
	infos := []ToolInfo{}
	finalInfos := []ToolInfo{}
	alternatives := []string{}
	for id, tool := range ToolRegistry {
		info := ToolInfo{
			Id:             id,
			Name:           tool.Name,
			IsMandatory:    tool.IsMandatory,
			MinVersion:     tool.MinVersion,
			IsInstalled:    false,
			IsUpToDate:     false,
			UseAlternative: false,
		}
		t, err := GetToolResolver().GetTool(id)
		if err == nil {
			info.Version = t.Version
			info.IsInstalled = t.IsInstalled
			info.IsUpToDate = t.IsUpToDate
			info.UseAlternative = t.ResolvedName != id
			info.ResolvedName = t.ResolvedName
		}
		if len(t.Tool.Alternatives) > 0 {
			alternatives = append(alternatives, t.Tool.Alternatives...)
		}
		infos = append(infos, info)
	}

	for _, info := range infos {
		if slices.Contains(alternatives, info.Id) {
			continue
		}
		finalInfos = append(finalInfos, info)
	}
	return finalInfos
}
