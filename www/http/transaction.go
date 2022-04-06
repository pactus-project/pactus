package http

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	b := s.capnp.GetTransaction(s.ctx, func(p capnp.ZarbServer_getTransaction_Params) error {
		vars := mux.Vars(r)
		id, err := hex.DecodeString(vars["id"])
		if err != nil {
			return err
		}
		return p.SetId(id)
	})

	t, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := t.Result()
	data, _ := res.Data()
	fmt.Printf("%x\n", data)
	trx, err := tx.FromBytes(data)
	if err != nil {
		s.writeError(w, err)
		return
	}
	tm := newTableMaker()
	txToTable(trx, tm)
	s.writeHTML(w, tm.html())

}

func txToTable(trx *tx.Tx, tm *tableMaker) {
	d, _ := trx.Bytes()

	tm.addRowTxID("ID", trx.ID().Bytes())
	tm.addRowBytes("Data", d)
	tm.addRowInt("Version", int(trx.Version()))
	tm.addRowBytes("Stamp", trx.Stamp().Bytes())
	tm.addRowInt("Sequence", int(trx.Sequence()))
	tm.addRowInt("Fee", int(trx.Fee()))
	tm.addRowString("Memo", trx.Memo())
	switch trx.Payload().Type() {
	case payload.PayloadTypeBond:
		tm.addRowString("Payload type", "Bond")
		tm.addRowAccAddress("Sender", trx.Payload().(*payload.BondPayload).Sender.String())
		tm.addRowValAddress("Validator address", trx.Payload().(*payload.BondPayload).PublicKey.Address().String())
		tm.addRowBytes("Validator PublicKey", trx.Payload().(*payload.BondPayload).PublicKey.Bytes())
		tm.addRowInt("Stake", int(trx.Payload().(*payload.BondPayload).Stake))

	case payload.PayloadTypeSend:
		tm.addRowString("Payload type", "Send")
		tm.addRowAccAddress("Sender", trx.Payload().(*payload.SendPayload).Sender.String())
		tm.addRowAccAddress("Receiver", trx.Payload().(*payload.SendPayload).Receiver.String())
		tm.addRowInt("Amount", int(trx.Payload().(*payload.SendPayload).Amount))

	case payload.PayloadTypeSortition:
		tm.addRowString("Payload type", "Sortition")
		tm.addRowValAddress("Address", trx.Payload().(*payload.SortitionPayload).Address.String())
		tm.addRowBytes("Proof", trx.Payload().(*payload.SortitionPayload).Proof[:])
	}
	if trx.PublicKey() != nil {
		tm.addRowBytes("PublicKey", trx.PublicKey().Bytes())
	}
	if trx.Signature() != nil {
		tm.addRowBytes("Signature", trx.Signature().Bytes())
	}
}
