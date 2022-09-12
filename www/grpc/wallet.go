package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
)

type walletServer struct {
	unlockedWallet *wallet.Wallet
	network        wallet.Network
	logger         *logger.Logger
}

func walletPath(name string) string {
	return util.MakeAbs(fmt.Sprintf("wallet%c%s", os.PathSeparator, name))
}

func (s *walletServer) GenerateMnemonic(ctx context.Context,
	request *pactus.GenerateMnemonicRequest) (*pactus.GenerateMnemonicResponse, error) {
	mnemonic := wallet.GenerateMnemonic(int(request.Entropy))

	return &pactus.GenerateMnemonicResponse{
		Mnemonic: mnemonic,
	}, nil
}

func (s *walletServer) CreateWallet(ctx context.Context,
	request *pactus.CreateWalletRequest) (*pactus.CreateWalletResponse, error) {
	if request.Name == "" {
		return nil, fmt.Errorf("wallet name is required")
	}

	path := walletPath(request.Name)
	w, err := wallet.Create(path, request.Mnemonic, request.Language, s.network)
	if err != nil {
		return nil, err
	}
	err = w.UpdatePassword("", request.Password)
	if err != nil {
		return nil, err
	}
	err = w.Save()
	if err != nil {
		return nil, err
	}
	return &pactus.CreateWalletResponse{}, nil
}
