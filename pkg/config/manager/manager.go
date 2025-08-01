package manager

import (
	"github.com/sh31k30ps/gikops/pkg/config/internal/convertors/v1alpha1"
	"github.com/sh31k30ps/gikops/pkg/config/internal/encoding"
)

var (
	DefaultConfigManager encoding.ConfigManager
)

func GetConfigManager() encoding.ConfigManager {
	if DefaultConfigManager == nil {
		DefaultConfigManager = encoding.NewConfigManager()
		DefaultConfigManager.AddConverter(v1alpha1.NewProjectConverter())
		DefaultConfigManager.AddConverter(v1alpha1.NewComponentConverter())
	}
	return DefaultConfigManager
}
