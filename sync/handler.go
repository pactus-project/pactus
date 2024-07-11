package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type messageHandler interface {
	ParseMessage(message.Message, peer.ID)
	PrepareBundle(message.Message) *bundle.Bundle
}
