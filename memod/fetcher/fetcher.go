package fetcher

/* fetches */

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/encoder"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/processor"

	"github.com/jasonlvhit/gocron"
)

// SetupMempoolFetcher sets up a periodic mempool fetch job
func SetupMempoolFetcher() {
	mempoolFetchInterval := config.GetInt("mempool.fetchInterval")
	s := gocron.NewScheduler()
	s.Every(uint64(mempoolFetchInterval)).Seconds().Do(getMempool)
	logger.Info.Println("Setup mempool fetcher. First fetch in", mempoolFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func getMempool() {

	body, err := fetchMempool()
	if err != nil {
		logger.Error.Printf("Could not fetch mempool: %v", err.Error())
	}

	mempool, err := encoder.DecodeFetchedMempoolBody(body) // decode fetched response body
	if err != nil {
		logger.Error.Printf("Failed to decode response body as JSON: %s", err.Error())
		return
	}

	processor.ProcessMempool(mempool)

}

func fetchMempool() (body []byte, err error) {
	if config.GetString("mempool.fetchInterface") == "REST" {

		body, err := fetchMempoolFromREST()
		if err != nil {
			return nil, fmt.Errorf("Could not fetch mempool from REST: %s", err.Error())
		}

		return body, nil

	} else if config.GetString("mempool.fetchInterface") == "JSON-RPC" {

		body, err := fetchMempoolFromJSONRPC()
		if err != nil {
			return nil, fmt.Errorf("Could not fetch mempool from JSON-RPC: %s", err.Error())
		}

		return body, nil

	} else {
		return nil, errors.New("Unknown interface " + config.GetString("mempool.fetchInterface"))
	}
}

// fetches the current mempool
func fetchMempoolFromREST() ([]byte, error) {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromREST()")

	resp, err := getMempoolContentsREST()
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
func getMempoolContentsREST() (*http.Response, error) {
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

func fetchMempoolFromJSONRPC() ([]byte, error) {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromJSONRPC()")

	rpcURL := config.GetString("bitcoind.jsonrpc.protocol") +
		"://" + config.GetString("bitcoind.jsonrpc.username") +
		":" + config.GetString("bitcoind.jsonrpc.password") +
		"@" + config.GetString("bitcoind.jsonrpc.host") +
		":" + config.GetString("bitcoind.jsonrpc.port")

	timeout := time.Duration(config.GetInt("bitcoind.jsonrpc.responseTimeout")) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	bodyReq := strings.NewReader("{\"jsonrpc\": \"1.0\", \"id\":\"memod-via-rpc\", \"method\": \"getrawmempool\", \"params\": [true] }")
	req, err := http.NewRequest("POST", rpcURL, bodyReq)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JSON-RPC Request failed with status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	// The JSON-RPC response is encapsulated in a JSON result object
	// like {"result":{...},"error":null,"id":"memod-via-rpc"}
	// The REST response isn't. We remove this here.
	body = body[10 : len(body)-36]

	return body, nil
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
