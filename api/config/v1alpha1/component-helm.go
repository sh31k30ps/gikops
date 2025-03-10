package v1alpha1

// HelmConfig contains Helm-specific configuration for a component
type HelmConfig struct {
	// Repo is the name of the Helm repository
	Repo string `json:"repo" yaml:"repo"`

	// RepoURL is the URL of the Helm repository
	RepoURL string `json:"repo-url" yaml:"repo-url"`

	// Version is the version of the Helm chart to use
	Version string `json:"version" yaml:"version"`

	// Chart is the name of the Helm chart
	Chart string `json:"chart" yaml:"chart"`

	// CRDChart is the name of a separate chart containing CRDs
	// If unset, CRDs are assumed to be in the main chart
	CRDsChart string `json:"crds-chart,omitempty" yaml:"crds-chart,omitempty"`

	// CRDsVersion is the version of the CRD chart
	// If unset and CRDChart is set, Version will be used
	CRDsVersion string `json:"crds-version,omitempty" yaml:"crds-version,omitempty"`

	// Before contains actions to perform before chart installation
	Before *HelmBeforeInitConfig `json:"before,omitempty" yaml:"before,omitempty"`

	// After contains actions to perform after chart installation
	After *HelmAfterInitConfig `json:"after,omitempty" yaml:"after,omitempty"`
}

// HelmBeforeInitConfig contains pre-installation actions for a Helm chart
type HelmBeforeInitConfig struct {
	// Uploads specifies files to be uploaded before installation
	Uploads []Upload `json:"uploads,omitempty" yaml:"uploads,omitempty"`
}

// HelmAfterInitConfig contains post-installation actions for a Helm chart
type HelmAfterInitConfig struct {
	// Uploads specifies files to be uploaded after installation
	Uploads []Upload `json:"uploads,omitempty" yaml:"uploads,omitempty"`

	// Resolves targeted files if they contains only an URL
	Resolves []string `json:"resolves,omitempty" yaml:"resolves,omitempty"`

	// Renames specifies files to be renamed after installation
	Renames []Rename `json:"renames,omitempty" yaml:"renames,omitempty"`

	// Concats specifies files to be concatenated after installation
	Concats []Concat `json:"concats,omitempty" yaml:"concats,omitempty"`
}
