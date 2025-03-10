package encoding

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/config"
)

// DefaultConfigManager implémente ConfigManager
type DefaultConfigManager struct {
	parsers    map[config.ConfigExtension]ConfigParser
	converters map[string]map[string]ConfigConverter
}

func NewConfigManager() ConfigManager {
	cm := &DefaultConfigManager{
		parsers:    make(map[config.ConfigExtension]ConfigParser),
		converters: make(map[string]map[string]ConfigConverter),
	}
	cm.parsers[config.ConfigExtensionJSON] = NewJSONParser(cm)
	cm.parsers[config.ConfigExtensionYAML] = NewYAMLParser(cm)
	return cm
}

func (cm *DefaultConfigManager) GetConverter(version, kind string) ConfigConverter {
	return cm.converters[version][kind]
}

func (cm *DefaultConfigManager) AddConverter(converter ConfigConverter) {
	if cm.converters[converter.GetVersion()] == nil {
		cm.converters[converter.GetVersion()] = make(map[string]ConfigConverter)
	}
	cm.converters[converter.GetVersion()][converter.GetKind()] = converter
}

// Load charge un projet depuis un fichier
func (cm *DefaultConfigManager) Load(path string) (ConfigObject, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	ext := config.GetConfigExtension(path)
	parser, ok := cm.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unknown extension: %s", ext)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return parser.Parse(raw)
}

// Save sauvegarde un projet dans un fichier
func (cm *DefaultConfigManager) Save(path string, cfg ConfigObject, version, kind string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}
	if kind == "" {
		return fmt.Errorf("config type is required")
	}
	if version == "" {
		return fmt.Errorf("version is required")
	}

	ext := config.GetConfigExtension(path)
	parser, ok := cm.parsers[ext]
	if !ok {
		return fmt.Errorf("unknown extension: %s", ext)
	}

	converter := cm.GetConverter(version, kind)
	if converter == nil {
		return fmt.Errorf("unknown converter: %s:%s", version, kind)
	}
	content, err := parser.Generate(cfg, converter)
	if err != nil {
		return fmt.Errorf("error generating project file content: %w", err)
	}

	return os.WriteFile(path, content, 0644)
}

func parseConfig(
	raw []byte,
	unmarshal func([]byte, interface{}) error,
	unmarshalStrict func([]byte, ConfigFile) error,
	converter ConfigConverter,
) (ConfigObject, error) {
	var tm TypeMeta
	if err := unmarshal(raw, &tm); err != nil {
		return nil, fmt.Errorf("error parsing document: %w", err)
	}
	if tm.APIVersion != converter.GetVersion() {
		return nil, fmt.Errorf("unknown api version %s", tm.APIVersion)
	}
	if tm.Kind != converter.GetKind() {
		return nil, fmt.Errorf("unknown api kind %s", tm.Kind)
	}

	tmpCfg := converter.GetConfigFile()
	if err := unmarshalStrict(raw, tmpCfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	cfg, err := converter.FromFile(tmpCfg)
	if err != nil {
		return nil, fmt.Errorf("error converting file: %w", err)
	}

	errs := converter.Validate(cfg)
	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid config: %v", errs)
	}

	return cfg, nil
}

func generateConfigFileContent(
	config ConfigObject,
	marshal func(interface{}) ([]byte, error),
	converter ConfigConverter,
) ([]byte, error) {
	if errs := converter.Validate(config); len(errs) > 0 {
		return nil, fmt.Errorf("invalid config: %v", errs)
	}
	cfg, err := converter.ToFile(config)
	if err != nil {
		return nil, fmt.Errorf("error converting config: %w", err)
	}
	return marshal(cfg)
}
