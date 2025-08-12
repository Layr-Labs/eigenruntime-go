package client

import (
	"context"
	"fmt"
	"io"

	"github.com/Layr-Labs/eigenruntime-go/pkg/common"
	"github.com/Layr-Labs/eigenruntime-go/pkg/manifest"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type ClientOptions struct {
	PlainHTTP bool
}

type Client struct {
	opts ClientOptions
}

func NewClient(opts ClientOptions) *Client {
	return &Client{
		opts: opts,
	}
}

func (c *Client) Pull(ctx context.Context, reference string) (*common.Artifact, error) {
	repo, err := c.createRepository(reference)
	if err != nil {
		return nil, err
	}

	store := memory.New()
	manifestDesc, err := oras.Copy(ctx, repo, reference, store, reference, oras.DefaultCopyOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to pull artifact: %w", err)
	}

	return c.fetchArtifact(ctx, store, manifestDesc)
}

func (c *Client) PullByDigest(ctx context.Context, registry, digestStr string) (*common.Artifact, error) {
	d, err := digest.Parse(digestStr)
	if err != nil {
		return nil, fmt.Errorf("invalid digest: %w", err)
	}

	reference := fmt.Sprintf("%s@%s", registry, d.String())
	return c.Pull(ctx, reference)
}

func (c *Client) FetchSpec(ctx context.Context, reference string) ([]byte, error) {
	art, err := c.Pull(ctx, reference)
	if err != nil {
		return nil, err
	}

	if len(art.Layers) == 0 {
		return nil, fmt.Errorf("no spec layer found in artifact")
	}

	return art.Layers[0].Content, nil
}

func (c *Client) createRepository(reference string) (oras.Target, error) {
	repo, err := remote.NewRepository(reference)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	repo.Client = &auth.Client{
		Credential: auth.StaticCredential("", auth.EmptyCredential),
	}
	repo.PlainHTTP = c.opts.PlainHTTP

	return repo, nil
}

func (c *Client) fetchArtifact(ctx context.Context, store *memory.Store, desc ocispec.Descriptor) (*common.Artifact, error) {
	manifestRC, err := store.Fetch(ctx, desc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest from store: %w", err)
	}
	defer manifestRC.Close()

	manifestBytes, err := io.ReadAll(manifestRC)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	m, err := manifest.ParseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}

	configRC, err := store.Fetch(ctx, m.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer configRC.Close()

	configBytes, err := io.ReadAll(configRC)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var layers []common.Layer
	for _, layerDesc := range m.Layers {
		layerRC, err := store.Fetch(ctx, layerDesc)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch layer: %w", err)
		}

		layerBytes, err := io.ReadAll(layerRC)
		layerRC.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read layer: %w", err)
		}

		layers = append(layers, common.Layer{
			Content:   layerBytes,
			MediaType: layerDesc.MediaType,
			Digest:    string(layerDesc.Digest),
			Size:      layerDesc.Size,
		})
	}

	return &common.Artifact{
		Manifest:     manifestBytes,
		Config:       configBytes,
		Layers:       layers,
		Digest:       string(desc.Digest),
		MediaType:    desc.MediaType,
		ArtifactType: m.ArtifactType,
	}, nil
}