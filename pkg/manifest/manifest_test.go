package manifest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCreateManifest(t *testing.T) {
	specContent := []byte("apiVersion: v1\nkind: Test\nname: test")
	config := []byte(`{"created":"2023-01-01T00:00:00Z"}`)
	
	opts := BuildOptions{
		Description: "Test artifact",
		Source:      "https://github.com/test/repo",
		Version:     "v1",
		Annotations: map[string]string{
			"custom.annotation": "value",
		},
	}
	
	manifest, err := CreateManifest(specContent, config, opts)
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}
	
	if manifest.SchemaVersion != 2 {
		t.Errorf("Expected schema version 2, got %d", manifest.SchemaVersion)
	}
	
	if manifest.MediaType != MediaTypeOCIManifest {
		t.Errorf("Expected media type %s, got %s", MediaTypeOCIManifest, manifest.MediaType)
	}
	
	if manifest.ArtifactType != MediaTypeEigenRuntimeManifest {
		t.Errorf("Expected artifact type %s, got %s", MediaTypeEigenRuntimeManifest, manifest.ArtifactType)
	}
	
	if manifest.Config.MediaType != MediaTypeEigenRuntimeConfig {
		t.Errorf("Expected config media type %s, got %s", MediaTypeEigenRuntimeConfig, manifest.Config.MediaType)
	}
	
	if len(manifest.Layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(manifest.Layers))
	}
	
	if manifest.Layers[0].MediaType != MediaTypeYAML {
		t.Errorf("Expected layer media type %s, got %s", MediaTypeYAML, manifest.Layers[0].MediaType)
	}
	
	if manifest.Annotations[AnnotationSpecVersion] != "v1" {
		t.Errorf("Expected spec version v1, got %s", manifest.Annotations[AnnotationSpecVersion])
	}
	
	if manifest.Annotations["custom.annotation"] != "value" {
		t.Errorf("Expected custom annotation value, got %s", manifest.Annotations["custom.annotation"])
	}
}

func TestManifestToJSON(t *testing.T) {
	specContent := []byte("test")
	config := []byte("config")
	
	manifest, err := CreateManifest(specContent, config, BuildOptions{})
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}
	
	jsonData, err := manifest.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert manifest to JSON: %v", err)
	}
	
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if parsed["schemaVersion"].(float64) != 2 {
		t.Errorf("Expected schema version 2 in JSON")
	}
}

func TestParseManifest(t *testing.T) {
	manifestJSON := `{
		"schemaVersion": 2,
		"mediaType": "application/vnd.oci.image.manifest.v1+json",
		"artifactType": "application/vnd.eigenruntime.manifest.v1",
		"config": {
			"mediaType": "application/vnd.eigenruntime.manifest.config.v1+json",
			"digest": "sha256:abc123",
			"size": 100
		},
		"layers": [
			{
				"mediaType": "text/yaml",
				"digest": "sha256:def456",
				"size": 200
			}
		],
		"annotations": {
			"test": "value"
		}
	}`
	
	manifest, err := ParseManifest([]byte(manifestJSON))
	if err != nil {
		t.Fatalf("Failed to parse manifest: %v", err)
	}
	
	if manifest.SchemaVersion != 2 {
		t.Errorf("Expected schema version 2, got %d", manifest.SchemaVersion)
	}
	
	if manifest.Config.Digest != "sha256:abc123" {
		t.Errorf("Expected config digest sha256:abc123, got %s", manifest.Config.Digest)
	}
	
	if len(manifest.Layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(manifest.Layers))
	}
	
	if manifest.Annotations["test"] != "value" {
		t.Errorf("Expected annotation test=value, got %s", manifest.Annotations["test"])
	}
}

func TestCreateManifestWithCustomTime(t *testing.T) {
	specContent := []byte("test")
	config := []byte("config")
	
	customTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	opts := BuildOptions{
		CreatedTime: &customTime,
	}
	
	manifest, err := CreateManifest(specContent, config, opts)
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}
	
	expectedTime := customTime.Format(time.RFC3339)
	if manifest.Annotations[AnnotationImageCreated] != expectedTime {
		t.Errorf("Expected created time %s, got %s", expectedTime, manifest.Annotations[AnnotationImageCreated])
	}
}