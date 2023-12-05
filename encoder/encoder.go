package encoder

import (
	"encoding/json"
	"time"

	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/types"
)

// DecodeFetchedMempoolBody decode the Body of the JSON response as a map of PartialMempoolEntry
func DecodeFetchedMempoolBody(body []byte) (map[string]types.PartialMempoolEntry, error) {
	defer logger.TrackTime(time.Now(), "decodeFetchedMempoolBody()")

	var mempool map[string]types.PartialMempoolEntry
	err := json.Unmarshal(body, &mempool)
	if err != nil {
		return nil, err
	}

	return mempool, nil
}
