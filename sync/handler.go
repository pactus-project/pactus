package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type messageHandler interface {
	ParseMessage(message.Message, peer.ID) error
	PrepareBundle(message.Message) *bundle.Bundle
}
