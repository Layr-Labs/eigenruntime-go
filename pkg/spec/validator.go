package spec

import (
	"fmt"
	"reflect"
)

type Validator interface {
	Validate(spec interface{}) error
	ValidateField(field string, value interface{}) error
}

type ValidationRule func(value interface{}) error

type validator struct {
	rules map[string]ValidationRule
}

func NewValidator() Validator {
	return &validator{
		rules: make(map[string]ValidationRule),
	}
}

func NewValidatorWithRules(rules map[string]ValidationRule) Validator {
	return &validator{
		rules: rules,
	}
}

func (v *validator) Validate(spec interface{}) error {
	if spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	
	val := reflect.ValueOf(spec)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("spec must be a struct or pointer to struct")
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		
		if rule, exists := v.rules[field.Name]; exists {
			if err := rule(fieldValue.Interface()); err != nil {
				return fmt.Errorf("validation failed for field %s: %w", field.Name, err)
			}
		}
	}
	
	return nil
}

func (v *validator) ValidateField(field string, value interface{}) error {
	if rule, exists := v.rules[field]; exists {
		return rule(value)
	}
	return nil
}

func RequiredString(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if s == "" {
		return fmt.Errorf("field is required")
	}
	return nil
}

func ValidateVersion(value interface{}) error {
	version, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if version == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}