package manifest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenruntime-go/pkg/artifact"
	"github.com/Layr-Labs/eigenruntime-go/pkg/common"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
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

func CreateMinimalConfig() []byte {
	config := map[string]interface{}{
		"created": time.Now().Format(time.RFC3339),
	}
	
	data, _ := json.Marshal(config)
	return data
}

func CreateManifest(specContent []byte, config []byte, opts BuildOptions) (*Manifest, error) {
	if opts.Annotations == nil {
		opts.Annotations = make(map[string]string)
	}
	
	createdTime := time.Now()
	if opts.CreatedTime != nil {
		createdTime = *opts.CreatedTime
	}
	
	version := common.DefaultSpecVersion
	if opts.Version != "" {
		version = opts.Version
	}
	
	opts.Annotations[common.AnnotationSpecVersion] = version
	opts.Annotations[common.AnnotationImageCreated] = createdTime.Format(time.RFC3339)
	
	if opts.Description != "" {
		opts.Annotations[common.AnnotationImageDescription] = opts.Description
	}
	
	if opts.Source != "" {
		opts.Annotations[common.AnnotationImageSource] = opts.Source
	}

	specDigest := artifact.ComputeDigest(specContent)
	configDigest := artifact.ComputeDigest(config)
	
	manifest := &Manifest{
		SchemaVersion: 2,
		MediaType:     common.MediaTypeOCIManifest,
		ArtifactType:  common.MediaTypeEigenRuntimeManifest,
		Config: ocispec.Descriptor{
			MediaType: common.MediaTypeEigenRuntimeConfig,
			Digest:    digest.Digest(configDigest),
			Size:      int64(len(config)),
		},
		Layers: []ocispec.Descriptor{
			{
				MediaType: common.MediaTypeYAML,
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