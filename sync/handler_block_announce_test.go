package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerBlockAnnounceBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert, nil)
	td.sync.broadcast(msg)

	msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	assert.Equal(t, msg.Certificate.Height(), msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height())
}

func TestHandlerBlockAnnounceCacheBlock(t *testing.T) {
	td := setup(t, nil)

	height := td.RandHeight()
	blk1, cert1 := td.GenerateTestBlock(height)
	blk2, cert2 := td.GenerateTestBlock(height)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1, nil)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2, nil)

	td.receivingNewMessage(td.sync, msg1, td.RandPeerID())
	td.receivingNewMessage(td.sync, msg2, td.RandPeerID())

	cachedBlock := td.sync.cache.GetBlock(height)
	cachedCert := td.sync.cache.GetCertificate(height)
	assert.Equal(t, blk1, cachedBlock)
	assert.Equal(t, cert1, cachedCert)
}

// An announcement of a later height can't be verified: it is cached as an announcement,
// and the blocks that we ask for replace it.
func TestHandlerBlockAnnounceLaterHeightIsNotTrusted(t *testing.T) {
	td := setup(t, nil)

	height := td.state.LastBlockHeight() + 2
	announced, announcedCert := td.GenerateTestBlock(height)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(announced, announcedCert, nil), td.RandPeerID())

	require.Equal(t, announced, td.sync.cache.GetBlock(height))
	assert.True(t, td.sync.cache.HasAnnouncedBlock(height))

	downloaded, _ := td.GenerateTestBlock(height)
	td.sync.cache.AddBlock(downloaded)
	assert.Equal(t, downloaded, td.sync.cache.GetBlock(height))
}

// setupCommittee replaces the committee of the fake state with the new validators,
// and returns their keys.
func (td *testData) setupCommittee() []*bls.ValidatorKey {
	const size = 4

	keys := make([]*bls.ValidatorKey, 0, size)
	vals := make([]*validator.Validator, 0, size)
	for i := range size {
		key := td.RandValKey()
		keys = append(keys, key)
		vals = append(vals, validator.NewValidator(key.PublicKey(), int32(i)))
	}
	td.state.FakeCommittee.FakeValidators = vals

	return keys
}

// certifyBlock makes a certificate for the block that is signed by all the committee members.
func (td *testData) certifyBlock(blk *block.Block, height types.Height, keys []*bls.ValidatorKey,
) *certificate.Certificate {
	cert := certificate.NewCertificate(height, td.RandRound())
	signBytes := cert.SignBytesPrecommit(blk.Hash())

	committers := make([]int32, 0, len(keys))
	sigs := make([]*bls.Signature, 0, len(keys))
	for i, key := range keys {
		committers = append(committers, int32(i))
		sigs = append(sigs, key.Sign(signBytes))
	}

	aggSig, err := bls.SignatureAggregate(sigs...)
	if err != nil {
		panic(err)
	}
	cert.SetSignature(committers, []int32{}, aggSig)

	return cert
}

// The next block can be verified, and it should be committed if the certificate is made by the committee.
func TestHandlerBlockAnnounceVerifiedNextBlock(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	height := td.state.LastBlockHeight() + 1
	blk, _ := td.GenerateTestBlock(height)
	cert := td.certifyBlock(blk, height, keys)

	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk, cert, nil), pid)

	assert.Equal(t, height, td.state.LastBlockHeight())
	assert.Equal(t, height, td.sync.peerSet.GetPeer(pid).Height)
}

// The announcements are not authenticated: anyone can send a block with a fake certificate.
// A fake certificate for the next block should not be cached, nor be used for updating the peer.
func TestHandlerBlockAnnounceFakeCertificate(t *testing.T) {
	td := setup(t, nil)
	td.setupCommittee()
	// These keys do not belong to the committee.
	fakeKeys := []*bls.ValidatorKey{td.RandValKey(), td.RandValKey(), td.RandValKey(), td.RandValKey()}

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	lastHeight := td.state.LastBlockHeight()
	height := lastHeight + 1
	blk, _ := td.GenerateTestBlock(height)
	cert := td.certifyBlock(blk, height, fakeKeys)

	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk, cert, nil), pid)

	assert.Equal(t, lastHeight, td.state.LastBlockHeight())
	assert.False(t, td.sync.cache.HasBlockInCache(height), "the block should not be cached")
	assert.Nil(t, td.sync.cache.GetCertificate(height), "the certificate should not be cached")
	assert.NotEqual(t, height, td.sync.peerSet.GetPeer(pid).Height, "the height of the peer should not be updated")
}

// A certificate that is made by the committee for another block should not be accepted.
func TestHandlerBlockAnnounceCertificateOfAnotherBlock(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	lastHeight := td.state.LastBlockHeight()
	height := lastHeight + 1
	blk, _ := td.GenerateTestBlock(height)
	otherBlk, _ := td.GenerateTestBlock(height)
	cert := td.certifyBlock(otherBlk, height, keys)

	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk, cert, nil), pid)

	assert.Equal(t, lastHeight, td.state.LastBlockHeight())
	assert.False(t, td.sync.cache.HasBlockInCache(height))
	assert.Nil(t, td.sync.cache.GetCertificate(height))
}

