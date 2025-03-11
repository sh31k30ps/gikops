package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// Tool represents an external tool requirement
type Tool struct {
	Name        string
	MinVersion  string
	VersionArgs []string
	VersionGet  func(string) string
	IsMandatory bool
}

// RequiredTools lists all external tools required by the application
var RequiredTools = []Tool{
	{
		Name:        "docker",
		MinVersion:  "28.0.0",
		VersionArgs: []string{"--version"},
		VersionGet: func(output string) string {
			parts := strings.Split(output, " ")
			if len(parts) >= 3 {
				return strings.TrimPrefix(parts[2], "version ")
			}
			return ""
		},
		IsMandatory: true,
	},
	{
		Name:        "kubectl",
		MinVersion:  "1.32.0",
		VersionArgs: []string{"version", "--client"},
		VersionGet: func(output string) string {
			parts := strings.Split(output, "Client Version: ")
			if len(parts) > 1 {
				version := strings.Split(parts[1], "\n")[0]
				return strings.TrimPrefix(version, "v")
			}
			return ""
		},
		IsMandatory: true,
	},
	{
		Name:        "helm",
		MinVersion:  "3.0.0",
		VersionArgs: []string{"version", "--short"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
		IsMandatory: false,
	},
	// {
	// 	Name:        "kind",
	// 	MinVersion:  "0.27.0",
	// 	VersionArgs: []string{"version"},
	// 	VersionGet: func(output string) string {
	// 		parts := strings.Split(strings.TrimSpace(output), " ")
	// 		if len(parts) >= 2 {
	// 			return strings.TrimPrefix(parts[1], "v")
	// 		}
	// 		return ""
	// 	},
	// 	IsMandatory: true,
	// },
	{
		Name:        "kustomize",
		MinVersion:  "5.5.0",
		VersionArgs: []string{"version"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
		IsMandatory: true,
	},
	{
		Name:        "cloud-provider-kind",
		MinVersion:  "0.6.0",
		VersionArgs: []string{"version"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
		IsMandatory: false,
	},
}

// VerifyTools checks if all required tools are installed and meet version requirements
func VerifyTools() error {
	for _, tool := range RequiredTools {
		if !tool.IsMandatory {
			continue
		}

		cmd := exec.Command(tool.Name, tool.VersionArgs...)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("%s n'est pas installé ou n'est pas accessible: %w", tool.Name, err)
		}

		version := tool.VersionGet(string(output))
		if version == "" || compareVersions(version, tool.MinVersion) < 0 {
			return fmt.Errorf(
				"%s version %s est inférieure à la version minimale requise %s",
				tool.Name,
				version,
				tool.MinVersion,
			)
		}
	}
	return nil
}

// GetToolVersions returns a map of tool names to their required and current versions
func GetTools() []struct {
	Name           string
	MinVersion     string
	CurrentVersion string
	IsInstalled    bool
	IsUpToDate     bool
	IsMandatory    bool
} {
	var tools []struct {
		Name           string
		MinVersion     string
		CurrentVersion string
		IsInstalled    bool
		IsUpToDate     bool
		IsMandatory    bool
	}

	for _, tool := range RequiredTools {
		cmd := exec.Command(tool.Name, tool.VersionArgs...)
		output, err := cmd.Output()
		currentVersion := ""
		isInstalled := false
		isUpToDate := false
		if err == nil {
			currentVersion = tool.VersionGet(string(output))
			isInstalled = true
			isUpToDate = compareVersions(currentVersion, tool.MinVersion) >= 0
		}

		tools = append(tools, struct {
			Name           string
			MinVersion     string
			CurrentVersion string
			IsInstalled    bool
			IsUpToDate     bool
			IsMandatory    bool
		}{
			Name:           tool.Name,
			MinVersion:     tool.MinVersion,
			CurrentVersion: currentVersion,
			IsInstalled:    isInstalled,
			IsUpToDate:     isUpToDate,
			IsMandatory:    tool.IsMandatory,
		})
	}

	return tools
}

// compareVersions compares two version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(strings.TrimPrefix(v1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(v2, "v"), ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		v1Num := parseInt(v1Parts[i])
		v2Num := parseInt(v2Parts[i])

		if v1Num < v2Num {
			return -1
		}
		if v1Num > v2Num {
			return 1
		}
	}

	if len(v1Parts) < len(v2Parts) {
		return -1
	}
	if len(v1Parts) > len(v2Parts) {
		return 1
	}

	return 0
}

func parseInt(s string) int {
	// Remove any non-numeric prefix/suffix
	numStr := strings.TrimFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})

	// Convert to int, default to 0 if empty or invalid
	num := 0
	fmt.Sscanf(numStr, "%d", &num)
	return num
}
