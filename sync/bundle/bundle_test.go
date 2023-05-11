package bundle

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	pid := network.TestRandomPeerID()
	msg := NewBundle(pid, message.NewQueryProposalMessage(100, 0))
	assert.Zero(t, msg.Flags)
	assert.Equal(t, msg.Initiator, pid)
}

func TestInvalidCBOR(t *testing.T) {
	d1, _ := hex.DecodeString("000000000000000000")
	d2, _ := hex.DecodeString("A401000242000003000440")
	m := new(Bundle)
	_, err := m.Decode(bytes.NewReader(d1))
	assert.Error(t, err)
	_, err = m.Decode(bytes.NewReader(d2))
	assert.Error(t, err)
}
func TestMessageCompress(t *testing.T) {
	var blocks = []*block.Block{}
	for i := 0; i < 10; i++ {
		b := block.GenerateTestBlock(nil, nil)
		blocks = append(blocks, b)
	}
	msg := message.NewBlocksResponseMessage(message.ResponseCodeBusy, 1234, 888, blocks, nil)
	bdl := NewBundle(network.TestRandomPeerID(), msg)
	bs0, err := bdl.Encode()
	assert.NoError(t, err)
	bdl.CompressIt()
	bs1, err := bdl.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs1)*100/(len(bs0)))
	fmt.Printf("Uncompressed len :%v\n", len(bs0))
	fmt.Printf("Compressed len :%v\n", len(bs1))
	msg2 := new(Bundle)
	msg3 := new(Bundle)
	_, err = msg2.Decode(bytes.NewReader(bs0))
	assert.NoError(t, err)
	_, err = msg3.Decode(bytes.NewReader(bs1))
	assert.NoError(t, err)
	assert.NoError(t, msg2.SanityCheck())
	assert.NoError(t, msg3.SanityCheck())
	assert.True(t, util.IsFlagSet(bdl.Flags, BundleFlagCompressed))
}

func TestDecodeVoteMessage(t *testing.T) {
	v, _ := vote.GenerateTestPrecommitVote(88, 0)
	msg := message.NewVoteMessage(v)
	bdl := NewBundle(network.TestRandomPeerID(), msg)
	bs0, err := bdl.Encode()
	assert.NoError(t, err)
	bdl.CompressIt()
	bs1, err := bdl.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs1)*100/(len(bs0)))
	fmt.Printf("Uncompressed len :%v\n", len(bs0))
	fmt.Printf("Compressed len :%v\n", len(bs1))
}

func TestDecodeVoteCBOR(t *testing.T) {
	d1, _ := hex.DecodeString(
		"a40119001102582212206e58a3dbd95357010000000000000000000000000000000000000000000000000307045877a101a6010202185803" +
			"000458205ffca0da6582ee795bdb73a518797bd4f2ccde1f8692e2b2a5ba0dd60f576410055501c94b4b3489c5370ae23923c2325cd80eee" +
			"749231065830a009f5f3ebe4fef05602813d099c539d13ba6ae209ecefe1ca72c55fd9b392ddb828d35a9d64abb3ca9694963e2d8338")
	// Compressed
	d2, _ := hex.DecodeString(
		"a40119011102582212206e58a3dbd953570100000000000000000000000000000000000000000000000003070458931f8b08000000000000" +
			"ff00770088ffa101a6010202185803000458205ffca0da6582ee795bdb73a518797bd4f2ccde1f8692e2b2a5ba0dd60f576410055501c94b" +
			"4b3489c5370ae23923c2325cd80eee749231065830a009f5f3ebe4fef05602813d099c539d13ba6ae209ecefe1ca72c55fd9b392ddb828d3" +
			"5a9d64abb3ca9694963e2d8338010000ffffcf174a7977000000")
	m1 := new(Bundle)
	m2 := new(Bundle)
	_, err := m1.Decode(bytes.NewReader(d1))
	assert.NoError(t, err)
	_, err = m2.Decode(bytes.NewReader(d2))
	assert.NoError(t, err)
	assert.NoError(t, m2.SanityCheck())

	assert.Equal(t, m1.Message, m2.Message)
}
