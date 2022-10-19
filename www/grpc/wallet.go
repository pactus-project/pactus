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
	req *pactus.GenerateMnemonicRequest) (*pactus.GenerateMnemonicResponse, error) {
	mnemonic := wallet.GenerateMnemonic(int(req.Entropy))

	return &pactus.GenerateMnemonicResponse{
		Mnemonic: mnemonic,
	}, nil
}

func (s *walletServer) CreateWallet(ctx context.Context,
	req *pactus.CreateWalletRequest) (*pactus.CreateWalletResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("wallet name is required")
	}

	path := walletPath(req.Name)
	w, err := wallet.Create(path, req.Mnemonic, req.Language, s.network)
	if err != nil {
		return nil, err
	}
	err = w.UpdatePassword("", req.Password)
	if err != nil {
		return nil, err
	}
	err = w.Save()
	if err != nil {
		return nil, err
	}
	return &pactus.CreateWalletResponse{}, nil
}
