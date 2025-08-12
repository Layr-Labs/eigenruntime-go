package spec

import (
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenruntime-go/pkg/common"
	"gopkg.in/yaml.v3"
)

func ParseYAML(data []byte) (*common.RuntimeSpec, error) {
	var spec common.RuntimeSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &spec, nil
}

func ParseJSON(data []byte) (*common.RuntimeSpec, error) {
	var spec common.RuntimeSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &spec, nil
}

func ToYAML(spec *common.RuntimeSpec) ([]byte, error) {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return data, nil
}

func ToJSON(spec *common.RuntimeSpec) ([]byte, error) {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return data, nil
}