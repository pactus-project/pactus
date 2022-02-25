package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

var (
	tSize     int
	tNetworks []*network
)

func init() {
	logger.InitLogger(logger.TestConfig())
	tSize = 5

	tNetworks = make([]*network, 5)

	util.TempFilePath()

	for i := 0; i < tSize; i++ {
		conf := TestConfig()

		if i == 0 {
			// bootstrap node
			conf.ListenAddress = []string{"/ip4/0.0.0.0/tcp/1347"}
			conf.NodeKeyFile = util.TempFilePath()
		} else {
			conf.Bootstrap.Addresses = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", tNetworks[0].SelfID().String())}
		}

		net, _ := NewNetwork(conf)
		if err := net.Start(); err != nil {
			panic(err)
		}

		if err := net.JoinGeneralTopic(); err != nil {
			panic(err)
		}

		tNetworks[i] = net.(*network)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Peers are connected")
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
	conf.Bootstrap.Addresses = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/1347/p2p/%s", tNetworks[0].SelfID())}

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
// func TestDisconnecting(t *testing.T) {
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
