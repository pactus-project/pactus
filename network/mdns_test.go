package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMDNS(t *testing.T) {
	conf1 := testConfig()
	conf1.Listens = []string{"/ip4/127.0.0.1/tcp/0"}
	conf1.EnableMdns = true
	net1, _ := newNetwork("test", conf1, nil)

	conf2 := testConfig()
	conf2.Listens = []string{"/ip4/127.0.0.1/tcp/0"}
	conf2.EnableMdns = true
	net2, _ := newNetwork("test", conf2, nil)

	assert.NoError(t, net1.Start())
	time.Sleep(250 * time.Millisecond)

	assert.NoError(t, net2.Start())
	time.Sleep(250 * time.Millisecond)

	msg := []byte("test-mdns")
	assert.NoError(t, net1.SendTo(msg, net2.SelfID()))
	e1 := shouldReceiveEvent(t, net2, EventTypeConnect).(*ConnectEvent)
	assert.Equal(t, e1.PeerID, net1.SelfID())
	e2 := shouldReceiveEvent(t, net2, EventTypeStream).(*StreamMessage)
	assert.Equal(t, e2.Source, net1.SelfID())
	assert.Equal(t, readData(t, e2.Reader, len(msg)), msg)
}
