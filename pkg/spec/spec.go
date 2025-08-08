package spec

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	ParseYAML(data []byte, v interface{}) error
	ParseJSON(data []byte, v interface{}) error
	ToYAML(v interface{}) ([]byte, error)
	ToJSON(v interface{}) ([]byte, error)
}

type parser struct{}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) ParseYAML(data []byte, v interface{}) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	return nil
}

func (p *parser) ParseJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func (p *parser) ToYAML(v interface{}) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return data, nil
}

func (p *parser) ToJSON(v interface{}) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return data, nil
}

func ParseYAML(data []byte, v interface{}) error {
	return NewParser().ParseYAML(data, v)
}

func ParseJSON(data []byte, v interface{}) error {
	return NewParser().ParseJSON(data, v)
}

func ToYAML(v interface{}) ([]byte, error) {
	return NewParser().ToYAML(v)
}

func ToJSON(v interface{}) ([]byte, error) {
	return NewParser().ToJSON(v)
}