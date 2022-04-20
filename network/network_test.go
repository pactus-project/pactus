package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
)

func testConfig() *Config {
	return &Config{
		Name:             "test-network",
		ListenAddress:    []string{"/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"},
		NodeKeyFile:      util.TempFilePath(),
		EnableNATService: false,
		EnableRelay:      false,
		EnableMdns:       true,
		EnableKademlia:   true,
		EnablePing:       false,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 4,
			MaxThreshold: 8,
			Period:       1 * time.Second,
		},
	}
}

func setup(t *testing.T, size int) []*network {
	nets := make([]*network, size)

	networkName := fmt.Sprintf("test-network-%d", util.RandInt32(10000))
	port := util.RandInt32(9999) + 10000

	for i := 0; i < size; i++ {
		conf := testConfig()
		conf.Name = networkName

		bootstrapAddr := ""
		if i == 0 {
			// bootstrap node
			conf.ListenAddress = []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)}
		} else {
			bootstrapAddr = fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/p2p/%s", port, nets[0].SelfID().String())
			conf.Bootstrap.Addresses = []string{bootstrapAddr}
		}

		net, err := NewNetwork(conf)
		require.NoError(t, err)
		require.NoError(t, net.Start())
		require.NoError(t, net.JoinGeneralTopic())

		nets[i] = net.(*network)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Peers are connected")
	return nets
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
	size := 2
	nets := setup(t, size)

	for i := 0; i < size; i++ {
		// Should stop without any error
		nets[i].Stop()
	}
}

func TestDHT(t *testing.T) {
	nets := setup(t, 4)
	conf := nets[1].config
	conf.EnableMdns = false

	net, err := NewNetwork(conf)
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
// 	nets := setup(t, 2)

// 	assert.NoError(t, nets[0].Start())
// 	assert.NoError(t, nets[1].Start())

// 	for {
// 		if nets[0].NumConnectedPeers() > 0 && nets[1].NumConnectedPeers() > 0 {
// 			break
// 		}
// 	}

// 	nets[0].CloseConnection(nets[1].SelfID())
// 	assert.Equal(t, nets[0].NumConnectedPeers(), 0)

// 	nets[0].Stop()
// 	nets[1].Stop()
// }
