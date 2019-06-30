package zmq

import (
	"encoding/binary"
	"strings"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/logger"

	"github.com/btcsuite/btcd/wire"
	"github.com/pebbe/zmq4"
)

func Start() {
	setupZMQ()
}

const rawBlock string = "rawblock"
const hashBlock string = "hashblock"
const rawTx string = "rawtx"
const rawTx2 string = "rawtx2"
const hashTx string = "hashtx"

func setupZMQ() {

	zmqHost := config.GetString("zmq.host")
	zmqPort := config.GetString("zmq.port")
	connectionString := "tcp://" + zmqHost + ":" + zmqPort

	subscriber, _ := zmq4.NewSocket(zmq4.SUB)
	subscriber.Connect(connectionString)

	if config.GetBool("zmq.subscribeTo.rawTx") {
		subscriber.SetSubscribe(rawTx)
	}
	if config.GetBool("zmq.subscribeTo.rawTx2") {
		subscriber.SetSubscribe(rawTx2)
	}
	if config.GetBool("zmq.subscribeTo.hashTx") {
		subscriber.SetSubscribe(hashTx)
	}
	if config.GetBool("zmq.subscribeTo.rawBlock") {
		subscriber.SetSubscribe(rawBlock)
	}
	if config.GetBool("zmq.subscribeTo.hashBlock") {
		subscriber.SetSubscribe(hashBlock)
	}

	defer subscriber.Close() // cancel subscribe

	loopZMQ(subscriber)
}

func loopZMQ(subscriber *zmq4.Socket) {
	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			logger.Error.Println(err)
		}
		handleZMQMessage(msg)
	}
}

func handleZMQMessage(zmqMessage []string) {
	topic := zmqMessage[0]
	payload := zmqMessage[1]

	switch topic {
	case rawBlock:
		go handleRawBlock(payload)
	case hashBlock:
		go handleHashBlock(payload)
	case rawTx:
		//go handleRawTx(payload)
	case rawTx2:
		go handleRawTxWithSizeAndFee(payload)
	case hashTx:
		go handleHashTx(payload)
	default:
		logger.Warning.Println("Unhandled ZMQ topic", topic)
	}
}

func handleRawBlock(payload string) {
	block, err := deserializeRawBlock(payload)
	if err != nil {
		logger.Error.Printf("Error handling raw block: %v", err)
	}

	height := getBlockHeightFromCoinbase(block.Transactions[0])
	numTx := len(block.Transactions)
	sizeWithWitness := block.SerializeSize()
	sizeWithoutWitness := block.SerializeSizeStripped()
	weight := sizeWithWitness + sizeWithoutWitness*3

	err = database.WriteNewBlockData(height, numTx, sizeWithWitness, weight)
	if err != nil {
		logger.Error.Printf("Error writing block to database: %v", err)
	}
	logger.Info.Printf("Success writing new block %d with %d transactions, size %d, weight %d", height, numTx, sizeWithWitness, weight)
}

func handleHashBlock(payload string) {
	//logger.Warning.Println("handleHashBlock() not Implemented")
}

func handleRawTxWithSizeAndFee(payload string) {
	payloadLength := len(payload)

	tx, err := deserializeRawTx(payload[0 : payloadLength-16])
	if err != nil {
		logger.Error.Printf("Error handling raw tx: %v", err)
	}

	sizeBytes := []byte(payload[payloadLength-16 : payloadLength-8])
	feeBytes := []byte(payload[payloadLength-8 : payloadLength])

	sizeInByte := int64(binary.LittleEndian.Uint64(sizeBytes))
	feeInSat := int64(binary.LittleEndian.Uint64(feeBytes))

	feerate := float64(feeInSat) / float64(sizeInByte)

	logger.Warning.Println(tx.TxHash(), feerate, tx.HasWitness())
}

func handleHashTx(payload string) {
	//logger.Warning.Println("handleHashTx() not Implemented")
}

func deserializeRawTx(rawTx string) (tx *wire.MsgTx, err error) {
	//defer logger.TrackTime(time.Now(), "deserializeRawTx()")
	tx = wire.NewMsgTx(1)
	r := strings.NewReader(rawTx)
	err = tx.Deserialize(r)
	if err != nil {
		return
	}
	return
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
	// base transaction starts with the height. The first byte sets
	// height length. (only true for blocks with BIP34, but this can
	// be ignored here, since we probably don't work with older blocks)
	heightLength := coinbase.TxIn[0].SignatureScript[0]

	// we get the bytes from pos 1 till pos heightLength + 1 since the
	// second parameter is exclusive in Go
	heightLE := coinbase.TxIn[0].SignatureScript[1 : heightLength+1]

	// since we want the block height in a int32 (4 byte) and big
	// endian we first add padding (at the end, since it's little
	// endian) and then convert it to big endian
	for i := heightLength; i < 4; i++ {
		heightLE = append(heightLE, 0x0)
	}

	height = int(binary.LittleEndian.Uint32(heightLE))
	return
}
