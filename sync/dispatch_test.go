package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"

	"github.com/zarbchain/zarb-go/vote"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

func TestSendTxs(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in his cache
	tAliceSync.cache.AddTransaction(trx1)
	tBobSync.cache.AddTransaction(trx2)

	tAliceBroadcastCh <- message.NewTxsReqMessage([]crypto.Hash{trx1.ID()})
	tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeTxsReq)

	tAliceSync.cache.AddTransaction(trx1)
	tAliceBroadcastCh <- message.NewTxsReqMessage([]crypto.Hash{trx1.ID(), trx2.ID()})
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeTxsReq)
	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeTxs)

	assert.NotNil(t, tAliceSync.cache.GetTransaction(trx2.ID()))
}

func TestSendVoteSet(t *testing.T) {
	setup(t)

	tAliceConsensus.HRS_ = hrs.NewHRS(100, 1, 1)
	tBobConsensus.HRS_ = hrs.NewHRS(100, 1, 1)
	v1, _ := vote.GenerateTestPrepareVote(100, 0)
	v2, _ := vote.GenerateTestPrepareVote(100, 1)
	v3, _ := vote.GenerateTestPrepareVote(100, 1)
	v4, _ := vote.GenerateTestPrepareVote(100, 1)
	v5, _ := vote.GenerateTestPrepareVote(101, 1)

	tAliceConsensus.Votes = []*vote.Vote{v2, v3}
	tBobBroadcastCh <- message.NewVoteSetMessage(100, 1, []crypto.Hash{v1.Hash(), v4.Hash()})
	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeVoteSet)
	tAliceNetAPI.shouldPublishThisMessage(t, message.NewVoteMessage(v2))
	tAliceNetAPI.shouldPublishThisMessage(t, message.NewVoteMessage(v3))

	tBobBroadcastCh <- message.NewVoteSetMessage(101, 1, []crypto.Hash{v5.Hash()})
	tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeVoteSet)
	tAliceNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeVote)

}
