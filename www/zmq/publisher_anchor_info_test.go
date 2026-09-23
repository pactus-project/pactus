package zmq

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestAnchorPublisherAbsentWhenUnset(t *testing.T) {
	td := setup(t, DefaultConfig())

	for _, pub := range td.server.publishers {
		require.NotEqual(t, TopicAnchorInfo.String(), pub.TopicName())
	}
}

func TestAnchorInfoPublisher(t *testing.T) {
	td, sub := subscribeAnchor(t)
	height := uint32(0x01020304)

	first, _ := td.RandBLSKeyPair()
	second, _ := td.RandBLSKeyPair()
	set32 := anchorTx(first.AccountAddress(), payload.AnchorActionSet, bytesRepeat(0x11, 32))
	setZero := anchorTx(first.AccountAddress(), payload.AnchorActionSet, make([]byte, 32))
	set64 := anchorTx(second.AccountAddress(), payload.AnchorActionSet, bytesRepeat(0x22, 64))
	del := anchorTx(first.AccountAddress(), payload.AnchorActionDelete, bytesRepeat(0x33, 32))
	updated := anchorTx(second.AccountAddress(), payload.AnchorActionSet, bytesRepeat(0x44, 32))
	transfer := td.GenerateTestTransferTx()

	plain, _ := td.GenerateTestBlock(1)
	td.pipe.Send(plain)
	td.pipe.Send(transfer)

	sendBlock(t, td, []*tx.Tx{set32})
	msg := mustRecv(t, sub)
	require.Equal(t, []byte{0x00, 0x05}, msg[:2])
	require.Equal(t, first.AccountAddress().Bytes(), msg[2:23])
	require.Equal(t, payload.AnchorActionSet, msg[23])
	require.Equal(t, height, binary.BigEndian.Uint32(msg[24:28]))
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, msg[24:28])
	require.Equal(t, byte(32), msg[28])
	require.Equal(t, bytesRepeat(0x11, 32), msg[29:61])
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(msg[61:]))
	require.Len(t, msg, 65)

	sendBlock(t, td, []*tx.Tx{del})
	msg = mustRecv(t, sub)
	require.Equal(t, payload.AnchorActionDelete, msg[23])
	require.Equal(t, byte(0), msg[28])
	require.Equal(t, uint32(1), binary.BigEndian.Uint32(msg[29:]))
	require.Len(t, msg, 33)

	sendBlock(t, td, []*tx.Tx{set32, transfer, del})
	firstMsg := mustRecv(t, sub)
	secondMsg := mustRecv(t, sub)
	require.Equal(t, uint32(2), binary.BigEndian.Uint32(firstMsg[len(firstMsg)-4:]))
	require.Equal(t, uint32(3), binary.BigEndian.Uint32(secondMsg[len(secondMsg)-4:]))
	require.Equal(t, payload.AnchorActionSet, firstMsg[23])
	require.Equal(t, payload.AnchorActionDelete, secondMsg[23])

	sendBlock(t, td, []*tx.Tx{setZero})
	msg = mustRecv(t, sub)
	require.Equal(t, byte(32), msg[28])
	require.Equal(t, make([]byte, 32), msg[29:61])
	require.Len(t, msg, 65)

	sendBlock(t, td, []*tx.Tx{set64})
	msg = mustRecv(t, sub)
	require.Equal(t, byte(64), msg[28])
	require.Equal(t, bytesRepeat(0x22, 64), msg[29:93])
	require.Equal(t, uint32(5), binary.BigEndian.Uint32(msg[93:]))
	require.Len(t, msg, 97)

	sendBlock(t, td, []*tx.Tx{set32, set64})
	one := mustRecv(t, sub)
	two := mustRecv(t, sub)
	require.Equal(t, first.AccountAddress().Bytes(), one[2:23])
	require.Equal(t, second.AccountAddress().Bytes(), two[2:23])
	require.NotEqual(t, one[2:23], two[2:23])

	others := []*tx.Tx{
		td.GenerateTestSubsidyTx(),
		td.GenerateTestTransferTx(),
		td.GenerateTestBondTx(),
		td.GenerateTestUnbondTx(),
		td.GenerateTestWithdrawTx(),
		td.GenerateTestSortitionTx(),
		td.GenerateTestBatchTransferTx(),
		updated,
	}
	sendBlock(t, td, others)
	msg = mustRecv(t, sub)
	require.Equal(t, second.AccountAddress().Bytes(), msg[2:23])
	require.Equal(t, payload.AnchorActionSet, msg[23])
	require.Equal(t, bytesRepeat(0x44, 32), msg[29:61])
	require.Equal(t, uint32(8), binary.BigEndian.Uint32(msg[len(msg)-4:]))

	require.NoError(t, sub.Close())
}

func TestAnchorInfoTwoSetsSameSender(t *testing.T) {
	td, sub := subscribeAnchor(t)
	sender, _ := td.RandBLSKeyPair()
	addr := sender.AccountAddress()
	first := anchorTx(addr, payload.AnchorActionSet, bytesRepeat(0x11, 32))
	second := anchorTx(addr, payload.AnchorActionSet, bytesRepeat(0x22, 32))

	plain, _ := td.GenerateTestBlock(1)
	td.pipe.Send(plain)
	sendBlock(t, td, []*tx.Tx{first, second})

	one := mustRecv(t, sub)
	two := mustRecv(t, sub)
	require.Equal(t, addr.Bytes(), one[2:23])
	require.Equal(t, addr.Bytes(), two[2:23])
	require.Equal(t, payload.AnchorActionSet, one[23])
	require.Equal(t, payload.AnchorActionSet, two[23])
	require.Equal(t, byte(32), one[28])
	require.Equal(t, bytesRepeat(0x11, 32), one[29:61])
	require.Equal(t, bytesRepeat(0x22, 32), two[29:61])
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(one[len(one)-4:]))
	require.Equal(t, uint32(1), binary.BigEndian.Uint32(two[len(two)-4:]))
	require.Len(t, one, 65)
	require.Len(t, two, 65)
	require.NoError(t, sub.Close())
}

