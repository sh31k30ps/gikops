package component

type ApplyMode string

const (
	ApplyModeAll       ApplyMode = "all"
	ApplyModeCRDs      ApplyMode = "crds"
	ApplyModeManifests ApplyMode = "manifests"
)

var (
	ApplyModes = []ApplyMode{
		ApplyModeAll,
		ApplyModeCRDs,
		ApplyModeManifests,
	}
	ApplyModesLabels = []string{
		string(ApplyModeAll),
		string(ApplyModeCRDs),
		string(ApplyModeManifests),
	}
)

func (m ApplyMode) String() string {
	return string(m)
}

func (m ApplyMode) IsValid() bool {
	for _, am := range ApplyModes {
		if m == am {
			return true
		}
	}
	return false
}

type ComponentType string

const (
	ComponentTypeInternal ComponentType = "internal"
	ComponentTypeGit      ComponentType = "git"
)

var (
	ComponentTypes = []ComponentType{
		ComponentTypeInternal,
		ComponentTypeGit,
	}
	ComponentTypesLabels = []string{
		string(ComponentTypeInternal),
		string(ComponentTypeGit),
	}
)

func (t ComponentType) String() string {
	return string(t)
}

func (t ComponentType) IsValid() bool {
	for _, ct := range ComponentTypes {
		if t == ct {
			return true
		}
	}
	return false
}

type Initializer interface {
	Init(name string, keepTmp bool) error
}

type Applier interface {
	Apply() error
}
