package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type rawBlockPub struct {
	basePub
}

func newRawBlockPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &rawBlockPub{
		basePub: basePub{
			topic:     TopicRawBlock,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (*rawBlockPub) onNewBlock(_ *block.Block) {
	// TODO implement me
	panic("implement me")
}
