package mempool

/* fetches */

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/jasonlvhit/gocron"
)

// SetupMempoolFetcher sets up a periodic mempool fetch job
func SetupMempoolFetcher() {
	mempoolFetchInterval := config.GetInt("mempool.fetchInterval")
	s := gocron.NewScheduler()
	s.Every(uint64(mempoolFetchInterval)).Seconds().Do(doWork)
	logger.Info.Println("Setup mempool fetcher. First fetch in", mempoolFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func doWork() {

	body, err := fetchMempoolFromREST() // fetch from REST API
	if err != nil {
		logger.Error.Printf("Could not fetch mempool from REST: %s", err.Error())
		return
	}

	mempool, err := decodeFetchedMempoolBody(body) // decode fetched response body
	if err != nil {
		logger.Error.Printf("Failed to decode response body as JSON: %s", err.Error())
		logger.Error.Println("Response body: ", string(body))
		return
	}

	processMempool(mempool)
}

// fetches the current mempool
func fetchMempoolFromREST() ([]byte, error) {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromREST()")

	resp, err := getMempoolContents()
	if err != nil {
		return nil, err
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// make a HTTP GET Request to the Bitcoin Core REST API
func getMempoolContents() (*http.Response, error) {
	timeout := time.Duration(config.GetInt("bitcoind.rest.responseTimeout")) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	urlPrefix := config.GetString("bitcoind.rest.protocol") +
		"://" + config.GetString("bitcoind.rest.host") +
		":" + config.GetString("bitcoind.rest.port")
	const urlSuffix = "/rest/mempool/contents.json"

	logger.Trace.Println("Fetching mempool contents from ", urlPrefix+urlSuffix)

	resp, err := client.Get(urlPrefix + urlSuffix)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// read the response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	defer resp.Body.Close()

	return body, nil
}
