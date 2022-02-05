package network

import (
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T, conf1 *Config, conf2 *Config) (*network, *network) {
	netName := fmt.Sprintf("net_%v", util.RandInt(0))
	conf1.Name = netName
	conf2.Name = netName

	net1, err := NewNetwork(conf1)
	assert.NoError(t, err)

	net2, err := NewNetwork(conf2)
	assert.NoError(t, err)

	return net1.(*network), net2.(*network)
}

func TestStoppingNetwork(t *testing.T) {
	net1, net2 := setup(t, TestConfig(), TestConfig())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	net1.Stop()
	net2.Stop()
}

func TestDHT(t *testing.T) {
	conf1 := TestConfig()
	conf2 := TestConfig()

	nodeKeyPath := util.TempFilePath()
	nodeKey, _ := loadOrCreateKey(nodeKeyPath)
	pid, _ := peer.IDFromPrivateKey(nodeKey)
	conf1.NodeKeyFile = nodeKeyPath

	conf1.EnableMdns = false
	conf2.EnableMdns = false
	conf1.ListenAddress = []string{"/ip4/0.0.0.0/tcp/1347"}
	conf2.Bootstrap.Addresses = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", pid)}

	net1, net2 := setup(t, conf1, conf2)

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

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
