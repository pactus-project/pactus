package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	res, err := s.transaction.GetTransaction(ctx,
		&pactus.GetTransactionRequest{
			Id:        vars["id"],
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO,
		},
	)
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := newTableMaker()
	txToTable(tm, res.Transaction)
	s.writeHTML(w, tm.html())
}

func txToTable(tm *tableMaker, trx *pactus.TransactionInfo) {
	if trx == nil {
		return
	}
	tm.addRowTxID("ID", trx.Id)
	tm.addRowInt("Version", int(trx.Version))
	tm.addRowInt("LockTime", int(trx.LockTime))
	tm.addRowAmount("Fee", amount.Amount(trx.Fee))
	tm.addRowString("Memo", trx.Memo)
	tm.addRowString("Payload type", trx.PayloadType.String())

	switch trx.PayloadType {
	case pactus.PayloadType_TRANSFER_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Transfer).Transfer
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowAccAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Amount", amount.Amount(pld.Amount))

	case pactus.PayloadType_BOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Bond).Bond
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowValAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Stake", amount.Amount(pld.Stake))

	case pactus.PayloadType_SORTITION_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Sortition).Sortition
		tm.addRowValAddress("Address", pld.Address)
		tm.addRowString("Proof", pld.Proof)

	case pactus.PayloadType_UNBOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Unbond).Unbond
		tm.addRowValAddress("Validator", pld.Validator)

	case pactus.PayloadType_WITHDRAW_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Withdraw).Withdraw
		tm.addRowValAddress("Sender", pld.ValidatorAddress)
		tm.addRowAccAddress("Receiver", pld.AccountAddress)
		tm.addRowAmount("Amount", amount.Amount(pld.Amount))

	case pactus.PayloadType_UNKNOWN:
		tm.addRowValAddress("error", "unknown payload type")
	}
	if trx.PublicKey != "" {
		tm.addRowString("PublicKey", trx.PublicKey)
	}
	if trx.Signature != "" {
		tm.addRowString("Signature", trx.Signature)
	}
}
