package v1alpha1

const (
	// ProjectVersion is the API version for the project configuration
	Version = "gikopsctl.config.k8s.io/v1alpha1"
)

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}
