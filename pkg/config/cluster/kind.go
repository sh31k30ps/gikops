package cluster

import (
	"errors"
	"fmt"
)

type KindConfigProvider string

const (
	DefaultKindConfigFile = "kind.yaml"
)

const (
	KindConfigProviderDocker  KindConfigProvider = "docker"
	KindConfigProviderPodman  KindConfigProvider = "podman"
	KindConfigProviderNerdctl KindConfigProvider = "nerdctl"
)

var KindConfigProvidersLabels = []string{
	string(KindConfigProviderDocker),
	string(KindConfigProviderPodman),
	string(KindConfigProviderNerdctl),
}

var KindConfigProviders = []KindConfigProvider{
	KindConfigProviderDocker,
	KindConfigProviderPodman,
	KindConfigProviderNerdctl,
}

type KindCluster struct {
	name       string
	kindConfig *KindConfig
}

func NewKindCluster() *KindCluster {
	return &KindCluster{
		name:       "",
		kindConfig: NewKindConfig(),
	}
}

func (c *KindCluster) Name() string {
	return c.name
}

func (c *KindCluster) Config() ClusterConfig {
	return c.kindConfig
}

func (c *KindCluster) SetName(name string) {
	c.name = name
}

func (c *KindCluster) GetClusterName() string {
	name := c.Config().(*KindConfig).ClusterName
	if name == "" {
		name = c.Name()
	}
	return name
}

func (c *KindCluster) GetContext() string {
	return fmt.Sprintf("kind-%s", c.GetClusterName())
}

func (c *KindCluster) SetConfig(config ClusterConfig) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if _, ok := config.(*KindConfig); !ok {
		return errors.New("config is not a KindConfig")
	}
	c.kindConfig = config.(*KindConfig)
	return nil
}

type KindConfig struct {
	ClusterName     string
	ConfigFile      string
	OverridesFolder []string
	Provider        KindConfigProvider
}

func NewKindConfig() *KindConfig {
	return &KindConfig{
		ClusterName:     "",
		OverridesFolder: []string{},
		Provider:        "",
		ConfigFile:      "",
	}
}
