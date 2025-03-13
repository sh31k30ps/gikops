package v1alpha1

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/api/config/v1alpha1"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
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
		cfg, err := ConvertV1Alpha1ToProject(config)
		if err != nil {
			return nil, err
		}
		errs := c.Validate(cfg)
		if len(errs) > 0 {
			return cfg, fmt.Errorf("invalid project configuration")
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("invalid project configuration")
}

func (c *ProjectConverter) Validate(config encoding.ConfigObject) []error {
	if config == nil {
		return []error{fmt.Errorf("no project configuration provided")}
	}
	if config, ok := config.(*project.Project); ok {
		return project.Validate(*config)
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

	p := project.NewConfig()
	p.Name = cfg.Metadata.Name

	if len(cfg.Components) > 0 {
		p.Components = make([]project.ProjectComponent, len(cfg.Components))
		for i, component := range cfg.Components {
			ConvertV1Alpha1ToProjectComponent(&component, &p.Components[i])
		}
	}

	if len(cfg.Clusters) > 0 {
		p.Clusters = make([]cluster.Cluster, len(cfg.Clusters))
		for i, c := range cfg.Clusters {
			if c.KindConfig != nil {
				p.Clusters[i] = cluster.NewKindCluster()
				ConvertV1Alpha1ToKindCluster(&c, p.Clusters[i].(*cluster.KindCluster))
			} else {
				p.Clusters[i] = cluster.NewBasicCluster()
				p.Clusters[i].(*cluster.BasicCluster).SetName(c.Name)
			}
		}
	}
	project.SetProjectDefaults(p)
	return p, nil
}

func ConvertV1Alpha1ToProjectComponent(in *v1alpha1.ProjectComponent, out *project.ProjectComponent) {
	out.Name = in.Name
	out.Require = in.Require
}

func ConvertV1Alpha1ToKindCluster(in *v1alpha1.Cluster, out *cluster.KindCluster) {
	out.SetName(in.Name)
	if in.KindConfig != nil {
		cfg := cluster.NewKindConfig()
		ConvertV1Alpha1ToKindConfig(in.KindConfig, cfg)
		out.SetConfig(cfg)
	}
}

func ConvertV1Alpha1ToKindConfig(in *v1alpha1.ClusterKindConfig, out *cluster.KindConfig) {
	out.ConfigFile = in.ConfigFile
	out.OverridesFolder = in.OverridesFolder
	out.Provider = cluster.KindConfigProvider(in.Provider)
}

func ConvertProjectToV1Alpha1(cfg *project.Project) (*v1alpha1.Project, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no project configuration provided")
	}
	project.SetProjectDefaults(cfg)
	p := &v1alpha1.Project{
		Metadata: &v1alpha1.ProjectMetadata{
			Name: cfg.Name,
		},
	}

	if len(cfg.Clusters) > 0 {
		p.Clusters = make([]v1alpha1.Cluster, len(cfg.Clusters))
		for i, c := range cfg.Clusters {
			ConvertClusterToV1Alpha1(c, &p.Clusters[i])
		}
	}

	if len(cfg.Components) > 0 {
		p.Components = make([]v1alpha1.ProjectComponent, len(cfg.Components))
		for i, c := range cfg.Components {
			ConvertProjectComponentToV1Alpha1(c, &p.Components[i])
		}
	}
	v1alpha1.SetProjectDefaults(p)
	return p, nil
}

func ConvertProjectComponentToV1Alpha1(in project.ProjectComponent, out *v1alpha1.ProjectComponent) {
	out.Name = in.Name
	out.Require = in.Require
}

func ConvertClusterToV1Alpha1(in cluster.Cluster, out *v1alpha1.Cluster) {
	switch c := in.(type) {
	case *cluster.KindCluster:
		ConvertKindClusterToV1Alpha1(c, out)
	case *cluster.BasicCluster:
		out.Name = c.Name()
	default:
	}
}

func ConvertKindClusterToV1Alpha1(in *cluster.KindCluster, out *v1alpha1.Cluster) {
	if in == nil {
		return
	}
	out.Name = in.Name()
	if in.Config() != nil {
		out.KindConfig = &v1alpha1.ClusterKindConfig{}
		ConvertKindConfigToV1Alpha1(in.Config().(*cluster.KindConfig), out.KindConfig)
	}
}

func ConvertKindConfigToV1Alpha1(in *cluster.KindConfig, out *v1alpha1.ClusterKindConfig) {
	out.ConfigFile = in.ConfigFile
	out.OverridesFolder = in.OverridesFolder
	out.Provider = string(in.Provider)
}
