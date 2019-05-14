package mempool

import (
	"time"

	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/logger"
)

/* wirtes data to database */

func writeCurrentMempoolData(feerateMapJSON string, mempoolSizeInByte int, megabyteMarkersJSON string) {
	defer logger.TrackTime(time.Now(), "writeCurrentMempoolData()")

	if database.Database != nil {
		sql := "UPDATE current_mempool SET byCount = ?, positionsInGreedyBlocks = ?, timestamp = UTC_TIMESTAMP, mempoolSize = ? WHERE id = 1"
		_, err := database.Database.Exec(sql, feerateMapJSON, megabyteMarkersJSON, mempoolSizeInByte)
		if err != nil {
			panic(err)
		}
	} else {
		logger.Error.Println("Database is nil")
	}

}
