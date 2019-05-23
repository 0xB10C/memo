package processor

/* processes */

import (
	"sort"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/encoder"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/types"
)

// cMEGABYTE: size of one megabyte in byte
const cMEGABYTE = 1000000

// cSATOSHIPERBITCOIN: satoshi per bitcoin
const cSATOSHIPERBITCOIN = 100000000

// ProcessMempool retives the mempool and starts various processing functions on it
func ProcessMempool(mempool map[string]types.PartialTransaction) {

	go historicalMempool(mempool)

	if config.GetBool("mempool.processing.processCurrentMempool") {
		go currentMempool(mempool) // start _current mempool_ stat generation in a goroutine
	}

}

func currentMempool(mempool map[string]types.PartialTransaction) {
	feerateMap, mempoolSizeInByte, megabyteMarkers := generateCurrentMempoolStats(mempool)
	feerateMapJSON, megabyteMarkersJSON, err := encoder.EncodeCurrentMempoolStatsToJSON(feerateMap, megabyteMarkers)
	if err != nil {
		logger.Error.Printf("Failed to encode generated data as JSON: %s", err.Error())
		return
	}

	err = database.WriteCurrentMempoolData(feerateMapJSON, mempoolSizeInByte, megabyteMarkersJSON)
	if err != nil {
		logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
		return
	}

	logger.Info.Println("Success writing Current Mempool to database.")
}

func historicalMempool(mempool map[string]types.PartialTransaction) {

	const timeframe2h = 1
	const timeframe12h = 2
	const timeframe48h = 3
	const timeframe7d = 4

	nu, err := database.ReadHistroricalMempoolNeedUpdate()
	if err != nil {
		logger.Error.Printf("Failed to get Needs Update data from database: %s", err.Error())
	}

	if nu.Update2h || nu.Update12h || nu.Update48h || nu.Update7d {

		countInBuckets := generateHistoricalMempoolStats(mempool)
		countInBucketsJSON, err := encoder.EncodeHistoricalStatsToJSON(countInBuckets)
		if err != nil {
			logger.Error.Printf("Failed to encode generated data as JSON: %s", err.Error())
			return
		}

		if nu.Update2h {
			logger.Info.Println("Writing 2h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBucketsJSON, timeframe2h)
			if err != nil {
				logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
				return
			}
		}

		if nu.Update12h {
			logger.Info.Println("Writing 12h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBucketsJSON, timeframe12h)
			if err != nil {
				logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
				return
			}
		}

		if nu.Update48h {
			logger.Info.Println("Writing 48h Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBucketsJSON, timeframe48h)
			if err != nil {
				logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
				return
			}
		}

		if nu.Update7d {
			logger.Info.Println("Writing 7d Historical Mempool data.")
			err = database.WriteHistoricalMempoolData(countInBucketsJSON, timeframe7d)
			if err != nil {
				logger.Error.Printf("Failed to write Current Mempool to database: %s", err.Error())
				return
			}
		}
		logger.Info.Println("Success writing Historical Mempool to database.")
	}

}

/* generateCurrentMempoolStats()
This function generates the _Current Mempool_ data. Which is:
	- The size of the transactions in the mempool `mempoolSizeInByte`
	- A map mapping the transaction count to the feerate (as a whole
		number). Named `feerateMap`.
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
			megabyteBucket = cMEGABYTE - entry.size               // start a new megabyte bucket minus the current entry.size
			megabyteMarkers = append(megabyteMarkers, memlistPos) // append current position to the megabyteMarkers list
			memlistPos--
		}
	}

	return feerateMap, mempoolSizeInByte, megabyteMarkers
}

var feerateBuckets = [40]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 18, 22, 27, 33, 41, 50, 62, 76, 93, 114, 140, 172, 212, 261, 321, 395, 486, 598, 736, 905, 1113, 1369, 1684, 2071, 2547, 3133, 3854, 3855}

// generates a list of counts of transactions representing the count in a feerate bucket
func generateHistoricalMempoolStats(mempool map[string]types.PartialTransaction) (countInBuckets []int) {

	countInBuckets = make([]int, len(feerateBuckets), len(feerateBuckets))

	for _, tx := range mempool {
		feerate := tx.Fee * cSATOSHIPERBITCOIN / float64(tx.Size)
		bucketIndex := findBucketForFeerate(feerate)
		countInBuckets[bucketIndex]++
	}

	return
}

// finds the bucket index for a given feerate in feerateBuckets. the last bucket is a catch all larger-equal
func findBucketForFeerate(feerate float64) int {
	i := sort.Search(len(feerateBuckets), func(i int) bool { return feerateBuckets[i] >= int(feerate) })
	if i < len(feerateBuckets) && feerateBuckets[i] >= int(feerate) {
		return i
	}
	return len(feerateBuckets) - 1
}
