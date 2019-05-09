package main

import (
    "os"
    "fmt"
    "github.com/google/logger"
    "github.com/spf13/viper"
  )

  const logPath = "memod.log"
  const configPath = "memod.toml"


    func loadConfig() {
        viper.SetDefault("ContentDir", "content")

        viper.SetConfigType("toml")
        viper.SetConfigName("config")   
        viper.AddConfigPath(".")        
        err := viper.ReadInConfig()
        if err != nil {
            panic(fmt.Errorf("Fatal error config file: %s ", err))
        }
    }

    func initLogger() (*logger.Logger) {
        lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
        if err!= nil { logger.Fatal("Could not open log file ", logPath) }
        return logger.Init("foo", true, true, lf)
    }

    
  func main() {
    
    loadConfig()
    defer initLogger().Close() // inits a new logger and defers Close()


  }