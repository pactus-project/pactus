package bundle

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pid := ts.RandPeerID()
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
	ts := testsuite.NewTestSuite(t)

	blocksData := [][]byte{}
	for i := 0; i < 10; i++ {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		d, _ := blk.Bytes()
		blocksData = append(blocksData, d)
	}
	msg := message.NewBlocksResponseMessage(message.ResponseCodeOK, message.ResponseCodeOK.String(),
		1234, 888, blocksData, nil)
	bdl := NewBundle(ts.RandPeerID(), msg)
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
	assert.NoError(t, msg2.BasicCheck())
	assert.NoError(t, msg3.BasicCheck())
	assert.True(t, util.IsFlagSet(bdl.Flags, BundleFlagCompressed))
}

func TestDecodeVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	v, _ := ts.GenerateTestPrecommitVote(88, 0)
	msg := message.NewVoteMessage(v)
	bdl := NewBundle(ts.RandPeerID(), msg)
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
		"a401000258221220aa8ece876ceabec6bd64430d0ad8d7bd5dbd72500f0a1642a13bd71b5b3522a90307045879a101a70101" +
			"02186403010458200264572d4d6bfcd2140d4f885fd5a32fe42fdbf40551e4ff89f3d235e32b4b92055501c0067d277f2dff" +
			"99943016d6a0f379cf09846c6f06f60758308ab7aecbe03c4ed5b688bcb7e848baffa62bcbf1a4021522c56693f0a7bbcc1f" +
			"e865277556ee59c1f63ba592acfe1b43")
	d2, _ := hex.DecodeString(
		"a4011901000258221220aa8ece876ceabec6bd64430d0ad8d7bd5dbd72500f0a1642a13bd71b5b3522a903070458951f8b08" +
			"000000000000ff00790086ffa101a7010102186403010458200264572d4d6bfcd2140d4f885fd5a32fe42fdbf40551e4ff89" +
			"f3d235e32b4b92055501c0067d277f2dff99943016d6a0f379cf09846c6f06f60758308ab7aecbe03c4ed5b688bcb7e848ba" +
			"ffa62bcbf1a4021522c56693f0a7bbcc1fe865277556ee59c1f63ba592acfe1b43010000ffff798ce7ec79000000")

	bdl1 := new(Bundle)
	bdl2 := new(Bundle)
	_, err := bdl1.Decode(bytes.NewReader(d1))
	require.NoError(t, err)
	_, err = bdl2.Decode(bytes.NewReader(d2))
	require.NoError(t, err)
	assert.NoError(t, bdl2.BasicCheck())

	assert.Equal(t, bdl1.Message, bdl2.Message)
}
