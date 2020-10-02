package txpool

import (
	"context"
	"fmt"
	"reflect"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/network"
	"gitlab.com/zarb-chain/zarb-go/tx"
	"gitlab.com/zarb-chain/zarb-go/txpool/message"
)

const (
	TxPoolChannel = byte(0x41)
)

type synchronizer struct {
	ctx    context.Context
	config *config.Config
	pool   *TxPool
	topic  *pubsub.Topic
	sub    *pubsub.Subscription
	self   peer.ID
	logger *logger.Logger
}

func newSynchronizer(conf *config.Config, pool *TxPool, net *network.Network, logger *logger.Logger) (*synchronizer, error) {
	syncer := &synchronizer{
		ctx:    context.Background(),
		config: conf,
		pool:   pool,
		logger: logger,
	}

	topic, err := net.JoinTopic("tx_pool")
	if err != nil {
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	syncer.self = net.ID()
	syncer.topic = topic
	syncer.sub = sub

	return syncer, nil
}

func (syncer *synchronizer) Start() error {
	go syncer.readLoop()
	return nil
}

func (syncer *synchronizer) Stop() error {
	syncer.ctx.Done()
	syncer.sub.Cancel()
	syncer.topic.Close()
	return nil
}

func (syncer *synchronizer) readLoop() {
	for {
		m, err := syncer.sub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}
		// only forward messages delivered by others
		if m.ReceivedFrom == syncer.self {
			continue
		}

		msg := new(message.Message)
		err = msg.UnmarshalCBOR(m.Data)
		if err != nil {
			syncer.logger.Error("Error decoding message", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
			continue
		}
		syncer.logger.Trace("Received a message", "from", m.ReceivedFrom.Pretty(), "message", msg)

		if err = msg.SanityCheck(); err != nil {
			syncer.logger.Error("Peer sent us invalid msg", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
			continue
		}

		switch msg.PayloadType() {
		case message.PayloadTypeTx:
			pld := msg.Payload.(*message.TxPayload)

			trx := new(tx.Tx)
			err := trx.Decode(pld.TxData)
			if err != nil {
				syncer.logger.Error("Received invalid transaction", "from", m.ReceivedFrom.Pretty(), "err", err)
				return
			}
			syncer.pool.AppendTx(trx)

		case message.PayloadTypeRequest:
			pld := msg.Payload.(*message.RequestPayload)

			if syncer.pool.HasTx(pld.Hash) {
				trx, _ := syncer.pool.PendingTx(pld.Hash)
				msg := message.NewTxMessage(trx)
				syncer.publishMessage(msg)
			}

		default:
			// don't punish (leave room for soft upgrades)
			syncer.logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}
	}
}

func (syncer *synchronizer) BroadcastTx(trx *tx.Tx) {
	msg := message.NewTxMessage(trx)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) BroadcastRequestTx(hash crypto.Hash) {
	msg := message.NewRequestMessage(hash)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) publishMessage(msg *message.Message) {
	bs, _ := msg.MarshalCBOR()
	if err := syncer.topic.Publish(syncer.ctx, bs); err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Trace("Publishing new message", "message", msg)
	}
}
