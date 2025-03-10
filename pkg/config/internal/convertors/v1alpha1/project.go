package v1alpha1

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/api/config/v1alpha1"
	"github.com/sh31k30ps/gikopsctl/pkg/config/internal/encoding"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

type ProjectConverter struct{}

func NewProjectConverter() encoding.ConfigConverter {
	return &ProjectConverter{}
}

func (c *ProjectConverter) ToFile(config encoding.ConfigObject) (encoding.ConfigFile, error) {
	if config == nil {
		return nil, fmt.Errorf("no project configuration provided")
	}
	if config, ok := config.(*project.Project); ok {
		return ConvertProjectToV1Alpha1(config)
	}
	return nil, fmt.Errorf("invalid project configuration")
}

func (c *ProjectConverter) FromFile(config encoding.ConfigFile) (encoding.ConfigObject, error) {
	if config == nil {
		return nil, fmt.Errorf("no project configuration provided")
	}
	if config, ok := config.(*v1alpha1.Project); ok {
		return ConvertV1Alpha1ToProject(config)
	}
	return nil, fmt.Errorf("invalid project configuration")
}

func (c *ProjectConverter) Validate(config encoding.ConfigObject) []error {
	if config == nil {
		return []error{fmt.Errorf("no project configuration provided")}
	}
	if config, ok := config.(*project.Project); ok {
		return project.ValidateProject(*config)
	}
	return []error{fmt.Errorf("invalid project configuration")}
}

func (c *ProjectConverter) GetVersion() string {
	return v1alpha1.Version
}

func (c *ProjectConverter) GetKind() string {
	return v1alpha1.ProjectKind
}

func (c *ProjectConverter) GetConfigFile() encoding.ConfigFile {
	return &v1alpha1.Project{}
}

func ConvertV1Alpha1ToProject(cfg *v1alpha1.Project) (*project.Project, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no project configuration provided")
	}
	cfg = cfg.DeepCopy()
	v1alpha1.SetProjectDefaults(cfg)

	p := project.NewProject()
	p.Name = cfg.Metadata.Name

	if cfg.Components != nil {
		ConvertV1Alpha1ToProjectComponents(cfg.Components, p.Components)
	}
	if cfg.ClusterLocal != nil {
		ConvertV1Alpha1ToKindCluster(cfg.ClusterLocal, p.LocalCluster)
	}
	if cfg.Environments != nil {
		p.Environments = cfg.Environments
	}

	return p, nil
}

func ConvertV1Alpha1ToProjectComponents(in *v1alpha1.ProjectComponents, out *project.ProjectComponents) {
	out.Folders = in.Folders
	out.FileName = in.FileName
}

func ConvertV1Alpha1ToKindCluster(in *v1alpha1.ClusterLocal, out *project.KindCluster) {
	out.Name = in.Name
	ConvertV1Alpha1ToKindConfig(in.KindConfig, out.KindConfig)
}

func ConvertV1Alpha1ToKindConfig(in *v1alpha1.ClusterLocalKindConfig, out *project.KindConfig) {
	out.ConfigFile = in.ConfigFile
	out.OverridesFolder = in.OverridesFolder
	out.Provider = project.KindConfigProvider(in.Provider)
}

func ConvertProjectToV1Alpha1(cfg *project.Project) (*v1alpha1.Project, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no project configuration provided")
	}

	p := &v1alpha1.Project{
		TypeMeta: v1alpha1.TypeMeta{
			APIVersion: v1alpha1.Version,
			Kind:       v1alpha1.ProjectKind,
		},
		Metadata: &v1alpha1.ProjectMetadata{
			Name: cfg.Name,
		},
	}
	if cfg.Components != nil {
		p.Components = &v1alpha1.ProjectComponents{}
		ConvertProjectComponentsToV1Alpha1(cfg.Components, p.Components)
	}
	if cfg.LocalCluster != nil {
		p.ClusterLocal = &v1alpha1.ClusterLocal{}
		ConvertKindClusterToV1Alpha1(cfg.LocalCluster, p.ClusterLocal)
	}
	v1alpha1.SetProjectDefaults(p)
	return p, nil
}

func ConvertProjectComponentsToV1Alpha1(in *project.ProjectComponents, out *v1alpha1.ProjectComponents) {
	out.Folders = in.Folders
	out.FileName = in.FileName
}

func ConvertKindClusterToV1Alpha1(in *project.KindCluster, out *v1alpha1.ClusterLocal) {
	if in.KindConfig != nil {
		out.KindConfig = &v1alpha1.ClusterLocalKindConfig{}
		ConvertKindConfigToV1Alpha1(in.KindConfig, out.KindConfig)
	}
}

func ConvertKindConfigToV1Alpha1(in *project.KindConfig, out *v1alpha1.ClusterLocalKindConfig) {
	out.ConfigFile = in.ConfigFile
	out.OverridesFolder = in.OverridesFolder
	out.Provider = string(in.Provider)
}
