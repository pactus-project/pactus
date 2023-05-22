package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMDNS(t *testing.T) {
	conf1 := testConfig()
	conf1.EnableMdns = true
	net1, err := NewNetwork(conf1)
	assert.NoError(t, err)

	conf2 := testConfig()
	conf2.EnableMdns = true
	net2, err := NewNetwork(conf2)
	assert.NoError(t, err)

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	assert.NoError(t, net1.JoinGeneralTopic())
	assert.NoError(t, net2.JoinGeneralTopic())

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, net1.SendTo([]byte("test"), net2.SelfID()))
}
