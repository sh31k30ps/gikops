package cluster

type Cluster interface {
	Name() string
	Config() ClusterConfig
	SetName(name string)
	SetConfig(config ClusterConfig) error
	GetClusterName() string
	GetContext() string
}

func DefaultKindCluster() *KindCluster {
	cluster := NewKindCluster()
	cluster.SetName("local")
	cluster.SetConfig(&KindConfig{
		Provider:   KindConfigProviderDocker,
		ConfigFile: "kind.yaml",
		OverridesFolder: []string{
			"overrides",
		},
	})
	return cluster
}
