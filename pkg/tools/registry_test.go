package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var execCommand = exec.Command

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func resetExecCommand() {
	execCommand = exec.Command
}

func setupMockCommands() {
	execCommand = mockExecCommand
}

var mockOutputs = map[string]struct {
	stdout string
	stderr string
	err    error
}{
	"docker --version": {
		stdout: "Docker version 28.10.12, build 20.10.12-0ubuntu1~20.04.1",
		err:    nil,
	},
	"docker-fail --version": {
		stdout: "",
		err:    errors.New("command not found"),
	},
	"kustomize version": {
		stdout: "v5.6.0",
		err:    nil,
	},
	"kustomize-fail version": {
		stdout: "",
		err:    errors.New("command not found"),
	},
	"kubectl version --client": {
		stdout: "Client Version: v1.32.2\nKustomize Version: v5.5.0",
		err:    nil,
	},
	"kubectl-fail version --client": {
		stdout: "",
		err:    errors.New("command not found"),
	},
	"critical-tool --version": {
		stdout: "critical-tool version 0.9.0",
		err:    nil,
	},
	"cached-tool --version": {
		stdout: "cached-tool version 1.0.0",
		err:    nil,
	},
}

var testRegistry = map[string]*Tool{
	"docker":      ToolRegistry["docker"],
	"docker-fail": ToolRegistry["docker"],
	"kustomize":   ToolRegistry["kustomize"],
	"kustomize-fail": {
		Name:        "kustomize-fail",
		MinVersion:  "5.5.0",
		VersionArgs: []string{"version"},
		VersionGet: func(output string) string {
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
		IsMandatory:  true,
		Alternatives: []string{"kubectl-kustomize"},
	},
	"kubectl": ToolRegistry["kubectl"],
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
		IsMandatory:   false,
		IsAlternative: true,
	},
	"kubectl-fail": ToolRegistry["kubectl"],
	"critical-tool": {
		Name:        "critical-tool",
		MinVersion:  "1.0.0",
		VersionArgs: []string{"--version"},
		VersionGet:  func(output string) string { return "0.9.0" },
		IsMandatory: true,
	},
	"cached-tool": {
		Name:        "cached-tool",
		MinVersion:  "1.0.0",
		VersionArgs: []string{"--version"},
		VersionGet:  func(output string) string { return "1.0.0" },
		IsMandatory: false,
	},
}

// func initRegistry() {
// 	if r, ok := testRegistry["docker-fail"]; ok {
// 		r.Name = "docker-fail"
// 		testRegistry["docker-fail"] = r
// 	}
// 	if r, ok := testRegistry["kustomize-fail"]; ok {
// 		r.Name = "kustomize-fail"
// 		testRegistry["kustomize-fail"] = r
// 	}
// 	if r, ok := testRegistry["kubectl-fail"]; ok {
// 		r.Name = "kubectl-fail"
// 		testRegistry["kubectl-fail"] = r
// 	}
// }

func TestGetTool_Success(t *testing.T) {
	setupMockCommands()
	// initRegistry()
	defer resetExecCommand()

	resolver := &ToolResolver{
		registry: testRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: execCommand,
	}

	tool, err := resolver.GetTool("docker")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := &ResolvedTool{
		Tool:         testRegistry["docker"],
		Version:      "28.10.12",
		IsInstalled:  true,
		IsUpToDate:   true,
		ResolvedName: "docker",
	}

	tool.Tool.VersionGet = nil
	expected.Tool.VersionGet = nil

	if !reflect.DeepEqual(tool, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tool)
	}
}

func TestGetTool_UseAlternative(t *testing.T) {
	setupMockCommands()
	// initRegistry()
	defer resetExecCommand()

	resolver := &ToolResolver{
		registry: testRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: execCommand,
	}

	tool, err := resolver.GetTool("kustomize-fail")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if tool.ResolvedName != "kustomize-fail" {
		t.Errorf("Expected resolved name to be 'kustomize', got %s", tool.ResolvedName)
	}

	if tool.Tool.Name != "kubectl" {
		t.Errorf("Expected tool name to be 'kubectl', got %s", tool.Tool.Name)
	}

	expected := []string{"kustomize"}
	got := tool.GetCmdArgs()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected cmd args to be %v, got %v", expected, got)
	}
}

func TestGetTool_NotFound(t *testing.T) {
	resolver := &ToolResolver{
		registry: testRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: execCommand,
	}

	_, err := resolver.GetTool("nonexistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expected := "tool nonexistent not found in registry"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestGetTool_MandatoryToolNotAvailable(t *testing.T) {
	setupMockCommands()
	// initRegistry()
	defer resetExecCommand()

	resolver := &ToolResolver{
		registry: testRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: execCommand,
	}

	_, err := resolver.GetTool("critical-tool")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expected := "tool critical-tool (or its alternatives) is not available or up to date"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestGetTool_CacheUsage(t *testing.T) {
	setupMockCommands()
	defer resetExecCommand()

	resolver := &ToolResolver{
		registry: testRegistry,
		cache:    make(map[string]*ResolvedTool),
		executor: execCommand,
	}

	tool1, _ := resolver.GetTool("cached-tool")

	tool2, _ := resolver.GetTool("cached-tool")

	if tool1.Version != tool2.Version {
		t.Errorf("Expected cached version %s, got %s", tool1.Version, tool2.Version)
	}

	if tool2.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", tool2.Version)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd := args[0]
	cmdArgs := args[1:]
	key := cmd
	for _, arg := range cmdArgs {
		key += " " + arg
	}

	if output, ok := mockOutputs[key]; ok {
		if output.err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", output.stderr)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", output.stdout)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Command not mocked: %s\n", key)
	os.Exit(1)
}
