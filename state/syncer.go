package state

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"gitlab.com/zarb-chain/zarb-go/block"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/network"
	"gitlab.com/zarb-chain/zarb-go/state/message"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/utils"
)

type synchronizer struct {
	ctx       context.Context
	config    *config.Config
	store     *store.Store
	state     *State
	topic     *pubsub.Topic
	sub       *pubsub.Subscription
	self      peer.ID
	blockPool map[int]*block.Block
	logger    *logger.Logger
}

func newSynchronizer(conf *config.Config, state *State, store *store.Store, net *network.Network, logger *logger.Logger) (*synchronizer, error) {
	syncer := &synchronizer{
		ctx:       context.Background(),
		config:    conf,
		state:     state,
		store:     store,
		blockPool: make(map[int]*block.Block),
		logger:    logger,
	}
	topic, err := net.JoinTopic("state")
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

	// Let's other peers know our height
	syncer.BroadcastStateInfo()

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

		ourHeight, _ := syncer.state.LastBlockInfo()
		switch msg.PayloadType() {
		case message.PayloadTypeStateInfo:
			pld := msg.Payload.(*message.StateInfoPayload)
			switch h := pld.Height; {
			case h == ourHeight:
			case h == ourHeight+1:
			case h == ourHeight-1:
				{
					// Do nothing
					// Consensus lagging?
				}

			case h > ourHeight+1:
				{

				}

			case h < ourHeight-1:
				{
					// Help peer to catch up
					from := h
					end := utils.Min(h+10, ourHeight)
					blocks := make([]block.Block, end-from)
					for h := from; h <= end; h++ {
						b, err := syncer.store.BlockByHeight(h)
						if err != nil {
							syncer.logger.Error("An error occured while retriveng a block", "err", err, "height", h)
							return
						}
						blocks[h-from] = *b
					}

					syncer.BroadcastBlocks(from, blocks)
				}
			}

		case message.PayloadTypeBlocks:
			pld := msg.Payload.(*message.BlocksPayload)
			height := pld.From
			for _, b := range pld.Blocks {
				if height > ourHeight {
					bp, has := syncer.blockPool[height]
					if has {
						if !bp.Hash().EqualsTo(b.Hash()) {
							syncer.logger.Error("We have recieved twoblock from same height but different hash", "from", m.ReceivedFrom.Pretty(), "height", height)
						}
					} else {
						syncer.blockPool[height] = &b
					}
				}
				height++
			}

		default:
			syncer.logger.Error("Unknown message type", "msg", msg)
		}
	}
}

func (syncer *synchronizer) BroadcastStateInfo() {
	height, hash := syncer.state.LastBlockInfo()
	msg := message.NewStateInfoMessage(height, hash)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) BroadcastBlocks(from int, blocks []block.Block) {
	msg := message.NewBlocksMessage(from, blocks)
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
