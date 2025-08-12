package common

const (
	MediaTypeEigenRuntimeManifest = "application/vnd.eigenruntime.manifest.v1"
	MediaTypeEigenRuntimeConfig   = "application/vnd.eigenruntime.manifest.config.v1+json"
	MediaTypeYAML                  = "text/yaml"
	MediaTypeOCIManifest          = "application/vnd.oci.image.manifest.v1+json"

	AnnotationSpecVersion      = "io.eigenruntime.spec.version"
	AnnotationImageCreated     = "org.opencontainers.image.created"
	AnnotationImageDescription = "org.opencontainers.image.description"
	AnnotationImageSource      = "org.opencontainers.image.source"

	DefaultSpecVersion = "v1"
)