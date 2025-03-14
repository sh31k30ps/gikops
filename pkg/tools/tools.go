package tools

import "strings"

var ToolRegistry = map[string]Tool{
	"docker": {
		Name:        "docker",
		MinVersion:  "28.0.0",
		VersionArgs: []string{"--version"},
		VersionGet: func(output string) string {
			parts := strings.Split(output, ",")
			return strings.TrimPrefix(parts[0], "Docker version ")
		},
		IsMandatory: false,
	},
	"podman": {
		Name:        "podman",
		MinVersion:  "5.4.0",
		VersionArgs: []string{"--version"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(output, "podman version ")
		},
		IsMandatory: false,
	},
	"nerdctl": {
		Name:        "nerdctl",
		MinVersion:  "2.0.0",
		VersionArgs: []string{"-v"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(output, "nerdctl version ")
		},
		IsMandatory: false,
	},
	"kubectl": {
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
	"helm": {
		Name:        "helm",
		MinVersion:  "3.0.0",
		VersionArgs: []string{"version", "--short"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.Split(output, "+")[0], "v")
		},
		IsMandatory: false,
	},
	"kustomize": {
		Name:        "kustomize",
		MinVersion:  "5.5.0",
		VersionArgs: []string{"version"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
		IsMandatory:  true,
		Alternatives: []string{"kubectl-kustomize"},
	},
	"kubectl-kustomize": {
		Name:        "kubectl",
		MinVersion:  "1.32.0",
		CmdArgs:     []string{"kustomize"},
		VersionArgs: []string{"version", "--client"},
		VersionGet: func(output string) string {
			parts := strings.Split(output, "Kustomize Version: ")
			if len(parts) > 1 {
				version := strings.Split(parts[1], "\n")[0]
				return strings.TrimPrefix(version, "v")
			}
			return ""
		},
		IsMandatory: false,
	},
	"git": {
		Name:        "git",
		MinVersion:  "2.39.0",
		VersionArgs: []string{"--version"},
		VersionGet: func(output string) string {
			return strings.Split(strings.TrimPrefix(output, "git version "), " ")[0]
		},
		IsMandatory: true,
	},
}
