package encoding

const (
	DefaultAPIVersion = "gikopsctl.config.k8s.io/v1alpha1"
)

type ConfigFile any
type ConfigObject any

type TypeMeta struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
}

// ProjectLoader définit l'interface pour charger et sauvegarder des projets
type ConfigManager interface {
	Load(path string) (ConfigObject, error)
	Save(path string, cfg ConfigObject, version, kind string) error
	GetConverter(version, kind string) ConfigConverter
	AddConverter(converter ConfigConverter)
}

// ProjectParser définit l'interface pour parser et générer du contenu
type ConfigParser interface {
	Parse(raw []byte) (ConfigObject, error)
	Generate(config ConfigObject, converter ConfigConverter) ([]byte, error)
}

type ConfigConverter interface {
	ToFile(config ConfigObject) (ConfigFile, error)
	FromFile(config ConfigFile) (ConfigObject, error)
	Validate(config ConfigObject) []error
	GetVersion() string
	GetKind() string
	GetConfigFile() ConfigFile
}
