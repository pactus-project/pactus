package http

import (
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
		return p.SetId([]byte(vars["id"]))
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

	tm.addRowBytes("ID", trx.ID().Bytes())
	tm.addRowBytes("Data", d)
	tm.addRowInt("Version", int(trx.Version()))
	tm.addRowBytes("Stamp", trx.Stamp().Bytes())
	tm.addRowInt("Sequence", int(trx.Sequence()))
	tm.addRowInt("Fee", int(trx.Fee()))
	tm.addRowString("Fee", trx.Memo())
	switch trx.Payload().Type() {
	case payload.PayloadTypeBond:
		tm.addRowString("Payload type", "Bond")
		tm.addRowString("Sender", trx.Payload().(*payload.BondPayload).Sender.String())
		tm.addRowString("Validator address", trx.Payload().(*payload.BondPayload).PublicKey.Address().String())
		tm.addRowBytes("Validator PublicKey", trx.Payload().(*payload.BondPayload).PublicKey.Bytes())
		tm.addRowInt("Stake", int(trx.Payload().(*payload.BondPayload).Stake))

	case payload.PayloadTypeSend:
		tm.addRowString("Payload type", "Send")
		tm.addRowString("Sender", trx.Payload().(*payload.SendPayload).Sender.String())
		tm.addRowString("Receiver", trx.Payload().(*payload.SendPayload).Receiver.String())
		tm.addRowInt("Amount", int(trx.Payload().(*payload.SendPayload).Amount))

	case payload.PayloadTypeSortition:
		tm.addRowString("Payload type", "Sortition")
		tm.addRowString("Address", trx.Payload().(*payload.SortitionPayload).Address.String())
		tm.addRowBytes("Proof", trx.Payload().(*payload.SortitionPayload).Proof[:])

	}
	if trx.PublicKey() != nil {
		tm.addRowBytes("PublicKey", trx.PublicKey().Bytes())
	}
	if trx.Signature() != nil {
		tm.addRowBytes("Signature", trx.Signature().Bytes())
	}
}
