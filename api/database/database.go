package database

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/0xb10c/memo/api/config"
	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

func newPool() *redis.Pool {
	dbUser := config.GetString("redis.user")
	dbPasswd := config.GetString("redis.passwd")
	dbHost := config.GetString("redis.host")
	dbPort := config.GetString("redis.port")
	dbConnection := config.GetString("redis.connection")

	connectionString := dbConnection + "://" + dbUser + ":" + dbPasswd + "@" + dbHost + ":" + dbPort

	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {

			c, err := redis.DialURL(connectionString)
			if err != nil {
				return c, err
			}
			return c, err
		},
	}
}

func SetupDatabase() error {
	Pool = newPool()

	c := Pool.Get() // get a new connection
	defer c.Close()

	_, err := c.Do("PING")
	if err != nil {
		return err
	}

	return nil
}

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

func GetRecentBlocks() (blocks []RecentBlock, err error) {
	c := Pool.Get()
	defer c.Close()

	reJSON, err := redis.Strings(c.Do("LRANGE", "recentBlocks", 0, 9))
	if err != nil {
		return
	}

	blocks = make([]RecentBlock, 0)
	for index := range reJSON {
		block := RecentBlock{}
		err = json.Unmarshal([]byte(reJSON[index]), &block)
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

func GetHistorical(timeframe int, by string) (hmds []HistoricalMempoolData, err error) {
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

	hmds = make([]HistoricalMempoolData, 0)
	for index := range hmdsJSON {
		hmd := HistoricalMempoolData{}
		err = json.Unmarshal([]byte(hmdsJSON[index]), &hmd)
		if err != nil {
			return
		}
		hmds = append(hmds, hmd)
	}

	return
}

// GetTimeInMempool gets the TimeInMempool data from the database
func GetTimeInMempool() (timestamp int64, timeAxis []int, feerateAxis []float64, err error) {
	c := Pool.Get()
	defer c.Close()

	prefix := "timeInMempool"

	response, err := redis.Strings(c.Do("MGET", prefix+":timeAxis", prefix+":feerateAxis", prefix+":utcTimestamp"))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(response[0]), &timeAxis)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(response[1]), &feerateAxis)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(response[2]), &timestamp)
	if err != nil {
		return
	}

	return
}

// GetTransactionStats gets the Transaction Stats data from the database
func GetTransactionStats() (tss []TransactionStat, err error) {
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

	tss = make([]TransactionStat, 0)
	for index := range tssJSON {
		ts := TransactionStat{}
		err = json.Unmarshal([]byte(tssJSON[index]), &ts)
		if err != nil {
			return
		}
		tss = append(tss, ts)
	}

	return
}
