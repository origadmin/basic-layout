package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"basic-layout/simple/simple_app/configs"
)

func main() {
	// 设置输出目录
	outputDir := filepath.Join("..", "..", "configs_new") // 输出到项目根目录的configs文件夹

	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 1. 生成 bootstrap.yaml
	generateBootstrapConfig(outputDir)

	// 2. 生成 app.yaml
	generateAppConfig(outputDir)

	// 3. 生成 server.yaml
	generateServerConfig(outputDir)

	// 4. 生成 databases.yaml
	generateDatabasesConfig(outputDir)

	// 5. 生成 logger.yaml
	generateLoggerConfig(outputDir)

	log.Println("配置文件生成成功!")
}

// generateBootstrapConfig 生成 bootstrap.yaml
func generateBootstrapConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 设置 bootstrap 配置
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
	}

	v.Set("sources", sources)

	// 写入文件
	outputFile := filepath.Join(outputDir, "bootstrap.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("生成 bootstrap.yaml 失败: %v", err)
	} else {
		log.Printf("生成配置文件: %s", outputFile)
	}
}

// generateAppConfig 生成 app.yaml
func generateAppConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 从默认配置中获取应用配置
	app := configs.DefaultApp()

	// 设置应用配置
	v.Set("app", map[string]interface{}{
		"id":      app.Id,
		"name":    app.Name,
		"version": app.Version,
		"env":     app.Env,
	})

	// 写入文件
	outputFile := filepath.Join(outputDir, "app.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("生成 app.yaml 失败: %v", err)
	} else {
		log.Printf("生成配置文件: %s", outputFile)
	}
}

// generateServerConfig 生成 server.yaml
func generateServerConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 从默认配置中获取服务器配置
	servers := configs.DefaultServers()

	// 转换服务器配置
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
					"network": srv.Http.Network,
					"addr":    srv.Http.Addr,
				}
			}
		case "grpc":
			if srv.Grpc != nil {
				srvCfg["grpc"] = map[string]interface{}{
					"network": srv.Grpc.Network,
					"addr":    srv.Grpc.Addr,
				}
			}
		}

		serverConfigs = append(serverConfigs, srvCfg)
	}

	v.Set("servers", map[string]interface{}{
		"configs": serverConfigs,
	})

	// 写入文件
	outputFile := filepath.Join(outputDir, "server.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("生成 server.yaml 失败: %v", err)
	} else {
		log.Printf("生成配置文件: %s", outputFile)
	}
}

// generateDatabasesConfig 生成 databases.yaml
func generateDatabasesConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 从默认配置中获取数据配置
	data := configs.DefaultData()

	// 转换数据库配置
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

	// 转换缓存配置
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

	// 设置数据配置
	v.Set("data", map[string]interface{}{
		"databases": map[string]interface{}{
			"configs": dbConfigs,
		},
		"caches": map[string]interface{}{
			"configs": cacheConfigs,
		},
	})

	// 写入文件
	outputFile := filepath.Join(outputDir, "databases.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("生成 databases.yaml 失败: %v", err)
	} else {
		log.Printf("生成配置文件: %s", outputFile)
	}
}

// generateLoggerConfig 生成 logger.yaml
func generateLoggerConfig(outputDir string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 设置默认日志配置
	v.Set("logger", map[string]interface{}{
		"level":  "info",
		"format": "text",
		"output": "stdout",
		"caller": true,
	})

	// 写入文件
	outputFile := filepath.Join(outputDir, "logger.yaml")
	if err := v.WriteConfigAs(outputFile); err != nil {
		log.Printf("生成 logger.yaml 失败: %v", err)
	} else {
		log.Printf("生成配置文件: %s", outputFile)
	}
}
