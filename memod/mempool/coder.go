package mempool

import (
	"encoding/json"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

/* decodes and encodes */

// decode the Body of the JSON response as a map of PartialTransactions
func decodeFetchedMempoolBody(body []byte) (map[string]PartialTransaction, error) {
	defer logger.TrackTime(time.Now(), "decodeFetchedMempoolBody()")

	var mempool map[string]PartialTransaction
	err := json.Unmarshal(body, &mempool)
	if err != nil {
		return nil, err
	}

	return mempool, nil
}

func encodeCurrentMempoolStatsToJSON(feerateMap map[int]int, megabyteMarkers []int) (string, string, error) {
	defer logger.TrackTime(time.Now(), "encodeCurrentMempoolStatsToJSON()")

	feerateMapJSON, err := json.Marshal(feerateMap)
	if err != nil {
		return "", "", err
	}

	megabyteMarkersJSON, err := json.Marshal(megabyteMarkers)
	if err != nil {
		return "", "", err
	}

	return string(feerateMapJSON), string(megabyteMarkersJSON), nil
}
