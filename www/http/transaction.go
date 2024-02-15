package http

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := hex.DecodeString(vars["id"])
	if err != nil {
		s.writeError(w, err)

		return
	}

	res, err := s.transaction.GetTransaction(r.Context(),
		&pactus.GetTransactionRequest{
			Id:        id,
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO,
		},
	)
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := newTableMaker()
	txToTable(res.Transaction, tm)
	s.writeHTML(w, tm.html())
}

func txToTable(trx *pactus.TransactionInfo, tm *tableMaker) {
	if trx == nil {
		return
	}
	tm.addRowTxID("ID", trx.Id)
	tm.addRowBytes("Data", trx.Data)
	tm.addRowInt("Version", int(trx.Version))
	tm.addRowInt("LockTime", int(trx.LockTime))
	tm.addRowInt("Fee", int(trx.Fee))
	tm.addRowString("Memo", trx.Memo)
	tm.addRowString("Payload type", trx.PayloadType.String())

	switch trx.PayloadType {
	case pactus.PayloadType_TRANSFER_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Transfer).Transfer
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowAccAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Amount", pld.Amount)

	case pactus.PayloadType_BOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Bond).Bond
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowValAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Stake", pld.Stake)

	case pactus.PayloadType_SORTITION_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Sortition).Sortition
		tm.addRowValAddress("Address", pld.Address)
		tm.addRowBytes("Proof", pld.Proof)

	case pactus.PayloadType_UNBOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Unbond).Unbond
		tm.addRowValAddress("Validator", pld.Validator)

	case pactus.PayloadType_WITHDRAW_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Withdraw).Withdraw
		tm.addRowValAddress("Sender", pld.From)
		tm.addRowAccAddress("Receiver", pld.To)
		tm.addRowAmount("Amount", pld.Amount)

	case pactus.PayloadType_UNKNOWN:
		tm.addRowValAddress("error", "unknown payload type")
	}
	if trx.PublicKey != "" {
		tm.addRowString("PublicKey", trx.PublicKey)
	}
	if trx.Signature != nil {
		tm.addRowBytes("Signature", trx.Signature)
	}
}
