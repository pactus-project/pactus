package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func (zs *zarbServer) GetBlockHeight(args ZarbServer_getBlockHeight) error {
	s, _ := args.Params.Hash()
	h, err := hash.FromString(string(s))
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
	if block.PrevCertificate() != nil {
		if err := clc.SetBlockHash(block.PrevCertificate().BlockHash().RawBytes()); err != nil {
			return err
		}
		clc.SetRound(uint32(block.PrevCertificate().Round()))
		if err := clc.SetSignature(block.PrevCertificate().Signature().RawBytes()); err != nil {
			return err
		}
		committers, _ := clc.NewCommitters(int32(len(block.PrevCertificate().Committers())))
		for i, num := range block.PrevCertificate().Committers() {
			committers.Set(i, int32(num))
		}
		absentees, _ := clc.NewAbsentees(int32(len(block.PrevCertificate().Absentees())))
		for i, num := range block.PrevCertificate().Absentees() {
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
	if err := ch.SetPrevBlockHash(block.Header().PrevBlockHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetPrevCertificateHash(block.Header().PrevCertificateHash().RawBytes()); err != nil {
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
