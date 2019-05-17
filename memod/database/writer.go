package database

import (
	"errors"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

/* wirtes data to database */

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
