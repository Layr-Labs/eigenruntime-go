package spec

import (
	"fmt"
	"testing"
)

type TestSpec struct {
	Name    string
	Version string
	Count   int
}

func TestValidator(t *testing.T) {
	rules := map[string]ValidationRule{
		"Name":    RequiredString,
		"Version": ValidateVersion,
		"Count": func(value interface{}) error {
			count, ok := value.(int)
			if !ok {
				return fmt.Errorf("expected int")
			}
			if count < 0 {
				return fmt.Errorf("count must be non-negative")
			}
			return nil
		},
	}
	
	v := NewValidatorWithRules(rules)
	
	tests := []struct {
		name    string
		spec    TestSpec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: TestSpec{
				Name:    "test",
				Version: "v1.0.0",
				Count:   5,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			spec: TestSpec{
				Name:    "",
				Version: "v1.0.0",
				Count:   5,
			},
			wantErr: true,
		},
		{
			name: "missing version",
			spec: TestSpec{
				Name:    "test",
				Version: "",
				Count:   5,
			},
			wantErr: true,
		},
		{
			name: "negative count",
			spec: TestSpec{
				Name:    "test",
				Version: "v1.0.0",
				Count:   -1,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateField(t *testing.T) {
	v := NewValidatorWithRules(map[string]ValidationRule{
		"Name": RequiredString,
	})
	
	if err := v.ValidateField("Name", "test"); err != nil {
		t.Errorf("ValidateField() should not error for valid field")
	}
	
	if err := v.ValidateField("Name", ""); err == nil {
		t.Errorf("ValidateField() should error for empty required field")
	}
	
	if err := v.ValidateField("UnknownField", "anything"); err != nil {
		t.Errorf("ValidateField() should not error for unknown field")
	}
}

func TestValidateNilSpec(t *testing.T) {
	v := NewValidator()
	if err := v.Validate(nil); err == nil {
		t.Errorf("Validate() should error for nil spec")
	}
}

func TestValidateNonStruct(t *testing.T) {
	v := NewValidator()
	if err := v.Validate("not a struct"); err == nil {
		t.Errorf("Validate() should error for non-struct")
	}
}

func TestRequiredString(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid string", "test", false},
		{"empty string", "", true},
		{"not a string", 123, true},
		{"nil", nil, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequiredString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequiredString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid version", "v1.0.0", false},
		{"simple version", "v1", false},
		{"empty version", "", true},
		{"not a string", 1.0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}