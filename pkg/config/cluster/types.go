package cluster

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
