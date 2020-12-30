package stats

import (
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
)

var tGenHash crypto.Hash
var tPeerID1 peer.ID
var tPeerID2 peer.ID
var tStats *Stats

func setup(t *testing.T) {
	tGenHash := crypto.GenerateTestHash()
	tPeerID1, _ = peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	tPeerID2, _ = peer.IDB58Decode("12D3KooWPKx9syCJSytRg8x6V27SZsc2kiF89suxETZDsUMiWexZ")
	tStats = NewStats(tGenHash)
}

func TestParsSalamWithGoodGenesisHash(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	msg1 := message.NewSalamMessage("kitty", pub, tPeerID1, tGenHash, 123, 0)
	bs, _ := msg1.Encode(false, nil)
	assert.True(t, tStats.mustGetPeer(tPeerID1).BelongsToSameNetwork(tGenHash))
	msg2 := tStats.ParsMessage(bs, tPeerID1)
	assert.Equal(t, msg1.SignBytes(), msg2.SignBytes())
	assert.True(t, tStats.getPeer(tPeerID1).BelongsToSameNetwork(tGenHash))
}

func TestParsSalamWithBadGenesisHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()
	msg1 := message.NewSalamMessage("kitty", pub, tPeerID1, invHash, 123, 0)
	bs, _ := msg1.Encode(false, nil)
	msg2 := tStats.ParsMessage(bs, tPeerID1)
	assert.Equal(t, msg1.SignBytes(), msg2.SignBytes())
	assert.False(t, tStats.getPeer(tPeerID1).BelongsToSameNetwork(tGenHash))
}
