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
	Components []ProjectComponent `json:"components,omitempty" yaml:"components,omitempty"`

	// ClusterLocal contains the local cluster configuration
	Clusters []Cluster `json:"clusters,omitempty" yaml:"clusters,omitempty"`
}

// ProjectMetadata contains metadata for the project
// It is used to identify the project
type ProjectMetadata struct {
	// Name is the name of the project
	Name string `json:"name" yaml:"name"`
}

// ProjectComponents contains the components for the project
type ProjectComponent struct {
	// Folders is a list of folders containing component configurations
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Require is a list of components that are required by the component
	Require []string `json:"require,omitempty" yaml:"require,omitempty"`
}
