package cluster

func SetKindClusterDefaults(c *KindCluster) {
	cfg := c.Config().(*KindConfig)
	if cfg == nil {
		cfg = NewKindConfig()
	}
	SetKindConfigDefaults(cfg)
	c.SetConfig(cfg)
}

func SetKindConfigDefaults(c *KindConfig) {
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
