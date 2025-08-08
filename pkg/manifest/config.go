package manifest

import (
	"encoding/json"
	"time"
)

type Config struct {
	Created      time.Time              `json:"created"`
	Architecture string                 `json:"architecture"`
	OS           string                 `json:"os"`
	Config       map[string]interface{} `json:"config,omitempty"`
	RootFS       RootFS                 `json:"rootfs"`
}

type RootFS struct {
	Type    string `json:"type"`
	DiffIDs []string `json:"diff_ids,omitempty"`
}

func CreateConfig(metadata map[string]interface{}) ([]byte, error) {
	config := Config{
		Created:      time.Now(),
		Architecture: "unknown",
		OS:           "unknown",
		Config:       metadata,
		RootFS: RootFS{
			Type:    "layers",
			DiffIDs: []string{},
		},
	}
	
	return json.Marshal(config)
}

func CreateMinimalConfig() ([]byte, error) {
	return CreateConfig(nil)
}