// The certificate should be for the height of the block.
func TestHandlerBlockAnnounceMismatchedHeights(t *testing.T) {
	td := setup(t, nil)

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	height := td.RandHeight()
	blk, _ := td.GenerateTestBlock(height)
	cert := td.GenerateTestCertificate(height + 1)

	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk, cert, nil), pid)

	assert.False(t, td.sync.cache.HasBlockInCache(height))
	assert.Nil(t, td.sync.cache.GetCertificate(height))
	assert.Nil(t, td.sync.cache.GetCertificate(height+1))
	assert.NotEqual(t, height, td.sync.peerSet.GetPeer(pid).Height)
}

// The attacker is faster than the honest peer: the fake block `S+1` is announced first.
// It should not occupy the cache; otherwise the honest announcement is ignored as "already processed".
func TestHandlerBlockAnnounceAttackerFirst(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()
	fakeKeys := []*bls.ValidatorKey{td.RandValKey(), td.RandValKey(), td.RandValKey(), td.RandValKey()}

	attacker := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	honest := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	lastHeight := td.state.LastBlockHeight()
	height := lastHeight + 1

	fakeBlk, _ := td.GenerateTestBlock(height)
	fakeCert := td.certifyBlock(fakeBlk, height, fakeKeys)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(fakeBlk, fakeCert, nil), attacker)
	require.Equal(t, lastHeight, td.state.LastBlockHeight())

	blk, _ := td.GenerateTestBlock(height)
	cert := td.certifyBlock(blk, height, keys)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk, cert, nil), honest)

	assert.Equal(t, height, td.state.LastBlockHeight())
	assert.Equal(t, cert, td.state.FakeStore.FakeCertificates[height])
	cBlk, err := td.state.CommittedBlock(height)
	require.NoError(t, err)
	assert.Equal(t, blk.Hash(), cBlk.BlockHash)
}

// A block of a far height can't be verified yet; it is cached, but it should not break the verified blocks.
// The attacker announces the block of the height `S+2` with a fake certificate for `S+1` inside it.
// The honest block of the height `S+1` should still be committed.
func TestHandlerBlockAnnouncePoisoningPreviousCertificate(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()

	attacker := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	honest := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	lastHeight := td.state.LastBlockHeight()

	// The attacker announces `S+2`, and the previous certificate of its block is a fake one.
	fakeBlk, fakeCert := td.GenerateTestBlock(lastHeight + 2)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(fakeBlk, fakeCert, nil), attacker)
	require.Equal(t, lastHeight, td.state.LastBlockHeight())
	require.Equal(t, fakeBlk.PrevCertificate(), td.sync.cache.GetCertificate(lastHeight+1))

	// The honest peer announces the block `S+1`.
	blk1, _ := td.GenerateTestBlock(lastHeight + 1)
	cert1 := td.certifyBlock(blk1, lastHeight+1, keys)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(blk1, cert1, nil), honest)

	// The block `S+1` is committed with the verified certificate.
	// The fake state doesn't validate the cached block `S+2`, a real state rejects it.
	assert.GreaterOrEqual(t, td.state.LastBlockHeight(), lastHeight+1)
	cBlk, err := td.state.CommittedBlock(lastHeight + 1)
	require.NoError(t, err)
	assert.Equal(t, blk1.Hash(), cBlk.BlockHash)
	assert.Equal(t, cert1, td.state.FakeStore.FakeCertificates[lastHeight+1])
}

// A fake announcement for a height that can't be verified yet must not mask
// the announcement of the committee, which arrives when the height is the next one.
func TestHandlerBlockAnnounceVerifiableReplacesAnnounced(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()

	height := td.state.LastBlockHeight() + 1
	fakeBlk, fakeCert := td.GenerateTestBlock(height)

	// The fake block is cached as an announcement: it is the next block here only
	// to keep the test short, it is added the way an unverifiable announcement is.
	td.sync.cache.AddAnnouncedBlock(fakeBlk, fakeCert)
	require.True(t, td.sync.cache.HasAnnouncedBlock(height))

	realBlk, _ := td.GenerateTestBlock(height)
	realCert := td.certifyBlock(realBlk, height, keys)
	td.receivingNewMessage(td.sync, message.NewBlockAnnounceMessage(realBlk, realCert, nil), td.RandPeerID())

	// The verified announcement is not ignored: it replaced the fake one, or it is committed already.
	if cached := td.sync.cache.GetBlock(height); cached != nil {
		assert.Equal(t, realBlk.Hash(), cached.Hash())
		assert.False(t, td.sync.cache.HasAnnouncedBlock(height))
	}
	assert.NotEqual(t, fakeCert, td.sync.cache.GetCertificate(height))
}
