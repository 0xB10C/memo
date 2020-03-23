package processor

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/types"
	"github.com/btcsuite/btcd/wire"

	"github.com/0xb10c/rawtx"
)

// This file includes the functions to handle zmq events
// and their subfunctions.

/* ----- RAW BLOCK handling ----- */

// HandleRawBlock handles a raw incoming zmq block
func HandleRawBlock(payload string, pool *database.RedisPool) {
	block, err := deserializeRawBlock(payload)
	if err != nil {
		logger.Error.Printf("Error handling raw block: %v", err)
	}

	height := getBlockHeightFromCoinbase(block.Transactions[0])
	numTx := len(block.Transactions)
	sizeWithWitness := block.SerializeSize()
	sizeWithoutWitness := block.SerializeSizeStripped()
	weight := sizeWithWitness + sizeWithoutWitness*3

	err = pool.WriteNewBlockData(height, numTx, sizeWithWitness, weight)
	if err != nil {
		logger.Error.Printf("Error writing block to database: %v", err)
	}

	// According to https://rusnak.io/longest-txid-prefix-collision-in-bitcoin/ (updated 2019)
	// the longest TXID collision is 15 hex bytes long. Choosing a short txid of 16 here should
	// be way more than enough. Compared to the blog post, TXIDs are here only compared to the 30k
	// transactions displayed in the Bitcoin Transaction Monitor.
	const shortTXIDLength int = 16
	shortTXIDs := make([]string, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		shortTXIDs = append(shortTXIDs, tx.TxHash().String()[0:shortTXIDLength])
	}

	err = pool.WriteNewBlockEntry(height, shortTXIDs)
	if err != nil {
		logger.Error.Printf("Error writing block entries to database: %v", err)
	}

	logger.Info.Printf("Success writing new block %d with %d transactions, size %d, weight %d", height, numTx, sizeWithWitness, weight)
}

func deserializeRawBlock(rawBlock string) (block *wire.MsgBlock, err error) {
	defer logger.TrackTime(time.Now(), "deserializeRawBlock()")
	blockHeader := wire.BlockHeader{}
	block = wire.NewMsgBlock(&blockHeader)
	r := strings.NewReader(rawBlock)
	err = block.Deserialize(r)
	if err != nil {
		return
	}
	return
}

func getBlockHeightFromCoinbase(coinbase *wire.MsgTx) (height int) {
	defer logger.TrackTime(time.Now(), "getBlockHeightFromCoinbase()")
	// To get the block height we look into the coinbase transaction
	// (the first transaction in a block). The scriptsig of the coin-
	// base transaction starts with the height. **The first byte sets
	// height length**. (only true for blocks with BIP34, but this can
	// be ignored here, since we probably don't work with older blocks)
	heightLength := coinbase.TxIn[0].SignatureScript[0]

	// we get the bytes from pos 1 till pos heightLength + 1 since the
	// second parameter is exclusive in Go
	heightLittleEndian := coinbase.TxIn[0].SignatureScript[1 : heightLength+1]

	// since we want the block height in a int32 (4 byte) and big
	// endian we first add padding (at the end, since it's little
	// endian) and then convert it to big endian.
	for i := heightLength; i < 4; i++ {
		heightLittleEndian = append(heightLittleEndian, 0x0)
	}

	height = int(binary.LittleEndian.Uint32(heightLittleEndian))
	return
}

/* ----- HASH BLOCK handling ----- */

// HandleHashBlock handles a incoming zmq block hash
func HandleHashBlock(payload string) {
	logger.Warning.Println("HandleHashBlock() not Implemented")
}

/* ----- RAW TX handling ----- */

// HandleRawTx handles a incoming zmq raw tx
func HandleRawTx(payload string) {
	logger.Warning.Println("HandleRawTx() not Implemented")
}

// HandleRawTxWithSizeAndFee handles the special rawtx2 zmq message
// which contains a tx, it's size and it's fee as 64 bit ints
func HandleRawTxWithSizeAndFee(payload string, pool *database.RedisPool) {
	payloadLength := len(payload)

	tx, err := rawtx.DeserializeRawTxBytes([]byte(payload[0 : payloadLength-16]))
	if err != nil {
		logger.Error.Printf("Error handling raw tx: %v", err)
	}

	sizeBytes := []byte(payload[payloadLength-16 : payloadLength-8])
	feeBytes := []byte(payload[payloadLength-8 : payloadLength])

	sizeInByte := int64(binary.LittleEndian.Uint64(sizeBytes))
	feeInSat := int64(binary.LittleEndian.Uint64(feeBytes))

	me := types.MempoolEntry{}
	me.EntryTime = time.Now().Unix()
	me.TxID = tx.HashString
	me.Fee = feeInSat
	me.Size = sizeInByte
	me.Version = tx.Version
	me.InputCount = tx.GetNumInputs()
	me.OutputCount = tx.GetNumOutputs()
	me.Locktime = tx.GetLocktime()
	me.OutputSum = tx.GetOutputSum()
	me.SpendsSegWit = tx.IsSpendingSegWit()
	me.SpendsMultisig = tx.IsSpendingMultisig()
	me.IsBIP69 = tx.IsBIP69Compliant()
	me.SignalsRBF = tx.IsExplicitlyRBFSignaling()

	me.Spends = make(map[string]int)
	me.PaysTo = make(map[string]int)

	if me.SpendsMultisig {
		me.Multisig = make(map[string]int)
	}

	for _, in := range tx.Inputs {
		inputType := in.GetType()
		me.Spends[inputType.String()]++
		switch inputType {
		case rawtx.InP2SH:
			isMultisig, m, n := in.GetP2SHRedeemScript().IsMultisigScript()
			if isMultisig {
				me.Multisig[strconv.Itoa(m)+"-of-"+strconv.Itoa(n)]++
			}
		case rawtx.InP2SH_P2WSH:
			isMultisig, m, n := in.GetNestedP2WSHRedeemScript().IsMultisigScript()
			if isMultisig {
				me.Multisig[strconv.Itoa(m)+"-of-"+strconv.Itoa(n)]++
			}
		case rawtx.InP2WSH:
			isMultisig, m, n := in.GetP2WSHRedeemScript().IsMultisigScript()
			if isMultisig {
				me.Multisig[strconv.Itoa(m)+"-of-"+strconv.Itoa(n)]++
			}
		case rawtx.InUNKNOWN:
			logger.Warning.Printf("Tx with the id %s categorized as Spends UNKNOWN.", me.TxID)
		}
	}

	for _, out := range tx.Outputs {
		outputType := out.GetType()
		me.PaysTo[outputType.String()]++

		switch outputType {
		case rawtx.OutOPRETURN:
			var isOPRETURN, data = out.GetOPReturnData()
			if isOPRETURN {
				me.OPReturnData = string(data.PushedData)
				me.OPReturnLength = len(data.PushedData)
			}
		case rawtx.OutUNKNOWN:
			logger.Warning.Printf("Tx with the id %s categorized as PaysTo UNKNOWN.", me.TxID)
		}
	}

	err = pool.WriteMempoolEntries(me)
	if err != nil {
		logger.Error.Printf("Error handling raw tx: %v", err)
	}
}

/* ----- HASH TX handling ----- */

// HandleHashTx handles incoming zmq tx hashes
func HandleHashTx(payload string) {
	logger.Warning.Println("HandleHashTx() not Implemented")
}
