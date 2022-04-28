package grpc

import (
	"context"
	"encoding/hex"

	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetTransaction(ctx context.Context, request *zarb.TransactionRequest) (*zarb.TransactionResponse, error) {
	id, err := hash.FromString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction ID: %v", err.Error())

	}
	trx := zs.state.Transaction(id)
	if trx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction not found")
	}

	return &zarb.TransactionResponse{
		Tranaction: transactionToProto(trx),
	}, nil

}

func (zs *zarbServer) SendRawTransaction(ctx context.Context, request *zarb.SendRawTransactionRequest) (*zarb.SendRawTransactionResponse, error) {
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

	return &zarb.SendRawTransactionResponse{
		Id: trx.ID().String(),
	}, nil
}

func transactionToProto(trx *tx.Tx) *zarb.TransactionInfo {
	transaction := &zarb.TransactionInfo{
		Id:       trx.ID().Bytes(),
		Version:  int32(trx.Version()),
		Stamp:    trx.Stamp().Bytes(),
		Sequence: trx.Sequence(),
		Fee:      trx.Fee(),
		Type:     zarb.PayloadType(trx.Payload().Type()),
		Memo:     trx.Memo(),
	}

	if trx.PublicKey() != nil {
		transaction.PublicKey = trx.PublicKey().Bytes()
	}

	if trx.Signature() != nil {
		transaction.Signature = trx.Signature().Bytes()
	}

	switch trx.Payload().Type() {
	case payload.PayloadTypeSend:
		pld := trx.Payload().(*payload.SendPayload)
		transaction.Payload = &zarb.TransactionInfo_Send{
			Send: &zarb.SEND_PAYLOAD{
				Sender:   pld.Sender.String(),
				Receiver: pld.Receiver.String(),
				Amount:   pld.Amount,
			},
		}
	case payload.PayloadTypeBond:
		pld := trx.Payload().(*payload.BondPayload)
		transaction.Payload = &zarb.TransactionInfo_Bond{
			Bond: &zarb.BOND_PAYLOAD{
				Sender:    pld.Sender.String(),
				Validator: pld.PublicKey.Bytes(),
				Stake:     pld.Stake,
			},
		}
	case payload.PayloadTypeSortition:
		pld := trx.Payload().(*payload.SortitionPayload)
		transaction.Payload = &zarb.TransactionInfo_Sortition{
			Sortition: &zarb.SORTITION_PAYLOAD{
				Address: pld.Address.String(),
				Proof:   pld.Proof[:],
			},
		}
	default:
		logger.Error("payload type not defined", "Type", trx.Payload().Type())
	}

	return transaction
}
