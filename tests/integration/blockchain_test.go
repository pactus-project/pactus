package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lastHeight() types.Height {
	res, err := tBlockchainClient.GetBlockchainInfo(tCtx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		panic(err)
	}

	return types.Height(res.LastBlockHeight)
}

func waitForNewBlocks(num uint32) {
	for i := uint32(0); i < num; i++ {
		height := lastHeight()
		if lastHeight() > height {
			break
		}
		time.Sleep(4 * time.Second)
	}
}

func lastBlock() *pactus.GetBlockResponse {
	return getBlockAt(lastHeight())
}

func getBlockAt(height types.Height) *pactus.GetBlockResponse {
	for i := 0; i < 120; i++ {
		res, err := tBlockchainClient.GetBlock(
			tCtx,
			&pactus.GetBlockRequest{
				Height:    uint32(height),
				Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_INFO,
			},
		)
		if err != nil {
			fmt.Printf("getBlockAt err: %s\n", err.Error())
			time.Sleep(1 * time.Second)

			continue
		}

		return res
	}
	logger.Panic("getBlockAt timeout", "height", height)

	return nil
}

func TestChainInfo(t *testing.T) {
	res, err := tBlockchainClient.GetBlockchainInfo(t.Context(), &pactus.GetBlockchainInfoRequest{})
	require.NoError(t, err)

	assert.Greater(t, res.LastBlockHeight, uint32(8))
}

func TestCommitteeInfo(t *testing.T) {
	res, err := tBlockchainClient.GetCommitteeInfo(t.Context(), &pactus.GetCommitteeInfoRequest{})
	require.NoError(t, err)

	assert.Greater(t, uint32(4), res.CommitteePower)
}
