package block_test

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/simplemerkle"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		data := ts.DecodingHex(
			"01" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"00") // Txs: Len

		blk, _ := block.FromBytes(data)

		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "no subsidy transaction",
		})
	})

	t.Run("Too many transactions", func(t *testing.T) {
		data := ts.DecodingHex(
			"02" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"e907") // Txs: Len (1001)

		_, err := block.FromBytes(data)
		assert.ErrorIs(t, err, block.ErrTooManyTransactions)
	})

	t.Run("Without the previous certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight(), testsuite.BlockWithPrevCert(nil))

		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid genesis block hash",
		})
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		data := ts.DecodingHex(
			"01" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"00000000" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"02" + // Tx[0]: Flags
				"01" + // Tx[0]: Version
				"01000000" + // Tx[0]: LockTime
				"01" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		blk, _ := block.FromBytes(data)

		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid certificate: height is not positive: 0",
		})
	})

	t.Run("Invalid transaction", func(t *testing.T) {
		data := ts.DecodingHex(
			"01" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"02" + // Tx[0]: Flags
				"00" + // Tx[0]: Version
				"00000000" + // Tx[0]: LockTime
				"00" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		blk, _ := block.FromBytes(data)

		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid transaction: invalid version: 0",
		})
	})

	t.Run("Invalid previous block hash", func(t *testing.T) {
		data := ts.DecodingHex(
			"01" + // Version
				"00000000" + // UnixTime
				"0000000000000000000000000000000000000000000000000000000000000000" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"00" + // Tx[0]: Flags
				"01" + // Tx[0]: Version
				"01000000" + // Tx[0]: LockTime
				"01" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		_, err := block.FromBytes(data)
		assert.Error(t, err)
	})

	t.Run("Invalid proposer address (type is 2)", func(t *testing.T) {
		data := ts.DecodingHex(
			"01" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"02AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"02" + // Tx[0]: Flags
				"01" + // Tx[0]: Version
				"01000000" + // Tx[0]: LockTime
				"01" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		blk, _ := block.FromBytes(data)
		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid proposer address: pc1z42424242424242424242424242424242klpmq4",
		})
	})

	t.Run("Invalid Version", func(t *testing.T) {
		data := ts.DecodingHex(
			"00" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"02" + // Tx[0]: Flags
				"01" + // Tx[0]: Version
				"01000000" + // Tx[0]: LockTime
				"00" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		blk, _ := block.FromBytes(data)
		err := blk.BasicCheck()
		assert.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid block version: 0",
		})
	})

	t.Run("Ok", func(t *testing.T) {
		data := ts.DecodingHex(
			"02" + // Version
				"00000000" + // UnixTime
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
				"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" + // StateRoot
				"333333333333333333333333333333333333333333333333" + // SortitionSeed
				"333333333333333333333333333333333333333333333333" +
				"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
				"04030201" + // PrevCert: Height
				"0100" + // PrevCert: Round
				"0401020304" + // PrevCert: Committers
				"0102" + // PrevCert: Absentees
				"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
				"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
				"01" + // Txs: Len
				"02" + // Tx[0]: Flags
				"01" + // Tx[0]: Version
				"01000000" + // Tx[0]: LockTime
				"00" + // Tx[0]: Fee
				"00" + // Tx[0]: Memo
				"01" + // Tx[0]: PayloadType
				"00" + // Tx[0]: Sender (treasury)
				"022222222222222222222222222222222222222222" + // Tx[0]: Receiver
				"01") // Tx[0]: Amount

		blk, _ := block.FromBytes(data)
		assert.NoError(t, blk.BasicCheck())
		assert.Zero(t, blk.Header().UnixTime())
		assert.Equal(t, protocol.ProtocolVersionLatest, blk.Header().Version())
	})
}

func TestCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk1, _ := ts.GenerateTestBlock(ts.RandHeight())
	bz1, err := cbor.Marshal(blk1)
	assert.NoError(t, err)
	var blk2 block.Block
	err = cbor.Unmarshal(bz1, &blk2)
	assert.NoError(t, err)
	assert.NoError(t, blk2.BasicCheck())
	assert.Equal(t, blk1.Hash(), blk2.Hash())

	assert.Equal(t, blk1.Hash(), blk2.Hash())

	err = cbor.Unmarshal([]byte{1}, &blk2)
	assert.Error(t, err)
}

func TestEncodingBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	length := blk.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, blk.Encode(w), "encode test %v failed", i)
	}
	writer := util.NewFixedWriter(length)
	assert.NoError(t, blk.Encode(writer))

	for i := 0; i < length; i++ {
		blk2 := new(block.Block)
		r := util.NewFixedReader(i, writer.Bytes())
		assert.Error(t, blk2.Decode(r), "decode test %v failed", i)
	}

	blk2 := new(block.Block)
	r := util.NewFixedReader(length, writer.Bytes())
	assert.NoError(t, blk2.Decode(r))
	assert.Equal(t, blk.Hash(), blk2.Hash())
	assert.Equal(t, blk.Header(), blk2.Header())
}

func TestTxFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	bs, _ := blk.Bytes()
	_, err := block.FromBytes(bs)
	assert.NoError(t, err)
	_, err = blk.Bytes()
	assert.NoError(t, err)

	// Invalid data
	_, err = block.FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestBlockHash(t *testing.T) {
	data, _ := hex.DecodeString(
		"01" + // Version
			"00000000" + // UnixTime
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // PrevBlockHash
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // StateRoot
			"333333333333333333333333333333333333333333333333" + // SortitionSeed
			"333333333333333333333333333333333333333333333333" +
			"01AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // ProposerAddress
			"04030201" + // PrevCert: Height
			"0100" + // PrevCert: Round
			"0401020304" + // PrevCert: Committers
			"0102" + // PrevCert: Absentees
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d" +
			"2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // PrevCert: Signature
			"01" + // Txs: Len
			"02" + // Tx[0]: Flags
			"01" + // Tx[0]: Version
			"01000000" + // Tx[0]: LockTime
			"01" + // Tx[0]: Fee
			"00" + // Tx[0]: Memo
			"01" + // Tx[0]: PayloadType
			"00" + // Tx[0]: Sender (treasury)
			"012222222222222222222222222222222222222222" + // Tx[0]: Receiver
			"01") // Tx[0]: Amount

	blk, err := block.FromBytes(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), blk.SerializeSize())

	blockData, _ := blk.Bytes()
	assert.Equal(t, data, blockData)

	headerSize := blk.Header().SerializeSize()
	headerData := data[:headerSize]
	certSize := blk.PrevCertificate().SerializeSize()
	certData := data[headerSize : headerSize+certSize]
	certHash := hash.CalcHash(certData)

	txHashes := make([]hash.Hash, 0)
	for _, trx := range blk.Transactions() {
		txHashes = append(txHashes, trx.ID())
	}
	txRoot := simplemerkle.NewTreeFromHashes(txHashes).Root()

	hashData := headerData
	hashData = append(hashData, certHash.Bytes()...)
	hashData = append(hashData, txRoot.Bytes()...)
	hashData = append(hashData, util.Int32ToSlice(int32(blk.Transactions().Len()))...)

	expected1 := hash.CalcHash(hashData)
	expected2, _ := hash.FromString("43399fa59adcfb7d8c515460ec9ca27b6a1cb865f5b7d9bde8fe56c18eaec9ab")
	assert.Equal(t, expected1, blk.Hash())
	assert.Equal(t, expected2, blk.Hash())
}

func TestMakeBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk0, _ := ts.GenerateTestBlock(ts.RandHeight())
	blk1 := block.MakeBlock(protocol.ProtocolVersion2, blk0.Header().Time(), blk0.Transactions(),
		blk0.Header().PrevBlockHash(),
		blk0.Header().StateRoot(),
		blk0.PrevCertificate(),
		blk0.Header().SortitionSeed(),
		blk0.Header().ProposerAddress())

	assert.Equal(t, blk0.Hash(), blk1.Hash())
}

func TestBlockHeight(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk1, _ := ts.GenerateTestBlock(1, testsuite.BlockWithPrevCert(nil), testsuite.BlockWithPrevHash(hash.UndefHash))
	blk2, _ := ts.GenerateTestBlock(2)

	assert.NoError(t, blk1.BasicCheck())
	assert.NoError(t, blk2.BasicCheck())

	assert.Equal(t, uint32(1), blk1.Height())
	assert.Equal(t, uint32(2), blk2.Height())
}
