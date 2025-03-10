package project

import "github.com/sh31k30ps/gikopsctl/pkg/config/component"

func SetProjectDefaults(p *Project) {
	if p == nil {
		p = NewProject()
	}
	SetProjectComponentsDefaults(p.Components)
	SetKindClusterDefaults(p.LocalCluster)
	if len(p.Environments) == 0 {
		p.Environments = []string{"local"}
	}
}

func SetProjectComponentsDefaults(c *ProjectComponents) {
	if c == nil {
		c = NewProjectComponents()
	}
	if c.FileName == "" {
		c.FileName = component.ComponentFileName
	}
	if len(c.Folders) == 0 {
		c.Folders = []string{
			"components",
			"core",
		}
	}
}

func SetKindClusterDefaults(c *KindCluster) {
	if c == nil {
		c = NewKindCluster()
	}
	SetKindConfigDefaults(c.KindConfig)
}

func SetKindConfigDefaults(c *KindConfig) {
	if c == nil {
		c = NewKindConfig()
	}
	if c.ConfigFile == "" {
		c.ConfigFile = DefaultKindConfigFile
	}

	if len(c.OverridesFolder) == 0 {
		c.OverridesFolder = []string{
			"overrides",
		}
	}

	if c.Provider == "" {
		c.Provider = KindConfigProviderDocker
	}
}
