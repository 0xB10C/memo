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
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.connection", "tcp")

	viper.SetDefault("bitcoind.rest.protocol", "http")
	viper.SetDefault("bitcoind.rest.host", "localhost")
	viper.SetDefault("bitcoind.rest.port", "8332")
	viper.SetDefault("bitcoind.rest.responseTimeout", 30)

	viper.SetDefault("mempool.fetchInterval", 60)

	viper.SetDefault("log.enableTrace", false)
	viper.SetDefault("log.colorizeOutput", true)
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
