package grpc

import (
	"context"
	"encoding/hex"

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
	request *pactus.TransactionRequest) (*pactus.TransactionResponse, error) {
	id, err := hash.FromString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction ID: %v", err.Error())
	}
	trx := zs.state.Transaction(id)
	if trx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	return &pactus.TransactionResponse{
		Transaction: transactionToProto(trx),
	}, nil
}

func (zs *transactionServer) SendRawTransaction(ctx context.Context,
	request *pactus.SendRawTransactionRequest) (*pactus.SendRawTransactionResponse, error) {
	data, err := hex.DecodeString(request.Data)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
	}
	trx, err := tx.FromBytes(data)
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
		Id: trx.ID().String(),
	}, nil
}

func transactionToProto(trx *tx.Tx) *pactus.TransactionInfo {
	transaction := &pactus.TransactionInfo{
		Id:       trx.ID().Bytes(),
		Version:  int32(trx.Version()),
		Stamp:    trx.Stamp().Bytes(),
		Sequence: trx.Sequence(),
		Fee:      trx.Fee(),
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
