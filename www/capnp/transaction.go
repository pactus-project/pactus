package capnp

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func (zs *zarbServer) GetTransaction(args ZarbServer_getTransaction) error {
	s, _ := args.Params.Id()
	h, err := crypto.HashFromString(string(s))
	if err != nil {
		return err
	}
	ctx, err := zs.store.Transaction(h)
	if err != nil {
		return err
	}

	res, _ := args.Results.NewResult()
	trxData, _ := ctx.Tx.Encode()
	if err := res.SetData(trxData); err != nil {
		return err
	}
	if err := res.SetId(ctx.Tx.ID().RawBytes()); err != nil {
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

//Send the raw transaction
func (zs *zarbServer) SendRawTransaction(args ZarbServer_sendRawTransaction) error {
	rawTx, _ := args.Params.RawTx()

	var tx tx.Tx

	if err := tx.Decode(rawTx); err != nil {
		return err
	}

	if err := tx.SanityCheck(); err != nil {
		return err
	}

	if err := zs.txPool.AppendTxAndBroadcast(&tx); err != nil {
		return err
	}

	res, _ := args.Results.NewResult()
	if err := res.SetId(tx.ID().RawBytes()); err != nil {
		return err
	}
	res.SetStatus(0)
	return nil

}
