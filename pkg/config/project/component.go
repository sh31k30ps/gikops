package project

type ProjectComponent struct {
	Name    string
	Require []string
}

func NewProjectComponent() *ProjectComponent {
	return &ProjectComponent{
		Name:    "",
		Require: []string{},
	}
}
