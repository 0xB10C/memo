package database

import (
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

type needsUpdate struct {
	Update2h   bool
	Update12h  bool
	Update48h  bool
	Update7d   bool
	Update30d  bool
	Update180d bool
}

func ReadHistroricalMempoolNeedUpdate() (nu needsUpdate, err error) {
	sqlStatement := `SELECT
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 1 ORDER BY timestamp DESC LIMIT 1) AS timediff2h,
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 2 ORDER BY timestamp DESC LIMIT 1) AS timediff12h,
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 3 ORDER BY timestamp DESC LIMIT 1) AS timediff48h,
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 4 ORDER BY timestamp DESC LIMIT 1) AS timediff7d,
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 5 ORDER BY timestamp DESC LIMIT 1) AS timediff30d,
	(select UTC_TIMESTAMP-timestamp from historicalMempool where timeframe = 6 ORDER BY timestamp DESC LIMIT 1) AS timediff180d;`

	row := Database.QueryRow(sqlStatement)
	var timediff2h, timediff12h, timediff48h, timediff7d, timediff30d, timediff180d int

	err = row.Scan(&timediff2h, &timediff12h, &timediff48h, &timediff7d, &timediff30d, &timediff180d)
	if err != nil {
		return nu, err
	}

	// Update 2h data every 4 minutes
	if time.Duration(timediff2h)*time.Second >= 4*time.Minute {
		nu.Update2h = true
		logger.Trace.Println("2h Historical Mempool data needs to be updated")
	}

	// Update 12h data every 24 minutes
	if time.Duration(timediff12h)*time.Second >= 24*time.Minute {
		nu.Update12h = true
		logger.Trace.Println("12h Historical Mempool data needs to be updated")
	}

	// Update 48h data every 96 minutes
	if time.Duration(timediff48h)*time.Second >= 96*time.Minute {
		nu.Update48h = true
		logger.Trace.Println("48h Historical Mempool data needs to be updated")
	}

	// Update 7d data every 336 minutes
	if time.Duration(timediff7d)*time.Second >= 336*time.Minute {
		nu.Update7d = true
		logger.Trace.Println("7d Historical Mempool data needs to be updated")
	}

	// Update 30d data every 1440 minutes
	if time.Duration(timediff30d)*time.Second >= 1440*time.Minute {
		nu.Update30d = true
		logger.Trace.Println("30d Historical Mempool data needs to be updated")
	}

	// Update 180d data every 8640 minutes
	if time.Duration(timediff180d)*time.Second >= 8640*time.Minute {
		nu.Update180d = true
		logger.Trace.Println("180d Historical Mempool data needs to be updated")
	}

	return nu, nil
}
