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

	committedTx, err := s.state.CommittedTx(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	res := &pactus.GetTransactionResponse{
		BlockHeight: committedTx.Height,
		BlockTime:   committedTx.BlockTime,
	}

	switch req.Verbosity {
	case pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_DATA:
		res.Transaction = &pactus.TransactionInfo{
			Id:   committedTx.TxID.String(),
			Data: hex.EncodeToString(committedTx.Data),
		}

	case pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO:
		trx, err := committedTx.ToTx()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err.Error())
		}
		lastBlockHeight := s.state.LastBlockHeight()
		confirmations := int(lastBlockHeight) - int(committedTx.Height)
		res.Transaction = transactionToProto(trx, committedTx.Height, confirmations)
	}

	return res, nil
}

func (s *transactionServer) BroadcastTransaction(_ context.Context,
	req *pactus.BroadcastTransactionRequest,
) (*pactus.BroadcastTransactionResponse, error) {
	trx, err := tx.FromString(req.SignedRawTransaction)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
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

func (s *transactionServer) GetRawBatchTransferTransaction(_ context.Context,
	req *pactus.GetRawBatchTransferTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	sender, err := crypto.AddressFromString(req.Sender)
	if err != nil {
		return nil, err
	}

	totalAmount := amount.Amount(0)

	recipients := make([]payload.BatchRecipient, 0, len(req.Recipients))
	for _, recipient := range req.Recipients {
		receiver, err := crypto.AddressFromString(recipient.Receiver)
		if err != nil {
			return nil, err
		}

		amt := amount.Amount(recipient.Amount)

		recipients = append(recipients, payload.BatchRecipient{
			To:     receiver,
			Amount: amt,
		})

		totalAmount += amt
	}

	fee := s.getFee(req.Fee, totalAmount)
	lockTime := s.getLockTime(req.LockTime)

	batchTransferTx := tx.NewBatchTransferTx(lockTime, sender, recipients, fee, tx.WithMemo(req.Memo))
	rawTx, err := batchTransferTx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.GetRawTransactionResponse{
		RawTransaction: hex.EncodeToString(rawTx),
	}, nil
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

func transactionToProto(trx *tx.Tx, blockHeight uint32, confirmations int) *pactus.TransactionInfo {
	trxInfo := &pactus.TransactionInfo{
		Id:            trx.ID().String(),
		Version:       int32(trx.Version()),
		LockTime:      trx.LockTime(),
		Fee:           trx.Fee().ToNanoPAC(),
		Value:         trx.Payload().Value().ToNanoPAC(),
		PayloadType:   pactus.PayloadType(trx.Payload().Type()),
		Memo:          trx.Memo(),
		BlockHeight:   blockHeight,
		Confirmed:     confirmations > 0,
		Confirmations: int32(confirmations),
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

		publicKeyStr := ""
		if pld.PublicKey != nil {
			publicKeyStr = pld.PublicKey.String()
		}

		trxInfo.Payload = &pactus.TransactionInfo_Bond{
			Bond: &pactus.PayloadBond{
				Sender:    pld.From.String(),
				Receiver:  pld.To.String(),
				Stake:     pld.Stake.ToNanoPAC(),
				PublicKey: publicKeyStr,
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

	case payload.TypeBatchTransfer:
		pld := trx.Payload().(*payload.BatchTransferPayload)
		recipients := make([]*pactus.Recipient, 0, len(pld.Recipients))
		for _, recipient := range pld.Recipients {
			recipients = append(recipients, &pactus.Recipient{
				Receiver: recipient.To.String(),
				Amount:   recipient.Amount.ToNanoPAC(),
			})
		}

		trxInfo.Payload = &pactus.TransactionInfo_BatchTransfer{
			BatchTransfer: &pactus.PayloadBatchTransfer{
				Sender:     pld.From.String(),
				Recipients: recipients,
			},
		}
	default:
		logger.Error("payload type not defined", "type", trx.Payload().Type())
	}

	return trxInfo
}

func (*transactionServer) DecodeRawTransaction(_ context.Context,
	req *pactus.DecodeRawTransactionRequest,
) (*pactus.DecodeRawTransactionResponse, error) {
	trx, err := tx.FromString(req.RawTransaction)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't decode transaction: %v", err.Error())
	}

	return &pactus.DecodeRawTransactionResponse{
		Transaction: transactionToProto(trx, 0, 0),
	}, nil
}
