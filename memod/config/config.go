package config

import (
	"fmt"

	"github.com/0xb10c/memo/memod/logger"
	"github.com/spf13/viper"
)

func init() {
	setDefaults()

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}

func setDefaults() {
	viper.SetDefault("logFile", "memod.log")
	viper.SetDefault("ContentDir", "content")
}

// GetInt returns a config property as int
func GetInt(property string) int {
	result := viper.GetInt(property)
	if result == 0 {
		logger.Warning.Println("Property " + property + " is not set. Using " + property + " = 0.")
	}
	return result
}

// GetString returns a config property as string
func GetString(property string) string {
	result := viper.GetString(property)
	if result == "" {
		logger.Warning.Println("Property " + property + " is not set. Using " + property + " = \"\"")
	}
	return result
}
