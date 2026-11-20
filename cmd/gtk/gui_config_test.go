//go:build gtk

package main

import (
	"testing"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/stretchr/testify/assert"
)

// The GUI doesn't change the gRPC server of the config: it could be public, without the wallet.
func TestGUIConfigKeepsTheGRPCServer(t *testing.T) {
	cfg := config.DefaultConfigForChain(genesis.Mainnet)
	cfg.GRPC.Enable = true
	cfg.GRPC.Listen = "0.0.0.0:50051"
	cfg.GRPC.EnableWallet = false
	want := *cfg.GRPC

	got := guiConfig(cfg)
	assert.Equal(t, want, *got.GRPC)
}
