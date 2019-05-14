package mempool

import (
	"encoding/json"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

/* decodes and encodes */

// decode the Body of the JSON response as a map of PartialTransactions
func decodeFetchedMempoolBody(body []byte) map[string]PartialTransaction {
	defer logger.TrackTime(time.Now(), "decodeFetchedMempoolBody()")

	var mempool map[string]PartialTransaction
	err := json.Unmarshal(body, &mempool)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	return mempool
}

func encodeCurrentMempoolStatsToJSON(feerateMap map[int]int, megabyteMarkers []int) (string, string) {
	defer logger.TrackTime(time.Now(), "encodeCurrentMempoolStatsToJSON()")
	feerateMapJSON, _ := json.Marshal(feerateMap)
	megabyteMarkersJSON, _ := json.Marshal(megabyteMarkers)
	return string(feerateMapJSON), string(megabyteMarkersJSON)
}
