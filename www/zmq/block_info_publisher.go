package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type blockInfoPub struct {
	seqNo uint32
	basePub
}

func newBlockInfoPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &blockInfoPub{
		basePub: basePub{
			topic:     BlockInfo,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (b *blockInfoPub) onNewBlock(blk *block.Block) {
	seq := b.seqNo + 1

	rawMsg := makeTopicMsg(
		b.topic,
		blk.Header().ProposerAddress(),
		blk.Header().UnixTime(),
		uint16(len(blk.Transactions())),
		blk.Height(),
		seq,
	)

	message := zmq4.NewMsg(rawMsg)

	if err := b.zmqSocket.Send(message); err != nil {
		b.logger.Error("zmq publish message error", "err", err, "publisher", b.TopicName())
	}

	b.logger.Debug("zmq published message success",
		"publisher", b.TopicName(),
		"block_height", blk.Height(),
	)

	b.seqNo = seq
}
