package v1alpha1

// ComponentExec contains execution configuration
type ComponentExec struct {
	// Pre contains commands to execute before deployment
	Before []string `json:"before,omitempty" yaml:"before,omitempty"`

	After []string `json:"after,omitempty" yaml:"after,omitempty"`
}
