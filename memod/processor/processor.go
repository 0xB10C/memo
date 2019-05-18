package processor

/* processes */

import (
	"sort"
	"time"

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

	// start _current mempool_ stat generation in a goroutine
	go currentMempool(mempool)

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
