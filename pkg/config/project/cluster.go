package project

type ProjectCluster interface {
	Name() string
	Config() ClusterConfig
	SetName(name string)
	SetConfig(config ClusterConfig) error
	GetClusterName() string
	GetContext() string
}

type ClusterConfig interface{}

type ClusterType string

const (
	ClusterTypeKind  ClusterType = "kind"
	ClusterTypeBasic ClusterType = "basic"
)

var ClusterTypes = []ClusterType{
	ClusterTypeKind,
	ClusterTypeBasic,
}

var ClusterTypesLabels = []string{
	string(ClusterTypeKind),
	string(ClusterTypeBasic),
}
