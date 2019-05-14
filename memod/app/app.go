package app

import (
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/mempool"
)

// Start starts the memo deamon
func Start() {

	database.Setup()
	defer database.Database.Close()
	mempool.SetupMempoolFetcher()

}
