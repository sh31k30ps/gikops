package v1alpha1

// ComponentFiles contains file-based configuration
type ComponentFiles struct {
	// CRD is the path to custom resource definitions
	CRDs string `json:"crds,omitempty" yaml:"crds,omitempty"`

	// SkipCRDs specifies if CRDs should be skipped
	SkipCRDs bool `json:"skipCRDs,omitempty" yaml:"skipCRDs,omitempty"`

	// Keep specifies files to preserve during operations
	Keep []string `json:"keep,omitempty" yaml:"keep,omitempty"`
}