func TestAnchorInfoThreeAnchorsSameSender(t *testing.T) {
	td, sub := subscribeAnchor(t)
	sender, _ := td.RandBLSKeyPair()
	addr := sender.AccountAddress()
	zero := anchorTx(addr, payload.AnchorActionSet, make([]byte, 32))
	wide := anchorTx(addr, payload.AnchorActionSet, bytesRepeat(0x22, 64))
	del := anchorTx(addr, payload.AnchorActionDelete, bytesRepeat(0x33, 32))

	plain, _ := td.GenerateTestBlock(1)
	td.pipe.Send(plain)
	sendBlock(t, td, []*tx.Tx{zero, wide, del})

	first := mustRecv(t, sub)
	second := mustRecv(t, sub)
	third := mustRecv(t, sub)
	require.Equal(t, addr.Bytes(), first[2:23])
	require.Equal(t, addr.Bytes(), second[2:23])
	require.Equal(t, addr.Bytes(), third[2:23])
	require.Equal(t, payload.AnchorActionSet, first[23])
	require.Equal(t, payload.AnchorActionSet, second[23])
	require.Equal(t, payload.AnchorActionDelete, third[23])
	require.Equal(t, byte(32), first[28])
	require.Equal(t, make([]byte, 32), first[29:61])
	require.Equal(t, byte(64), second[28])
	require.Equal(t, bytesRepeat(0x22, 64), second[29:93])
	require.Equal(t, byte(0), third[28])
	require.Len(t, first, 65)
	require.Len(t, second, 97)
	require.Len(t, third, 33)
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(first[len(first)-4:]))
	require.Equal(t, uint32(1), binary.BigEndian.Uint32(second[len(second)-4:]))
	require.Equal(t, uint32(2), binary.BigEndian.Uint32(third[len(third)-4:]))
	require.NoError(t, sub.Close())
}

func TestAnchorInfoInterleavedSenders(t *testing.T) {
	td, sub := subscribeAnchor(t)
	firstKey, _ := td.RandBLSKeyPair()
	secondKey, _ := td.RandBLSKeyPair()
	first := firstKey.AccountAddress()
	second := secondKey.AccountAddress()
	firstSet := anchorTx(first, payload.AnchorActionSet, bytesRepeat(0x11, 32))
	middleSet := anchorTx(second, payload.AnchorActionSet, bytesRepeat(0x22, 32))
	secondSet := anchorTx(first, payload.AnchorActionSet, bytesRepeat(0x33, 32))

	plain, _ := td.GenerateTestBlock(1)
	td.pipe.Send(plain)
	sendBlock(t, td, []*tx.Tx{firstSet, middleSet, secondSet})

	one := mustRecv(t, sub)
	two := mustRecv(t, sub)
	three := mustRecv(t, sub)
	require.Equal(t, first.Bytes(), one[2:23])
	require.Equal(t, second.Bytes(), two[2:23])
	require.Equal(t, first.Bytes(), three[2:23])
	require.Equal(t, bytesRepeat(0x11, 32), one[29:61])
	require.Equal(t, bytesRepeat(0x22, 32), two[29:61])
	require.Equal(t, bytesRepeat(0x33, 32), three[29:61])
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(one[len(one)-4:]))
	require.Equal(t, uint32(1), binary.BigEndian.Uint32(two[len(two)-4:]))
	require.Equal(t, uint32(2), binary.BigEndian.Uint32(three[len(three)-4:]))
	require.NoError(t, sub.Close())
}

func subscribeAnchor(t *testing.T) (*testData, zmq4.Socket) {
	t.Helper()

	port := testsuite.FindFreePort()
	conf := DefaultConfig()
	conf.ZmqPubAnchorInfo = fmt.Sprintf("tcp://127.0.0.1:%d", port)
	td := setup(t, conf)

	sub := zmq4.NewSub(t.Context(), zmq4.WithAutomaticReconnect(false), zmq4.WithTimeout(200*time.Millisecond))
	require.NoError(t, sub.SetOption(zmq4.OptionSubscribe, string(TopicAnchorInfo.Bytes())))
	require.NoError(t, sub.Dial(conf.ZmqPubAnchorInfo))
	time.Sleep(100 * time.Millisecond)

	return td, sub
}

func sendBlock(t *testing.T, td *testData, txs []*tx.Tx) {
	t.Helper()

	blk, _ := td.GenerateTestBlock(0x01020304, testsuite.BlockWithTransactions(block.Txs(txs)))
	require.Equal(t, types.Height(0x01020304), blk.Height())
	td.pipe.Send(blk)
}

func mustRecv(t *testing.T, sub zmq4.Socket) []byte {
	t.Helper()

	received, err := sub.Recv()
	require.NoError(t, err)
	require.NotEmpty(t, received.Frames)

	return received.Frames[0]
}

func anchorTx(from crypto.Address, action byte, root []byte) *tx.Tx {
	deposit := amount.Amount(1)
	if action == payload.AnchorActionDelete {
		deposit = 0
	}

	return tx.NewAnchorTx(1, from, action, root, "", 0, deposit, 1)
}

func bytesRepeat(fill byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = fill
	}

	return out
}
