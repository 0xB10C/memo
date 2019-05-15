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
	mempoolFetchInterval := int(5)
	s := gocron.NewScheduler()
	s.Every(uint64(mempoolFetchInterval)).Seconds().Do(doWork)
	logger.Trace.Println("Setup mempool fetcher. First fetch in", mempoolFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func doWork() {

	// fetch from REST API
	body, err := fetchMempoolFromREST()
	if err != nil {
		logger.Error.Printf("Could not fetch mempool from REST: %s", err.Error())
		return
	}

	// decode fetched response body
	mempool, err := decodeFetchedMempoolBody(body)
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

	timeout := time.Duration(30 * time.Second) // TODO: make the timeout configurable
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get("http://localhost:18332/rest/mempool/contents.json")
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
