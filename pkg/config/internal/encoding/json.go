package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JSONParser implémente ProjectParser pour JSON
type JSONParser struct {
	unmarshal       func([]byte, interface{}) error
	unmarshalStrict func([]byte, ConfigFile) error
	marshal         func(interface{}) ([]byte, error)
	manager         ConfigManager
}

func NewJSONParser(manager ConfigManager) ConfigParser {
	return &JSONParser{
		unmarshal:       json.Unmarshal,
		unmarshalStrict: jsonUnmarshalStrict,
		marshal:         json.Marshal,
		manager:         manager,
	}
}

func (p *JSONParser) Parse(raw []byte) (ConfigObject, error) {
	var tm TypeMeta
	if err := p.unmarshal(raw, &tm); err != nil {
		return nil, fmt.Errorf("error parsing document: %w", err)
	}
	converter := p.manager.GetConverter(tm.APIVersion, tm.Kind)
	if converter == nil {
		return nil, fmt.Errorf("converter not found for  %s/%s", tm.APIVersion, tm.Kind)
	}
	return parseConfig(raw, p.unmarshal, p.unmarshalStrict, converter)
}

func (p *JSONParser) Generate(config ConfigObject, converter ConfigConverter) ([]byte, error) {
	if converter == nil {
		return nil, fmt.Errorf("converter is required")
	}
	return generateConfigFileContent(config, p.marshal, converter)
}

func jsonUnmarshalStrict(raw []byte, c ConfigFile) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	return d.Decode(c)
}
