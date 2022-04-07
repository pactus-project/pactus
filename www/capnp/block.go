package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
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

func (zs zarbServer) ToVerboseBlock(b *block.Block, res *BlockResult) error {
	cb, _ := res.NewBlock()
	ch, _ := cb.NewHeader()
	ctxs, _ := cb.NewTxs(int32(b.Transactions().Len()))
	clc, _ := cb.NewPrevCert()

	// previous certificate
	if b.PrevCertificate() != nil {
		clc.SetRound(b.PrevCertificate().Round())
		if err := clc.SetSignature(b.PrevCertificate().Signature().Bytes()); err != nil {
			return err
		}
		committers, _ := clc.NewCommitters(int32(len(b.PrevCertificate().Committers())))
		for i, num := range b.PrevCertificate().Committers() {
			committers.Set(i, num)
		}
		absentees, _ := clc.NewAbsentees(int32(len(b.PrevCertificate().Absentees())))
		for i, num := range b.PrevCertificate().Absentees() {
			absentees.Set(i, num)
		}
	}
	// header
	ch.SetVersion(b.Header().Version())
	ch.SetTime(int32(b.Header().Time().Unix()))
	if err := ch.SetStateRoot(b.Header().StateRoot().Bytes()); err != nil {
		return err
	}
	if err := ch.SetPrevBlockHash(b.Header().PrevBlockHash().Bytes()); err != nil {
		return err
	}
	if err := ch.SetProposerAddress(b.Header().ProposerAddress().String()); err != nil {
		return err
	}
	// Transactions
	for i, trx := range b.Transactions() {
		d, _ := trx.Bytes()
		if err := ctxs.Set(i, d); err != nil {
			return err
		}
	}

	return nil
}
