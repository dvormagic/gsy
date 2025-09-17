package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	"github.com/dvormagic/gsy/secretstring"
)

type Config struct {
	APIKey   secretstring.SecretString `yaml:"api_key"`
	Database secretstring.SecretString `yaml:"database_url"`
	OtherKey string                    `yaml:"other_key"`
}

func main() {
	// Set to local mode for this demo (plain strings and refs as-is)
	secretstring.SetEnv("local")

	// YAML config for local dev
	yamlData := []byte(`
api_key: "sk-123-local-api-key"  # Plain string
database_url: 
  secret: "projects/my-project/secrets/db-url/versions/latest"  # GCP ref (used as string in local)
other_key: "regular-value"
	`)

	var cfg Config
	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resolved Config (Local Mode):")
	fmt.Printf("API Key: %s\n", string(cfg.APIKey))  // sk-123-local-api-key
	fmt.Printf("DB URL: %s\n", string(cfg.Database)) // projects/my-project/secrets/db-url/versions/latest
	fmt.Printf("Other: %s\n", cfg.OtherKey)          // regular-value

	// In prod mode (uncomment and set GCP auth):
	// secretstring.SetEnv("prod")
	// // Use same YAMLData
	// // Unmarshal again: DB URL would fetch real secret
	// fmt.Printf("DB URL (Prod): %s\n", string(cfg.Database)) // Actual secret value
}
