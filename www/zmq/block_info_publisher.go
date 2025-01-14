package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type blockInfoPub struct {
	basePub
}

func newBlockInfoPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &blockInfoPub{
		basePub: basePub{
			topic:     BlockInfo,
			seqNo:     0,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (b *blockInfoPub) onNewBlock(blk *block.Block) {
	rawMsg := b.makeTopicMsg(
		blk.Header().ProposerAddress(),
		blk.Header().UnixTime(),
		uint16(len(blk.Transactions())),
		blk.Height(),
	)

	message := zmq4.NewMsg(rawMsg)

	if err := b.zmqSocket.Send(message); err != nil {
		b.logger.Error("zmq publish message error", "err", err, "publisher", b.TopicName())
	}

	b.logger.Debug("zmq published message success",
		"publisher", b.TopicName(),
		"block_height", blk.Height())

	b.seqNo++
}
