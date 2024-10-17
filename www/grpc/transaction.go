package grpc

import (
	"context"
	"encoding/hex"

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
	id, err := hash.FromString(req.Id)
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

	switch req.Verbosity {
	case pactus.TransactionVerbosity_TRANSACTION_DATA:
		res.Transaction = &pactus.TransactionInfo{
			Id:   committedTx.TxID.String(),
			Data: hex.EncodeToString(committedTx.Data),
		}

	case pactus.TransactionVerbosity_TRANSACTION_INFO:
		trx, err := committedTx.ToTx()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err.Error())
		}
		res.Transaction = transactionToProto(trx)
	}

	return res, nil
}

func (s *transactionServer) BroadcastTransaction(_ context.Context,
	req *pactus.BroadcastTransactionRequest,
) (*pactus.BroadcastTransactionResponse, error) {
	b, err := hex.DecodeString(req.SignedRawTransaction)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid signed transaction")
	}

	trx, err := tx.FromBytes(b)
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
		Id: trx.ID().String(),
	}, nil
}

func (s *transactionServer) CalculateFee(_ context.Context,
	req *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	amt := amount.Amount(req.Amount)
	fee := s.state.CalculateFee(amt, payload.Type(req.PayloadType))

	if req.FixedAmount {
		amt -= fee
	}

	return &pactus.CalculateFeeResponse{
		Amount: amt.ToNanoPAC(),
		Fee:    fee.ToNanoPAC(),
	}, nil
}

func (s *transactionServer) GetRawTransaction(_ context.Context,
	req *pactus.GetRawTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	lockTime := s.getLockTime(req.LockTime)

	var trx *tx.Tx
	var err error
	switch pld := req.Payload.(type) {
	case *pactus.GetRawTransactionRequest_Transfer:
		trx, err = s.handleRawTransfer(lockTime, req.Memo, req.Fee, pld.Transfer)

	case *pactus.GetRawTransactionRequest_Bond:
		trx, err = s.handleRawBond(lockTime, req.Memo, req.Fee, pld.Bond)

	case *pactus.GetRawTransactionRequest_Unbond:
		trx, err = s.handleRawUnbond(lockTime, req.Memo, pld.Unbond)

	case *pactus.GetRawTransactionRequest_Withdraw:
		trx, err = s.handleRawWithdraw(lockTime, req.Memo, req.Fee, pld.Withdraw)

	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction type")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	data, err := trx.Bytes()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(data),
		Id:             trx.ID().String(),
	}, err
}

// Deprecated: GetRawTransferTransaction is deprecated.
// Use GetRawTransaction API instead.
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

	amt := amount.Amount(req.Amount)
	fee := s.getFee(req.Fee, amt)
	lockTime := s.getLockTime(req.LockTime)

	transferTx := tx.NewTransferTx(lockTime, sender, receiver, amt, fee, tx.WithMemo(req.Memo))
	rawTx, err := transferTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(rawTx),
	}, nil
}

// Deprecated: GetRawBondTransaction is deprecated.
// Use GetRawTransaction API instead.
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

	amt := amount.Amount(req.Stake)
	fee := s.getFee(req.Fee, amt)
	lockTime := s.getLockTime(req.LockTime)

	bondTx := tx.NewBondTx(lockTime, sender, receiver, publicKey, amt, fee, tx.WithMemo(req.Memo))
	rawTx, err := bondTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(rawTx),
	}, nil
}

// Deprecated: GetRawUnbondTransaction is deprecated.
// Use GetRawTransaction API instead.
func (s *transactionServer) GetRawUnbondTransaction(_ context.Context,
	req *pactus.GetRawUnbondTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	validatorAddr, err := crypto.AddressFromString(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	lockTime := s.getLockTime(req.LockTime)

	unbondTx := tx.NewUnbondTx(lockTime, validatorAddr, tx.WithMemo(req.Memo))
	rawTx, err := unbondTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(rawTx),
	}, nil
}

// Deprecated: GetRawWithdrawTransaction is deprecated.
// Use GetRawTransaction API instead.
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

	amt := amount.Amount(req.Amount)
	fee := s.getFee(req.Fee, amt)
	lockTime := s.getLockTime(req.LockTime)

	withdrawTx := tx.NewWithdrawTx(lockTime, validatorAddr, accountAddr, amt, fee, tx.WithMemo(req.Memo))
	rawTx, err := withdrawTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(rawTx),
	}, nil
}

func (s *transactionServer) handleRawTransfer(lockTime uint32, memo string, feeInt int64,
	pld *pactus.PayloadTransfer,
) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(pld.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := crypto.AddressFromString(pld.Receiver)
	if err != nil {
		return nil, err
	}

	amt := amount.Amount(pld.Amount)
	fee := s.getFee(feeInt, amt)

	transferTx := tx.NewTransferTx(lockTime, sender, receiver, amt, fee, tx.WithMemo(memo))

	return transferTx, nil
}

