package config

import (
	"fmt"
	"github.com/khivuksergey/webserver"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server  webserver.ServerConfig
	Router  webserver.RouterConfig
	Swagger *webserver.SwaggerConfig
	Logger  *LoggerConfig
	DB      DBConfig
}

type DBConfig struct {
	ConnectionString string
	TablePrefix      string
}

type LoggerConfig struct {
	LogLevel string
}

func LoadConfiguration(path string) *Configuration {
	config := &Configuration{}

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config from file %s: %v\n", path, err)
		return defaultConfiguration()
	}
	if err := viper.Unmarshal(config); err != nil {
		fmt.Printf("error unmarshalling configuration: %v\n", err)
		return defaultConfiguration()
	}

	return config
}

func defaultConfiguration() *Configuration {
	fmt.Println("loading default configuration...")
	return &Configuration{
		Server: webserver.DefaultServerConfig,
		Router: webserver.DefaultRouterConfig,
	}
}
