package component

import "strings"

type Component struct {
	Name                    string
	Namespace               string
	Disabled                bool
	DependsOn               []string
	Helm                    *HelmConfig
	Files                   *ComponentFiles
	Exec                    *ComponentExec
	EnvironmentAvailability []string
}

func NewComponent() *Component {
	return &Component{
		Name:                    "",
		Namespace:               "",
		Disabled:                false,
		DependsOn:               []string{},
		Files:                   NewComponentFiles(),
		EnvironmentAvailability: []string{},
	}
}

type ComponentFiles struct {
	CRDs     string
	SkipCRDs bool
	Keep     []string
}

func NewComponentFiles() *ComponentFiles {
	return &ComponentFiles{
		CRDs:     "",
		SkipCRDs: false,
		Keep:     []string{},
	}
}

type ComponentExec struct {
	Before []string
	After  []string
}

func NewComponentExec() *ComponentExec {
	return &ComponentExec{
		Before: []string{},
		After:  []string{},
	}
}

type HelmConfig struct {
	Chart     *HelmChart
	CRDsChart *HelmChart
	Before    *HelmInitHooks
	After     *HelmInitHooks
	Repo      string
	URL       string
}

func NewHelmConfig() *HelmConfig {
	return &HelmConfig{
		Chart: NewHelmChart(),
		Repo:  "",
		URL:   "",
	}
}

type HelmChart struct {
	Chart   string
	Version string
}

func NewHelmChart() *HelmChart {
	return &HelmChart{
		Chart:   "",
		Version: "",
	}
}

type HelmInitHooks struct {
	Uploads  []HelmHookUpload
	Resolves []string
	Renames  []HelmHookRename
	Concats  []HelmHookConcat
}

func NewHelmInitHooks() *HelmInitHooks {
	return &HelmInitHooks{
		Uploads:  []HelmHookUpload{},
		Resolves: []string{},
		Renames:  []HelmHookRename{},
		Concats:  []HelmHookConcat{},
	}
}

type HelmHookUpload struct {
	Name string
	URL  string
}

func NewHelmHookUpload() *HelmHookUpload {
	return &HelmHookUpload{
		Name: "",
		URL:  "",
	}
}

type HelmHookRename struct {
	Original string
	Renamed  string
}

func NewHelmHookRename() *HelmHookRename {
	return &HelmHookRename{
		Original: "",
		Renamed:  "",
	}
}

type HelmHookConcat struct {
	Folder   string
	Includes []string
	Output   string
}

func NewHelmHookConcat() *HelmHookConcat {
	return &HelmHookConcat{
		Folder:   "",
		Includes: []string{},
		Output:   "",
	}
}

func GetComponentPrefix(component string) string {
	if idx := strings.LastIndex(component, "/"); idx >= 0 {
		return component[:idx]
	}
	return ""
}
