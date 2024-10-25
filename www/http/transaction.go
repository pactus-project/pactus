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

func txToTable(tmk *tableMaker, trx *pactus.TransactionInfo) {
	if trx == nil {
		return
	}
	tmk.addRowTxID("ID", trx.Id)
	tmk.addRowInt("Version", int(trx.Version))
	tmk.addRowInt("LockTime", int(trx.LockTime))
	tmk.addRowAmount("Fee", amount.Amount(trx.Fee))
	tmk.addRowString("Memo", trx.Memo)
	tmk.addRowString("Payload type", trx.PayloadType.String())

	switch trx.PayloadType {
	case pactus.PayloadType_TRANSFER_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Transfer).Transfer
		tmk.addRowAccAddress("Sender", pld.Sender)
		tmk.addRowAccAddress("Receiver", pld.Receiver)
		tmk.addRowAmount("Amount", amount.Amount(pld.Amount))

	case pactus.PayloadType_BOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Bond).Bond
		tmk.addRowAccAddress("Sender", pld.Sender)
		tmk.addRowValAddress("Receiver", pld.Receiver)
		tmk.addRowAmount("Stake", amount.Amount(pld.Stake))

	case pactus.PayloadType_SORTITION_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Sortition).Sortition
		tmk.addRowValAddress("Address", pld.Address)
		tmk.addRowString("Proof", pld.Proof)

	case pactus.PayloadType_UNBOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Unbond).Unbond
		tmk.addRowValAddress("Validator", pld.Validator)

	case pactus.PayloadType_WITHDRAW_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Withdraw).Withdraw
		tmk.addRowValAddress("Sender", pld.ValidatorAddress)
		tmk.addRowAccAddress("Receiver", pld.AccountAddress)
		tmk.addRowAmount("Amount", amount.Amount(pld.Amount))

	case pactus.PayloadType_UNKNOWN:
		tmk.addRowValAddress("error", "unknown payload type")
	}
	if trx.PublicKey != "" {
		tmk.addRowString("PublicKey", trx.PublicKey)
	}
	if trx.Signature != "" {
		tmk.addRowString("Signature", trx.Signature)
	}
}
