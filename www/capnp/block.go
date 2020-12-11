package capnp

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/store"
)

type factory struct {
	store  store.StoreReader
	logger *logger.Logger
}

func (f factory) BlockToBlockInfo(block *block.Block, height int, cbi *BlockInfo) error {
	data, _ := block.Encode()

	cb, _ := cbi.NewBlock()
	ch, _ := cb.NewHeader()
	ctxs, _ := cb.NewTxs()
	clc, _ := cb.NewLastCommit()
	cbi.SetHash(block.Hash().RawBytes())
	cbi.SetHeight(uint32(height))
	cbi.SetData(data)
	// last commit
	if block.LastCommit() != nil {
		clc.SetRound(uint32(block.LastCommit().Round()))
		clc.SetSignature(block.LastCommit().Signature().RawBytes())
		clcc, _ := clc.NewCommitters(int32(len(block.LastCommit().Committers())))
		for i, committer := range block.LastCommit().Committers() {
			c := clcc.At(i)
			c.SetAddress(committer.Address.RawBytes())
			c.SetStatus(int32(committer.Status))
		}
	}
	// header
	ch.SetVersion(int32(block.Header().Version()))
	ch.SetTime(block.Header().Time().Unix())
	ch.SetTxsHash(block.Header().TxsHash().RawBytes())
	ch.SetStateHash(block.Header().StateHash().RawBytes())
	ch.SetCommittersHash(block.Header().CommittersHash().RawBytes())
	ch.SetLastBlockHash(block.Header().LastBlockHash().RawBytes())
	ch.SetLastReceiptsHash(block.Header().LastReceiptsHash().RawBytes())
	ch.SetLastCommitHash(block.Header().LastCommitHash().RawBytes())
	ch.SetProposerAddress(block.Header().ProposerAddress().RawBytes())
	// Transactions
	ctxHashes, _ := ctxs.NewHashes(int32(block.TxHashes().Count()))
	for i, hash := range block.TxHashes().Hashes() {
		ctxHashes.Set(i, hash.RawBytes())
	}

	return nil
}

func (f factory) BlockAt(b ZarbServer_blockAt) error {
	defer func() {
		if r := recover(); r != nil {
			f.logger.Error("Block method recovered from a panic", "r", r)
		}
	}()

	height := b.Params.Height()
	block, err := f.store.BlockByHeight(int(height))
	if err != nil {
		f.logger.Error("Error on retriving block", "height", height, "err", err)
		return err
	}

	bi, _ := b.Results.NewBlockInfo()

	return f.BlockToBlockInfo(block, int(height), &bi)
}

func (f factory) Block(b ZarbServer_block) error {
	defer func() {
		if r := recover(); r != nil {
			f.logger.Error("Block method recovered from a panic", "r", r)
		}
	}()

	h, _ := b.Params.Hash()
	hash, err := crypto.HashFromRawBytes(h)
	if err != nil {
		f.logger.Error("Error on retriving block", "hash", h, "err", err)
		return err
	}
	block, height, err := f.store.BlockByHash(hash)
	if err != nil {
		f.logger.Error("Error on retriving block", "hash", h, "err", err)
		return err
	}
	bi, _ := b.Results.NewBlockInfo()

	return f.BlockToBlockInfo(block, height, &bi)
}
