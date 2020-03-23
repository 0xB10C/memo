package fetcher

/* fetches */

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/encoder"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/processor"

	"github.com/jasonlvhit/gocron"
)

// SetupMempoolFetcher sets up a periodic mempool fetch job
func SetupMempoolFetcher(redisPool *database.RedisPool) {
	mempoolFetchInterval := config.GetInt("mempool.fetchInterval")
	s := gocron.NewScheduler()
	s.Every(uint64(mempoolFetchInterval)).Seconds().Do(getMempool, redisPool)
	logger.Info.Println("Setup mempool fetcher. First fetch in", mempoolFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func getMempool(redisPool *database.RedisPool) {

	body, err := fetchMempool()
	if err != nil {
		logger.Error.Printf("Could not fetch mempool: %v", err.Error())
	}

	mempool, err := encoder.DecodeFetchedMempoolBody(body) // decode fetched response body
	if err != nil {
		logger.Error.Printf("Failed to decode response body as JSON: %s", err.Error())
		return // we return here to stop the execution
	}

	processor.ProcessMempool(mempool, redisPool)

	if config.GetBool("mempool.callSaveMempool") {
		if rand.Intn(100) <= 25 { // Only call savemempool every 4th call
			err = saveMempoolJSONRPC()
			if err != nil {
				logger.Error.Printf("Failed to save mempool: %s", err.Error())
			}
		}
	}
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

// fetches the current mempool from the Bitcoin Core REST API
func fetchMempoolFromREST() (body []byte, err error) {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromREST()")

	urlPrefix := config.GetString("bitcoind.rest.protocol") +
		"://" + config.GetString("bitcoind.rest.host") +
		":" + config.GetString("bitcoind.rest.port")
	const urlSuffix = "/rest/mempool/contents.json"

	logger.Trace.Println("Fetching mempool contents from ", urlPrefix+urlSuffix)

	body, err = makeHTTPGETReq(urlPrefix+urlSuffix, config.GetInt("bitcoind.rest.responseTimeout"))
	if err != nil {
		return
	}

	return
}

func fetchMempoolFromJSONRPC() ([]byte, error) {
	defer logger.TrackTime(time.Now(), "fetchMempoolFromJSONRPC()")
	logger.Trace.Println("Fetching mempool via JSON-RPC")

	rpcURL := getJSONRPCURL()

	timeout := time.Duration(config.GetInt("bitcoind.jsonrpc.responseTimeout")) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	bodyReq := strings.NewReader("{\"jsonrpc\": \"1.0\", \"id\":\"memod-via-rpc\", \"method\": \"getrawmempool\", \"params\": [true] }")
	req, err := http.NewRequest("POST", rpcURL, bodyReq)
	if err != nil {
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

func getJSONRPCURL() (rpcURL string) {
	return config.GetString("bitcoind.jsonrpc.protocol") +
		"://" + config.GetString("bitcoind.jsonrpc.username") +
		":" + config.GetString("bitcoind.jsonrpc.password") +
		"@" + config.GetString("bitcoind.jsonrpc.host") +
		":" + config.GetString("bitcoind.jsonrpc.port")
}

func saveMempoolJSONRPC() (err error) {
	defer logger.TrackTime(time.Now(), "saveMempoolJSONRPC()")
	logger.Trace.Println("Saving mempool via JSON-RPC")

	client := http.Client{Timeout: 5 * time.Second}
	rpcURL := getJSONRPCURL()

	bodyReq := strings.NewReader("{\"jsonrpc\": \"1.0\", \"id\":\"memod-via-rpc\", \"method\": \"savemempool\" }")
	req, err := http.NewRequest("POST", rpcURL, bodyReq)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JSON-RPC Request failed with status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}

	if string(body) != "{\"result\":null,\"error\":null,\"id\":\"memod-via-rpc\"}\n" {
		return fmt.Errorf("JSON-RPC Request failed with response %s", string(body))
	}

	return nil
}
