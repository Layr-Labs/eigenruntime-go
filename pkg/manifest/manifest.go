package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	MediaTypeEigenRuntimeManifest       = "application/vnd.eigenruntime.manifest.v1"
	MediaTypeEigenRuntimeConfig         = "application/vnd.eigenruntime.manifest.config.v1+json"
	MediaTypeYAML                       = "text/yaml"
	MediaTypeOCIManifest               = "application/vnd.oci.image.manifest.v1+json"
	
	AnnotationSpecVersion     = "io.eigenruntime.spec.version"
	AnnotationImageCreated    = "org.opencontainers.image.created"
	AnnotationImageDescription = "org.opencontainers.image.description"
	AnnotationImageSource     = "org.opencontainers.image.source"
	
	DefaultSpecVersion = "v1"
)

type BuildOptions struct {
	Description string
	Source      string
	Version     string
	Annotations map[string]string
	CreatedTime *time.Time
}

type Manifest struct {
	SchemaVersion int                    `json:"schemaVersion"`
	MediaType     string                 `json:"mediaType"`
	ArtifactType  string                 `json:"artifactType"`
	Config        ocispec.Descriptor     `json:"config"`
	Layers        []ocispec.Descriptor   `json:"layers"`
	Annotations   map[string]string      `json:"annotations,omitempty"`
}

func computeDigest(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:]))
}

func CreateManifest(specContent []byte, config []byte, opts BuildOptions) (*Manifest, error) {
	if opts.Annotations == nil {
		opts.Annotations = make(map[string]string)
	}
	
	createdTime := time.Now()
	if opts.CreatedTime != nil {
		createdTime = *opts.CreatedTime
	}
	
	version := DefaultSpecVersion
	if opts.Version != "" {
		version = opts.Version
	}
	
	opts.Annotations[AnnotationSpecVersion] = version
	opts.Annotations[AnnotationImageCreated] = createdTime.Format(time.RFC3339)
	
	if opts.Description != "" {
		opts.Annotations[AnnotationImageDescription] = opts.Description
	}
	
	if opts.Source != "" {
		opts.Annotations[AnnotationImageSource] = opts.Source
	}

	specDigest := computeDigest(specContent)
	configDigest := computeDigest(config)
	
	manifest := &Manifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeOCIManifest,
		ArtifactType:  MediaTypeEigenRuntimeManifest,
		Config: ocispec.Descriptor{
			MediaType: MediaTypeEigenRuntimeConfig,
			Digest:    digest.Digest(configDigest),
			Size:      int64(len(config)),
		},
		Layers: []ocispec.Descriptor{
			{
				MediaType: MediaTypeYAML,
				Digest:    digest.Digest(specDigest),
				Size:      int64(len(specContent)),
			},
		},
		Annotations: opts.Annotations,
	}
	
	return manifest, nil
}

func (m *Manifest) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func ParseManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}
	return &manifest, nil
}