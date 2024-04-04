package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMDNS(t *testing.T) {
	conf1 := testConfig()
	conf1.ListenAddrStrings = []string{
		"/ip6/::1/tcp/0", "/ip6/::1/udp/0/quic-v1",
		"/ip4/127.0.0.1/tcp/0", "/ip4/127.0.0.1/udp/0/quic-v1",
	}
	conf1.EnableMdns = true
	net1 := makeTestNetwork(t, conf1, nil)

	conf2 := testConfig()
	conf2.ListenAddrStrings = []string{
		"/ip6/::1/tcp/0", "/ip6/::1/udp/0/quic-v1",
		"/ip4/127.0.0.1/tcp/0", "/ip4/127.0.0.1/udp/0/quic-v1",
	}
	conf2.EnableMdns = true
	net2 := makeTestNetwork(t, conf2, nil)

	assert.NoError(t, net1.Start())
	time.Sleep(250 * time.Millisecond)

	assert.NoError(t, net2.Start())
	time.Sleep(250 * time.Millisecond)

	msg := []byte("test-mdns")
	assert.NoError(t, net1.SendTo(msg, net2.SelfID()))

	se := shouldReceiveEvent(t, net2, EventTypeStream).(*StreamMessage)
	assert.Equal(t, se.From, net1.SelfID())
	assert.Equal(t, readData(t, se.Reader, len(msg)), msg)

	net1.Stop()
	net2.Stop()
}
