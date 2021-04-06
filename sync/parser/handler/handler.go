package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type Handler interface {
	ParsPayload(payload.Payload, peer.ID) error
}
