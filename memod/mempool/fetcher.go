package mempool

/* fetches */

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/0xb10c/memo/memod/logger"
	"github.com/jasonlvhit/gocron"
)

// SetupMempoolFetcher sets up a periodic mempool fetch job
func SetupMempoolFetcher() {
	mempoolFetchInterval := int(60)
	s := gocron.NewScheduler()
	s.Every(uint64(mempoolFetchInterval)).Seconds().Do(doWork)
	logger.Trace.Println("Setup mempool fetcher. First fetch in", mempoolFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func doWork() {
	body := fetchMempoolFromREST()
	mempool := decodeFetchedMempoolBody(body)
	processMempool(mempool)
}

// fetches the current mempool and gives it to processMempool()
func fetchMempoolFromREST() []byte {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromREST()")

	// make a HTTP GET Request to the Bitcoin Core REST API
	resp, err := http.Get("http://localhost:8332/rest/mempool/contents.json")
	if err != nil {
		logger.Error.Println(err.Error())
	}

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	defer resp.Body.Close()

	return body
}
