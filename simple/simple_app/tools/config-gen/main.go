package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"basic-layout/simple/simple_app/configs"
	// Import the generated protobuf code
	"basic-layout/simple/simple_app/internal/conf"
)

func main() {
	// Paths
	outputDir := filepath.Join(".", "configs")

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create a default Bootstrap config
	bootstrap := &conf.Bootstrap{
		// Initialize with default values as needed
		App:         configs.DefaultApp(),
		Servers:     configs.DefaultServers(),
		Discoveries: configs.DefaultDiscoveries(),
		Trace:       configs.DefaultTrace(),
		Middlewares: configs.DefaultMiddlewares(),
		Data:        configs.DefaultData(),
		// Add other fields as needed...
	}

	// Convert to map for easier manipulation
	jsonData, err := json.Marshal(bootstrap)
	if err != nil {
		log.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &configMap); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Define the config files to generate with their corresponding paths
	configs := map[string]string{
		"app":       "app",
		"bootstrap": "", // Root level config
		"databases": "data",
		"logger":    "logger",
		"server":    "servers",
	}

	// Generate each config file
	for name, path := range configs {
		var configData interface{}
		if path == "" {
			configData = configMap // Use the whole config for bootstrap
		} else {
			// Get the nested config
			if nested, ok := getNestedValue(configMap, path).(map[string]interface{}); ok {
				configData = nested
			} else {
				log.Printf("Warning: No config found for path '%s', creating empty config", path)
				continue
			}
		}

		// Convert to YAML
		yamlData, err := yaml.Marshal(configData)
		if err != nil {
			log.Printf("Failed to marshal %s to YAML: %v", name, err)
			continue
		}

		// Write to file
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", name))
		if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
			log.Printf("Failed to write config file %s: %v", outputFile, err)
		} else {
			log.Printf("Generated config file: %s", outputFile)
		}
	}

	log.Println("Configuration files generated successfully!")
}

// getNestedValue gets a value from a nested map using dot notation (e.g., "data.database")
func getNestedValue(m map[string]interface{}, path string) interface{} {
	keys := splitPath(path)
	for i, key := range keys {
		if val, ok := m[key]; ok {
			if i == len(keys)-1 {
				return val
			}
			if nested, ok := val.(map[string]interface{}); ok {
				m = nested
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return m
}

// splitPath splits a dot-separated path into individual keys
func splitPath(path string) []string {
	var result []string
	var current []rune
	inQuotes := false
	for _, r := range path {
		switch r {
		case '.':
			if !inQuotes && len(current) > 0 {
				result = append(result, string(current))
				current = nil
			} else {
				current = append(current, r)
			}
		case '"':
			inQuotes = !inQuotes
		default:
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}
