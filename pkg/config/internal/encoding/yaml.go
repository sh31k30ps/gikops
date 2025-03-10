package encoding

import (
	"bytes"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

// YAMLParser implémente ProjectParser pour YAML
type YAMLParser struct {
	unmarshal       func([]byte, interface{}) error
	unmarshalStrict func([]byte, ConfigFile) error
	marshal         func(interface{}) ([]byte, error)
	manager         ConfigManager
}

func NewYAMLParser(manager ConfigManager) ConfigParser {
	return &YAMLParser{
		unmarshal:       yaml.Unmarshal,
		unmarshalStrict: yamlUnmarshalStrict,
		marshal:         yaml.Marshal,
		manager:         manager,
	}
}

func (p *YAMLParser) Parse(raw []byte) (ConfigObject, error) {
	var tm TypeMeta
	if err := p.unmarshal(raw, &tm); err != nil {
		return nil, fmt.Errorf("error parsing document: %w", err)
	}
	converter := p.manager.GetConverter(tm.APIVersion, tm.Kind)

	return parseConfig(raw, p.unmarshal, p.unmarshalStrict, converter)
}

func (p *YAMLParser) Generate(config ConfigObject, converter ConfigConverter) ([]byte, error) {
	if converter == nil {
		return nil, fmt.Errorf("converter is required")
	}
	return generateConfigFileContent(config, p.marshal, converter)
}

func yamlUnmarshalStrict(raw []byte, c ConfigFile) error {
	d := yaml.NewDecoder(bytes.NewReader(raw))
	d.KnownFields(true)
	return d.Decode(c)
}
