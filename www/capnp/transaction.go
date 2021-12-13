package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
)

func (zs *zarbServer) GetTransaction(args ZarbServer_getTransaction) error {
	s, _ := args.Params.Id()
	h, err := hash.FromString(string(s))
	if err != nil {
		return fmt.Errorf("invalid transaction id: %s", err)
	}
	trx := zs.state.Transaction(h)
	if trx == nil {
		return fmt.Errorf("transaction not found")
	}

	res, _ := args.Results.NewResult()
	trxData, _ := trx.Encode()
	if err := res.SetData(trxData); err != nil {
		return err
	}
	if err := res.SetId(trx.ID().RawBytes()); err != nil {
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

	if err := zs.state.AddPendingTxAndBroadcast(&tx); err != nil {
		return err
	}

	res, _ := args.Results.NewResult()
	if err := res.SetId(tx.ID().RawBytes()); err != nil {
		return err
	}
	res.SetStatus(0)
	return nil

}
