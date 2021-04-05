package grpc

import (
	"context"
	"encoding/hex"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetTransaction(ctx context.Context, request *zarb.TransactionRequest) (*zarb.TransactionResponse, error) {
	id, err := crypto.HashFromString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid transaction ID: %v", err.Error())

	}
	tx := zs.state.Transaction(id)
	if tx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Transaction not found")
	}

	data, err := tx.Tx.Encode()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var json string
	if request.Verbosity == 1 {
		bz, err := tx.Tx.MarshalJSON()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		json = string(bz)
	}
	res := &zarb.TransactionResponse{
		Data: data,
		Json: json,
	}

	return res, nil

}

func (zs *zarbServer) SendRawTransaction(ctx context.Context, request *zarb.SendRawTransactionRequest) (*zarb.SendRawTransactionResponse, error) {
	var tx tx.Tx

	hexDecoded, err := hex.DecodeString(request.Data)
	if err != nil {
		zs.logger.Error("Invalid transaction", "err", err, "type", "hex decode")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid transaction: Couldn't decode transaction")
	}
	if err := tx.Decode(hexDecoded); err != nil {
		zs.logger.Error("Invalid transaction", "err", err, "type", "decode")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid transaction: Couldn't decode transaction")
	}

	if err := tx.SanityCheck(); err != nil {
		zs.logger.Error("Invalid transaction", "err", err, "type", "sanity")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid transaction: Couldn't Verify Transaction")
	}

	if err := zs.state.AddPendingTxAndBroadcast(&tx); err != nil {
		zs.logger.Error("Couldn't add trx to pool", "err", err)
		return nil, status.Errorf(codes.Canceled, "Couldn't add to Pending pool")
	}

	return &zarb.SendRawTransactionResponse{
		Id: tx.ID().String(),
	}, nil
}
