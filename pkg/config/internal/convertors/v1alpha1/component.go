package v1alpha1

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/api/config/v1alpha1"
	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/internal/encoding"
)

type ComponentConverter struct{}

func NewComponentConverter() encoding.ConfigConverter {
	return &ComponentConverter{}
}

func (c *ComponentConverter) ToFile(config encoding.ConfigObject) (encoding.ConfigFile, error) {
	if config == nil {
		return nil, fmt.Errorf("no component configuration provided")
	}
	if config, ok := config.(*component.Component); ok {
		return ConvertComponentToV1Alpha1(config)
	}
	return nil, fmt.Errorf("invalid component configuration")
}

func (c *ComponentConverter) FromFile(config encoding.ConfigFile) (encoding.ConfigObject, error) {
	if config == nil {
		return nil, fmt.Errorf("no component configuration provided")
	}
	if config, ok := config.(*v1alpha1.Component); ok {
		return ConvertV1Alpha1ToComponent(config)
	}
	return nil, fmt.Errorf("invalid component configuration")
}

func (c *ComponentConverter) GetVersion() string {
	return v1alpha1.Version
}

func (c *ComponentConverter) GetKind() string {
	return v1alpha1.ComponentKind
}

func (c *ComponentConverter) GetConfigFile() encoding.ConfigFile {
	return &v1alpha1.Component{}
}

func (c *ComponentConverter) Validate(config encoding.ConfigObject) []error {
	if config == nil {
		return []error{fmt.Errorf("no component configuration provided")}
	}
	if config, ok := config.(*component.Component); ok {
		return component.ValidateComponent(*config)
	}
	return []error{fmt.Errorf("invalid component configuration")}
}

func ConvertV1Alpha1ToComponent(cfg *v1alpha1.Component) (*component.Component, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no component configuration provided")
	}
	cfg = cfg.DeepCopy()
	v1alpha1.SetComponentDefaults(cfg)
	if cfg.Metadata == nil {
		cfg.Metadata = &v1alpha1.ComponentMetadata{}
	}

	c := &component.Component{
		Name:      cfg.Metadata.Name,
		Namespace: cfg.Metadata.Namespace,
		Disabled:  cfg.Metadata.Disabled,
		DependsOn: cfg.DependsOn,
	}

	if cfg.Helm != nil {
		c.Helm = &component.HelmConfig{}
		ConvertV1Alpha1ToHelmConfig(cfg.Helm, c.Helm)
	}
	if cfg.Files != nil {
		c.Files = &component.ComponentFiles{}
		ConvertV1Alpha1ToComponentFiles(cfg.Files, c.Files)
	}
	if cfg.Exec != nil {
		c.Exec = &component.ComponentExec{}
		ConvertV1Alpha1ToComponentExec(cfg.Exec, c.Exec)
	}

	return c, nil
}

func ConvertV1Alpha1ToHelmConfig(in *v1alpha1.HelmConfig, out *component.HelmConfig) {
	out.Repo = in.Repo
	out.URL = in.RepoURL
	if in.Chart != "" {
		out.Chart = &component.HelmChart{
			Chart:   in.Chart,
			Version: in.Version,
		}
	}
	if in.CRDsChart != "" {
		out.CRDsChart = &component.HelmChart{
			Chart:   in.CRDsChart,
			Version: in.CRDsVersion,
		}
	}
	if in.Before != nil {
		out.Before = &component.HelmInitHooks{}
		ConvertV1Alpha1ToComponentHelmInitHooks(in.Before, out.Before)
	}
	if in.After != nil {
		out.After = &component.HelmInitHooks{}
		ConvertV1Alpha1ToComponentHelmInitHooks(in.After, out.After)
	}
}

func ConvertV1Alpha1ToComponentHelmInitHooks(in interface{}, out *component.HelmInitHooks) {
	if in, ok := in.(*v1alpha1.HelmBeforeInitConfig); ok {
		if len(in.Uploads) > 0 {
			out.Uploads = ConvertV1Alpha1ToComponentHelmInitHooksUploads(in.Uploads)
		}
	}
	if in, ok := in.(*v1alpha1.HelmAfterInitConfig); ok {
		if len(in.Uploads) > 0 {
			out.Uploads = ConvertV1Alpha1ToComponentHelmInitHooksUploads(in.Uploads)
		}
		if len(in.Renames) > 0 {
			out.Renames = ConvertV1Alpha1ToComponentHelmInitHooksRenames(in.Renames)
		}
		if len(in.Concats) > 0 {
			out.Concats = ConvertV1Alpha1ToComponentHelmInitHooksConcats(in.Concats)
		}
		out.Resolves = in.Resolves
	}
}

func ConvertV1Alpha1ToComponentHelmInitHooksUploads(in []v1alpha1.Upload) []component.HelmHookUpload {
	out := []component.HelmHookUpload{}
	for _, u := range in {
		out = append(out, component.HelmHookUpload{
			Name: u.Name,
			URL:  u.URL,
		})
	}
	return out
}

