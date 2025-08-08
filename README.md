# eigenruntime-go

A Go package for handling EigenRuntime OCI artifacts. This library provides tools to create, push, pull, and manipulate EigenRuntime specifications as OCI artifacts.

## Features

- **OCI Artifact Creation**: Build EigenRuntime artifacts from YAML specifications
- **Registry Operations**: Push and pull artifacts to/from OCI registries
- **Flexible Spec Handling**: Parse and validate YAML/JSON specifications
- **Authentication Support**: Multiple authentication methods for registry access
- **Standard OCI Compliance**: Follows OCI artifact specifications

## Installation

```bash
go get github.com/Layr-Labs/eigenruntime-go
```

## Quick Start

### Building and Pushing an Artifact

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/Layr-Labs/eigenruntime-go/pkg/artifact"
)

func main() {
    // Read your spec file
    specContent, err := os.ReadFile("runtime-spec.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a builder
    builder := artifact.NewBuilder()
    
    // Build and push the artifact
    digest, err := builder.BuildAndPush(
        context.Background(),
        specContent,
        artifact.BuildOptions{
            Description: "My EigenRuntime specification",
            Source:      "https://github.com/myorg/myrepo",
        },
        "ghcr.io/myorg/myartifact:v1.0.0",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Pushed artifact with digest: %s", digest)
}
```

### Pulling an Artifact

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Layr-Labs/eigenruntime-go/pkg/client"
)

func main() {
    // Create a client
    c := client.NewClient(client.ClientOptions{})
    
    // Pull the artifact and get the spec
    specContent, err := c.FetchSpec(
        context.Background(),
        "ghcr.io/myorg/myartifact:v1.0.0",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(specContent))
}
```

## Package Structure

- `pkg/artifact/` - Core OCI artifact handling
  - `types.go` - Media types and constants
  - `builder.go` - Artifact builder implementation
  - `digest.go` - Digest computation utilities

- `pkg/client/` - OCI registry client
  - `client.go` - Client for pulling artifacts
  - `auth.go` - Authentication handling

- `pkg/manifest/` - OCI manifest management
  - `manifest.go` - Manifest creation and parsing
  - `config.go` - Config blob handling

- `pkg/spec/` - Specification handling
  - `spec.go` - YAML/JSON parsing
  - `validator.go` - Spec validation

## Authentication

The client uses the default credential chain for authentication. Configure your Docker credentials using `docker login` or your cloud provider's credential helper.

## Examples

See the `examples/` directory for complete examples:

- `examples/push/` - Example of building and pushing an artifact
- `examples/pull/` - Example of pulling an artifact

### Running Examples

Push an artifact:
```bash
go run examples/push/main.go \
    -spec myspec.yaml \
    -registry ghcr.io/myorg/myartifact \
    -tag v1.0.0 \
    -description "My artifact" \
    -source "https://github.com/myorg/myrepo"
```

Pull an artifact:
```bash
go run examples/pull/main.go \
    -ref ghcr.io/myorg/myartifact:v1.0.0 \
    -output spec.yaml
```

## OCI Artifact Structure

EigenRuntime artifacts follow this structure:

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "artifactType": "application/vnd.eigenruntime.manifest.v1",
  "config": {
    "mediaType": "application/vnd.eigenruntime.manifest.config.v1+json",
    "digest": "sha256:...",
    "size": 319
  },
  "layers": [
    {
      "mediaType": "text/yaml",
      "digest": "sha256:...",
      "size": 437
    }
  ],
  "annotations": {
    "io.eigenruntime.spec.version": "v1",
    "org.opencontainers.image.created": "2025-08-06T05:12:29Z",
    "org.opencontainers.image.description": "EigenRuntime specification",
    "org.opencontainers.image.source": "https://github.com/..."
  }
}
```

## Testing

Run tests with:
```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Add license information here]