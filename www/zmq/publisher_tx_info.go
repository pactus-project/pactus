package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type txInfoPub struct {
	basePub
}

func newTxInfoPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &txInfoPub{
		basePub: basePub{
			topic:     TopicTransactionInfo,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (t *txInfoPub) onNewBlock(blk *block.Block) {
	for _, txn := range blk.Transactions() {
		rawMsg := t.makeTopicMsg(txn.ID().Bytes(), blk.Height())
		message := zmq4.NewMsg(rawMsg)

		if err := t.zmqSocket.Send(message); err != nil {
			t.logger.Error("zmq publish message error", "err", err, "publisher", t.TopicName())

			continue
		}

		t.logger.Debug("ZMQ published the message successfully",
			"publisher", t.TopicName(),
			"block_height", blk.Height(),
			"tx_hash", txn.ID().String(),
		)

		t.seqNo++
	}
}
