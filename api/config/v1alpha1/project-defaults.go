package v1alpha1

func SetProjectDefaults(p *Project) {
	if p == nil {
		return
	}
	p.TypeMeta = TypeMeta{
		APIVersion: Version,
		Kind:       ProjectKind,
	}
	if len(p.Environments) == 0 {
		p.Environments = []string{"local"}
	}
	if p.Components == nil {
		p.Components = &ProjectComponents{}
	}
	SetProjectComponentsDefaults(p.Components)
	if p.ClusterLocal == nil {
		p.ClusterLocal = &ClusterLocal{}
	}
	SetKindClusterDefaults(p.ClusterLocal)
}

func SetProjectComponentsDefaults(c *ProjectComponents) {
	if len(c.Folders) == 0 {
		c.Folders = []string{
			"components",
			"core",
		}
	}
	if c.FileName == "" {
		c.FileName = "gikcpnt"
	}
}

func SetKindClusterDefaults(c *ClusterLocal) {
	if c.KindConfig == nil {
		c.KindConfig = &ClusterLocalKindConfig{}
	}
	SetKindConfigDefaults(c.KindConfig)
}

func SetKindConfigDefaults(c *ClusterLocalKindConfig) {
	if c.ConfigFile == "" {
		c.ConfigFile = "kind.yaml"
	}
	if len(c.OverridesFolder) == 0 {
		c.OverridesFolder = []string{"overrides"}
	}
	if c.Provider == "" {
		c.Provider = "docker"
	}
}
