package tests

import (
	"time"

	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
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
	h, _ := hash.FromRawBytes(data)
	return h
}

func lastHeight() int {
	res := tCapnpServer.GetBlockchainInfo(tCtx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		panic(err)
	}

	return int(st.LastBlockHeight())
}

func waitForNewBlock() {
	getBlockAt(lastHeight() + 1)
}

func lastBlock() *capnp.BlockResult {
	return getBlockAt(lastHeight())
}

func getBlockAt(height int) *capnp.BlockResult {
	for i := 0; i < 120; i++ {
		hashRes, _ := tCapnpServer.GetBlockHash(tCtx, func(p capnp.ZarbServer_getBlockHash_Params) error {
			p.SetHeight(uint64(height))
			return nil
		}).Struct()

		blockRes := tCapnpServer.GetBlock(tCtx, func(p capnp.ZarbServer_getBlock_Params) error {
			data, _ := hashRes.Result()
			p.SetHash(data)
			p.SetVerbosity(0)
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
