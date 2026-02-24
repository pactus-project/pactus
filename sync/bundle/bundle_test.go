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

func TestInvalidCBOR(t *testing.T) {
	data1, _ := hex.DecodeString("00")
	data2, _ := hex.DecodeString("A3")
	data3, _ := hex.DecodeString("A3010002000340")
	bdl := new(Bundle)
	_, err := bdl.Decode(bytes.NewReader(data1))
	assert.Error(t, err)
	_, err = bdl.Decode(bytes.NewReader(data2))
	assert.Error(t, err)
	_, err = bdl.Decode(bytes.NewReader(data3))
	assert.Error(t, err)
}

func TestMessageCompress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blocksData := [][]byte{}
	for i := 0; i < ts.RandIntMax(40); i++ {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		data, _ := blk.Bytes()
		blocksData = append(blocksData, data)
	}
	msg1 := message.NewBlocksResponseMessage(message.ResponseCodeOK, message.ResponseCodeOK.String(),
		1234, 888, blocksData, nil)
	bdl := NewBundle(msg1)
	bs0, err := bdl.Encode()
	assert.NoError(t, err)
	assert.False(t, util.IsFlagSet(bdl.Flags, BundleFlagCompressed))

	bdl.CompressIt()
	bs1, err := bdl.Encode()
	assert.NoError(t, err)
	assert.True(t, util.IsFlagSet(bdl.Flags, BundleFlagCompressed))

	fmt.Printf("Compressed :%v%%\n", 100-len(bs1)*100/(len(bs0)))
	fmt.Printf("Uncompressed len :%v\n", len(bs0))
	fmt.Printf("Compressed len :%v\n", len(bs1))

	msg2 := new(Bundle)
	bytesRead1, err := msg2.Decode(bytes.NewReader(bs0))
	assert.NoError(t, err)
	assert.Equal(t, len(bs0), bytesRead1)
	assert.NoError(t, msg2.BasicCheck())

	msg3 := new(Bundle)
	bytesRead2, err := msg3.Decode(bytes.NewReader(bs1))
	assert.NoError(t, err)
	assert.Equal(t, len(bs1), bytesRead2)
	assert.NoError(t, msg3.BasicCheck())
}

func TestDecodeVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	v, _ := ts.GenerateTestPrecommitVote(88, 0)
	msg := message.NewVoteMessage(v)
	bdl := NewBundle(msg)
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
	dat1, _ := hex.DecodeString(
		"a4" + // Map(4)
			"0100" + // Flags = 0
			"0207" + // Message Type = 7 (TypeVote)
			"035879" + // Message + Len
			"" + "a101a7010102186403010458200264572d4d6bfcd2140d4f885fd5a32fe42fdb" + // Vote Message Uncompressed
			"" + "f40551e4ff89f3d235e32b4b92055501c0067d277f2dff99943016d6a0f379cf" +
			"" + "09846c6f06f60758308ab7aecbe03c4ed5b688bcb7e848baffa62bcbf1a40215" +
			"" + "22c56693f0a7bbcc1fe865277556ee59c1f63ba592acfe1b43" +
			"041a00001234") // Consensus Height (0x00001234)
	data2, _ := hex.DecodeString(
		"a4" + // Map(4)
			"01190100" + // Flags = 0x0100 (compressed)
			"0207" + // Message Type = 7 (TypeVote)
			"035895" + // Message + Len
			"" + "1f8b08000000000000ff00790086ffa101a7010102186403010458200264572d" + // Vote Uncompressed
			"" + "4d6bfcd2140d4f885fd5a32fe42fdbf40551e4ff89f3d235e32b4b92055501c0" +
			"" + "067d277f2dff99943016d6a0f379cf09846c6f06f60758308ab7aecbe03c4ed5" +
			"" + "b688bcb7e848baffa62bcbf1a4021522c56693f0a7bbcc1fe865277556ee59c1" +
			"" + "f63ba592acfe1b43010000ffff798ce7ec79000000" +
			"041a00001234") // Consensus Height (0x00001234)

	bdl1 := new(Bundle)
	bdl2 := new(Bundle)
	bytesRead1, err := bdl1.Decode(bytes.NewReader(dat1))
	require.NoError(t, err)
	assert.Equal(t, len(dat1), bytesRead1)
	assert.NoError(t, bdl1.BasicCheck())

	bytesRead2, err := bdl2.Decode(bytes.NewReader(data2))
	require.NoError(t, err)
	assert.Equal(t, len(data2), bytesRead2)
	assert.NoError(t, bdl2.BasicCheck())

	assert.Equal(t, bdl1.Message, bdl2.Message)
	assert.Equal(t, 0x0000, bdl1.Flags)
	assert.Equal(t, 0x0100, bdl2.Flags)
	assert.Contains(t, bdl1.LogString(), "vote")

	assert.Equal(t, uint32(0x1234), bdl1.ConsensusHeight)
	assert.Equal(t, uint32(0x1234), bdl2.ConsensusHeight)
}

func TestEncodingData(t *testing.T) {
	t.Run("Encoding non-consensus message", func(t *testing.T) {
		msg := message.NewBlocksRequestMessage(0x00, 0x12, 0x13)
		bdl := NewBundle(msg)
		data, _ := bdl.Encode()

		expectedData := "a4" + // Map(3)
			"0100" + // Flags = 0
			"0209" + // Message Type = 9 (TypeBlocksRequest)
			"0347" + // Message + Len
			"" + "a3" +
			"" + "0100" +
			"" + "0212" +
			"" + "0313" +
			"041a00000000" // Consensus height (0x00000000)
		assert.Equal(t, expectedData, hex.EncodeToString(data))
		assert.Equal(t, uint32(0x00), bdl.ConsensusHeight)
	})

	t.Run("Encoding consensus message", func(t *testing.T) {
		ts := testsuite.NewTestSuite(t)

		rndAddr := ts.RandValAddress()
		msg := message.NewQueryVoteMessage(0x12, 0x01, rndAddr)
		bdl := NewBundle(msg)
		data, _ := bdl.Encode()

		expectedData := "a4" + // Map(3)
			"0100" + // Flags = 0
			"0206" + // Message Type = 6 (TypeQueryVote)
			"03581c" + // Message + Len
			"" + "a3" +
			"" + "0112" +
			"" + "0201" +
			"" + "0355" + hex.EncodeToString(rndAddr.Bytes()) +
			"041a00000012" // Consensus height (0x00000012)
		assert.Equal(t, expectedData, hex.EncodeToString(data))
		assert.Equal(t, uint32(0x12), bdl.ConsensusHeight)
	})
}

func TestCBORLengthAttack(t *testing.T) {
	tests := []struct {
		data string
		msg  string
	}{
		{"9A00010001", "exceeded max number of elements 65536"},        // Major type 4 (100x xxxx): Array
		{"BA00010001", "exceeded max number of key-value pairs 65536"}, // Major type 5 (101x xxxx): Map
	}

	for _, tt := range tests {
		data, _ := hex.DecodeString(tt.data)
		bdl := new(Bundle)
		_, err := bdl.Decode(bytes.NewReader(data))

		assert.ErrorContains(t, err, tt.msg)
	}
}
