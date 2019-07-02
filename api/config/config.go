package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// logger/logger.go can't be used here to log,
// since logger.go itself depends on the config
// to read `log.enableTrace`.

func init() {
	setDefaults()

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("ERROR: Error reading config file: ", err)
	}
}

func setDefaults() {
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.connection", "redis")
	viper.SetDefault("api.port", "23485")
	viper.SetDefault("api.production", true)
}

// GetInt returns a config property as int
func GetInt(property string) int {
	result := viper.GetInt(property)
	if result == 0 {
		fmt.Println("WARN: Property " + property + " is 0. Is this not set?")
	}
	return result
}

// GetString returns a config property as string
func GetString(property string) string {
	result := viper.GetString(property)
	if result == "" {
		fmt.Println("WARN: Property " + property + " is \"\". Is this not set?")
	}
	return result
}

// GetBool returns a config property as bool
func GetBool(property string) bool {
	return viper.GetBool(property)
}
