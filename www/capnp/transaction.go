package capnp

import "github.com/zarbchain/zarb-go/crypto"

func (f factory) GetTransaction(args ZarbServer_getTransaction) error {
	s, _ := args.Params.Hash()
	h, err := crypto.HashFromString(string(s))
	if err != nil {
		return err
	}
	ctx, err := f.store.Transaction(h)
	if err != nil {
		return err
	}

	res, _ := args.Results.NewResult()
	trxData, _ := ctx.Tx.Encode()
	if err := res.SetData(trxData); err != nil {
		return err
	}
	if err := res.SetHash(ctx.Tx.ID().RawBytes()); err != nil {
		return err
	}
	rec, _ := res.NewReceipt()
	recData, _ := ctx.Receipt.Encode()
	if err := rec.SetData(recData); err != nil {
		return err
	}
	if err := rec.SetHash(ctx.Receipt.Hash().RawBytes()); err != nil {
		return err
	}
	return nil
}
