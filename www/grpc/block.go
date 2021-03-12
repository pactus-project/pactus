package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetBlockHeight(ctx context.Context, request *zarb.BlockHeightRequest) (*zarb.BlockHeightResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented yet!")
}

func (zs *zarbServer) GetBlock(ctx context.Context, request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	height := request.GetHeight()
	block := zs.state.Block(int(height))
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Block not found")
	}
	hash := block.Hash().String()
	timestamp, err := ptypes.TimestampProto(block.Header().Time())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
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
