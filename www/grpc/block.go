package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (zs *zarbServer) GetBlockHeight(ctx context.Context, request *zarb.BlockHeightRequest) (*zarb.BlockHeightResponse, error) {
	h, err := crypto.HashFromString(request.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Hash provided is not Valid")
	}
	height := zs.state.BlockHeight(h)
	if height == -1 {
		return nil, status.Errorf(codes.NotFound, "No block found with the Hash provided")
	}
	return &zarb.BlockHeightResponse{
		Height: int64(height),
	}, nil
}

func (zs *zarbServer) GetBlock(ctx context.Context, request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	height := request.GetHeight()
	block := zs.state.Block(int(height))
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Block not found")
	}
	hash := block.Hash().String()
	timestamp := timestamppb.New(block.Header().Time())
	data, err := block.Encode()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var json string
	if request.Verbosity == 1 {
		bz, err := block.MarshalJSON()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		json = string(bz)
	}
	res := &zarb.BlockResponse{
		Hash:      hash,
		BlockTime: timestamp,
		Data:      data,
		Json:      json,
	}

	return res, nil
}
