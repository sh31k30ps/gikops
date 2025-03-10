package setup

import "github.com/sh31k30ps/gikopsctl/pkg/config/project"

type SetupConfig struct {
	Name       string
	Provider   project.KindConfigProvider
	Components map[string][]string
}

func NewSetupConfig() *SetupConfig {
	return &SetupConfig{
		Components: make(map[string][]string),
	}
}
