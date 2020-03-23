package cache

import (
	"encoding/json"
	"log"

	"github.com/0xb10c/memo/api/database"
	"github.com/jasonlvhit/gocron"
)

// SetupCache sets up caching scripts
func SetupCache() {
	SetupMempoolEntriesCacher()
}

// SetupMempoolEntriesCacher sets up a periodic GetMempoolEntries() fetch job
func SetupMempoolEntriesCacher() {
	cacheMempoolEntries()

	fetchInterval := uint64(30)
	s := gocron.NewScheduler()
	s.Every(fetchInterval).Seconds().Do(cacheMempoolEntries)
	log.Printf("Setup GetMempoolEntries() cacher to run every %d seconds.\n", fetchInterval)
	<-s.Start()
	defer s.Clear()
}

func cacheMempoolEntries() {
	log.Printf("Caching mempool entries.\n")
	mes, err := database.GetMempoolEntries()
	if err != nil {
		log.Printf("Error getting mempool entries %v.\n", err)
	}

	mesJSON, err := json.Marshal(mes)
	if err != nil {
		log.Printf("Error marshalling mempool entries %v.\n", err)
	}

	err = database.SetMempoolEntriesCache(string(mesJSON))
	if err != nil {
		log.Printf("Could not cache mempool entries %v.\n", err)
	}

}
