package v1alpha1

const (
	// ComponentKind is the kind for the component configuration
	ComponentKind = "Component"
)

// Component contains the configuration for a Kubernetes component
type Component struct {
	TypeMeta `json:",inline" yaml:",inline"`

	// Metadata contains the component metadata
	Metadata *ComponentMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// Helm contains the Helm chart configuration for this component
	// If unset, the component will not use Helm for deployment
	Helm *HelmConfig `json:"helm,omitempty" yaml:"helm,omitempty"`

	// Files contains the file-based configuration for this component
	// If unset, the component will not use direct file management
	Files *ComponentFiles `json:"files,omitempty" yaml:"files,omitempty"`

	// Exec contains commands to execute before component deployment
	// If unset, no commands will be executed
	Exec *ComponentExec `json:"exec,omitempty" yaml:"exec,omitempty"`

	// DependsOn is a list of components that must be deployed before this component
	DependsOn []string `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`

	// EnvironmentAvailability is the environment availability for the component
	EnvironmentAvailability []string `json:"environmentAvailability,omitempty" yaml:"environmentAvailability,omitempty"`
}

// ComponentMetadata contains metadata for the component
type ComponentMetadata struct {
	// Name is the name of the component
	Name string `json:"name" yaml:"name"`

	// Namespace is the Kubernetes namespace where the component will be deployed
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Disabled indicates whether this component should be skipped during deployment
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	// DependsOn is a list of components that must be deployed before this component
}
