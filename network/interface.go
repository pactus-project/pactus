package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
)

type CallbackFn func([]byte, peer.ID)

type Network interface {
	Start() error
	Stop()
	PublishMessage(msg *message.Message) error
	JoinTopics(CallbackFn) error
	JoinDownloadTopic() error
	LeaveDownloadTopic()
	CloseConnection(pid peer.ID)
	SelfID() peer.ID
}
