package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/wallet"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

type walletServer struct {
	wallet *wallet.Wallet
	logger *logger.Logger
}

func (s *walletServer) GenerateMnemonic(ctx context.Context,
	request *zarb.GenerateMnemonicRequest) (*zarb.GenerateMnemonicResponse, error) {
	mnemonic := wallet.GenerateMnemonic(int(request.Entropy))

	return &zarb.GenerateMnemonicResponse{
		Mnemonic: mnemonic,
	}, nil
}
