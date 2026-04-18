package tests

import (
	"fmt"
	"time"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
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
		res, err := tBlockchainClient.GetBlock(tCtx,
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
