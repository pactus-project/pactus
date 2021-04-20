package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type payloadHandler interface {
	ParsPayload(payload.Payload, peer.ID) error
	PrepareMessage(payload.Payload) *message.Message
}
