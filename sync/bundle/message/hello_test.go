package message

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloType(t *testing.T) {
	msg := &HelloMessage{}
	assert.Equal(t, TypeHello, msg.Type())
}

func TestHelloMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid signature", func(t *testing.T) {
		valKey := ts.RandValKey()
		msg := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		msg.Sign([]*bls.ValidatorKey{valKey})
		msg.Signature = ts.RandBLSSignature()

		err := msg.BasicCheck()
		require.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Signature is nil", func(t *testing.T) {
		valKey := ts.RandValKey()
		msg := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		msg.Sign([]*bls.ValidatorKey{valKey})
		msg.Signature = nil

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{"no signature"})
	})

	t.Run("PublicKeys are empty", func(t *testing.T) {
		valKey := ts.RandValKey()
		msg := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		msg.Sign([]*bls.ValidatorKey{valKey})
		msg.PublicKeys = make([]*bls.PublicKey, 0)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{"no public key"})
	})

	t.Run("Invalid PublicKey", func(t *testing.T) {
		valKey := ts.RandValKey()
		msg := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		msg.Sign([]*bls.ValidatorKey{valKey})
		msg.PublicKeys = []*bls.PublicKey{{}}

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{"short buffer"})
	})

	t.Run("Check hello message time", func(t *testing.T) {
		time1 := time.Now()
		msg := NewHelloMessage(ts.RandPeerID(), "Alice", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		time2 := time.Now()

		assert.GreaterOrEqual(t, msg.MyTime().UnixMilli(), time1.UnixMilli())
		assert.LessOrEqual(t, msg.MyTime().UnixMilli(), time2.UnixMilli())
	})

	t.Run("Ok", func(t *testing.T) {
		valKey := ts.RandValKey()
		msg := NewHelloMessage(ts.RandPeerID(), "Alice", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		msg.Sign([]*bls.ValidatorKey{valKey})

		require.NoError(t, msg.BasicCheck())
		assert.Zero(t, msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "Alice")
		assert.Contains(t, msg.LogString(), "FULL")
	})
}
