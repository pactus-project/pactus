package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func lastHeight() int {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		panic(err)
	}

	return int(st.Height())
}

func waitForNewBlock() *block.Block {
	return getBlockAt(lastHeight() + 1)
}

func lastBlock() *block.Block {
	return getBlockAt(lastHeight())
}

func getBlockAt(height int) *block.Block {
	for i := 0; i < 120; i++ {
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
		err = b.Decode(d)
		if err != nil {
			panic(err)
		}
		return b
	}
	panic("get block timeout")
}

func TestGetBlock(t *testing.T) {
	require.NotNil(t, lastBlock())
}
