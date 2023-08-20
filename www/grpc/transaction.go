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
	// TODO: Use RawTransaction here
	storedTx := s.state.StoredTx(id)
	if storedTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	res := &pactus.GetTransactionResponse{
		BlockHeight: storedTx.Height,
		BlockTime:   storedTx.BlockTime,
	}

	if req.Verbosity > pactus.TransactionVerbosity_TRANSACTION_DATA {
		res.Transaction = transactionToProto(storedTx.ToTx())
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
		Sequence:    trx.Sequence(),
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
	case payload.PayloadTypeTransfer:
		pld := trx.Payload().(*payload.TransferPayload)
		transaction.Payload = &pactus.TransactionInfo_Transfer{
			Transfer: &pactus.PayloadTransfer{
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
	case payload.PayloadTypeUnbond:
		pld := trx.Payload().(*payload.UnbondPayload)
		transaction.Payload = &pactus.TransactionInfo_Unbond{
			Unbond: &pactus.PayloadUnbond{
				Validator: pld.Validator.String(),
			},
		}
	case payload.PayloadTypeWithdraw:
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
