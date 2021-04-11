package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetValidator(ctx context.Context, request *zarb.ValidatorRequest) (*zarb.ValidatorResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Validator Address:%s", err.Error())
	}
	validator := zs.state.Validator(addr)
	if validator == nil {
		return nil, status.Errorf(codes.NotFound, "NotFound Validator Address")
	}

	data, err := validator.Encode()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var json string
	if request.Verbosity == 1 {
		bz, err := validator.MarshalJSON()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		json = string(bz)
	}

	return &zarb.ValidatorResponse{
		Data: data,
		Json: json,
	}, nil
}

func (zs *zarbServer) GetValidatorByNumber(ctx context.Context, request *zarb.ValidatorByNumberRequest) (*zarb.ValidatorResponse, error) {
	validator := zs.state.ValidatorByNumber(int(request.Number))
	if validator == nil {
		return nil, status.Errorf(codes.NotFound, "NotFound Validator Address")
	}

	data, err := validator.Encode()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var json string
	if request.Verbosity == 1 {
		bz, err := validator.MarshalJSON()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		json = string(bz)
	}

	return &zarb.ValidatorResponse{
		Data: data,
		Json: json,
	}, nil
}
