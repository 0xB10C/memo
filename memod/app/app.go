package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/fetcher"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/zmq"
)

// Run starts the memo deamon
func Run() {

	exitSignals := make(chan os.Signal, 1)
	shouldExit := make(chan bool, 1)

	signal.Notify(exitSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go handleExitSig(exitSignals, shouldExit)

	db, err := database.Setup()
	if err != nil {
		logger.Error.Printf("Failed to setup database connection: %s", err.Error())
		shouldExit <- true
	} else {
		defer db.Close()
	}

	startWorkers()

	<-shouldExit // wait till memod should exit
	logger.Info.Println("Memod exiting")
}

// handles exit signals
func handleExitSig(exitSignals chan os.Signal, shouldExit chan bool) {
	sig := <-exitSignals
	logger.Info.Println("Received signal", sig)
	shouldExit <- true
}

func startWorkers() {
	// run the mempool fetcher in a goroutine
	go fetcher.SetupMempoolFetcher()

	// starts the ZMQ listener
	go zmq.Start()
}
