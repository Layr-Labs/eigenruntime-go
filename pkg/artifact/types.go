package artifact

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

type Artifact struct {
	Manifest     []byte
	Config       []byte
	Layers       []Layer
	Digest       string
	MediaType    string
	ArtifactType string
}

type Layer struct {
	Content   []byte
	MediaType string
	Digest    string
	Size      int64
}