func (s *transactionServer) handleRawBond(lockTime uint32, memo string, feeInt int64,
	pld *pactus.PayloadBond,
) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(pld.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := crypto.AddressFromString(pld.Receiver)
	if err != nil {
		return nil, err
	}

	var publicKey *bls.PublicKey
	if pld.PublicKey != "" {
		publicKey, err = bls.PublicKeyFromString(pld.PublicKey)
		if err != nil {
			return nil, err
		}
	} else {
		publicKey = nil
	}

	amt := amount.Amount(pld.Stake)
	fee := s.getFee(feeInt, amt)

	bondTx := tx.NewBondTx(lockTime, sender, receiver, publicKey, amt, fee, tx.WithMemo(memo))

	return bondTx, nil
}

func (*transactionServer) handleRawUnbond(lockTime uint32, memo string, pld *pactus.PayloadUnbond) (*tx.Tx, error) {
	validatorAddr, err := crypto.AddressFromString(pld.Validator)
	if err != nil {
		return nil, err
	}

	unbondTx := tx.NewUnbondTx(lockTime, validatorAddr, tx.WithMemo(memo))

	return unbondTx, nil
}

func (s *transactionServer) handleRawWithdraw(lockTime uint32, memo string, feeInt int64,
	pld *pactus.PayloadWithdraw,
) (*tx.Tx, error) {
	validatorAddr, err := crypto.AddressFromString(pld.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	accountAddr, err := crypto.AddressFromString(pld.AccountAddress)
	if err != nil {
		return nil, err
	}

	amt := amount.Amount(pld.Amount)
	fee := s.getFee(feeInt, amt)

	withdrawTx := tx.NewWithdrawTx(lockTime, validatorAddr, accountAddr, amt, fee, tx.WithMemo(memo))

	return withdrawTx, nil
}

func (s *transactionServer) getFee(f int64, amt amount.Amount) amount.Amount {
	fee := amount.Amount(f)
	if fee == 0 {
		fee = s.state.CalculateFee(amt, payload.TypeTransfer)
	}

	return fee
}

func (s *transactionServer) getLockTime(lockTime uint32) uint32 {
	if lockTime == 0 {
		lockTime = s.state.LastBlockHeight()
	}

	return lockTime
}

func transactionToProto(trx *tx.Tx) *pactus.TransactionInfo {
	trxInfo := &pactus.TransactionInfo{
		Id:          trx.ID().String(),
		Version:     int32(trx.Version()),
		LockTime:    trx.LockTime(),
		Fee:         trx.Fee().ToNanoPAC(),
		Value:       trx.Payload().Value().ToNanoPAC(),
		PayloadType: pactus.PayloadType(trx.Payload().Type()),
		Memo:        trx.Memo(),
	}

	if trx.PublicKey() != nil {
		trxInfo.PublicKey = trx.PublicKey().String()
	}

	if trx.Signature() != nil {
		trxInfo.Signature = trx.Signature().String()
	}

	switch trx.Payload().Type() {
	case payload.TypeTransfer:
		pld := trx.Payload().(*payload.TransferPayload)
		trxInfo.Payload = &pactus.TransactionInfo_Transfer{
			Transfer: &pactus.PayloadTransfer{
				Sender:   pld.From.String(),
				Receiver: pld.To.String(),
				Amount:   pld.Amount.ToNanoPAC(),
			},
		}
	case payload.TypeBond:
		pld := trx.Payload().(*payload.BondPayload)
		trxInfo.Payload = &pactus.TransactionInfo_Bond{
			Bond: &pactus.PayloadBond{
				Sender:   pld.From.String(),
				Receiver: pld.To.String(),
				Stake:    pld.Stake.ToNanoPAC(),
			},
		}
	case payload.TypeSortition:
		pld := trx.Payload().(*payload.SortitionPayload)
		trxInfo.Payload = &pactus.TransactionInfo_Sortition{
			Sortition: &pactus.PayloadSortition{
				Address: pld.Validator.String(),
				Proof:   hex.EncodeToString(pld.Proof[:]),
			},
		}
	case payload.TypeUnbond:
		pld := trx.Payload().(*payload.UnbondPayload)
		trxInfo.Payload = &pactus.TransactionInfo_Unbond{
			Unbond: &pactus.PayloadUnbond{
				Validator: pld.Validator.String(),
			},
		}
	case payload.TypeWithdraw:
		pld := trx.Payload().(*payload.WithdrawPayload)
		trxInfo.Payload = &pactus.TransactionInfo_Withdraw{
			Withdraw: &pactus.PayloadWithdraw{
				ValidatorAddress: pld.From.String(),
				AccountAddress:   pld.To.String(),
				Amount:           pld.Amount.ToNanoPAC(),
			},
		}
	default:
		logger.Error("payload type not defined", "type", trx.Payload().Type())
	}

	return trxInfo
}
