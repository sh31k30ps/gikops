package v1alpha1

func SetKindConfigDefaults(c *ClusterKindConfig) {
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
