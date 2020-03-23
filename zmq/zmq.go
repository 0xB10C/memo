package zmq

import (
	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/processor"

	"github.com/pebbe/zmq4"
)

const rawBlock string = "rawblock"
const hashBlock string = "hashblock"
const rawTx string = "rawtx"
const rawTx2 string = "rawtx2"
const hashTx string = "hashtx"

func SetupZMQ(pool *database.RedisPool) {

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

	loopZMQ(subscriber, pool)
}

func loopZMQ(subscriber *zmq4.Socket, pool *database.RedisPool) {
	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			logger.Error.Println(err)
		}
		handleZMQMessage(msg, pool)
	}
}

func handleZMQMessage(zmqMessage []string, pool *database.RedisPool) {
	topic := zmqMessage[0]
	payload := zmqMessage[1]

	switch topic {
	case rawBlock:
		go processor.HandleRawBlock(payload, pool)
	case hashBlock:
		go processor.HandleHashBlock(payload)
	case rawTx:
		go processor.HandleRawTx(payload)
	case rawTx2:
		go processor.HandleRawTxWithSizeAndFee(payload, pool)
	case hashTx:
		go processor.HandleHashTx(payload)
	default:
		logger.Warning.Println("Unhandled ZMQ topic", topic)
	}
}
