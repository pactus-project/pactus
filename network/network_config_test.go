package network

import (
	"path/filepath"
	"testing"

	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The network is created, started and stopped with each optional service.
func TestNetworkWithOptionalServices(t *testing.T) {
	tests := []struct {
		name   string
		modify func(conf *Config)
	}{
		{"metrics", func(conf *Config) { conf.EnableMetrics = true }},
		{"NAT service", func(conf *Config) { conf.EnableNATService = true }},
		{"UPnP", func(conf *Config) { conf.EnableUPnP = true }},
		{"no UDP", func(conf *Config) { conf.EnableUDP = false }},
		{"relay service", func(conf *Config) { conf.EnableRelayService = true }},
		{"bootstrapper", func(conf *Config) { conf.IsBootstrapper = true }},
		{"public address", func(conf *Config) { conf.PublicAddrString = "/ip4/8.8.8.8/tcp/21888" }},
		{"public network", func(conf *Config) {
			conf.ForcePrivateNetwork = false
			conf.PublicAddrString = "/ip4/8.8.8.8/tcp/21888"
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := testConfig()
			tt.modify(conf)

			net := makeTestNetwork(t, conf, nil)
			assert.NotEmpty(t, net.SelfID())
			assert.Equal(t, conf.NetworkName, net.Name())

			net.Stop()
		})
	}
}

// A network key that can't be read or created is an error, and nothing is left running.
func TestNetworkWithInvalidKeyPath(t *testing.T) {
	file := util.TempFilePath()
	require.NoError(t, util.WriteFile(file, []byte("not a directory")))

	conf := testConfig()
	conf.NetworkKey = filepath.Join(file, "network_key")

	pipe := pipeline.New[Event](t.Context())
	log := logger.NewSubLogger(t.Context(), "_network", nil)
	net, err := makeNetwork(t.Context(), conf, log, pipe, nil)

	var libp2pErr LibP2PError
	require.ErrorAs(t, err, &libp2pErr)
	assert.Nil(t, net)
}
