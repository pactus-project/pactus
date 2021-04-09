package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var (
	tAliceConfig      *Config
	tBobConfig        *Config
	tAliceState       *state.MockState
	tBobState         *state.MockState
	tAliceConsensus   *consensus.MockConsensus
	tBobConsensus     *consensus.MockConsensus
	tAliceNet         *network.MockNetwork
	tBobNet           *network.MockNetwork
	tAliceSync        *synchronizer
	tBobSync          *synchronizer
	tAliceBroadcastCh chan payload.Payload
	tBobBroadcastCh   chan payload.Payload
	tAlicePeerID      peer.ID
	tBobPeerID        peer.ID
)

type OverrideFingerprint struct {
	sync Synchronizer
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.sync.Fingerprint()
}

func init() {
	logger.InitLogger(logger.TestConfig())
	tAliceConfig = TestConfig()
	tBobConfig = TestConfig()

	tAliceConfig.Moniker = "alice"
	tBobConfig.Moniker = "bob"
}

func setup(t *testing.T) {
	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	aliceSigner := crypto.NewSigner(priv1)
	bobSigner := crypto.NewSigner(priv2)

	committee, _ := committee.GenerateTestCommittee()
	tAlicePeerID = util.RandomPeerID()
	tBobPeerID = util.RandomPeerID()
	tAliceState = state.MockingState(committee)
	tBobState = state.MockingState(committee)
	tAliceConsensus = consensus.MockingConsensus(tAliceState)
	tBobConsensus = consensus.MockingConsensus(tBobState)
	tAliceBroadcastCh = make(chan payload.Payload, 1000)
	tBobBroadcastCh = make(chan payload.Payload, 1000)
	tAliceNet = network.MockingNetwork(tAlicePeerID)
	tBobNet = network.MockingNetwork(tBobPeerID)

	tBobState.GenHash = tAliceState.GenHash

	// Apply 20 blocks for both Alice and Bob
	lastBlockHash := crypto.Hash{}
	for i := 0; i < 21; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		lastBlockHash = b.Hash()

		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCertificate = c

		tBobState.AddBlock(i+1, b, trxs)
		tBobState.LastBlockCertificate = c
	}

	tAliceSync = &synchronizer{ctx: context.Background()}
	aliceSync, err := NewSynchronizer(tAliceConfig,
		aliceSigner,
		tAliceState,
		tAliceConsensus,
		tAliceNet,
		tAliceBroadcastCh,
	)
	assert.NoError(t, err)
	tAliceSync = aliceSync.(*synchronizer)
	tAliceSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "alice: ", sync: tAliceSync})

	tBobSync = &synchronizer{ctx: context.Background()}
	bobSync, err := NewSynchronizer(tBobConfig,
		bobSigner,
		tBobState,
		tBobConsensus,
		tBobNet,
		tBobBroadcastCh,
	)
	assert.NoError(t, err)
	tBobSync = bobSync.(*synchronizer)
	tBobSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "bob: ", sync: tBobSync})

	tAliceNet.OtherNet = tBobNet
	tBobNet.OtherNet = tAliceNet

	assert.NoError(t, tAliceSync.Start())
	assert.NoError(t, tBobSync.Start())

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeSalam)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeSalam)

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeAleyk)

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
}

func shouldPublishPayloadWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.PayloadType) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("ShouldPublishPayloadWithThisType %v: Timeout", payloadType))
			return
		case msg := <-net.BroadcastCh:
			logger.Info("shouldPublishPayloadWithThisType", "msg", msg, "type", payloadType.String())
			net.SendMessageToOthePeer(msg)
			logger.Info("Nessage sent to other peer", "msg", msg)

			if msg.Payload.Type() == payloadType {
				return
			}
		}
	}
}

func shouldPublishPayloadWithThisTypeAndResponseCode(t *testing.T, net *network.MockNetwork, payloadType payload.PayloadType, code payload.ResponseCode) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("ShouldPublishPayloadWithThisType %v: Timeout", payloadType))
			return
		case msg := <-net.BroadcastCh:
			logger.Info("shouldPublishPayloadWithThisType", "msg", msg, "type", payloadType.String())
			net.SendMessageToOthePeer(msg)
			logger.Info("Nessage sent to other peer", "msg", msg)

			if msg.Payload.Type() == payloadType {
				switch payloadType {
				case payload.PayloadTypeAleyk:
					pld := msg.Payload.(*payload.AleykPayload)
					assert.Equal(t, pld.ResponseCode, code)

				case payload.PayloadTypeDownloadResponse:
					pld := msg.Payload.(*payload.DownloadResponsePayload)
					assert.Equal(t, pld.ResponseCode, code)

				case payload.PayloadTypeLatestBlocksResponse:
					pld := msg.Payload.(*payload.LatestBlocksResponsePayload)
					assert.Equal(t, pld.ResponseCode, code)
				}
				return
			}
		}
	}
}

func shouldNotPublishPayloadWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.PayloadType) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-net.BroadcastCh:
			logger.Info("shouldNotPublishPayloadWithThisType", "msg", msg, "type", payloadType.String())
			net.SendMessageToOthePeer(msg)
			logger.Info("Nessage sent to other peer", "msg", msg)

			assert.NotEqual(t, msg.Payload.Type(), payloadType)
		}
	}
}

func addMoreBlocksForBob(t *testing.T, count int) {
	lastBlockHash := tBobState.LastBlockHash()
	for i := 0; i < count; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		lastBlockHash = b.Hash()

		tBobState.AddBlock(tBobState.LastBlockHeight()+1, b, trxs)
		tBobState.LastBlockCertificate = c
	}
}

func addMoreBlocksForBobAndAnnounceLastBlock(t *testing.T, count int) {
	addMoreBlocksForBob(t, count)

	pld := payload.NewBlockAnnouncePayload(
		tBobState.LastBlockHeight(),
		tBobState.Store.Blocks[tBobState.LastBlockHeight()],
		tBobState.LastBlockCertificate)

	tBobBroadcastCh <- pld
}

func disableHeartbeat(t *testing.T) {
	tAliceSync.heartBeatTicker.Stop()
	tBobSync.heartBeatTicker.Stop()
}

func joinAliceToTheSet(t *testing.T) {
	val := validator.NewValidator(tAliceSync.signer.PublicKey(), 4, tAliceState.LastBlockHeight())
	val.UpdateLastJoinedHeight(tAliceState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func joinBobToTheSet(t *testing.T) {
	val := validator.NewValidator(tBobSync.signer.PublicKey(), 5, tBobState.LastBlockHeight())
	val.UpdateLastJoinedHeight(tBobState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func TestAccessors(t *testing.T) {
	setup(t)

	assert.Equal(t, tAliceSync.SelfID(), tAlicePeerID)
	assert.Equal(t, len(tAliceSync.Peers()), 1)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tAliceSync.Stop()
	tBobSync.Stop()
}
