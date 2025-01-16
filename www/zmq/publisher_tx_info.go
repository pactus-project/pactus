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

func (*txInfoPub) onNewBlock(_ *block.Block) {
	// TODO implement me
	panic("implement me")
}
