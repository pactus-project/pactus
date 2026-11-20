package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A BlocksResponse must only be able to act on a session that belongs to the
// peer who sent it. Today the handler never compares the sender with the
// session owner, so any peer can close or reject the session of another peer.
func TestBlocksResponseFromAnotherPeerClosesSession(t *testing.T) {
	codes := []message.ResponseCode{
		message.ResponseCodeNoMoreBlocks,
		message.ResponseCodeSynced,
		message.ResponseCodeRejected,
	}

	for _, code := range codes {
		t.Run(code.String(), func(t *testing.T) {
			td := setup(t, nil)

			owner := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
			intruder := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

			// Prevent updateBlockchain from wiping the sessions.
			otherOwner := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
			td.sync.peerSet.OpenSession(otherOwner, 1000, 10)

			sid := td.sync.peerSet.OpenSession(owner, td.state.LastBlockHeight()+1, 10)

			msg := message.NewBlocksResponseMessage(code, code.String(), sid, 0, nil, nil)
			td.receivingNewMessage(td.sync, msg, intruder)

			var found *session.Session
			for _, ssn := range td.sync.peerSet.Sessions() {
				if ssn.SessionID == sid {
					found = ssn
				}
			}
			require.NotNil(t, found)
			assert.Equal(t, session.Open, found.Status,
				"a response sent by another peer must not change the status of the session")
			assert.Equal(t, 0, td.sync.peerSet.GetPeer(owner).CompletedSessions,
				"the owner must not be credited with a session it did not complete")
		})
	}
}

// Blocks that do not belong to any session we opened must not be accepted in the
// cache: they are not verified until commit time, and a cached block at height H
// makes sendBlockRequestToRandomPeer skip downloading H.
func TestBlocksResponseUnsolicitedBlockIsCached(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.state.LastBlockHeight()
	blk, _ := td.GenerateTestBlock(lastHeight + 1)
	data, err := blk.Bytes()
	require.NoError(t, err)

	stranger := td.RandPeerID() // never connected, no session
	msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
		message.ResponseCodeMoreBlocks.String(), 987654, lastHeight+1, [][]byte{data}, nil)

	td.receivingNewMessage(td.sync, msg, stranger)

	assert.False(t, td.sync.cache.HasBlockInCache(lastHeight+1),
		"an unsolicited block from an unknown peer must not enter the cache")
}

// The cached (unverified) block makes the node skip the download of this height.
func TestBlocksResponseUnsolicitedBlockSkipsDownload(t *testing.T) {
	td := setup(t, nil)
	td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

	lastHeight := td.state.LastBlockHeight()
	blk, _ := td.GenerateTestBlock(lastHeight + 1)
	data, err := blk.Bytes()
	require.NoError(t, err)

	stranger := td.RandPeerID()
	msg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
		message.ResponseCodeMoreBlocks.String(), 987654, lastHeight+1, [][]byte{data}, nil)
	td.receivingNewMessage(td.sync, msg, stranger)

	td.sync.sendBlockRequestToRandomPeer(lastHeight+1, 10, false)

	bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	req := bdl.Message.(*message.BlocksRequestMessage)
	assert.Equal(t, lastHeight+1, req.From,
		"the node must still request the block that follows its height")
}
