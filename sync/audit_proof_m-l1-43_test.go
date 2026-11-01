package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// M-L1-43: a requester's BlocksResponse chunks and the final response ("NoMoreBlocks"/
// "Synced") are sent by the responder via independent calls to handler.respond, each of
// which goes through synchronizer.sendTo -> network.SendTo. In the real network
// implementation, network.SendTo (network/network.go:379-386) spawns a brand-new
// goroutine AND a brand-new libp2p stream per call, so nothing guarantees the final
// message reaches the requester after its preceding chunks.
//
// This test plays the requester's side of that race: it feeds the final
// "NoMoreBlocks" response before the "MoreBlocks" chunk it is supposed to follow, and
// shows that handler_blocks_response.go's updateSession (handler_blocks_response.go:69-72)
// closes the session unconditionally on the final response, with no check that the
// requested range was actually received -- exactly the condition the finding says lets
// the node believe it is done and move on (or re-download the same range from another
// peer) despite missing data.
func TestAuditBlocksResponseFinalMessageBeforeChunkClosesSessionEarly(t *testing.T) {
	td := setup(t, nil)
	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

	lastHeight := td.state.LastBlockHeight()
	from := lastHeight + 1
	count := uint32(1)

	sid := td.sync.peerSet.OpenSession(pid, from, count)
	require.True(t, td.sync.peerSet.HasOpenSession(pid))
	require.Nil(t, td.sync.cache.GetBlock(from), "precondition: the requested block hasn't arrived yet")

	// Simulate the final message winning the race against its own chunk.
	finalMsg := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
		message.ResponseCodeNoMoreBlocks.String(), sid, 0, nil, nil)
	td.receivingNewMessage(td.sync, finalMsg, pid)

	// BUG: updateSession closes (and updateBlockchain then discards) the original
	// session unconditionally on the final response, even though the chunk carrying
	// the actual block for `from` was never processed. The sync module now behaves as
	// if that height were already handled.
	for _, ssn := range td.sync.peerSet.Sessions() {
		assert.NotEqual(t, sid, ssn.SessionID,
			"the original session for the still-missing block was discarded by the lone final message")
	}
	assert.Nil(t, td.sync.cache.GetBlock(from),
		"the requested block is still missing even though its session was already closed")

	// Concrete consequence: believing the download for `from` is done, updateBlockchain
	// immediately opens a *new* session that re-requests the very same height `from`
	// from a peer -- a duplicate download of a range already (wrongly) marked complete.
	redownloaded := false
	for _, ssn := range td.sync.peerSet.Sessions() {
		if ssn.From == from {
			redownloaded = true
		}
	}
	assert.True(t, redownloaded,
		"expected a new session re-requesting the same height that the premature final message had already closed")

	// The original chunk now arrives late, for the same session ID. M-H-12's fix
	// (handler_blocks_response.go now gates every response on
	// peerSet.HasOpenSessionFor(pid, sid)) rejects it as a side effect: the session
	// is no longer Open, so the late chunk is correctly ignored rather than accepted
	// after the node had already (wrongly) moved on. The core M-L1-43 finding above
	// -- the premature close and duplicate re-download from a lone final message --
	// is unchanged.
	blk, cert := td.GenerateTestBlock(from)
	data, _ := blk.Bytes()
	chunkMsg := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
		message.ResponseCodeMoreBlocks.String(), sid, from, [][]byte{data}, cert)
	td.receivingNewMessage(td.sync, chunkMsg, pid)
	assert.Nil(t, td.sync.cache.GetBlock(from),
		"a chunk for a session that is no longer open must not be accepted either")
}
