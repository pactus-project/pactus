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
	if err := cbi.SetHash(block.Hash().RawBytes()); err != nil {
		return err
	}
	cbi.SetHeight(uint32(height))
	if err := cbi.SetData(data); err != nil {
		return err
	}
	// last commit
	if block.LastCommit() != nil {
		clc.SetRound(uint32(block.LastCommit().Round()))
		if err := clc.SetSignature(block.LastCommit().Signature().RawBytes()); err != nil {
			return err
		}
		clcc, _ := clc.NewCommitters(int32(len(block.LastCommit().Committers())))
		for i, commiter := range block.LastCommit().Committers() {
			c := clcc.At(i)
			if err := c.SetAddress(commiter.Address.RawBytes()); err != nil {
				return err
			}
			c.SetStatus(int32(commiter.Status))
		}
	}
	// header
	ch.SetVersion(int32(block.Header().Version()))
	ch.SetTime(block.Header().Time().Unix())
	if err := ch.SetTxsHash(block.Header().TxsHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetStateHash(block.Header().StateHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetCommittersHash(block.Header().CommittersHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetLastBlockHash(block.Header().LastBlockHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetLastReceiptsHash(block.Header().LastReceiptsHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetLastCommitHash(block.Header().LastCommitHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetProposerAddress(block.Header().ProposerAddress().RawBytes()); err != nil {
		return err
	}
	// Transactions
	ctxHashes, _ := ctxs.NewHashes(int32(block.TxHashes().Count()))
	for i, hash := range block.TxHashes().Hashes() {
		if err := ctxHashes.Set(i, hash.RawBytes()); err != nil {
			return err
		}
	}

	return nil
}

func (f factory) BlockAt(b ZarbServer_blockAt) error {
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
