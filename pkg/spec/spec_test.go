package spec

import (
	"testing"
)

func TestParseYAML(t *testing.T) {
	yamlData := `
apiVersion: v1
kind: Test
metadata:
  name: example
  labels:
    app: test
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: app
        image: example:latest`
	
	var result map[string]interface{}
	err := ParseYAML([]byte(yamlData), &result)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}
	
	if result["apiVersion"] != "v1" {
		t.Errorf("Expected apiVersion v1, got %v", result["apiVersion"])
	}
	
	if result["kind"] != "Test" {
		t.Errorf("Expected kind Test, got %v", result["kind"])
	}
	
	metadata, ok := result["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("metadata is not a map")
	}
	
	if metadata["name"] != "example" {
		t.Errorf("Expected name example, got %v", metadata["name"])
	}
}

func TestParseJSON(t *testing.T) {
	jsonData := `{
		"apiVersion": "v1",
		"kind": "Test",
		"metadata": {
			"name": "example"
		}
	}`
	
	var result map[string]interface{}
	err := ParseJSON([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if result["apiVersion"] != "v1" {
		t.Errorf("Expected apiVersion v1, got %v", result["apiVersion"])
	}
}

func TestToYAML(t *testing.T) {
	data := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name": "example",
		},
	}
	
	yamlBytes, err := ToYAML(data)
	if err != nil {
		t.Fatalf("Failed to convert to YAML: %v", err)
	}
	
	var result map[string]interface{}
	err = ParseYAML(yamlBytes, &result)
	if err != nil {
		t.Fatalf("Failed to parse generated YAML: %v", err)
	}
	
	if result["apiVersion"] != "v1" {
		t.Errorf("Expected apiVersion v1, got %v", result["apiVersion"])
	}
}

func TestToJSON(t *testing.T) {
	data := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Test",
	}
	
	jsonBytes, err := ToJSON(data)
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}
	
	var result map[string]interface{}
	err = ParseJSON(jsonBytes, &result)
	if err != nil {
		t.Fatalf("Failed to parse generated JSON: %v", err)
	}
	
	if result["apiVersion"] != "v1" {
		t.Errorf("Expected apiVersion v1, got %v", result["apiVersion"])
	}
}

func TestParserInterface(t *testing.T) {
	p := NewParser()
	
	yamlData := `name: test`
	var yamlResult map[string]interface{}
	if err := p.ParseYAML([]byte(yamlData), &yamlResult); err != nil {
		t.Fatalf("Parser.ParseYAML failed: %v", err)
	}
	
	jsonData := `{"name":"test"}`
	var jsonResult map[string]interface{}
	if err := p.ParseJSON([]byte(jsonData), &jsonResult); err != nil {
		t.Fatalf("Parser.ParseJSON failed: %v", err)
	}
	
	if yamlResult["name"] != jsonResult["name"] {
		t.Errorf("YAML and JSON parsing results don't match")
	}
}