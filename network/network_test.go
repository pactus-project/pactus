package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
)

func TestStoppingNetwork(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	net, err := NewNetwork(TestConfig())
	assert.NoError(t, err)
	assert.NoError(t, net.Start())
	net.Stop()
}

func TestDHT(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	conf1 := TestConfig()
	conf1.EnableMDNS = false
	conf1.ListenAddress = []string{"/ip4/0.0.0.0/tcp/1347"}
	net1, err := NewNetwork(conf1)
	assert.NoError(t, err)
	assert.NoError(t, net1.Start())

	conf2 := TestConfig()
	conf2.EnableMDNS = false
	conf2.Bootstrap.Addresses = append(conf2.Bootstrap.Addresses, fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", net1.SelfID()))
	net2, err := NewNetwork(conf2)
	require.NoError(t, err)
	require.NoError(t, net2.Start())

	for {
		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
			break
		}
	}

	net1.Stop()
	net2.Stop()
}

// TODO: Fix me

// func TestDisconrecting(t *testing.T) {
// 	logger.InitLogger(logger.TestConfig())

// 	net1, err := NewNetwork(TestConfig())
// 	assert.NoError(t, err)
// 	assert.NoError(t, net1.Start())

// 	net2, err := NewNetwork(TestConfig())
// 	require.NoError(t, err)
// 	require.NoError(t, net2.Start())

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
