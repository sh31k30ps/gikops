package v1alpha1

type KustomizeConfig struct {
	URLs []string `json:"urls,omitempty" yaml:"urls,omitempty"`
}
