package zmq

import (
	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/processor"

	"github.com/pebbe/zmq4"
)

const rawBlock string = "rawblock"
const hashBlock string = "hashblock"
const rawTx string = "rawtx"
const rawTx2 string = "rawtx2"
const hashTx string = "hashtx"

func SetupZMQ() {

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
		go processor.HandleRawBlock(payload)
	case hashBlock:
		go processor.HandleHashBlock(payload)
	case rawTx:
		go processor.HandleRawTx(payload)
	case rawTx2:
		go processor.HandleRawTxWithSizeAndFee(payload)
	case hashTx:
		go processor.HandleHashTx(payload)
	default:
		logger.Warning.Println("Unhandled ZMQ topic", topic)
	}
}
