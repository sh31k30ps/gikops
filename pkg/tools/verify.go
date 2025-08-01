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
	Id           string
	Name         string
	Version      string
	MinVersion   string
	IsMandatory  bool
	IsInstalled  bool
	IsUpToDate   bool
	ResolvedName string
	Alternative  string
}

func ListTools() []ToolInfo {
	infos := []ToolInfo{}

	for id, tool := range ToolRegistry {
		if tool.IsAlternative {
			continue
		}

		info := ToolInfo{
			Id:          id,
			Name:        tool.Name,
			IsMandatory: tool.IsMandatory,
			MinVersion:  tool.MinVersion,
			IsInstalled: false,
			IsUpToDate:  false,
			Alternative: "",
		}
		t, err := GetToolResolver().GetTool(id)

		if err == nil {
			info.Version = t.Version
			info.IsInstalled = t.IsInstalled
			info.IsUpToDate = t.IsUpToDate
			if t.Tool.IsAlternative {
				info.Alternative = t.Tool.Name
			}
			info.ResolvedName = t.ResolvedName
		}

		infos = append(infos, info)
	}

	return infos
}
