package tests

import (
	"time"

	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func lastHash() hash.Hash {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
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

func lastHeight() int32 {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		panic(err)
	}

	return st.LastBlockHeight()
}

func waitForNewBlocks(num int32) {
	height := lastHeight() + num
	for i := int32(0); i < num; i++ {
		if lastHeight() > height {
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func lastBlock() *capnp.BlockResult {
	return getBlockAt(lastHeight())
}

func getBlockAt(height int32) *capnp.BlockResult {
	for i := 0; i < 120; i++ {
		hashRes, _ := tCapnpServer.GetBlockHash(tCtx, func(p capnp.ZarbServer_getBlockHash_Params) error {
			p.SetHeight(height)
			return nil
		}).Struct()

		blockRes := tCapnpServer.GetBlock(tCtx, func(p capnp.ZarbServer_getBlock_Params) error {
			data, _ := hashRes.Result()
			p.SetVerbosity(0)
			return p.SetHash(data)
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
