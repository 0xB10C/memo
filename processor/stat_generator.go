package processor

/* processes the mempool and generates statistics */

import (
	"sort"
	"time"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/types"
)

// MEGABYTE size of one megabyte in byte
const MEGABYTE = 1000000

// COIN satoshi per bitcoin
const COIN = 100000000

// ProcessMempool starts various processing functions on the passed mempool map
func ProcessMempool(mempool map[string]types.PartialTransaction, redisPool *database.RedisPool) {
	if config.GetBool("mempool.processing.processHistoricalMempool") {
		go historicalMempool(mempool, redisPool)
	}

	if config.GetBool("mempool.processing.processCurrentMempool") {
		go currentMempool(mempool, redisPool) // start _current mempool_ stat generation in a goroutine
	}

	if config.GetBool("mempool.processing.processTransactionStats") {
		go transactionStatsMempool(mempool, redisPool)
	}
}

func currentMempool(mempool map[string]types.PartialTransaction, redisPool *database.RedisPool) {
	feerateMap, mempoolSizeInByte, megabyteMarkers := generateCurrentMempoolStats(mempool)

	err := redisPool.WriteCurrentMempoolData(feerateMap, mempoolSizeInByte, megabyteMarkers)
	if err != nil {
		logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
	}

	logger.Info.Println("Success writing Current Mempool to database.")
}

func historicalMempool(mempool map[string]types.PartialTransaction, redisPool *database.RedisPool) {

	const timeframe2h = 1
	const timeframe12h = 2
	const timeframe48h = 3
	const timeframe7d = 4
	const timeframe30d = 5
	const timeframe180d = 6

	needsUpdate, err := redisPool.ReadHistoricalMempoolNeedUpdate()
	if err != nil {
		logger.Error.Printf("Failed to get Needs Update data from database: %s", err.Error())
	}

	if needsUpdate.Update2h || needsUpdate.Update12h || needsUpdate.Update48h || needsUpdate.Update7d || needsUpdate.Update30d || needsUpdate.Update180d {

		countInBuckets, feeInBuckets, sizeInBuckets := generateHistoricalMempoolStats(mempool)

		if needsUpdate.Update2h {
			logger.Info.Println("Writing 2h Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe2h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update12h {
			logger.Info.Println("Writing 12h Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe12h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update48h {
			logger.Info.Println("Writing 48h Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe48h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update7d {
			logger.Info.Println("Writing 7d Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe7d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update30d {
			logger.Info.Println("Writing 30d Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe30d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update180d {
			logger.Info.Println("Writing 180d Historical Mempool data.")
			err = redisPool.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe180d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		logger.Info.Println("Success writing Historical Mempool to database.")
	}
}

func transactionStatsMempool(mempool map[string]types.PartialTransaction, redisPool *database.RedisPool) {
	segwitCount, rbfCount, txCount := generateTransactionStats(mempool)

	err := redisPool.WriteCurrentTransactionStats(segwitCount, rbfCount, txCount)
	if err != nil {
		logger.Error.Printf("Failed to write Transaction Stats to database: %s", err.Error())
		return
	}

	logger.Info.Println("Success writing Transaction Stats to database.")
}

/* generateCurrentMempoolStats()
This function generates the _Current Mempool_ data. Which is:
	- The size of the transactions in the mempool `mempoolSizeInByte`
	- A map mapping the transaction count to the feerate (as a whole
		needsUpdatember). Named `feerateMap`.
	- A list positions in the mempool (tx count), when sorted by
		feerate, which each mark one megabyte worth of transactions.
		Positions starting from the top. Named `megabyteMarkers`.
*/
func generateCurrentMempoolStats(mempool map[string]types.PartialTransaction) (map[int]int, int, []int) {
	defer logger.TrackTime(time.Now(), "generateCurrentMempoolStats()")

	// this represents a entry in a mempool list (memlist).
	// The memlist can be sorted by feerate making it to a
	// memqueue.
	type memlistEntry struct {
		feerate float64
		size    int
	}

	mempoolPos := 0
	mempoolSizeInByte := 0
	var megabyteMarkers []int
	feerateMap := make(map[int]int)
	memlist := make([]memlistEntry, len(mempool))

	for _, tx := range mempool {
		mempoolSizeInByte += tx.Size
		feerate := tx.Fee * COIN / float64(tx.Size)
		feerateMap[int(feerate)]++

		memlist[mempoolPos] = memlistEntry{feerate: feerate, size: tx.Size}
		mempoolPos++
	}

	// sort the list of memlistEntry's by highest feerate first
	sort.Slice(memlist, func(i, j int) bool {
		return memlist[i].feerate > memlist[j].feerate
	})

	memlistPos := len(memlist)
	megabyteBucket := MEGABYTE
	for _, entry := range memlist {
		if megabyteBucket-entry.size > 0 { // if entry.size fits in the bucket
			megabyteBucket = megabyteBucket - entry.size
			memlistPos--
		} else { // if entry.size doesn't fit in the bucket
			megabyteBucket = MEGABYTE - entry.size                // start a new megabyte bucket mineedsUpdates the current entry.size
			megabyteMarkers = append(megabyteMarkers, memlistPos) // append current position to the megabyteMarkers list
			memlistPos--
		}
	}

	return feerateMap, mempoolSizeInByte, megabyteMarkers
}

var feerateBuckets = [40]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 18, 22, 27, 33, 41, 50, 62, 76, 93, 114, 140, 172, 212, 261, 321, 395, 486, 598, 736, 905, 1113, 1369, 1684, 2071, 2547, 3133, 3854, 3855}

// generates a list of counts of transactions representing the count in a feerate bucket
func generateHistoricalMempoolStats(mempool map[string]types.PartialTransaction) (countInBuckets []int, feeInBuckets []float64, sizeInBuckets []int) {

	countInBuckets = make([]int, len(feerateBuckets), len(feerateBuckets))
	feeInBuckets = make([]float64, len(feerateBuckets), len(feerateBuckets))
	sizeInBuckets = make([]int, len(feerateBuckets), len(feerateBuckets))

	for _, tx := range mempool {
		feerate := tx.Fee * COIN / float64(tx.Size)
		bucketIndex := findBucketForFeerate(feerate)
		countInBuckets[bucketIndex]++
		feeInBuckets[bucketIndex] += tx.Fee
		sizeInBuckets[bucketIndex] += tx.Size
	}

	return
}

// finds the bucket index for a given feerate in feerateBuckets. the last bucket is a catch all larger-equal.
// given a feerate of 19 it gives the bucket index for 22. (since it's bigger than 18)
// buckets should be read like "in between or equal feerate [index-1] [index]"
func findBucketForFeerate(feerate float64) int {
	i := sort.Search(len(feerateBuckets), func(i int) bool { return feerateBuckets[i] >= int(feerate) })
	if i < len(feerateBuckets) && feerateBuckets[i] >= int(feerate) {
		return i
	}
	return len(feerateBuckets) - 1
}

func generateTransactionStats(mempool map[string]types.PartialTransaction) (segwitCount int, rbfCount int, txCount int) {

	for txid, tx := range mempool {
		if txid != tx.Wtxid {
			segwitCount++
		}
		if tx.Bip125Replaceable {
			rbfCount++
		}
	}

	txCount = len(mempool)

	return
}
