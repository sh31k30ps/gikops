package v1alpha1

const (
	// ProjectKind is the kind for the project configuration
	ProjectKind = "Project"
)

// Project contains the configuration for a project
// A project is a collection of components that are deployed together
// It is the root of the configuration
type Project struct {
	TypeMeta `json:",inline" yaml:",inline"`

	// Metadata contains the project metadata
	Metadata *ProjectMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// ComponentsFolders is a list of folders containing component configurations
	Components *ProjectComponents `json:"components,omitempty" yaml:"components,omitempty"`

	// ClusterLocal contains the local cluster configuration
	ClusterLocal *ClusterLocal `json:"clusterLocal,omitempty" yaml:"clusterLocal,omitempty"`

	// Environments is a list of environments for the project
	Environments []string `json:"environments,omitempty" yaml:"environments,omitempty"`
}

// ProjectMetadata contains metadata for the project
// It is used to identify the project
type ProjectMetadata struct {
	// Name is the name of the project
	Name string `json:"name" yaml:"name"`
}

// ProjectComponents contains the components for the project
type ProjectComponents struct {
	// Folders is a list of folders containing component configurations
	Folders []string `json:"folders,omitempty" yaml:"folders,omitempty"`
	// FileName is the name of the file containing the component configurations
	FileName string `json:"fileName,omitempty" yaml:"fileName,omitempty"`
}
