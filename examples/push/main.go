package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenruntime-go/pkg/artifact"
)

func main() {
	var (
		specFile    = flag.String("spec", "", "Path to the spec YAML file")
		registry    = flag.String("registry", "", "Registry URL (e.g., ghcr.io/myorg/myartifact)")
		tag         = flag.String("tag", "latest", "Tag for the artifact")
		description = flag.String("description", "", "Description for the artifact")
		source      = flag.String("source", "", "Source URL for the artifact")
	)
	flag.Parse()

	if *specFile == "" || *registry == "" {
		flag.Usage()
		log.Fatal("spec and registry are required")
	}

	specContent, err := os.ReadFile(*specFile)
	if err != nil {
		log.Fatalf("Failed to read spec file: %v", err)
	}

	reference := fmt.Sprintf("%s:%s", *registry, *tag)
	
	digest, err := artifact.BuildAndPush(
		context.Background(),
		specContent,
		artifact.BuildOptions{
			Description: *description,
			Source:      *source,
		},
		reference,
	)
	if err != nil {
		log.Fatalf("Failed to build and push artifact: %v", err)
	}

	fmt.Printf("Successfully pushed artifact to %s\n", reference)
	fmt.Printf("Digest: %s\n", digest)
}