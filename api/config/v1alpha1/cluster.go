package v1alpha1

// ClusterLocal contains the local cluster configuration
type ClusterLocal struct {
	// Name is the name of the cluster
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// KindConfig contains the configuration for a local cluster
	KindConfig *ClusterLocalKindConfig `json:"kindConfig,omitempty" yaml:"kindConfig,omitempty"`
}

// ClusterLocalKindConfig contains the configuration for a local cluster
type ClusterLocalKindConfig struct {
	// ConfigFile is the path to the Kind configuration file
	ConfigFile string `json:"configFile,omitempty" yaml:"configFile,omitempty"`

	// OverridesFolder is the path to the folder containing overrides
	OverridesFolder []string `json:"overridesFolder,omitempty" yaml:"overridesFolder,omitempty"`

	// Provider is the provider for the cluster
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}
