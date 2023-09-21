package grpc

import (
	"context"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transactionServer struct {
	state  state.Facade
	logger *logger.SubLogger
}

func (s *transactionServer) GetTransaction(_ context.Context,
	req *pactus.GetTransactionRequest,
) (*pactus.GetTransactionResponse, error) {
	id, err := hash.FromBytes(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction ID: %v", err.Error())
	}

	committedTx := s.state.CommittedTx(id)
	if committedTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	res := &pactus.GetTransactionResponse{
		BlockHeight: committedTx.Height,
		BlockTime:   committedTx.BlockTime,
	}

	if req.Verbosity == pactus.TransactionVerbosity_TRANSACTION_DATA {
		res.Transaction = &pactus.TransactionInfo{
			Data: committedTx.Data,
			Id:   committedTx.TxID.Bytes(),
		}
	} else {
		trx, err := committedTx.ToTx()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		res.Transaction = transactionToProto(trx)
	}

	return res, nil
}

func (s *transactionServer) SendRawTransaction(_ context.Context,
	req *pactus.SendRawTransactionRequest,
) (*pactus.SendRawTransactionResponse, error) {
	trx, err := tx.FromBytes(req.Data)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
	}

	if err := trx.BasicCheck(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't verify transaction: %v", err.Error())
	}

	if err := s.state.AddPendingTxAndBroadcast(trx); err != nil {
		return nil, status.Errorf(codes.Canceled, "couldn't add to transaction pool: %v", err.Error())
	}

	return &pactus.SendRawTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}

func (s *transactionServer) CalculateFee(_ context.Context,
	req *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	fee, err := s.state.CalculateFee(req.Amount, payload.Type(req.PayloadType))
	if err != nil {
		return nil, err
	}

	return &pactus.CalculateFeeResponse{
		Fee: fee,
	}, nil
}

func transactionToProto(trx *tx.Tx) *pactus.TransactionInfo {
	data, _ := trx.Bytes()
	transaction := &pactus.TransactionInfo{
		Id:          trx.ID().Bytes(),
		Data:        data,
		Version:     int32(trx.Version()),
		Stamp:       trx.Stamp().Bytes(),
		LockTime:    trx.LockTime(),
		Fee:         trx.Fee(),
		Value:       trx.Payload().Value(),
		PayloadType: pactus.PayloadType(trx.Payload().Type()),
		Memo:        trx.Memo(),
	}

	if trx.PublicKey() != nil {
		transaction.PublicKey = trx.PublicKey().String()
	}

	if trx.Signature() != nil {
		transaction.Signature = trx.Signature().Bytes()
	}

	switch trx.Payload().Type() {
	case payload.TypeTransfer:
		pld := trx.Payload().(*payload.TransferPayload)
		transaction.Payload = &pactus.TransactionInfo_Transfer{
			Transfer: &pactus.PayloadTransfer{
				Sender:   pld.From.String(),
				Receiver: pld.To.String(),
				Amount:   pld.Amount,
			},
		}
	case payload.TypeBond:
		pld := trx.Payload().(*payload.BondPayload)
		transaction.Payload = &pactus.TransactionInfo_Bond{
			Bond: &pactus.PayloadBond{
				Sender:   pld.From.String(),
				Receiver: pld.To.String(),
				Stake:    pld.Stake,
			},
		}
	case payload.TypeSortition:
		pld := trx.Payload().(*payload.SortitionPayload)
		transaction.Payload = &pactus.TransactionInfo_Sortition{
			Sortition: &pactus.PayloadSortition{
				Address: pld.Validator.String(),
				Proof:   pld.Proof[:],
			},
		}
	case payload.TypeUnbond:
		pld := trx.Payload().(*payload.UnbondPayload)
		transaction.Payload = &pactus.TransactionInfo_Unbond{
			Unbond: &pactus.PayloadUnbond{
				Validator: pld.Validator.String(),
			},
		}
	case payload.TypeWithdraw:
		pld := trx.Payload().(*payload.WithdrawPayload)
		transaction.Payload = &pactus.TransactionInfo_Withdraw{
			Withdraw: &pactus.PayloadWithdraw{
				From:   pld.From.String(),
				To:     pld.To.String(),
				Amount: pld.Amount,
			},
		}
	default:
		logger.Error("payload type not defined", "type", trx.Payload().Type())
	}

	return transaction
}
