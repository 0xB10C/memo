package database

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

// WriteCurrentMempoolData writes the current mempool data into the database
func WriteCurrentMempoolData(feerateMap map[int]int, mempoolSizeInByte int, megabyteMarkers []int) error {
	defer logger.TrackTime(time.Now(), "writeCurrentMempoolData()")
	c := Pool.Get()
	defer c.Close()

	feerateMapJSON, err := json.Marshal(feerateMap)
	if err != nil {
		return err
	}

	megabyteMarkersJSON, err := json.Marshal(megabyteMarkers)
	if err != nil {
		return err
	}

	prefix := "currentMempool"

	c.Send("MULTI")
	c.Send("SET", prefix+":feerateMap", feerateMapJSON)
	c.Send("SET", prefix+":mempoolSizeInByte", mempoolSizeInByte)
	c.Send("SET", prefix+":megabyteMarkers", megabyteMarkersJSON)
	c.Send("SET", prefix+":utcTimestamp", time.Now().Unix())
	_, err = c.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// WriteNewBlockData writes data for a new block into the database
func WriteNewBlockData(height int, numTx int, sizeWithWitness int, weight int) error {
	defer logger.TrackTime(time.Now(), "writeNewBlockData()")
	c := Pool.Get()
	defer c.Close()
	listName := "recentBlocks"

	rb := recentBlock{height, sizeWithWitness, time.Now().Unix(), numTx, weight}

	rbJSON, err := json.Marshal(rb)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", listName, rbJSON)
	if err != nil {
		return err
	}

	return nil
}

// WriteHistoricalMempoolData writes the histoical mempool data into the database
func WriteHistoricalMempoolData(countInBuckets []int, feeInBuckets []float64, sizeInBuckets []int, timeframe int) error {
	defer logger.TrackTime(time.Now(), "WriteHistoricalMempoolData()")
	c := Pool.Get()
	defer c.Close()

	countInBucketsJSON, err := json.Marshal(historicalMempoolData{countInBuckets, time.Now().Unix()})
	if err != nil {
		return err
	}

	feeInBucketsJSON, err := json.Marshal(historicalMempoolData{feeInBuckets, time.Now().Unix()})
	if err != nil {
		return err
	}

	sizeInBucketsJSON, err := json.Marshal(historicalMempoolData{sizeInBuckets, time.Now().Unix()})
	if err != nil {
		return err
	}

	listName := "historicalMempool" + strconv.Itoa(timeframe)

	_, err = c.Do("LPUSH", listName+":countInBuckets", countInBucketsJSON)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", listName+":feeInBuckets", feeInBucketsJSON)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", listName+":sizeInBuckets", sizeInBucketsJSON)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", listName+":lastUpdated", time.Now().Unix())

	return nil
}

// WriteTimeInMempoolData writes the time-in-mempool data into the database
func WriteTimeInMempoolData(timeAxis []int, feerateAxis []float64) error {
	defer logger.TrackTime(time.Now(), "WriteTimeInMempoolData()")
	c := Pool.Get()
	defer c.Close()
	prefix := "timeInMempool"

	timeAxisJSON, err := json.Marshal(timeAxis)
	if err != nil {
		return err
	}

	feerateAxisJSON, err := json.Marshal(feerateAxis)
	if err != nil {
		return err
	}

	c.Send("MULTI")
	c.Send("SET", prefix+":timeAxis", timeAxisJSON)
	c.Send("SET", prefix+":feerateAxis", feerateAxisJSON)
	c.Send("SET", prefix+":utcTimestamp", time.Now().Unix())
	_, err = c.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// WriteCurrentTransactionStats writes the current transaction stats into the database
func WriteCurrentTransactionStats(segwitCount int, rbfCount int, txCount int) error {
	defer logger.TrackTime(time.Now(), "WriteCurrentTransactionStats()")
	c := Pool.Get()
	defer c.Close()

	ts := transactionStat{segwitCount, rbfCount, txCount, time.Now().Unix()}
	tsJSON, err := json.Marshal(ts)
	if err != nil {
		return err
	}

	listName := "transactionStats"

	_, err = c.Do("LPUSH", listName, tsJSON)
	if err != nil {
		return err
	}

	return nil
}
