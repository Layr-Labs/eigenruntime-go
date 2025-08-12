package runtime

import "github.com/Layr-Labs/eigenruntime-go/pkg/common"

// EigenRuntime defines the interface that all AVS runtime implementations must satisfy.
// Runtimes are loaded dynamically and orchestrate the deployment and lifecycle of AVS components.
type EigenRuntime interface {
	// APIVersion returns the API version this runtime implements (e.g., "eigenruntime/v1alpha1")
	APIVersion() string

	// Kind returns the type of AVS architecture this runtime manages (e.g., "Hourglass")
	Kind() string

	// Run runs the AVS components based on the provided spec.
	// This should be idempotent and handle partial failures gracefully.
	Run(spec common.RuntimeSpec) error

	// Remove removes the AVS components based on the provided spec.
	// This should gracefully handle component removal and cleanup.
	Remove(spec common.RuntimeSpec) error
}

// Context provides runtime execution context including AVS metadata and Docker configuration.
type Context struct {
	// AVS identification
	AVSAddress    string
	OperatorSetID string
	ReleaseID     string

	// Docker configuration
	DockerHost string // Docker socket path or host

	// Working directory for runtime operations
	WorkDir string

	// Environment variables to pass to containers
	Env map[string]string
}

// Runtime loading contract constants
const (
	// RuntimeAPIVersion defines the current runtime API version
	RuntimeAPIVersion = "eigenruntime/v1alpha1"
)
