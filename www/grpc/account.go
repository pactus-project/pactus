package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetAccount(ctx context.Context, request *zarb.AccountRequest) (*zarb.AccountResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %v", err)

	}
	acc := zs.state.Account(addr)
	if acc == nil {
		return nil, status.Errorf(codes.InvalidArgument, "account not found")
	}
	res := &zarb.AccountResponse{
		Account: accountToProto(acc),
	}

	return res, nil
}

func accountToProto(acc *account.Account) *zarb.AccountInfo {
	return &zarb.AccountInfo{
		Address:  acc.Address().String(),
		Number:   acc.Number(),
		Sequence: acc.Sequence(),
		Balance:  acc.Balance(),
	}
}
