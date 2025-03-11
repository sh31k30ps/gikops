package project

type BasicCluster struct {
	name string
}

func NewBasicCluster() *BasicCluster {
	return &BasicCluster{
		name: "",
	}
}

func (c *BasicCluster) Name() string {
	return c.name
}

func (c *BasicCluster) Config() ClusterConfig {
	return nil
}

func (c *BasicCluster) SetName(name string) {
	c.name = name
}

func (c *BasicCluster) SetConfig(config ClusterConfig) error {
	return nil
}

func (c *BasicCluster) GetClusterName() string {
	return c.Name()
}

func (c *BasicCluster) GetContext() string {
	return c.GetClusterName()
}
