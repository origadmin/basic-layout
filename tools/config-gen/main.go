package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"

	"basic-layout/simple/simple_app/configs"
)

func main() {
	// Set output directory
	outputDir := filepath.Join("configs_new") // Output to the configs folder in the project root directory

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 1. Generate bootstrap.yaml
	generateBootstrapConfig(outputDir)

	// 2. Generate app.yaml
	generateAppConfig(outputDir)

	// 3. Generate server.yaml
	generateServerConfig(outputDir)

	// 4. Generate databases.yaml
	generateDatabasesConfig(outputDir)

	// 5. Generate logger.yaml
	generateLoggerConfig(outputDir)

	// 6. Generate middlewares.yaml
	generateMiddlewaresConfig(outputDir)

	// 7. Generate trace.yaml
	generateTraceConfig(outputDir)

	log.Println("Configuration files generated successfully!")
}

func generateTraceConfig(dir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get trace configuration from default config
	trace := configs.DefaultTrace()

	// Serialize trace configuration using protojson
	m := protojson.MarshalOptions{
		EmitUnpopulated: false, // Don't output empty values
	}

	// Serialize to JSON bytes
	jsonBytes, err := m.Marshal(trace)
	if err != nil {
		log.Printf("Failed to serialize trace configuration: %v", err)
		return
	}

	// Deserialize to map
	var traceMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &traceMap); err != nil {
		log.Printf("Failed to deserialize trace configuration: %v", err)
		return
	}

	// Set trace configuration
	v.Set("trace", traceMap)

	// Write to file
	outputFile := filepath.Join(dir, "trace.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate trace.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

func generateMiddlewaresConfig(dir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get middleware configuration from default config
	middlewares := configs.DefaultMiddlewares()

	// Serialize middleware configuration using protojson
	m := protojson.MarshalOptions{
		EmitUnpopulated: false, // Don't output empty values
	}

	// Serialize to JSON bytes
	jsonBytes, err := m.Marshal(middlewares)
	if err != nil {
		log.Printf("Failed to serialize middleware configuration: %v", err)
		return
	}

	// Deserialize to map
	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		log.Printf("Failed to deserialize middleware configuration: %v", err)
		return
	}

	// Set configuration
	v.Set("middlewares", configMap)

	// Write to file
	outputFile := filepath.Join(dir, "middlewares.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate middlewares.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateBootstrapConfig generates bootstrap.yaml
func generateBootstrapConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Set bootstrap configuration
	sources := []map[string]interface{}{
		{
			"type": "file",
			"file": map[string]string{
				"path": "app.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "server.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "databases.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "logger.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "middlewares.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "trace.yaml",
			},
		},
	}

	v.Set("sources", sources)

	// Write to file
	outputFile := filepath.Join(outputDir, "bootstrap.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate bootstrap.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateAppConfig generates app.yaml
func generateAppConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get application configuration from default config
	app := configs.DefaultApp()

	// Set application configuration
	v.Set("app", map[string]interface{}{
		"id":      app.Id,
		"name":    app.Name,
		"version": app.Version,
		"env":     app.Env,
	})

	// Write to file
	outputFile := filepath.Join(outputDir, "app.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate app.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateServerConfig generates server.yaml
func generateServerConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get server configuration from default config
	servers := configs.DefaultServers()

	// Convert server configuration
	var serverConfigs []map[string]interface{}
	for _, srv := range servers.Configs {
		srvCfg := map[string]interface{}{
			"name":     srv.Name,
			"protocol": srv.Protocol,
		}

		switch srv.Protocol {
		case "http":
			if srv.Http != nil {
				srvCfg["http"] = map[string]interface{}{
					"network":     srv.Http.Network,
					"addr":        srv.Http.Addr,
					"middlewares": srv.Http.Middlewares,
				}
			}
		case "grpc":
			if srv.Grpc != nil {
				srvCfg["grpc"] = map[string]interface{}{
					"network":     srv.Grpc.Network,
					"addr":        srv.Grpc.Addr,
					"middlewares": srv.Grpc.Middlewares,
				}
			}
		}

		serverConfigs = append(serverConfigs, srvCfg)
	}

	v.Set("servers", map[string]interface{}{
		"configs": serverConfigs,
	})

	// Write to file
	outputFile := filepath.Join(outputDir, "server.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate server.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateDatabasesConfig generates databases.yaml
func generateDatabasesConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get data configuration from default config
	data := configs.DefaultData()

	// Convert database configuration
	var dbConfigs []map[string]interface{}
	for _, db := range data.Databases.Configs {
		dbCfg := map[string]interface{}{
			"name":    db.Name,
			"dialect": db.Dialect,
		}
		if db.Source != "" {
			dbCfg["source"] = db.Source
		}
		dbConfigs = append(dbConfigs, dbCfg)
	}

	// Convert cache configuration
	var cacheConfigs []map[string]interface{}
	for _, cache := range data.Caches.Configs {
		cacheCfg := map[string]interface{}{
			"name":   cache.Name,
			"driver": cache.Driver,
		}
		if cache.Driver != "" {
			cacheCfg["source"] = cache.Driver
		}
		cacheConfigs = append(cacheConfigs, cacheCfg)
	}

	// Set data configuration
	v.Set("data", map[string]interface{}{
		"databases": map[string]interface{}{
			"configs": dbConfigs,
		},
		"caches": map[string]interface{}{
			"configs": cacheConfigs,
		},
	})

	// Write to file
	outputFile := filepath.Join(outputDir, "databases.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate databases.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateLoggerConfig generates logger.yaml
func generateLoggerConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Set default logger configuration
	v.Set("logger", map[string]interface{}{
		"level":  "info",
		"format": "text",
		"output": "stdout",
		"caller": true,
	})

	// Write to file
	outputFile := filepath.Join(outputDir, "logger.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate logger.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}
