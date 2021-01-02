package network_api

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/message"
)

type NetworkAPI interface {
	Start() error
	Stop()
	PublishMessage(msg *message.Message) error
	JoinDownloadTopic() error
	LeaveDownloadTopic()
	SelfID() peer.ID
}
