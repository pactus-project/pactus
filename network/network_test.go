package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
)

var (
	tConfig1  *Config
	tConfig2  *Config
	tNetwork1 *network
	tNetwork2 *network
)

func init() {
	logger.InitLogger(logger.TestConfig())

	tConfig1 = TestConfig()
	tConfig2 = TestConfig()
	tConfig1.ListenAddress = []string{"/ip4/0.0.0.0/tcp/1347"}

	net1, _ := NewNetwork(tConfig1)
	net2, _ := NewNetwork(tConfig2)

	tNetwork1 = net1.(*network)
	tNetwork2 = net2.(*network)

	err := tNetwork1.Start()
	if err != nil {
		panic(err)
	}
	err = tNetwork2.Start()
	if err != nil {
		panic(err)
	}

	err = tNetwork1.JoinGeneralTopic()
	if err != nil {
		panic(err)
	}
	err = tNetwork2.JoinGeneralTopic()
	if err != nil {
		panic(err)
	}

	for {
		if tNetwork1.NumConnectedPeers() > 0 && tNetwork2.NumConnectedPeers() > 0 {
			break
		}
	}

	time.Sleep(1 * time.Second)
}

func shouldReceiveEvent(t *testing.T, net *network) Event {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldReceiveEvent Timeout, test: %v", t.Name()))
			return nil
		case e := <-net.EventChannel():
			return e
		}
	}
}

func TestStoppingNetwork(t *testing.T) {
	net, err := NewNetwork(TestConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	// Should stop without error
	net.Stop()
}

func TestDHT(t *testing.T) {
	conf := TestConfig()
	conf.EnableMdns = false
	conf.Bootstrap.Addresses = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", tNetwork1.SelfID())}

	net, err := NewNetwork(TestConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())

	for {
		if net.NumConnectedPeers() > 0 {
			break
		}
	}

	net.Stop()
}

// TODO: Fix me
// func TestDisconrecting(t *testing.T) {
// 	net1, net2 := setup(t, TestConfig(), TestConfig())

// 	assert.NoError(t, net1.Start())
// 	assert.NoError(t, net2.Start())

// 	for {
// 		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
// 			break
// 		}
// 	}

// 	net1.CloseConnection(net2.SelfID())
// 	assert.Equal(t, net1.NumConnectedPeers(), 0)

// 	net1.Stop()
// 	net2.Stop()
// }
