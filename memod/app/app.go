package app

import (
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/mempool"
)

// Start starts the memo deamon
func Start() {

	db, err := database.Setup()
	if err != nil {
		if err != nil {
			logger.Error.Printf("Failed to setup database connection: %s", err.Error())
			return
		}
	}
	defer db.Close()

	mempool.SetupMempoolFetcher()

}
