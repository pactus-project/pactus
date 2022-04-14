package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetValidatorByNumber(ctx context.Context, request *zarb.ValidatorByNumberRequest) (*zarb.ValidatorResponse, error) {
	val := zs.state.ValidatorByNumber(request.Number)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &zarb.ValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (zs *zarbServer) GetValidator(ctx context.Context, request *zarb.ValidatorRequest) (*zarb.ValidatorResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %s", err.Error())
	}
	val := zs.state.ValidatorByAddress(addr)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &zarb.ValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}
func (zs *zarbServer) GetValidators(ctx context.Context, request *zarb.ValidatorsRequest) (*zarb.ValidatorsResponse, error) {
	validators := zs.state.CommitteeValidators()
	validatorsResp := make([]*zarb.ValidatorInfo, 0)
	for _, val := range validators {
		// TODO: make a function
		// proto validator from native validator
		validatorsResp = append(validatorsResp, validatorToProto(val))
	}
	return &zarb.ValidatorsResponse{Validators: validatorsResp}, nil
}

func validatorToProto(val *validator.Validator) *zarb.ValidatorInfo {
	return &zarb.ValidatorInfo{
		PublicKey:         val.PublicKey().Bytes(),
		Address:           val.Address().String(),
		Number:            val.Number(),
		Sequence:          val.Sequence(),
		Stake:             val.Stake(),
		LastBondingHeight: val.LastBondingHeight(),
		LastJoinedHeight:  val.LastJoinedHeight(),
		UnbondingHeight:   val.UnbondingHeight(),
	}
}
