package artifact

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenruntime-go/pkg/manifest"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
)

type Builder interface {
	Build(ctx context.Context, specContent []byte, opts manifest.BuildOptions) (*Artifact, error)
	Push(ctx context.Context, artifact *Artifact, reference string) (string, error)
	BuildAndPush(ctx context.Context, specContent []byte, opts manifest.BuildOptions, reference string) (string, error)
}

type builder struct {
	store *memory.Store
}

func NewBuilder() Builder {
	return &builder{
		store: memory.New(),
	}
}

func (b *builder) Build(ctx context.Context, specContent []byte, opts manifest.BuildOptions) (*Artifact, error) {
	config, err := manifest.CreateMinimalConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}
	
	m, err := manifest.CreateManifest(specContent, config, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest: %w", err)
	}
	
	manifestJSON, err := m.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	artifact := &Artifact{
		Manifest:     manifestJSON,
		Config:       config,
		MediaType:    manifest.MediaTypeOCIManifest,
		ArtifactType: manifest.MediaTypeEigenRuntimeManifest,
		Digest:       ComputeDigest(manifestJSON),
		Layers: []Layer{
			{
				Content:   specContent,
				MediaType: manifest.MediaTypeYAML,
				Digest:    ComputeDigest(specContent),
				Size:      int64(len(specContent)),
			},
		},
	}
	
	if err := b.storeArtifact(ctx, artifact); err != nil {
		return nil, fmt.Errorf("failed to store artifact: %w", err)
	}
	
	return artifact, nil
}

func (b *builder) Push(ctx context.Context, artifact *Artifact, reference string) (string, error) {
	repo, err := remote.NewRepository(reference)
	if err != nil {
		return "", fmt.Errorf("failed to create repository: %w", err)
	}
	
	desc := ocispec.Descriptor{
		MediaType: artifact.MediaType,
		Digest:    digest.Digest(artifact.Digest),
		Size:      int64(len(artifact.Manifest)),
	}
	
	_, err = oras.Copy(ctx, b.store, reference, repo, reference, oras.DefaultCopyOptions)
	if err != nil {
		return "", fmt.Errorf("failed to push artifact: %w", err)
	}
	
	return string(desc.Digest), nil
}

func (b *builder) BuildAndPush(ctx context.Context, specContent []byte, opts manifest.BuildOptions, reference string) (string, error) {
	artifact, err := b.Build(ctx, specContent, opts)
	if err != nil {
		return "", err
	}
	
	return b.Push(ctx, artifact, reference)
}

func (b *builder) storeArtifact(ctx context.Context, artifact *Artifact) error {
	configDesc := ocispec.Descriptor{
		MediaType: manifest.MediaTypeEigenRuntimeConfig,
		Digest:    digest.Digest(ComputeDigest(artifact.Config)),
		Size:      int64(len(artifact.Config)),
	}
	if err := b.store.Push(ctx, configDesc, bytes.NewReader(artifact.Config)); err != nil {
		return fmt.Errorf("failed to store config: %w", err)
	}
	
	for _, layer := range artifact.Layers {
		layerDesc := ocispec.Descriptor{
			MediaType: layer.MediaType,
			Digest:    digest.Digest(layer.Digest),
			Size:      layer.Size,
		}
		if err := b.store.Push(ctx, layerDesc, bytes.NewReader(layer.Content)); err != nil {
			return fmt.Errorf("failed to store layer: %w", err)
		}
	}
	
	manifestDesc := ocispec.Descriptor{
		MediaType: artifact.MediaType,
		Digest:    digest.Digest(artifact.Digest),
		Size:      int64(len(artifact.Manifest)),
	}
	if err := b.store.Push(ctx, manifestDesc, bytes.NewReader(artifact.Manifest)); err != nil {
		return fmt.Errorf("failed to store manifest: %w", err)
	}
	
	return nil
}