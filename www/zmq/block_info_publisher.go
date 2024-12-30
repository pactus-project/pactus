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
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (*blockInfoPub) onNewBlock(_ *block.Block) {
	// TODO implement me
	panic("implement me")
}
