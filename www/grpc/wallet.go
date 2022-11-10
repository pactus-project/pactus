package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type walletServer struct {
	unlockedWallet *wallet.Wallet
	network        wallet.Network
	logger         *logger.Logger
}

func walletPath(name string) string {
	return util.MakeAbs(fmt.Sprintf("wallet%c%s", os.PathSeparator, name))
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

func (s *walletServer) LoadWallet(ctx context.Context,
	req *pactus.LoadWalletRequest) (*pactus.LoadWalletResponse, error) {
	return nil, nil
}

func (s *walletServer) UnloadWallet(ctx context.Context,
	req *pactus.UnloadWalletRequest) (*pactus.UnloadWalletResponse, error) {
	return nil, nil
}

func (s *walletServer) LockWallet(ctx context.Context,
	req *pactus.LockWalletRequest) (*pactus.LockWalletResponse, error) {
	return nil, nil
}

func (s *walletServer) UnlockWallet(ctx context.Context,
	req *pactus.UnlockWalletRequest) (*pactus.UnlockWalletResponse, error) {
	return nil, nil
}
