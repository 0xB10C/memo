package database

import (
	"errors"
	"strconv"
	"time"

	"github.com/0xb10c/memo/api/types"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetMempool gets the current mempool from the database
func GetMempool() (timestamp time.Time, feerateMapJSON string, megabyteMarkersJSON string, mempoolSize int, err error) {
	c := Pool.Get()
	defer c.Close()

	prefix := "currentMempool"

	response, err := redis.Strings(c.Do("MGET", prefix+":utcTimestamp", prefix+":feerateMap", prefix+":megabyteMarkers", prefix+":mempoolSizeInByte"))
	if err != nil {
		return
	}

	if n, err := strconv.Atoi(response[0]); err == nil {
		timestamp = time.Unix(int64(n), 0)
	} else {
		timestamp = time.Unix(0, 0)
	}

	feerateMapJSON = response[1]
	megabyteMarkersJSON = response[2]

	if n, err := strconv.Atoi(response[3]); err == nil {
		mempoolSize = n
	} else {
		mempoolSize = 0
	}

	return
}

// GetRecentBlocks returns the 10 most recent blocks.
func GetRecentBlocks() (blocks []types.RecentBlock, err error) {
	c := Pool.Get()
	defer c.Close()

	reJSON, err := redis.Strings(c.Do("LRANGE", "recentBlocks", 0, 9))
	if err != nil {
		return
	}

	blocks = make([]types.RecentBlock, 0)
	for index := range reJSON {
		block := types.RecentBlock{}
		err = json.Unmarshal([]byte(reJSON[index]), &block)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}

	return
}

// GetBlockEntries returns the 20 most recent blocks with short TXIDs.
func GetBlockEntries() (blocks []types.BlockEntry, err error) {
	c := Pool.Get()
	defer c.Close()

	beJSON, err := redis.Strings(c.Do("LRANGE", "blockEntries", 0, 20))
	if err != nil {
		return
	}

	blocks = make([]types.BlockEntry, 0)
	for index := range beJSON {
		block := types.BlockEntry{}
		err = json.Unmarshal([]byte(beJSON[index]), &block)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}

	return
}

type MempoolState struct {
	time              time.Time
	Timestamp         int64 `json:"timestamp"`
	dataInBucketsJSON string
	DataInBuckets     []float64 `json:"dataInBuckets"`
}

func GetHistorical(timeframe int, by string) (hmds []types.HistoricalMempoolData, err error) {
	c := Pool.Get()
	defer c.Close()

	var timeSelector string
	switch timeframe {
	case 1:
		timeSelector = "historicalMempool1"
	case 2:
		timeSelector = "historicalMempool2"
	case 3:
		timeSelector = "historicalMempool3"
	case 4:
		timeSelector = "historicalMempool4"
	case 5:
		timeSelector = "historicalMempool5"
	case 6:
		timeSelector = "historicalMempool6"
	default:
		return hmds, errors.New("Invalid input")
	}

	var bySelector string
	switch by {
	case "byCount":
		bySelector = "countInBuckets"
	case "byFee":
		bySelector = "feeInBuckets"
	case "bySize":
		bySelector = "sizeInBuckets"
	default:
		return hmds, errors.New("Invalid input")
	}

	hmdsJSON, err := redis.Strings(c.Do("LRANGE", timeSelector+":"+bySelector, 0, 29))
	if err != nil {
		return
	}

	hmds = make([]types.HistoricalMempoolData, 0)
	for index := range hmdsJSON {
		hmd := types.HistoricalMempoolData{}
		err = json.Unmarshal([]byte(hmdsJSON[index]), &hmd)
		if err != nil {
			return
		}
		hmds = append(hmds, hmd)
	}

	return
}

// GetTransactionStats gets the Transaction Stats data from the database
func GetTransactionStats() (tss []types.TransactionStat, err error) {
	c := Pool.Get()
	defer c.Close()

	tssJSON, err := redis.Strings(c.Do("LRANGE", "transactionStats", 0, 180))
	if err != nil {
		return
	}

	type transactionStat struct {
		SegwitCount int   `json:"segwitCount"`
		RbfCount    int   `json:"rbfCount"`
		TxCount     int   `json:"txCount"`
		Timestamp   int64 `json:"timestamp"`
	}

	tss = make([]types.TransactionStat, 0)
	for index := range tssJSON {
		ts := types.TransactionStat{}
		err = json.Unmarshal([]byte(tssJSON[index]), &ts)
		if err != nil {
			return
		}
		tss = append(tss, ts)
	}

	return
}

// GetMempoolEntries gets the last x mempool Entries from the database
func GetMempoolEntries() (mes []types.MempoolEntry, err error) {
	c := Pool.Get()
	defer c.Close()

	// gets recent entries from 0 to 19999 (20k)
	mesJSON, err := redis.Strings(c.Do("ZREVRANGE", "mempoolEntries", 0, 29999))
	if err != nil {
		return
	}

	mes = make([]types.MempoolEntry, 0)
	for index := range mesJSON {
		me := types.MempoolEntry{}
		err = json.Unmarshal([]byte(mesJSON[index]), &me)
		if err != nil {
			return
		}
		mes = append(mes, me)
	}

	return
}

// SetMempoolEntriesCache SETs the response of a recent GetMempoolEntries() as a cache
func SetMempoolEntriesCache(mesJSON string) (err error) {
	c := Pool.Get()
	defer c.Close()

	_, err = c.Do("SET", "cache:mempoolEntries", mesJSON)
	if err != nil {
		return err
	}
	return nil
}

// GetMempoolEntriesCache GET the cached response of a recent GetMempoolEntries() call
func GetMempoolEntriesCache() (mesJSON string, err error) {
	c := Pool.Get()
	defer c.Close()

	mesJSON, err = redis.String(c.Do("GET", "cache:mempoolEntries"))
	if err != nil {
		return
	}
	return
}

// GetRecentFeerateAPIEntries returns the recent feeRate API entries from Redis
func GetRecentFeerateAPIEntries() (entries []types.FeeRateAPIEntry, err error) {
	c := Pool.Get()
	defer c.Close()

	entriesJSON, err := redis.Strings(c.Do("LRANGE", "feerateAPIEntries", 0, 400))
	if err != nil {
		return
	}

	entries = make([]types.FeeRateAPIEntry, 0)
	for index := range entriesJSON {
		entry := types.FeeRateAPIEntry{}
		err = json.Unmarshal([]byte(entriesJSON[index]), &entry)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}

	return
}
