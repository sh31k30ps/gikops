package component

type Component struct {
	Name      string
	Namespace string
	Disabled  bool
	DependsOn []string
	Helm      *HelmConfig
	Kustomize *KustomizeConfig
	Files     *FilesConfig
	Exec      *ExecConfig
	Clusters  []string
}

func NewComponent() *Component {
	return &Component{
		Name:      "",
		Namespace: "",
		Disabled:  false,
		DependsOn: []string{},
		Files:     NewFilesConfig(),
	}
}
