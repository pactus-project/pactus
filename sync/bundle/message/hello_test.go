package message

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestHelloType(t *testing.T) {
	m := &HelloMessage{}
	assert.Equal(t, TypeHello, m.Type())
}

func TestHelloMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid signature", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.Signature = ts.RandBLSSignature()

		assert.ErrorIs(t, crypto.ErrInvalidSignature, m.BasicCheck())
	})

	t.Run("Signature is nil", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.Signature = nil

		assert.Equal(t, errors.ErrInvalidSignature, errors.Code(m.BasicCheck()))
	})

	t.Run("PublicKeys are empty", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.PublicKeys = make([]*bls.PublicKey, 0)

		assert.Equal(t, errors.ErrInvalidPublicKey, errors.Code(m.BasicCheck()))
	})

	t.Run("MyTimeUnixMilli of time1 is less or equal than hello message time", func(t *testing.T) {
		time1 := time.Now()
		myTimeUnixMilli := time1.UnixMilli()

		m := NewHelloMessage(ts.RandPeerID(), "Alice", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())

		assert.LessOrEqual(t, m.MyTimeUnixMilli, time.Now().UnixMilli())
		assert.GreaterOrEqual(t, m.MyTimeUnixMilli, myTimeUnixMilli)
	})

	t.Run("Ok", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Alice", service.New(service.FullNode),
			ts.RandHeight(), ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "Alice")
		assert.Contains(t, m.String(), "FULL")
	})
}
