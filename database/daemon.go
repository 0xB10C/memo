package database

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/0xb10c/memo/config"

	"github.com/gomodule/redigo/redis"

	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/types"
)

// WriteCurrentMempoolData writes the current mempool data into the database
func (p *RedisPool) WriteCurrentMempoolData(feerateMap map[int]int, mempoolSizeInByte int, megabyteMarkers []int) error {
	defer logger.TrackTime(time.Now(), "writeCurrentMempoolData()")
	c := p.Get()
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
func (p *RedisPool) WriteNewBlockData(height int, numTx int, sizeWithWitness int, weight int) error {
	defer logger.TrackTime(time.Now(), "writeNewBlockData()")
	c := p.Get()
	defer c.Close()
	listName := "recentBlocks"

	rb := types.RecentBlock{Height: height, Size: sizeWithWitness, Timestamp: time.Now().Unix(), TxCount: numTx, Weight: weight}

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

// WriteNewBlockEntry writes a BlockEntry into the Redis database.
func (p *RedisPool) WriteNewBlockEntry(height int, shortTXIDs []string) error {
	defer logger.TrackTime(time.Now(), "WriteNewBlockEntry()")
	c := p.Get()
	defer c.Close()
	listName := "blockEntries"

	be := types.BlockEntry{Height: height, Timestamp: time.Now().Unix(), ShortTXIDs: shortTXIDs}

	beJSON, err := json.Marshal(be)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", listName, beJSON)
	if err != nil {
		return err
	}

	// only keep the last ~200 blocks
	// check 10% of all insertions
	if rand.Intn(10) == 4 {
		_, err := c.Do("LTRIM", listName, "0", 200)
		if err != nil {
			return fmt.Errorf("could not do LTRIM on %s: %s", listName, err.Error())
		}
	}

	return nil
}

// WriteHistoricalMempoolData writes the histoical mempool data into the database
func (p *RedisPool) WriteHistoricalMempoolData(countInBuckets []int, feeInBuckets []float64, sizeInBuckets []int, timeframe int) error {
	defer logger.TrackTime(time.Now(), "WriteHistoricalMempoolData()")
	c := p.Get()
	defer c.Close()

	countInBucketsJSON, err := json.Marshal(types.HistoricalMempoolData{DataInBuckets: countInBuckets, Timestamp: time.Now().Unix()})
	if err != nil {
		return err
	}

	feeInBucketsJSON, err := json.Marshal(types.HistoricalMempoolData{DataInBuckets: feeInBuckets, Timestamp: time.Now().Unix()})
	if err != nil {
		return err
	}

	sizeInBucketsJSON, err := json.Marshal(types.HistoricalMempoolData{DataInBuckets: sizeInBuckets, Timestamp: time.Now().Unix()})
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

// WriteCurrentTransactionStats writes the current transaction stats into the database
func (p *RedisPool) WriteCurrentTransactionStats(segwitCount int, rbfCount int, txCount int) error {
	defer logger.TrackTime(time.Now(), "WriteCurrentTransactionStats()")
	c := p.Get()
	defer c.Close()

	ts := types.TransactionStat{SegwitCount: segwitCount, RbfCount: rbfCount, TxCount: txCount, Timestamp: time.Now().Unix()}
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

// WriteMempoolEntries writes a txid and it's feerate to the database
func (p *RedisPool) WriteMempoolEntries(me types.MempoolEntry) error {
	//defer logger.TrackTime(time.Now(), "WriteMempoolEntries()")
	c := p.Get()
	defer c.Close()

	meJSON, err := json.Marshal(me)
	if err != nil {
		return fmt.Errorf("could not marshal the mempoolEntry to JSON: %s", err.Error())
	}

	listName := "mempoolEntries"

	// insert the mempool entry into a redis sorted list
	// the list is sorted by timestamps in ascending order.
	_, err = c.Do("ZADD", listName, me.EntryTime, meJSON)
	if err != nil {
		return fmt.Errorf("could not ZADD the mempoolEntry JSON to %s: %s", listName, err.Error())
	}

	if config.GetBool("zmq.saveMempoolEntries.enable") {
		err = writeMempoolEntriesSQLite(me)
		if err != nil {
			return fmt.Errorf("could not write the mempoolEntry to SQLite %s", err.Error())
		}
	}

	// only keep the last ~500k transactions
	// check every 0.1% of all insertions
	if rand.Intn(1000) == 42 {
		count, err := redis.Int64(c.Do("ZCOUNT", listName, "-inf", "+inf"))
		if err != nil {
			return fmt.Errorf("could not do ZCOUNT on %s: %s", listName, err.Error())
		}
		if count > 500000 {
			removeIndex := count - 500000
			_, err := c.Do("ZREMRANGEBYRANK", listName, "0", removeIndex)
			if err != nil {
				return fmt.Errorf("could not do ZREMRANGEBYRANK on %s: %s", listName, err.Error())
			}
		}
	}

	return nil
}

func writeMempoolEntriesSQLite(me types.MempoolEntry) error {
	spendsJSON, err := json.Marshal(me.Spends)
	if err != nil {
		return fmt.Errorf("could not marshal mempoolEntry.Spends to JSON: %s", err.Error())
	}

	paystoJSON, err := json.Marshal(me.PaysTo)
	if err != nil {
		return fmt.Errorf("could not marshal mempoolEntry.PaysTo to JSON: %s", err.Error())
	}

	multisigJSON, err := json.Marshal(me.Multisig)
	if err != nil {
		return fmt.Errorf("could not marshal mempoolEntry.Multisig to JSON: %s", err.Error())
	}

	_, err = SQLiteDB.Exec("INSERT INTO mempoolEntries(entryTime, txid, fee, size, version, inputs, outputs, locktime, outSum, spendsSegWit, spendsMultisig, bip69compliant, signalsRBF, spends, paysto, multisigs, opreturndata) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		me.EntryTime, me.TxID, me.Fee, me.Size, me.Version, me.InputCount, me.OutputCount, me.Locktime, me.OutputSum, me.SpendsSegWit, me.SpendsMultisig, me.IsBIP69, me.SignalsRBF, spendsJSON, paystoJSON, multisigJSON, me.OPReturnData)
	if err != nil {
		return fmt.Errorf("error while writing to SQLite: %s", err.Error())
	}
	return nil
}

// WriteFeerateAPIEntry writes a new feerate API entry into the database
func (p *RedisPool) WriteFeerateAPIEntry(entry types.FeeRateAPIEntry) error {
	defer logger.TrackTime(time.Now(), "WriteFeerateAPIEntry()")
	c := p.Get()
	defer c.Close()
	listName := "feerateAPIEntries"

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", listName, entryJSON)
	if err != nil {
		return err
	}

	return nil
}

// NeedsUpdate holds information about the timeframes that need an update
type NeedsUpdate struct {
	Update2h   bool
	Update12h  bool
	Update48h  bool
	Update7d   bool
	Update30d  bool
	Update180d bool
}

// ReadHistoricalMempoolNeedUpdate checks which time frame is due (needs an update)
func (p *RedisPool) ReadHistoricalMempoolNeedUpdate() (nu NeedsUpdate, err error) {
	defer logger.TrackTime(time.Now(), "ReadHistroricalMempoolNeedUpdate()")
	c := p.Get()
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
