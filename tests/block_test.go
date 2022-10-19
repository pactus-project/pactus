package tests

import (
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/www/capnp"
)

func lastHash() hash.Hash {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.PactusServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		panic(err)
	}

	data, _ := st.LastBlockHash()
	h, _ := hash.FromBytes(data)
	return h
}

func lastHeight() uint32 {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.PactusServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		panic(err)
	}

	return st.LastBlockHeight()
}

func waitForNewBlocks(num uint32) {
	height := lastHeight() + num
	for i := uint32(0); i < num; i++ {
		if lastHeight() > height {
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func lastBlock() *capnp.BlockResult {
	return getBlockAt(lastHeight())
}

func getBlockAt(height uint32) *capnp.BlockResult {
	for i := 0; i < 120; i++ {
		blockRes := tCapnpServer.GetBlock(tCtx, func(p capnp.PactusServer_getBlock_Params) error {
			p.SetVerbosity(0)
			p.SetHeight(height)
			return nil
		}).Result()

		st, err := blockRes.Struct()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return &st
	}
	logger.Panic("get block timeout", "height", height)
	return nil
}
