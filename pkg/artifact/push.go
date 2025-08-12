package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenruntime-go/pkg/common"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
)

type BuildOptions struct {
	Description string
	Source      string
	Version     string
	Annotations map[string]string
	CreatedTime *time.Time
}

func BuildAndPush(ctx context.Context, specContent []byte, opts BuildOptions, reference string) (string, error) {
	// Create minimal config
	config := map[string]interface{}{
		"created": time.Now().Format(time.RFC3339),
	}
	configData, _ := json.Marshal(config)
	
	// Create manifest
	manifest, err := createManifest(specContent, configData, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create manifest: %w", err)
	}
	
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	// Create memory store and push all components
	store := memory.New()
	
	// Store config
	configDesc := ocispec.Descriptor{
		MediaType: common.MediaTypeEigenRuntimeConfig,
		Digest:    digest.Digest(ComputeDigest(configData)),
		Size:      int64(len(configData)),
	}
	if err := store.Push(ctx, configDesc, bytes.NewReader(configData)); err != nil {
		return "", fmt.Errorf("failed to store config: %w", err)
	}
	
	// Store spec layer
	specDesc := ocispec.Descriptor{
		MediaType: common.MediaTypeYAML,
		Digest:    digest.Digest(ComputeDigest(specContent)),
		Size:      int64(len(specContent)),
	}
	if err := store.Push(ctx, specDesc, bytes.NewReader(specContent)); err != nil {
		return "", fmt.Errorf("failed to store spec: %w", err)
	}
	
	// Store manifest
	manifestDesc := ocispec.Descriptor{
		MediaType: common.MediaTypeOCIManifest,
		Digest:    digest.Digest(ComputeDigest(manifestJSON)),
		Size:      int64(len(manifestJSON)),
	}
	if err := store.Push(ctx, manifestDesc, bytes.NewReader(manifestJSON)); err != nil {
		return "", fmt.Errorf("failed to store manifest: %w", err)
	}
	
	// Push to registry
	repo, err := remote.NewRepository(reference)
	if err != nil {
		return "", fmt.Errorf("failed to create repository: %w", err)
	}
	
	_, err = oras.Copy(ctx, store, reference, repo, reference, oras.DefaultCopyOptions)
	if err != nil {
		return "", fmt.Errorf("failed to push to registry: %w", err)
	}
	
	return string(manifestDesc.Digest), nil
}

func createManifest(specContent []byte, config []byte, opts BuildOptions) (interface{}, error) {
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

	specDigest := ComputeDigest(specContent)
	configDigest := ComputeDigest(config)
	
	manifest := map[string]interface{}{
		"schemaVersion": 2,
		"mediaType":     common.MediaTypeOCIManifest,
		"artifactType":  common.MediaTypeEigenRuntimeManifest,
		"config": map[string]interface{}{
			"mediaType": common.MediaTypeEigenRuntimeConfig,
			"digest":    configDigest,
			"size":      len(config),
		},
		"layers": []map[string]interface{}{
			{
				"mediaType": common.MediaTypeYAML,
				"digest":    specDigest,
				"size":      len(specContent),
			},
		},
		"annotations": opts.Annotations,
	}
	
	return manifest, nil
}