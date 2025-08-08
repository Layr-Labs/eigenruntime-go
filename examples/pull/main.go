package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenruntime-go/pkg/client"
)

func main() {
	var (
		reference  = flag.String("ref", "", "Artifact reference (e.g., ghcr.io/myorg/myartifact:latest)")
		registry   = flag.String("registry", "", "Registry URL (used with -digest)")
		digest     = flag.String("digest", "", "Artifact digest (used with -registry)")
		outputFile = flag.String("output", "", "Output file for the spec (optional)")
		plainHTTP  = flag.Bool("plain-http", false, "Use plain HTTP instead of HTTPS")
	)
	flag.Parse()

	if *reference == "" && (*registry == "" || *digest == "") {
		flag.Usage()
		log.Fatal("Either -ref or both -registry and -digest are required")
	}

	c := client.NewClient(client.ClientOptions{
		PlainHTTP: *plainHTTP,
	})

	var specContent []byte
	var err error
	
	ctx := context.Background()
	
	if *reference != "" {
		fmt.Printf("Pulling artifact from %s...\n", *reference)
		specContent, err = c.FetchSpec(ctx, *reference)
	} else {
		fmt.Printf("Pulling artifact from %s@%s...\n", *registry, *digest)
		artifact, err := c.PullByDigest(ctx, *registry, *digest)
		if err == nil && len(artifact.Layers) > 0 {
			specContent = artifact.Layers[0].Content
		}
	}
	
	if err != nil {
		log.Fatalf("Failed to pull artifact: %v", err)
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, specContent, 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Spec written to %s\n", *outputFile)
	} else {
		fmt.Println("Spec content:")
		fmt.Println(string(specContent))
	}
}