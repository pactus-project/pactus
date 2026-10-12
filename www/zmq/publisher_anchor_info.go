package zmq

import (
	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type anchorInfoPub struct {
	basePub
}

func newAnchorInfoPub(socket zmq4.Socket, logger *logger.SubLogger) Publisher {
	return &anchorInfoPub{
		basePub: basePub{
			topic:     TopicAnchorInfo,
			zmqSocket: socket,
			logger:    logger,
		},
	}
}

func (p *anchorInfoPub) onNewBlock(blk *block.Block) {
	for _, txn := range blk.Transactions() {
		pld, ok := txn.Payload().(*payload.AnchorPayload)
		if !ok {
			continue
		}

		hashLen := byte(0)
		var root []byte
		if pld.Action == payload.AnchorActionSet {
			hashLen = byte(len(pld.RootHash))
			root = pld.RootHash
		}

		parts := []any{pld.From, []byte{pld.Action}, uint32(blk.Height()), []byte{hashLen}}
		if hashLen > 0 {
			parts = append(parts, root)
		}

		rawMsg := p.makeTopicMsg(parts...)
		if err := p.zmqSocket.Send(zmq4.NewMsg(rawMsg)); err != nil {
			p.logger.Error("zmq publish message error", "err", err, "publisher", p.TopicName())

			continue
		}

		p.seqNo++
	}
}