func ConvertV1Alpha1ToComponentHelmInitHooksRenames(in []v1alpha1.Rename) []component.HelmHookRename {
	out := []component.HelmHookRename{}
	for _, r := range in {
		out = append(out, component.HelmHookRename{
			Original: r.Original,
			Renamed:  r.Renamed,
		})
	}
	return out
}

func ConvertV1Alpha1ToComponentHelmInitHooksConcats(in []v1alpha1.Concat) []component.HelmHookConcat {
	out := []component.HelmHookConcat{}
	for _, c := range in {
		out = append(out, component.HelmHookConcat{
			Folder:   c.Folder,
			Includes: c.Includes,
			Output:   c.Output,
		})
	}
	return out
}

func ConvertV1Alpha1ToComponentFiles(in *v1alpha1.ComponentFiles, out *component.ComponentFiles) {
	out.CRDs = in.CRDs
	out.Keep = in.Keep
}

func ConvertV1Alpha1ToComponentExec(in *v1alpha1.ComponentExec, out *component.ComponentExec) {
	out.Before = in.Before
	out.After = in.After
}

func ConvertComponentToV1Alpha1(cfg *component.Component) (*v1alpha1.Component, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no component configuration provided")
	}
	c := &v1alpha1.Component{
		Metadata: &v1alpha1.ComponentMetadata{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
			Disabled:  cfg.Disabled,
		},
	}

	if cfg.Helm != nil {
		c.Helm = &v1alpha1.HelmConfig{}
		ConvertHelmConfigToV1Alpha1(cfg.Helm, c.Helm)
	}
	if cfg.Files != nil {
		c.Files = &v1alpha1.ComponentFiles{}
		ConvertComponentFilesToV1Alpha1(cfg.Files, c.Files)
	}
	if cfg.Exec != nil {
		c.Exec = &v1alpha1.ComponentExec{}
		ConvertComponentExecToV1Alpha1(cfg.Exec, c.Exec)
	}
	if len(cfg.DependsOn) > 0 {
		c.DependsOn = cfg.DependsOn
	}
	return c, nil
}

func ConvertHelmConfigToV1Alpha1(in *component.HelmConfig, out *v1alpha1.HelmConfig) {
	out.Repo = in.Repo
	out.RepoURL = in.URL
	if in.Chart != nil {
		out.Chart = in.Chart.Chart
		out.Version = in.Chart.Version
	}
	if in.CRDsChart != nil {
		out.CRDsChart = in.CRDsChart.Chart
		out.CRDsVersion = in.CRDsChart.Version
	}
	if in.Before != nil {
		out.Before = &v1alpha1.HelmBeforeInitConfig{}
		ConvertComponentHelmInitHooksToV1Alpha1(in.Before, out.Before)
	}
	if in.After != nil {
		out.After = &v1alpha1.HelmAfterInitConfig{}
		ConvertComponentHelmInitHooksToV1Alpha1(in.After, out.After)
	}
}

func ConvertComponentHelmInitHooksToV1Alpha1(in *component.HelmInitHooks, out interface{}) {
	if out, ok := out.(*v1alpha1.HelmBeforeInitConfig); ok {
		if len(in.Uploads) > 0 {
			ConvertComponentHookUploadsToV1Alpha1(in.Uploads, out.Uploads)
		}
	}
	if out, ok := out.(*v1alpha1.HelmAfterInitConfig); ok {
		if len(in.Uploads) > 0 {
			ConvertComponentHookUploadsToV1Alpha1(in.Uploads, out.Uploads)
		}
		if len(in.Renames) > 0 {
			ConvertComponentHookRenamesToV1Alpha1(in.Renames, out.Renames)
		}
		if len(in.Concats) > 0 {
			ConvertComponentHookConcatsToV1Alpha1(in.Concats, out.Concats)
		}
		if len(in.Resolves) > 0 {
			out.Resolves = in.Resolves
		}
	}

}

func ConvertComponentHookUploadsToV1Alpha1(in []component.HelmHookUpload, out []v1alpha1.Upload) {
	for _, u := range in {
		out = append(out, v1alpha1.Upload{
			Name: u.Name,
			URL:  u.URL,
		})
	}
}

func ConvertComponentHookRenamesToV1Alpha1(in []component.HelmHookRename, out []v1alpha1.Rename) {
	for _, r := range in {
		out = append(out, v1alpha1.Rename{
			Original: r.Original,
			Renamed:  r.Renamed,
		})
	}
}

func ConvertComponentHookConcatsToV1Alpha1(in []component.HelmHookConcat, out []v1alpha1.Concat) {
	for _, c := range in {
		out = append(out, v1alpha1.Concat{
			Folder:   c.Folder,
			Includes: c.Includes,
			Output:   c.Output,
		})
	}
}

func ConvertComponentFilesToV1Alpha1(in *component.ComponentFiles, out *v1alpha1.ComponentFiles) {
	out.CRDs = in.CRDs
	out.Keep = in.Keep
}

func ConvertComponentExecToV1Alpha1(in *component.ComponentExec, out *v1alpha1.ComponentExec) {
	out.Before = in.Before
	out.After = in.After
}
