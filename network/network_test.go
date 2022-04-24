package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
)

var (
	tNetworksSize   int
	tNetworks       []*network
	tBootstrapAddrs []string
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
			Addresses:    tBootstrapAddrs,
			MinThreshold: 4,
			MaxThreshold: 8,
			Period:       1 * time.Second,
		},
	}
}

func init() {
	tNetworksSize = 8
	tNetworks = make([]*network, tNetworksSize)

	port := util.RandInt32(9999) + 10000

	for i := 0; i < tNetworksSize; i++ {
		conf := testConfig()
		if i == 0 {
			// bootstrap node
			conf.ListenAddress = []string{
				fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", port),
				fmt.Sprintf("/ip6/::/tcp/%v", port),
			}
		}

		net, err := NewNetwork(conf)
		if err != nil {
			panic(err)
		}
		err = net.Start()
		if err != nil {
			panic(err)
		}
		err = net.JoinGeneralTopic()
		if err != nil {
			panic(err)
		}

		if i == 0 {
			tBootstrapAddrs = []string{
				fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", port, net.SelfID().String()),
				fmt.Sprintf("/ip6/::1/tcp/%v/p2p/%v", port, net.SelfID().String()),
			}
		}

		tNetworks[i] = net.(*network)
		time.Sleep(100 * time.Millisecond)

		fmt.Printf("peer %v id:%v\n", i, net.SelfID().String())
	}

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("Peers are connected")
}

func shouldReceiveEvent(t *testing.T, net *network) Event {
	timeout := time.NewTimer(2 * time.Second)

	for {
		net.logger.Debug("network connections", "NumConnectedPeers", net.NumConnectedPeers())
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldReceiveEvent Timeout, test: %v id:%s", t.Name(), net.SelfID().String()))
			return nil
		case e := <-net.EventChannel():
			net.logger.Debug("received an event", "event", e, "id", net.SelfID().String())
			return e
		}
	}
}

func TestStoppingNetwork(t *testing.T) {
	net, err := NewNetwork(testConfig())
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic())

	// Should stop peacefully
	net.Stop()
}

func TestDHT(t *testing.T) {
	conf := testConfig()
	conf.EnableMdns = false
	net, err := NewNetwork(conf)
	assert.NoError(t, err)

	assert.NoError(t, net.Start())
	assert.NoError(t, net.JoinGeneralTopic())

	net, err = NewNetwork(conf)
	assert.NoError(t, err)

	assert.NoError(t, net.Start())

	for {
		if net.NumConnectedPeers() > 0 {
			break
		}
	}

	net.Stop()
}
