package capnp

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
)

type factory struct {
	state  state.StateReader
	store  store.StoreReader
	txPool txpool.TxPool
	sync   sync.Synchronizer
	logger *logger.Logger
}

func (f factory) GetBlockHeight(args ZarbServer_getBlockHeight) error {
	s, _ := args.Params.Hash()
	h, err := crypto.HashFromString(string(s))
	if err != nil {
		return err
	}
	num, err := f.store.BlockHeight(h)
	if err != nil {
		f.logger.Debug("Error on retriving block number", "err", err)
		return err
	}
	args.Results.SetResult(uint64(num))
	return nil
}

func (f factory) GetBlock(args ZarbServer_getBlock) error {
	h := args.Params.Height()
	v := args.Params.Verbosity()
	b, err := f.store.Block(int(h))
	if err != nil {
		f.logger.Debug("Error on retriving block", "height", h, "err", err)
		return err
	}

	res, _ := args.Results.NewResult()
	d, _ := b.Encode()
	if err = res.SetData(d); err != nil {
		return err
	}
	if err := res.SetHash(b.Hash().RawBytes()); err != nil {
		return err
	}
	if v == 1 {
		if err := f.ToVerboseBlock(b, &res); err != nil {
			return err
		}
	}
	return nil
}

func (f factory) ToVerboseBlock(block *block.Block, res *BlockResult) error {
	cb, _ := res.NewBlock()
	ch, _ := cb.NewHeader()
	ctxs, _ := cb.NewTxs()
	clc, _ := cb.NewLastCommit()

	// last commit
	if block.LastCommit() != nil {
		if err := clc.SetBlockHash(block.LastCommit().BlockHash().RawBytes()); err != nil {
			return err
		}
		clc.SetRound(uint32(block.LastCommit().Round()))
		if err := clc.SetSignature(block.LastCommit().Signature().RawBytes()); err != nil {
			return err
		}
		clcc, _ := clc.NewCommitters(int32(len(block.LastCommit().Committers())))
		for i, commiter := range block.LastCommit().Committers() {
			c := clcc.At(i)
			c.SetNumber(int32(commiter.Number))
			c.SetStatus(int32(commiter.Status))
		}
	}
	// header
	ch.SetVersion(int32(block.Header().Version()))
	ch.SetTime(block.Header().Time().Unix())
	if err := ch.SetTxsHash(block.Header().TxIDsHash().RawBytes()); err != nil {
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
	cTxIDs, _ := ctxs.NewHashes(int32(block.TxIDs().Len()))
	for i, hash := range block.TxIDs().IDs() {
		if err := cTxIDs.Set(i, hash.RawBytes()); err != nil {
			return err
		}
	}

	return nil
}
