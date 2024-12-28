package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type Publisher interface {
	Address() string
	TopicName() string

	onNewBlock(blk *block.Block)
}

type basePub struct {
	topic     Topic
	zmqSocket zmq4.Socket
	logger    *logger.SubLogger
}

func (b *basePub) Address() string {
	return b.zmqSocket.Addr().String()
}

func (b *basePub) TopicName() string {
	return b.topic.String()
}
