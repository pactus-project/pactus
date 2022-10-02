package http

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/www/capnp"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	b := s.capnp.GetTransaction(s.ctx, func(p capnp.PactusServer_getTransaction_Params) error {
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
	case payload.PayloadTypeSend:
		tm.addRowString("Payload type", "Send")
		tm.addRowAccAddress("Sender", trx.Payload().(*payload.SendPayload).Sender.String())
		tm.addRowAccAddress("Receiver", trx.Payload().(*payload.SendPayload).Receiver.String())
		tm.addRowAmount("Amount", trx.Payload().(*payload.SendPayload).Amount)

	case payload.PayloadTypeBond:
		tm.addRowString("Payload type", "Bond")
		tm.addRowAccAddress("Sender", trx.Payload().(*payload.BondPayload).Sender.String())
		tm.addRowValAddress("Receiver", trx.Payload().(*payload.BondPayload).Receiver.String())
		tm.addRowAmount("Stake", trx.Payload().(*payload.BondPayload).Stake)

	case payload.PayloadTypeSortition:
		tm.addRowString("Payload type", "Sortition")
		tm.addRowValAddress("Address", trx.Payload().(*payload.SortitionPayload).Address.String())
		tm.addRowBytes("Proof", trx.Payload().(*payload.SortitionPayload).Proof[:])

	case payload.PayloadTypeUnbond:
		tm.addRowString("Payload type", "Unbond")
		tm.addRowValAddress("Validator", trx.Payload().(*payload.UnbondPayload).Validator.String())

	case payload.PayloadTypeWithdraw:
		tm.addRowString("Payload type", "Withdraw")
		tm.addRowValAddress("Sender", trx.Payload().(*payload.WithdrawPayload).From.String())
		tm.addRowAccAddress("Receiver", trx.Payload().(*payload.WithdrawPayload).To.String())
		tm.addRowAmount("Amount", trx.Payload().(*payload.WithdrawPayload).Amount)
	}
	if trx.PublicKey() != nil {
		tm.addRowBytes("PublicKey", trx.PublicKey().Bytes())
	}
	if trx.Signature() != nil {
		tm.addRowBytes("Signature", trx.Signature().Bytes())
	}
}
