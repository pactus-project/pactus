package bundle

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// makeAllMessages returns a valid message of each type, with all the types of the votes and justifications.
func makeAllMessages(ts *testsuite.TestSuite) map[string]message.Message {
	valKey := ts.RandValKey()
	addr := valKey.Address()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := ts.GenerateTestCertificate(height)
	blk, blkCert := ts.GenerateTestBlock(height)
	blockHash := ts.RandHash()

	signed := func(v *vote.Vote) *vote.Vote {
		ts.HelperSignVote(valKey, v)

		return v
	}

	conflict := &vote.JustMainVoteConflict{JustNo: &vote.JustInitNo{QCert: cert}, JustYes: &vote.JustInitYes{}}

	return map[string]message.Message{
		"hello": message.NewHelloMessage(ts.RandPeerID(), "moniker",
			service.New(service.FullNode), height, blockHash, blockHash),
		"hello-ack":      message.NewHelloAckMessage(message.ResponseCodeOK, "ok", height),
		"transactions":   message.NewTransactionsMessage([]*tx.Tx{ts.GenerateTestTransferTx(), ts.GenerateTestBondTx()}),
		"query-proposal": message.NewQueryProposalMessage(height, round, addr),
		"proposal":       message.NewProposalMessage(ts.GenerateTestProposal(height, round)),
		"query-vote":     message.NewQueryVoteMessage(height, round, addr),
		"prepare":        message.NewVoteMessage(signed(vote.NewPrepareVote(blockHash, height, round, addr))),
		"precommit":      message.NewVoteMessage(signed(vote.NewPrecommitVote(blockHash, height, round, addr))),
		"cp-init-no": message.NewVoteMessage(signed(vote.NewCPPreVote(blockHash, height, round, 0,
			vote.CPValueNo, &vote.JustInitNo{QCert: cert}, addr))),
		"cp-init-yes": message.NewVoteMessage(signed(vote.NewCPPreVote(blockHash, height, round, 0,
			vote.CPValueYes, &vote.JustInitYes{}, addr))),
		"cp-pre-vote-hard": message.NewVoteMessage(signed(vote.NewCPPreVote(blockHash, height, round, 1,
			vote.CPValueYes, &vote.JustPreVoteHard{QCert: cert}, addr))),
		"cp-pre-vote-soft": message.NewVoteMessage(signed(vote.NewCPPreVote(blockHash, height, round, 1,
			vote.CPValueYes, &vote.JustPreVoteSoft{QCert: cert}, addr))),
		"cp-main-vote-no-conflict": message.NewVoteMessage(signed(vote.NewCPMainVote(blockHash, height, round, 1,
			vote.CPValueNo, &vote.JustMainVoteNoConflict{QCert: cert}, addr))),
		"cp-main-vote-conflict": message.NewVoteMessage(signed(vote.NewCPMainVote(blockHash, height, round, 1,
			vote.CPValueAbstain, conflict, addr))),
		"cp-decided": message.NewVoteMessage(signed(vote.NewCPDecidedVote(blockHash, height, round, 1,
			vote.CPValueYes, &vote.JustDecided{QCert: cert}, addr))),
		"block-announce": message.NewBlockAnnounceMessage(blk, blkCert, cert),
		"blocks-request": message.NewBlocksRequestMessage(3, height, 10),
		"blocks-response": message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks, "more", 3, height,
			[][]byte{ts.RandBytes(50), ts.RandBytes(70)}, cert),
	}
}

func mutate(rng *rand.Rand, ts *testsuite.TestSuite, data []byte, kind int) []byte {
	mutated := append([]byte(nil), data...)

	switch kind % 4 {
	case 0: // Change one to three bytes
		for i := 0; i <= rng.Intn(3); i++ {
			mutated[rng.Intn(len(mutated))] ^= byte(1 + rng.Intn(255))
		}
	case 1: // Cut the end
		mutated = mutated[:rng.Intn(len(mutated))]
	case 2: // Use special values
		mutated[rng.Intn(len(mutated))] = []byte{0x00, 0xff, 0x80, 0x1b}[rng.Intn(4)]
	case 3: // Insert some bytes
		pos := rng.Intn(len(mutated))
		inserted := ts.RandBytes(1 + rng.Intn(8))
		mutated = append(mutated[:pos], append(inserted, mutated[pos:]...)...)
	}

	return mutated
}

// The data of the messages come from any peer of the network.
// Decoding, checking, or printing a corrupted message should never panic.
func TestDecodeMutatedBundlesDoesNotPanic(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	rng := rand.New(rand.NewSource(42))

	const mutationsPerMessage = 400

	for name, msg := range makeAllMessages(ts) {
		data, err := NewBundle(msg).Encode()
		require.NoError(t, err, name)

		for i := 0; i < mutationsPerMessage; i++ {
			mutated := mutate(rng, ts, data, i)

			require.NotPanics(t, func() {
				bdl := new(Bundle)
				if _, err := bdl.Decode(bytes.NewReader(mutated)); err == nil {
					_ = bdl.BasicCheck()
					_ = bdl.LogString()
				}
			}, "message %s, input: %x", name, mutated)
		}
	}
}
