package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type messageHandler interface {
	ParsMessage(message.Message, peer.ID) error
	PrepareBundle(message.Message) *bundle.Bundle
}
