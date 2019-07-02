package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/0xb10c/memo/memod/logger"
	"github.com/gomodule/redigo/redis"
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
	defer logger.TrackTime(time.Now(), "ReadHistroricalMempoolNeedUpdate()")
	c := Pool.Get()
	defer c.Close()

	lastUpdatedTimesStrings, err := redis.Strings(c.Do("MGET", "historicalMempool1:lastUpdated", "historicalMempool2:lastUpdated", "historicalMempool3:lastUpdated", "historicalMempool4:lastUpdated", "historicalMempool5:lastUpdated", "historicalMempool6:lastUpdated"))
	if err != nil {
		return
	}

	// convert responses from string to int64 unix timestamp and
	// calculate the time difference between now and the last updated time
	lastUpdatedTimeDiffs := make([]time.Duration, 0)
	for _, lastUpdatedString := range lastUpdatedTimesStrings {
		if n, err := strconv.Atoi(lastUpdatedString); err == nil {
			lastUpdatedTime := time.Unix(int64(n), 0)
			timeDiff := time.Duration(time.Now().Unix()-lastUpdatedTime.Unix()) * time.Second
			lastUpdatedTimeDiffs = append(lastUpdatedTimeDiffs, timeDiff)
		} else {
			fmt.Println(lastUpdatedString, "is not an integer.")
			lastUpdatedTimeDiffs = append(lastUpdatedTimeDiffs, time.Duration(1000)*time.Hour)
		}
	}

	// Update 2h data every 4 minutes
	if lastUpdatedTimeDiffs[0] >= 4*time.Minute {
		nu.Update2h = true
		logger.Trace.Println("2h Historical Mempool data needs to be updated")
	}

	// Update 12h data every 24 minutes
	if lastUpdatedTimeDiffs[1] >= 24*time.Minute {
		nu.Update12h = true
		logger.Trace.Println("12h Historical Mempool data needs to be updated")
	}

	// Update 48h data every 96 minutes
	if lastUpdatedTimeDiffs[2] >= 96*time.Minute {
		nu.Update48h = true
		logger.Trace.Println("48h Historical Mempool data needs to be updated")
	}

	// Update 7d data every 336 minutes
	if lastUpdatedTimeDiffs[3] >= 336*time.Minute {
		nu.Update7d = true
		logger.Trace.Println("7d Historical Mempool data needs to be updated")
	}

	// Update 30d data every 1440 minutes
	if lastUpdatedTimeDiffs[4] >= 1440*time.Minute {
		nu.Update30d = true
		logger.Trace.Println("30d Historical Mempool data needs to be updated")
	}

	// Update 180d data every 8640 minutes
	if lastUpdatedTimeDiffs[5] >= 8640*time.Minute {
		nu.Update180d = true
		logger.Trace.Println("180d Historical Mempool data needs to be updated")
	}

	return
}
