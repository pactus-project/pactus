package sync

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

type GeneralTopic struct {
	config    *Config
	selfID    peer.ID
	publicKey crypto.PublicKey
	peerSet   *peerset.PeerSet
	publishFn PublishMessageFn
	state     state.StateReader
	logger    *logger.Logger
}

func NewGeneralTopic(
	conf *Config,
	selfID peer.ID,
	publicKey crypto.PublicKey,
	peerSet *peerset.PeerSet,
	state state.StateReader,
	logger *logger.Logger,
	publishFn PublishMessageFn) *GeneralTopic {
	return &GeneralTopic{
		config:    conf,
		selfID:    selfID,
		publicKey: publicKey,
		peerSet:   peerSet,
		state:     state,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (gt *GeneralTopic) BroadcastSalam() {
	flags := 0
	if gt.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	msg := message.NewSalamMessage(
		gt.config.Moniker,
		gt.publicKey,
		gt.selfID,
		gt.state.GenesisHash(),
		gt.state.LastBlockHeight(),
		flags)

	gt.publishFn(msg)
}

func (gt *GeneralTopic) BroadcastAleyk(resStatus int, resMsg string) {
	flags := 0
	if gt.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	msg := message.NewAleykMessage(
		gt.config.Moniker,
		gt.publicKey,
		gt.selfID,
		gt.state.LastBlockHeight(),
		flags,
		resStatus,
		resMsg)

	gt.publishFn(msg)
}

func (gt *GeneralTopic) ProcessSalamPayload(pld *payload.SalamPayload) {
	gt.logger.Trace("Process salam payload", "pld", pld)

	if !pld.GenesisHash.EqualsTo(gt.state.GenesisHash()) {
		gt.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash)
		// Reply salam
		gt.BroadcastAleyk(payload.SalamResponseCodeRejected, "Invalid genesis hash")
		return
	}

	p := gt.peerSet.MustGetPeer(pld.PeerID)
	p.UpdateMoniker(pld.Moniker)
	p.UpdateHeight(pld.Height)
	p.UpdateVersion(pld.NodeVersion)
	p.UpdatePublicKey(pld.PublicKey)

	// Reply salam
	gt.BroadcastAleyk(payload.SalamResponseCodeOK, "Welcome!")
}

func (gt *GeneralTopic) ProcessAleykPayload(pld *payload.AleykPayload) {
	gt.logger.Trace("Process Aleyk payload", "pld", pld)
}
