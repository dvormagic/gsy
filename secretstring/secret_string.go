// Package secretstring provides a SecretString type that implements YAML unmarshaling
// for secure configuration management. It supports plain strings for local development
// and references to Google Cloud Secret Manager for production environments.
//
// Usage:
//
// Set the environment mode with SetEnv("local") or SetEnv("prod").
// Embed SecretString in your config structs, and unmarshal YAML as usual.
//
// For local: Use plain strings in YAML.
// For prod: Use maps like {secret: "projects/my-project/secrets/my-secret/versions/latest"}.
package secretstring

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"gopkg.in/yaml.v3"
)

// SecretString is a string type that implements yaml.Unmarshaler.
// It resolves secrets from YAML based on the current environment.
type SecretString string

var secretStringEnv = "local"

// SetEnv sets the environment mode for secret resolution.
// Valid values: "local" (use plain strings), "prod" (fetch from GCP Secret Manager).
func SetEnv(env string) {
	secretStringEnv = env
}

// UnmarshalYAML implements yaml.Unmarshaler for SecretString.
func (s *SecretString) UnmarshalYAML(value *yaml.Node) error {
	var str string
	// Try to unmarshal as string
	if err := value.Decode(&str); err == nil {
		*s = SecretString(str)
		return nil
	}
	// Try to unmarshal as map
	var m map[string]string
	if err := value.Decode(&m); err == nil {
		if secretRef, ok := m["secret"]; ok {
			if secretStringEnv == "local" {
				*s = SecretString(secretRef)
				return nil
			}

			val, err := FetchSecretFromGCP(secretRef)
			if err != nil {
				return err
			}
			*s = SecretString(val)
			return nil
		}
	}
	return fmt.Errorf("invalid secret format")
}

// FetchSecretFromGCP retrieves a secret from Google Cloud Secret Manager.
// The secretName should be the full resource name, e.g.,
// "projects/my-project/secrets/my-secret/versions/latest".
func FetchSecretFromGCP(secretName string) (string, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}