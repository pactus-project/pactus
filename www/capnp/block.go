package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func (zs *zarbServer) GetBlockHash(args ZarbServer_getBlockHash) error {
	height := args.Params.Height()
	hash := zs.state.BlockHash(int(height))
	args.Results.SetResult(hash.RawBytes())
	return nil
}

func (zs *zarbServer) GetBlock(args ZarbServer_getBlock) error {
	data, _ := args.Params.Hash()
	h, err := hash.FromRawBytes(data)
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
	d, _ := b.Encode()
	if err := res.SetData(d); err != nil {
		return err
	}
	// TODO: Set height?? Get it from store
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

func (zs zarbServer) ToVerboseBlock(b *block.Block, res *BlockResult) error {
	cb, _ := res.NewBlock()
	ch, _ := cb.NewHeader()
	ctxs, _ := cb.NewTxs(int32(b.Transactions().Len()))
	clc, _ := cb.NewPrevCert()

	// previous certificate
	if b.PrevCertificate() != nil {
		if err := clc.SetBlockHash(b.PrevCertificate().BlockHash().RawBytes()); err != nil {
			return err
		}
		clc.SetRound(uint32(b.PrevCertificate().Round()))
		if err := clc.SetSignature(b.PrevCertificate().Signature().RawBytes()); err != nil {
			return err
		}
		committers, _ := clc.NewCommitters(int32(len(b.PrevCertificate().Committers())))
		for i, num := range b.PrevCertificate().Committers() {
			committers.Set(i, int32(num))
		}
		absentees, _ := clc.NewAbsentees(int32(len(b.PrevCertificate().Absentees())))
		for i, num := range b.PrevCertificate().Absentees() {
			absentees.Set(i, int32(num))
		}
	}
	// header
	ch.SetVersion(int32(b.Header().Version()))
	ch.SetTime(b.Header().Time().Unix())
	if err := ch.SetTxsRoot(b.Header().TxsRoot().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetStateRoot(b.Header().StateRoot().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetPrevBlockHash(b.Header().PrevBlockHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetPrevCertHash(b.Header().PrevCertificateHash().RawBytes()); err != nil {
		return err
	}
	if err := ch.SetProposerAddress(b.Header().ProposerAddress().RawBytes()); err != nil {
		return err
	}
	// Transactions
	for i, trx := range b.Transactions() {
		d, _ := trx.Encode()
		if err := ctxs.Set(i, d); err != nil {
			return err
		}
	}

	return nil
}
