package database

import (
	"log"

	"github.com/jasonlvhit/gocron"
)

// SetupMempoolEntriesCache sets up a periodic GetMempoolEntries() fetch job
func SetupMempoolEntriesCache(pool *RedisPool) {
	cacheMempoolEntries(pool)

	fetchInterval := uint64(30)
	s := gocron.NewScheduler()
	s.Every(fetchInterval).Seconds().Do(cacheMempoolEntries, pool)
	log.Printf("Setup GetMempoolEntries() cache job to run every %d seconds.\n", fetchInterval)
	<-s.Start()
	defer s.Clear()
}

func cacheMempoolEntries(pool *RedisPool) {
	log.Printf("Caching mempool entries.\n")
	mes, err := pool.GetMempoolEntries()
	if err != nil {
		log.Printf("Error getting mempool entries %v.\n", err)
	}

	mesJSON, err := json.Marshal(mes)
	if err != nil {
		log.Printf("Error marshalling mempool entries %v.\n", err)
	}

	err = pool.SetMempoolEntriesCache(string(mesJSON))
	if err != nil {
		log.Printf("Could not cache mempool entries %v.\n", err)
	}

}
