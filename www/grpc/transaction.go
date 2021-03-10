package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetTransaction(ctx context.Context, request *zarb.TransactionRequest) (*zarb.TransactionResponse, error) {

	hash, err := crypto.HashFromString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID not found")

	}
	tx, err := zs.store.Transaction(hash)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}

	data, err := tx.Tx.Encode()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var json string
	bz, err := tx.Tx.MarshalJSON()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	json = string(bz)

	res := &zarb.TransactionResponse{
		Data: data,
		Json: json,
	}

	return res, nil

}

func (zs *zarbServer) SendRawTransaction(ctx context.Context, request *zarb.SendRawTransactionRequest) (*zarb.SendRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented yet!")

}
