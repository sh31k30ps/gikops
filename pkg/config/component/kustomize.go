package component

type KustomizeConfig struct {
	URLs []string
}

func NewKustomizeConfig() *KustomizeConfig {
	return &KustomizeConfig{
		URLs: []string{},
	}
}
