package component

type ExecConfig struct {
	Before []string
	After  []string
}

func NewExecConfig() *ExecConfig {
	return &ExecConfig{
		Before: []string{},
		After:  []string{},
	}
}
