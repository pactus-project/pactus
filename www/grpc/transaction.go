package grpc

import (
	"context"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transactionServer struct {
	state  state.Facade
	logger *logger.Logger
}

func (zs *transactionServer) GetTransaction(ctx context.Context,
	req *pactus.TransactionRequest) (*pactus.TransactionResponse, error) {
	id, err := hash.FromBytes(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction ID: %v", err.Error())
	}
	// TODO: Use RawTransaction here
	storedTx := zs.state.StoredTx(id)
	if storedTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	res := &pactus.TransactionResponse{}

	if req.Verbosity > pactus.TransactionVerbosity_TRANSACTION_DATA {
		res.Transaction = transactionToProto(storedTx.ToTx())
	}

	res.Transaction.BlockHeight = storedTx.Height
	res.Transaction.BlockTime = storedTx.BlockTime
	return res, nil
}

func (zs *transactionServer) SendRawTransaction(ctx context.Context,
	req *pactus.SendRawTransactionRequest) (*pactus.SendRawTransactionResponse, error) {
	trx, err := tx.FromBytes(req.Data)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
	}

	if err := trx.SanityCheck(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't verify transaction: %v", err.Error())
	}

	if err := zs.state.AddPendingTxAndBroadcast(trx); err != nil {
		return nil, status.Errorf(codes.Canceled, "couldn't add to transaction pool: %v", err.Error())
	}

	return &pactus.SendRawTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}

func transactionToProto(trx *tx.Tx) *pactus.TransactionInfo {
	data, _ := trx.Bytes()
	transaction := &pactus.TransactionInfo{
		Id:       trx.ID().Bytes(),
		Data:     data,
		Version:  int32(trx.Version()),
		Stamp:    trx.Stamp().Bytes(),
		Sequence: trx.Sequence(),
		Fee:      trx.Fee(),
		Value:    trx.Payload().Value(),
		Type:     pactus.PayloadType(trx.Payload().Type()),
		Memo:     trx.Memo(),
	}

	if trx.PublicKey() != nil {
		transaction.PublicKey = trx.PublicKey().String()
	}

	if trx.Signature() != nil {
		transaction.Signature = trx.Signature().Bytes()
	}

	switch trx.Payload().Type() {
	case payload.PayloadTypeSend:
		pld := trx.Payload().(*payload.SendPayload)
		transaction.Payload = &pactus.TransactionInfo_Send{
			Send: &pactus.PayloadSend{
				Sender:   pld.Sender.String(),
				Receiver: pld.Receiver.String(),
				Amount:   pld.Amount,
			},
		}
	case payload.PayloadTypeBond:
		pld := trx.Payload().(*payload.BondPayload)
		transaction.Payload = &pactus.TransactionInfo_Bond{
			Bond: &pactus.PayloadBond{
				Sender:   pld.Sender.String(),
				Receiver: pld.Receiver.String(),
				Stake:    pld.Stake,
			},
		}
	case payload.PayloadTypeSortition:
		pld := trx.Payload().(*payload.SortitionPayload)
		transaction.Payload = &pactus.TransactionInfo_Sortition{
			Sortition: &pactus.PayloadSortition{
				Address: pld.Address.String(),
				Proof:   pld.Proof[:],
			},
		}
	default:
		logger.Error("payload type not defined", "Type", trx.Payload().Type())
	}

	return transaction
}
