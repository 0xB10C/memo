package database

import (
	"errors"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

// WriteCurrentMempoolData writes the current mempool data into the database
func WriteCurrentMempoolData(feerateMapJSON string, mempoolSizeInByte int, megabyteMarkersJSON string) error {
	defer logger.TrackTime(time.Now(), "writeCurrentMempoolData()")

	if Database != nil {
		sql := "UPDATE current_mempool SET byCount = ?, positionsInGreedyBlocks = ?, timestamp = UTC_TIMESTAMP, mempoolSize = ? WHERE id = 1"
		_, err := Database.Exec(sql, feerateMapJSON, megabyteMarkersJSON, mempoolSizeInByte)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Database pointer is nil")

}

// WriteNewBlockData writes data for a new block into the database
func WriteNewBlockData(height int, numTx int, sizeWithWitness int, weight int) error {
	defer logger.TrackTime(time.Now(), "writeNewBlockData()")

	if Database != nil {
		sql := "REPLACE INTO pastBlocks SET height = ?, txCount = ?, size= ?, weight = ?, timestamp = UTC_TIMESTAMP"

		_, err := Database.Exec(sql, height, numTx, sizeWithWitness, weight)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Database pointer is nil")
}

// WriteHistoricalMempoolData writes the histoical mempool data into the database
func WriteHistoricalMempoolData(countInBucketsJSON string, feeInBucketsJSON string, sizeInBucketsJSON string, timeframe int) error {
	defer logger.TrackTime(time.Now(), "WriteHistoricalMempoolData()")

	if Database != nil {
		sql := "INSERT INTO historicalMempool(timeframe, timestamp, countInBuckets, feeInBuckets, sizeInBuckets) VALUES (?, UTC_TIMESTAMP, ?, ?, ?)"
		_, err := Database.Exec(sql, timeframe, countInBucketsJSON, feeInBucketsJSON, sizeInBucketsJSON)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Database pointer is nil")

}
