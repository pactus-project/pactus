package state

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitBlockRefusesAFutureHeight(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	blk, _ := td.makeBlockAndCertificate(t, 0)
	cert := td.GenerateTestCertificate(td.state.LastBlockHeight() + 2)

	err := td.state.CommitBlock(blk, cert)
	require.ErrorIs(t, err, FutureBlockError{Height: cert.Height(), LastHeight: td.state.LastBlockHeight()})
	assert.Zero(t, td.state.LastBlockHeight())
}

// The aggregated signature doesn't cover the absentees: a peer can change them.
func TestCommitBlockRefusesTamperedAbsentees(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	blk, cert := td.makeBlockAndCertificate(t, 0)

	// The signature is still valid, but the absentees are not in the list of committers.
	tampered := certificate.NewCertificate(cert.Height(), cert.Round())
	tampered.SetSignature(cert.Committers(), []int32{3, 3, 1000}, cert.Signature())

	require.Error(t, td.state.CommitBlock(blk, tampered))
	assert.Zero(t, td.state.LastBlockHeight())

	require.NoError(t, td.state.CommitBlock(blk, cert))
}

func TestScoreStartHeight(t *testing.T) {
	tests := []struct {
		end      types.Height
		expected types.Height
	}{
		{0, 2},
		{1, 2},
		{59_999, 2},
		{60_000, 2},
		{60_001, 2}, // The first block has no certificate.
		{60_002, 2},
		{60_003, 3},
		{100_000, 40_000},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, scoreStartHeight(tt.end, 60_000), "end height %d", tt.end)
	}
}
