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
			topic:     RawTransaction,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (*rawTxPub) onNewBlock(_ *block.Block) {
	// TODO implement me
	panic("implement me")
}
