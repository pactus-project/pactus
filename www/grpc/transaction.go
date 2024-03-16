package grpc

import (
	"context"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transactionServer struct {
	*Server
}

func newTransactionServer(server *Server) *transactionServer {
	return &transactionServer{
		Server: server,
	}
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

func (s *transactionServer) BroadcastTransaction(_ context.Context,
	req *pactus.BroadcastTransactionRequest,
) (*pactus.BroadcastTransactionResponse, error) {
	trx, err := tx.FromBytes(req.SignedRawTransaction)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
	}

	if err := trx.BasicCheck(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't verify transaction: %v", err.Error())
	}

	if err := s.state.AddPendingTxAndBroadcast(trx); err != nil {
		return nil, status.Errorf(codes.Canceled, "couldn't add to transaction pool: %v", err.Error())
	}

	return &pactus.BroadcastTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}

func (s *transactionServer) CalculateFee(_ context.Context,
	req *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	amt, err := amount.NewAmount(req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err.Error())
	}

	fee := s.state.CalculateFee(amt, payload.Type(req.PayloadType))

	if req.FixedAmount {
		amt -= fee
	}

	return &pactus.CalculateFeeResponse{
		Amount: amt.ToPAC(),
		Fee:    fee.ToPAC(),
	}, nil
}

func (s *transactionServer) GetRawTransferTransaction(_ context.Context,
	req *pactus.GetRawTransferTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	sender, err := crypto.AddressFromString(req.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := crypto.AddressFromString(req.Receiver)
	if err != nil {
		return nil, err
	}

	amt, err := s.getAmount(req.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := s.getFee(req.Fee, amt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid fee: %v", err.Error())
	}

	lockTime := req.LockTime
	if lockTime == 0 {
		lockTime = s.state.LastBlockHeight()
	}

	transferTx := tx.NewTransferTx(lockTime, sender, receiver, amt, fee, req.Memo)
	rawTx, err := transferTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: rawTx,
	}, nil
}

func (s *transactionServer) GetRawBondTransaction(_ context.Context,
	req *pactus.GetRawBondTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	sender, err := crypto.AddressFromString(req.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := crypto.AddressFromString(req.Receiver)
	if err != nil {
		return nil, err
	}

	var publicKey *bls.PublicKey
	if req.PublicKey != "" {
		publicKey, err = bls.PublicKeyFromString(req.PublicKey)
		if err != nil {
			return nil, err
		}
	} else {
		publicKey = nil
	}

	amt, err := s.getAmount(req.Stake)
	if err != nil {
		return nil, err
	}

	fee, err := s.getFee(req.Fee, amt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid fee: %v", err.Error())
	}

	lockTime := req.LockTime
	if lockTime == 0 {
		lockTime = s.state.LastBlockHeight()
	}

	bondTx := tx.NewBondTx(lockTime, sender, receiver, publicKey, amt, fee, req.Memo)
	rawTx, err := bondTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: rawTx,
	}, nil
}

func (s *transactionServer) GetRawUnbondTransaction(_ context.Context,
	req *pactus.GetRawUnbondTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	validatorAddr, err := crypto.AddressFromString(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	lockTime := req.LockTime
	if lockTime == 0 {
		lockTime = s.state.LastBlockHeight()
	}

	unbondTx := tx.NewUnbondTx(lockTime, validatorAddr, req.Memo)
	rawTx, err := unbondTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: rawTx,
	}, nil
}

func (s *transactionServer) GetRawWithdrawTransaction(_ context.Context,
	req *pactus.GetRawWithdrawTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	validatorAddr, err := crypto.AddressFromString(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	accountAddr, err := crypto.AddressFromString(req.AccountAddress)
	if err != nil {
		return nil, err
	}

	amt, err := s.getAmount(req.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := s.getFee(req.Fee, amt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid fee: %v", err.Error())
	}

	lockTime := req.LockTime
	if lockTime == 0 {
		lockTime = s.state.LastBlockHeight()
	}

	withdrawTx := tx.NewWithdrawTx(lockTime, validatorAddr, accountAddr, amt, fee, req.Memo)
	rawTx, err := withdrawTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: rawTx,
	}, nil
}

func (s *transactionServer) getAmount(a float64) (amount.Amount, error) {
	amt, err := amount.NewAmount(a)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err.Error())
	}
	if amt == 0 {
		return 0, status.Errorf(codes.InvalidArgument, "amount is zero")
	}

	return amt, nil
}

func (s *transactionServer) getFee(f float64, amt amount.Amount) (amount.Amount, error) {
	fee, err := amount.NewAmount(f)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "invalid fee: %v", err.Error())
	}

	if fee == 0 {
		fee = s.state.CalculateFee(amt, payload.TypeTransfer)
	}

	return fee, nil
}

func transactionToProto(trx *tx.Tx) *pactus.TransactionInfo {
	data, _ := trx.Bytes()
	transaction := &pactus.TransactionInfo{
		Id:          trx.ID().Bytes(),
		Data:        data,
		Version:     int32(trx.Version()),
		LockTime:    trx.LockTime(),
		Fee:         trx.Fee().ToPAC(),
		Value:       trx.Payload().Value().ToPAC(),
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
				Amount:   pld.Amount.ToPAC(),
			},
		}
	case payload.TypeBond:
		pld := trx.Payload().(*payload.BondPayload)
		transaction.Payload = &pactus.TransactionInfo_Bond{
			Bond: &pactus.PayloadBond{
				Sender:   pld.From.String(),
				Receiver: pld.To.String(),
				Stake:    pld.Stake.ToPAC(),
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
				Amount: pld.Amount.ToPAC(),
			},
		}
	default:
		logger.Error("payload type not defined", "type", trx.Payload().Type())
	}

	return transaction
}
