package processor

/* processes */

import (
	"sort"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/types"
)

// cMEGABYTE: size of one megabyte in byte
const cMEGABYTE = 1000000

// cSATOSHIPERBITCOIN: satoshi per bitcoin
const cSATOSHIPERBITCOIN = 100000000

// ProcessMempool retives the mempool and starts various processing functions on it
func ProcessMempool(mempool map[string]types.PartialTransaction) {
	if config.GetBool("mempool.processing.processHistoricalMempool") {
		go historicalMempool(mempool)
	}

	if config.GetBool("mempool.processing.processCurrentMempool") {
		go currentMempool(mempool) // start _current mempool_ stat generation in a goroutine
	}

	if config.GetBool("mempool.processing.processTimeInMempool") {
		go timeInMempool(mempool)
	}

	if config.GetBool("mempool.processing.processTransactionStats") {
		go transactionStatsMempool(mempool)
	}

}

func currentMempool(mempool map[string]types.PartialTransaction) {
	feerateMap, mempoolSizeInByte, megabyteMarkers := generateCurrentMempoolStats(mempool)

	err := database.WriteCurrentMempoolData(feerateMap, mempoolSizeInByte, megabyteMarkers)
	if err != nil {
		logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
	}

	logger.Info.Println("Success writing Current Mempool to database.")
}

func historicalMempool(mempool map[string]types.PartialTransaction) {

	const timeframe2h = 1
	const timeframe12h = 2
	const timeframe48h = 3
	const timeframe7d = 4
	const timeframe30d = 5
	const timeframe180d = 6

	needsUpdate, err := database.ReadHistroricalMempoolNeedUpdate()
	if err != nil {
		logger.Error.Printf("Failed to get Needs Update data from database: %s", err.Error())
	}

	if needsUpdate.Update2h || needsUpdate.Update12h || needsUpdate.Update48h || needsUpdate.Update7d || needsUpdate.Update30d || needsUpdate.Update180d {

		countInBuckets, feeInBuckets, sizeInBuckets := generateHistoricalMempoolStats(mempool)

		if needsUpdate.Update2h {
			logger.Info.Println("Writing 2h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe2h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update12h {
			logger.Info.Println("Writing 12h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe12h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update48h {
			logger.Info.Println("Writing 48h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe48h)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update7d {
			logger.Info.Println("Writing 7d Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe7d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update30d {
			logger.Info.Println("Writing 30d Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe30d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		if needsUpdate.Update180d {
			logger.Info.Println("Writing 180d Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBuckets, feeInBuckets, sizeInBuckets, timeframe180d)
			if err != nil {
				logger.Error.Printf("Failed to write Historical Mempool to database: %s", err.Error())
				return
			}
		}

		logger.Info.Println("Success writing Historical Mempool to database.")
	}
}

func timeInMempool(mempool map[string]types.PartialTransaction) {
	timeAxis, feerateAxis := generateTimeInMempoolStats(mempool)

	err := database.WriteTimeInMempoolData(timeAxis, feerateAxis)
	if err != nil {
		logger.Error.Printf("Failed to write Time in Mempool to database: %s", err.Error())
		return
	}

	logger.Info.Println("Success writing Time in Mempool to database.")
}

func transactionStatsMempool(mempool map[string]types.PartialTransaction) {
	segwitCount, rbfCount, txCount := generateTransactionStats(mempool)

	err := database.WriteCurrentTransactionStats(segwitCount, rbfCount, txCount)
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
		feerate := tx.Fee * cSATOSHIPERBITCOIN / float64(tx.Size)
		feerateMap[int(feerate)]++

		memlist[mempoolPos] = memlistEntry{feerate: feerate, size: tx.Size}
		mempoolPos++
	}

	// sort the list of memlistEntry's by highest feerate first
	sort.Slice(memlist, func(i, j int) bool {
		return memlist[i].feerate > memlist[j].feerate
	})

	memlistPos := len(memlist)
	megabyteBucket := cMEGABYTE
	for _, entry := range memlist {
		if megabyteBucket-entry.size > 0 { // if entry.size fits in the bucket
			megabyteBucket = megabyteBucket - entry.size
			memlistPos--
		} else { // if entry.size doesn't fit in the bucket
			megabyteBucket = cMEGABYTE - entry.size               // start a new megabyte bucket mineedsUpdates the current entry.size
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
		feerate := tx.Fee * cSATOSHIPERBITCOIN / float64(tx.Size)
		bucketIndex := findBucketForFeerate(feerate)
		countInBuckets[bucketIndex]++
		feeInBuckets[bucketIndex] += tx.Fee
		sizeInBuckets[bucketIndex] += tx.Size
	}

	return
}

// finds the bucket index for a given feerate in feerateBuckets. the last bucket is a catch all larger-equal.
// given a feerate of 19 it gives the bucket index for 22. (since it's bigger than 18)
// buckets should be read like "inbetween or equal feerate [index-1] [index]"
func findBucketForFeerate(feerate float64) int {
	i := sort.Search(len(feerateBuckets), func(i int) bool { return feerateBuckets[i] >= int(feerate) })
	if i < len(feerateBuckets) && feerateBuckets[i] >= int(feerate) {
		return i
	}
	return len(feerateBuckets) - 1
}

func generateTimeInMempoolStats(mempool map[string]types.PartialTransaction) ([]int, []float64) {
	/* We two slices:
	- one with the timestamp the transaction entered the mempool
	- one with the feerate it paid
	*/

	timeAxis := make([]int, 0)
	feerateAxis := make([]float64, 0)

	for _, tx := range mempool {
		feerate := tx.Fee * cSATOSHIPERBITCOIN / float64(tx.Size)
		feerateTruncated := float64(int(feerate*1000)) / 1000 // same as toFixed(3)

		feerateAxis = append(feerateAxis, feerateTruncated)
		timeAxis = append(timeAxis, tx.Time)
	}

	return timeAxis, feerateAxis
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
	//segwitPercentage = float64(int(float64(segwitCount)/float64(txCount)*1000)) / 1000
	//rbfPercentage = float64(int(float64(rbfCount)/float64(txCount)*1000)) / 1000

	return
}
