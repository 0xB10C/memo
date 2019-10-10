package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

// logger/logger.go can't be used here to log,
// since logger.go itself depends on the config
// to read `log.enableTrace`.

func init() {
	// do not try to load a config file when running tests
	if flag.Lookup("test.v") != nil {
		return
	}

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

	viper.SetDefault("bitcoind.rest.protocol", "http")
	viper.SetDefault("bitcoind.rest.host", "localhost")
	viper.SetDefault("bitcoind.rest.port", "8332")
	viper.SetDefault("bitcoind.rest.responseTimeout", 30)

	viper.SetDefault("bitcoind.jsonrpc.protocol", "http")
	viper.SetDefault("bitcoind.jsonrpc.host", "localhost")
	viper.SetDefault("bitcoind.jsonrpc.port", "8332")
	viper.SetDefault("bitcoind.jsonrpc.responseTimeout", 30)

	viper.SetDefault("mempool.enable", false)
	viper.SetDefault("mempool.fetchInterface", "REST")
	viper.SetDefault("mempool.callSaveMempool", true)
	viper.SetDefault("mempool.processing.processCurrentMempool", false)
	viper.SetDefault("mempool.processing.processHistoricalMempool", false)
	viper.SetDefault("mempool.processing.processTransactionStats", false)
	viper.SetDefault("mempool.fetchInterval", 60)

	viper.SetDefault("feeratefetcher.enable", false)
	viper.SetDefault("feeratefetcher.fetchInterval", 180)

	viper.SetDefault("log.enableTrace", false)
	viper.SetDefault("log.colorizeOutput", true)

	viper.SetDefault("zmq.enable", false)
	viper.SetDefault("zmq.host", "localhost")
	viper.SetDefault("zmq.port", 28332)
	viper.SetDefault("zmq.subscribeTo.rawTx", false)
	viper.SetDefault("zmq.subscribeTo.rawBlock", false)
	viper.SetDefault("zmq.subscribeTo.hashTx", false)
	viper.SetDefault("zmq.subscribeTo.hashBlock", false)
	viper.SetDefault("zmq.saveMempoolEntries.enable", false) // saves mempool entries to a SQLite database
	viper.SetDefault("zmq.saveMempoolEntries.dbPath", "/tmp/mempool-entries.sqlite")
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
