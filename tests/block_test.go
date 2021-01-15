package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func lastHeight(t *testing.T) int {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	assert.NoError(t, err)

	return int(st.Height())
}

func waitForNewBlock(t *testing.T) *block.Block {
	return getBlockAt(t, lastHeight(t)+1)
}

func lastBlock(t *testing.T) *block.Block {
	return getBlockAt(t, lastHeight(t))
}

func getBlockAt(t *testing.T, height int) *block.Block {
	for i := 0; i < 10; i++ {
		res := tCapnpServer.GetBlock(tCtx, func(p capnp.ZarbServer_getBlock_Params) error {
			p.SetHeight(uint64(height))
			p.SetVerbosity(0)
			return nil
		}).Result()

		st, err := res.Struct()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		d, _ := st.Data()
		b := new(block.Block)
		assert.NoError(t, b.Decode(d))
		return b
	}
	require.NoError(t, fmt.Errorf("timeout"))
	return nil
}
