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

func getBlockAt(t *testing.T, height int) *block.Block {
	for i := 0; i < 10; i++ {
		res := tCapnpServer.GetBlock(tCtx, func(p capnp.ZarbServer_getBlock_Params) error {
			p.SetHeight(uint64(height))
			p.SetVerbosity(0)
			return nil
		}).Result()

		st, err := res.Struct()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
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

func TestGeneratingBlocks(t *testing.T) {
	res := getBlockAt(t, 1)
	assert.Contains(t, res, "0000000000000000000000000000000000000000000000000000000000000000")
	fmt.Println(res)
}
