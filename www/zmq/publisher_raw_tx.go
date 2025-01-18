package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type rawTxPub struct {
	basePub
}

func newRawTxPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &rawTxPub{
		basePub: basePub{
			topic:     TopicRawTransaction,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (r *rawTxPub) onNewBlock(blk *block.Block) {
	for _, tx := range blk.Transactions() {
		buf, err := tx.Bytes()
		if err != nil {
			r.logger.Error("failed to serializing raw tx", "err", err, "topic", r.TopicName())

			return
		}

		rawMsg := r.makeTopicMsg(buf, blk.Height())
		message := zmq4.NewMsg(rawMsg)

		if err := r.zmqSocket.Send(message); err != nil {
			r.logger.Error("zmq publish message error", "err", err, "publisher", r.TopicName())

			return
		}

		r.logger.Debug("ZMQ published the message successfully",
			"publisher", r.TopicName(),
			"block_height", blk.Height())

		r.seqNo++
	}
}
