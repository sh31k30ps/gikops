package tools

import (
	"fmt"
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
	Name           string
	Version        string
	MinVersion     string
	IsMandatory    bool
	IsInstalled    bool
	IsUpToDate     bool
	UseAlternative bool
	ResolvedName   string
}

func ListTools() []ToolInfo {
	infos := []ToolInfo{}
	for id, tool := range ToolRegistry {
		info := ToolInfo{
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
			info.IsInstalled = true
			info.IsUpToDate = t.IsUpToDate
			info.UseAlternative = t.ResolvedName != id
			info.ResolvedName = t.ResolvedName
		}
		infos = append(infos, info)
	}
	return infos
}
