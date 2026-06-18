package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockAnnounceType(t *testing.T) {
	smg := &BlockAnnounceMessage{}
	assert.Equal(t, TypeBlockAnnounce, smg.Type())
}

func TestBlockAnnounceMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Block is not set", func(t *testing.T) {
		cert := certificate.NewCertificate(0, 0)
		msg := NewBlockAnnounceMessage(nil, cert, nil)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{
			Reason: "block is not set",
		})
		assert.Zero(t, msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "nil")
	})

	t.Run("Certificate is not set", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		msg := NewBlockAnnounceMessage(blk, nil, nil)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{
			Reason: "certificate is not set",
		})
		assert.Zero(t, msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "0")
	})

	t.Run("Invalid block", func(t *testing.T) {
		blk, cert := ts.GenerateTestBlock(ts.RandHeight(),
			testsuite.BlockWithPrevCert(nil))
		msg := NewBlockAnnounceMessage(blk, cert, nil)
		err := msg.BasicCheck()

		require.ErrorIs(t, err, block.BasicCheckError{
			Reason: "invalid genesis block hash",
		})
	})

	t.Run("Invalid block certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		cert := certificate.NewCertificate(0, 0)
		msg := NewBlockAnnounceMessage(blk, cert, nil)
		err := msg.BasicCheck()

		require.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
		assert.Zero(t, msg.ConsensusHeight())
	})

	t.Run("Invalid proof", func(t *testing.T) {
		height := ts.RandHeight()
		blk, cert := ts.GenerateTestBlock(height)
		proof := ts.GenerateTestCertificate(0)
		msg := NewBlockAnnounceMessage(blk, cert, proof)
		err := msg.BasicCheck()

		require.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("OK", func(t *testing.T) {
		blk, cert := ts.GenerateTestBlock(ts.RandHeight())
		msg := NewBlockAnnounceMessage(blk, cert, nil)

		require.NoError(t, msg.BasicCheck())
		assert.Equal(t, cert.Height(), msg.Height())
		assert.Equal(t, cert.Height(), msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), blk.Hash().LogString())
	})
}
