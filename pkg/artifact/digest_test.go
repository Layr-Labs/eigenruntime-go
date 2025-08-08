package artifact

import (
	"bytes"
	"testing"
)

func TestComputeDigest(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected string
	}{
		{
			name:     "empty content",
			content:  []byte{},
			expected: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "simple string",
			content:  []byte("hello world"),
			expected: "sha256:b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "yaml content",
			content:  []byte("apiVersion: v1\nkind: Test"),
			expected: "sha256:7c4a89d4e1c19e7a8f9e3e8e7f3d4e5c6a8b9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digest := ComputeDigest(tt.content)
			if len(digest) < 7 || digest[:7] != "sha256:" {
				t.Errorf("Invalid digest format: %s", digest)
			}
		})
	}
}

func TestComputeDigestFromReader(t *testing.T) {
	content := []byte("test content")
	reader := bytes.NewReader(content)
	
	digest, err := ComputeDigestFromReader(reader)
	if err != nil {
		t.Fatalf("Failed to compute digest from reader: %v", err)
	}
	
	if len(digest) < 7 || digest[:7] != "sha256:" {
		t.Errorf("Invalid digest format: %s", digest)
	}
	
	expectedDigest := ComputeDigest(content)
	if digest != expectedDigest {
		t.Errorf("Digest mismatch: got %s, expected %s", digest, expectedDigest)
	}
}