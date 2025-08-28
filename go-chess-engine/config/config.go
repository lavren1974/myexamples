package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config holds all application configuration.
type Config struct {
	BoardRepresentation string `json:"board_representation"`
}

// AppConfig is the global configuration instance.
var AppConfig Config

// init runs before main() and loads the configuration from disk.
func init() {
	configFile := "config.json"

	// Attempt to read the file.
	f, err := os.Open(configFile)
	if err != nil {
		// If it doesn't exist, create a default one.
		log.Println("INFO: config.json not found, creating a default one.")
		createDefaultConfig(configFile)
		// Try opening again
		f, err = os.Open(configFile)
		if err != nil {
			log.Fatalf("FATAL: Could not create or open config file: %v", err)
		}
	}
	defer f.Close()

	// Decode the JSON into our AppConfig struct
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("FATAL: Could not parse config.json: %v", err)
	}
	log.Printf("INFO: Loaded config. Board representation set to '%s'", AppConfig.BoardRepresentation)
}

func createDefaultConfig(filename string) {
	defaultConfig := Config{
		BoardRepresentation: "array", // Default to the stable version
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("FATAL: Could not create default config file: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(defaultConfig)
	if err != nil {
		log.Fatalf("FATAL: Could not write to default config file: %v", err)
	}
}
