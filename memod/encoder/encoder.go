package encoder

import (
	"encoding/json"
	"time"

	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/types"
)

/* decodes and encodes */

// DecodeFetchedMempoolBody decode the Body of the JSON response as a map of PartialTransactions
func DecodeFetchedMempoolBody(body []byte) (map[string]types.PartialTransaction, error) {
	defer logger.TrackTime(time.Now(), "decodeFetchedMempoolBody()")

	var mempool map[string]types.PartialTransaction
	err := json.Unmarshal(body, &mempool)
	if err != nil {
		return nil, err
	}

	return mempool, nil
}

// EncodeCurrentMempoolStatsToJSON encodes the current mempool stats to JSON
func EncodeCurrentMempoolStatsToJSON(feerateMap map[int]int, megabyteMarkers []int) (string, string, error) {
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
