package dockerclient

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types/registry"
)

// EncodeRegistryAuth converts a registry.AuthConfig into the base64-encoded JSON expected by the Docker Engine API.
func EncodeRegistryAuth(cfg registry.AuthConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
