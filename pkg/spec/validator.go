package spec

import (
	"fmt"

	"github.com/Layr-Labs/eigenruntime-go/pkg/common"
)

func ValidateRuntimeSpec(spec *common.RuntimeSpec) error {
	if spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}

	if spec.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}

	if spec.Kind == "" {
		return fmt.Errorf("kind is required")
	}

	if spec.Name == "" {
		return fmt.Errorf("name is required")
	}

	if spec.Version == "" {
		return fmt.Errorf("version is required")
	}

	if spec.Spec == nil || len(spec.Spec) == 0 {
		return fmt.Errorf("spec must contain at least one component")
	}

	for name, component := range spec.Spec {
		if err := ValidateComponent(name, &component); err != nil {
			return fmt.Errorf("validation failed for component %s: %w", name, err)
		}
	}

	return nil
}

func ValidateComponent(name string, component *common.Component) error {
	if component.Registry == "" {
		return fmt.Errorf("registry is required")
	}

	if component.Digest == "" {
		return fmt.Errorf("digest is required")
	}

	for _, env := range component.Env {
		if env.Name == "" {
			return fmt.Errorf("environment variable name cannot be empty")
		}
	}

	return nil
}