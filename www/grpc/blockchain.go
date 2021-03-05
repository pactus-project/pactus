package grpc

import (
	"context"

	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetBlockchainInfo(ctx context.Context, request *zarb.BlockchainInfoRequest) (*zarb.BlockchainInfoResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented yet!")
}

func (zs *zarbServer) GetNetworkInfo(ctx context.Context, request *zarb.NetworkInfoRequest) (*zarb.NetworkInfoResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented yet!")
}
