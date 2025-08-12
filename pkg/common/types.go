package common

type RuntimeSpec struct {
	APIVersion string               `yaml:"apiVersion" json:"apiVersion"`
	Kind       string               `yaml:"kind" json:"kind"`
	Name       string               `yaml:"name" json:"name"`
	Version    string               `yaml:"version" json:"version"`
	Spec       map[string]Component `yaml:"spec" json:"spec"`
}

type Component struct {
	Registry  string     `yaml:"registry" json:"registry"`
	Digest    string     `yaml:"digest" json:"digest"`
	Command   []string   `yaml:"command,omitempty" json:"command,omitempty"`
	Env       []EnvVar   `yaml:"env,omitempty" json:"env,omitempty"`
	Resources *Resources `yaml:"resources,omitempty" json:"resources,omitempty"`
}

type EnvVar struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type,omitempty" json:"type,omitempty"`
	Required bool   `yaml:"required,omitempty" json:"required,omitempty"`
}

type Resources struct {
	TEEEnabled bool `yaml:"teeEnabled" json:"teeEnabled"`
}

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