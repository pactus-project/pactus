package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
)

func (zs *zarbServer) GetBlockHeight(args ZarbServer_getBlockHeight) error {
	s, _ := args.Params.Hash()
	h, err := crypto.HashFromString(string(s))
	if err != nil {
		return err
	}
	num := zs.state.BlockHeight(h)
	args.Results.SetResult(uint64(num))
	return nil
}

func (zs *zarbServer) GetBlock(args ZarbServer_getBlock) error {
	h := args.Params.Height()
	v := args.Params.Verbosity()
	b := zs.state.Block(int(h))
	if b == nil {
		return fmt.Errorf("block not found")
	}

	res, _ := args.Results.NewResult()
	d, _ := b.Encode()
	if err := res.SetData(d); err != nil {
		return err
	}
	if err := res.SetHash(b.Hash().RawBytes()); err != nil {
		return err
	}
	if v == 1 {
		if err := zs.ToVerboseBlock(b, &res); err != nil {
			return err
		}
	}
	return nil
}

func (zs zarbServer) ToVerboseBlock(block *block.Block, res *BlockResult) error {
	cb, _ := res.NewBlock()
	ch, _ := cb.NewHeader()
	ctxs, _ := cb.NewTxs()
	clc, _ := cb.NewLastCertificate()

	// last commit
	if block.LastCertificate() != nil {
		if err := clc.SetBlockHash(block.LastCertificate().BlockHash().RawBytes()); err != nil {
			return err
		}
		clc.SetRound(uint32(block.LastCertificate().Round()))
		if err := clc.SetSignature(block.LastCertificate().Signature().RawBytes()); err != nil {
			return err
		}
		committers, _ := clc.NewCommitters(int32(len(block.LastCertificate().Committers())))
		for i, num := range block.LastCertificate().Committers() {
			committers.Set(i, int32(num))
		}
		absentees, _ := clc.NewAbsentees(int32(len(block.LastCertificate().Absentees())))
		for i, num := range block.LastCertificate().Absentees() {
			absentees.Set(i, int32(num))
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
	if err := ch.SetLastBlockHash(block.Header().LastBlockHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetLastCertificateHash(block.Header().LastCertificateHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetProposerAddress(block.Header().ProposerAddress().RawBytes()); err != nil {
		return err
	}
	// Transactions
	cTxIDs, _ := ctxs.NewHashes(int32(block.TxIDs().Len()))
	for i, id := range block.TxIDs().IDs() {
		if err := cTxIDs.Set(i, id.RawBytes()); err != nil {
			return err
		}
	}

	return nil
}
