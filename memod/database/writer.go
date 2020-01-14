package database

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/0xb10c/memo/memod/config"

	"github.com/gomodule/redigo/redis"

	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/types"
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
func WriteNewBlockEntry(height int, shortTXIDs []string) error {
	defer logger.TrackTime(time.Now(), "WriteNewBlockEntry()")
	c := Pool.Get()
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
func WriteHistoricalMempoolData(countInBuckets []int, feeInBuckets []float64, sizeInBuckets []int, timeframe int) error {
	defer logger.TrackTime(time.Now(), "WriteHistoricalMempoolData()")
	c := Pool.Get()
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
func WriteCurrentTransactionStats(segwitCount int, rbfCount int, txCount int) error {
	defer logger.TrackTime(time.Now(), "WriteCurrentTransactionStats()")
	c := Pool.Get()
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
func WriteMempoolEntries(me types.MempoolEntry) error {
	//defer logger.TrackTime(time.Now(), "WriteMempoolEntries()")
	c := Pool.Get()
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
func WriteFeerateAPIEntry(entry types.FeeRateAPIEntry) error {
	defer logger.TrackTime(time.Now(), "WriteFeerateAPIEntry()")
	c := Pool.Get()
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
