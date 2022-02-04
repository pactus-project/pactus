package network

import (
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

var (
	tConf1 *Config
	tConf2 *Config
	tNet1  *network
	tNet2  *network
)

func init() {
	logger.InitLogger(logger.TestConfig())
	tConf1 = TestConfig()
	tConf2 = TestConfig()

	netName := fmt.Sprintf("net_%v", util.RandInt(0))
	tConf1.Name = netName
	tConf2.Name = netName
}

func setup(t *testing.T) {
	net1, err := NewNetwork(tConf1)
	assert.NoError(t, err)

	net2, err := NewNetwork(tConf2)
	assert.NoError(t, err)

	tNet1 = net1.(*network)
	tNet2 = net2.(*network)
}

func TestStoppingNetwork(t *testing.T) {
	setup(t)

	assert.NoError(t, tNet1.Start())
	assert.NoError(t, tNet2.Start())

	tNet1.Stop()
	tNet2.Stop()
}

func TestDHT(t *testing.T) {
	nodeKeyPath := util.TempFilePath()
	nodeKey, _ := loadOrCreateKey(nodeKeyPath)
	pid, _ := peer.IDFromPrivateKey(nodeKey)
	tConf1.NodeKeyFile = nodeKeyPath

	tConf1.EnableMDNS = false
	tConf2.EnableMDNS = false
	tConf1.ListenAddress = []string{"/ip4/0.0.0.0/tcp/1347"}
	tConf2.Bootstrap.Addresses = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", pid)}

	setup(t)

	assert.NoError(t, tNet1.Start())
	assert.NoError(t, tNet2.Start())

	for {
		if tNet1.NumConnectedPeers() > 0 && tNet2.NumConnectedPeers() > 0 {
			break
		}
	}

	tNet1.Stop()
	tNet2.Stop()
}

// TODO: Fix me
// func TestDisconrecting(t *testing.T) {
// 	setup(t)

// 	assert.NoError(t, tNet1.Start())
// 	assert.NoError(t, tNet2.Start())

// 	for {
// 		if tNet1.NumConnectedPeers() > 0 && tNet2.NumConnectedPeers() > 0 {
// 			break
// 		}
// 	}

// 	tNet1.CloseConnection(tNet2.SelfID())
// 	assert.Equal(t, tNet1.NumConnectedPeers(), 0)

// 	tNet1.Stop()
// 	tNet2.Stop()
// }
