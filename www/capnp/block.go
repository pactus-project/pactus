package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/types/block"
)

func (zs *zarbServer) GetBlockHash(args ZarbServer_getBlockHash) error {
	height := args.Params.Height()
	hash := zs.state.BlockHash(height)
	return args.Results.SetResult(hash.Bytes())
}

func (zs *zarbServer) GetBlock(args ZarbServer_getBlock) error {
	data, _ := args.Params.Hash()
	h, err := hash.FromBytes(data)
	if err != nil {
		return err
	}
	v := args.Params.Verbosity()
	b := zs.state.Block(h)
	if b == nil {
		return fmt.Errorf("block not found")
	}

	res, _ := args.Results.NewResult()
	// TODO: Get it from store
	d, _ := b.Bytes()
	if err := res.SetData(d); err != nil {
		return err
	}
	// TODO: Set height?? Get it from store
	if err := res.SetHash(b.Hash().Bytes()); err != nil {
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
	capBlock, _ := res.NewBlock()
	capHeader, _ := capBlock.NewHeader()
	capTrxs, _ := capBlock.NewTxs(int32(block.Transactions().Len()))
	capPrevCert, _ := capBlock.NewPrevCert()

	// previous certificate
	if block.PrevCertificate() != nil {
		capPrevCert.SetRound(block.PrevCertificate().Round())
		if err := capPrevCert.SetSignature(block.PrevCertificate().Signature().Bytes()); err != nil {
			return err
		}
		capCommitters, _ := capPrevCert.NewCommitters(
			int32(len(block.PrevCertificate().Committers())))
		for i, num := range block.PrevCertificate().Committers() {
			capCommitters.Set(i, num)
		}
		capAbsentees, _ := capPrevCert.NewAbsentees(
			int32(len(block.PrevCertificate().Absentees())))
		for i, num := range block.PrevCertificate().Absentees() {
			capAbsentees.Set(i, num)
		}
	}
	// header
	capHeader.SetVersion(block.Header().Version())
	capHeader.SetTime(int32(block.Header().Time().Unix()))
	err := capHeader.SetStateRoot(block.Header().StateRoot().Bytes())
	if err != nil {
		return err
	}
	err = capHeader.SetPrevBlockHash(block.Header().PrevBlockHash().Bytes())
	if err != nil {
		return err
	}
	err = capHeader.SetProposerAddress(block.Header().ProposerAddress().String())
	if err != nil {
		return err
	}
	// Transactions
	for i, trx := range block.Transactions() {
		d, _ := trx.Bytes()
		if err := capTrxs.Set(i, d); err != nil {
			return err
		}
	}

	return nil
}
