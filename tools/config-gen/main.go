package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/config/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	target string
)

func init() {
	flag.StringVar(&target, "target", "", "The target project to generate configs for (e.g., 'simple' or 'multiple')")
}

func main() {
	flag.Parse()

	var outputDir string
	switch target {
	case "simple":
		outputDir = filepath.Join("..", "simple", "simple_app", "resources", "configs")
	case "multiple":
		outputDir = filepath.Join("..", "multiple", "multiple_sample", "resources", "configs")
	default:
		fmt.Println("Error: Please specify a valid target with -target flag. (e.g., -target=simple or -target=multiple)")
		os.Exit(1)
	}

	log.Printf("Generating configuration for target: '%s' into directory: '%s'", target, outputDir)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 1. Generate bootstrap.yaml
	generateBootstrapConfig(outputDir)

	// 2. Generate server.yaml
	generateServerConfig(outputDir)

	// 3. Generate clients.yaml
	generateClientsConfig(outputDir)

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
	trace := DefaultTrace()

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
	middlewares := DefaultMiddlewares()

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
				"path": "server.yaml",
			},
		},
		{
			"type": "file",
			"file": map[string]string{
				"path": "clients.yaml", // Added clients.yaml here
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
		{
			"type": "env", // Added env source
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

// generateServerConfig generates server.yaml
func generateServerConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get server configuration from default config
	servers := DefaultServers()

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

// generateClientsConfig generates clients.yaml
func generateClientsConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get client configuration from default config
	clients := DefaultClients()

	// Convert client configuration
	var clientConfigs []map[string]interface{}
	for _, cli := range clients.Configs {
		cliCfg := map[string]interface{}{
			"name":     cli.Name,
			"protocol": cli.Protocol,
		}

		switch cli.Protocol {
		case "http":
			if cli.Http != nil {
				cliCfg["http"] = map[string]interface{}{
					"endpoint":    cli.Http.Endpoint,
					"middlewares": cli.Http.Middlewares,
				}
			}
		case "grpc":
			if cli.Grpc != nil {
				cliCfg["grpc"] = map[string]interface{}{
					"endpoint":    cli.Grpc.Endpoint,
					"middlewares": cli.Grpc.Middlewares,
				}
			}
		}
		clientConfigs = append(clientConfigs, cliCfg)
	}

	v.Set("clients", map[string]interface{}{
		"configs": clientConfigs,
	})

	// Write to file
	outputFile := filepath.Join(outputDir, "clients.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("Failed to generate clients.yaml: %v", err)
	} else {
		log.Printf("Generate the profile: %s", outputFile)
	}
}

// generateDatabasesConfig generates databases.yaml
func generateDatabasesConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Get data configuration from default config
	data := DefaultData()

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

func DefaultServers() *transportv1.Servers {
	return &transportv1.Servers{
		Configs: []*transportv1.Server{
			{
				Name:     "grpc_server",
				Protocol: "grpc",
				Grpc: &grpcv1.Server{
					Network:     "tcp",
					Addr:        "0.0.0.0:${GRPC_PORT:9090}",
					Middlewares: []string{"recovery", "logger"},
				},
			},
			{
				Name:     "http_server",
				Protocol: "http",
				Http: &httpv1.Server{
					Network:     "tcp",
					Addr:        "0.0.0.0:${HTTP_PORT:8080}",
					Middlewares: []string{"recovery", "logger"},
				},
			},
		},
	}
}

func DefaultClients() *transportv1.Clients {
	return &transportv1.Clients{
		Configs: []*transportv1.Client{
			{
				Name:     "client.user",
				Protocol: "grpc",
				Grpc: &grpcv1.Client{
					Endpoint:    "localhost:${GRPC_USER_PORT:9091}",
					Middlewares: []string{"recovery", "logger"},
				},
			},
			{
				Name:     "client.order",
				Protocol: "grpc",
				Grpc: &grpcv1.Client{
					Endpoint:    "localhost:${GRPC_ORDER_PORT:9092}",
					Middlewares: []string{"recovery", "logger"},
				},
			},
		},
	}
}

func DefaultData() *datav1.Data {
	return &datav1.Data{
		Databases: &datav1.Databases{
			Configs: []*databasev1.DatabaseConfig{
				{
					Name:    "default",
					Dialect: "sqlite3",
					Source:  "file:./test.db?cache=shared&mode=memory&_fk=1",
				},
			},
		},
		Caches: &datav1.Caches{
			Configs: []*cachev1.CacheConfig{
				{
					Name:   "default",
					Driver: "memory",
				},
			},
		},
	}
}

func DefaultMiddlewares() *middlewarev1.Middlewares {
	return &middlewarev1.Middlewares{
		Configs: []*middlewarev1.Middleware{
			{
				Name:     "recovery",
				Type:     "recovery",
				Recovery: &middlewarev1.Recovery{},
			},
			{
				Name:    "logging",
				Type:    "logging",
				Logging: &middlewarev1.Logging{},
			},
		},
	}
}

func DefaultTrace() *tracev1.Trace {
	return &tracev1.Trace{
		Name:        "jaeger",
		Endpoint:    "localhost:6831",
		ServiceName: "test-app-name",
	}
}
