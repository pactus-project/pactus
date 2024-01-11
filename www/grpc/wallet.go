package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loadedWallet struct {
	wallet *wallet.Wallet
}

type walletServer struct {
	*Server
	wallets   map[string]*loadedWallet
	chainType genesis.ChainType
}

func walletPath(name string) string {
	return util.MakeAbs(fmt.Sprintf("wallet%c%s", os.PathSeparator, name))
}

func (s *walletServer) CreateWallet(_ context.Context,
	req *pactus.CreateWalletRequest,
) (*pactus.CreateWalletResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("wallet name is required")
	}

	path := walletPath(req.Name)
	w, err := wallet.Create(path, req.Mnemonic, req.Language, s.chainType)
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

func (s *walletServer) LoadWallet(_ context.Context,
	req *pactus.LoadWalletRequest,
) (*pactus.LoadWalletResponse, error) {
	_, ok := s.wallets[req.Name]
	if !ok {
		return nil, status.Errorf(codes.AlreadyExists, "wallet already loaded")
	}

	path := walletPath(req.Name)
	w, err := wallet.Open(path, true)
	if err != nil {
		return nil, err
	}

	s.wallets[req.Name] = &loadedWallet{wallet: w}

	return &pactus.LoadWalletResponse{
		Name: req.Name,
	}, nil
}

func (s *walletServer) UnloadWallet(_ context.Context,
	req *pactus.UnloadWalletRequest,
) (*pactus.UnloadWalletResponse, error) {
	_, ok := s.wallets[req.Name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	delete(s.wallets, req.Name)

	return &pactus.UnloadWalletResponse{
		Name: req.Name,
	}, nil
}

func (s *walletServer) LockWallet(_ context.Context,
	_ *pactus.LockWalletRequest,
) (*pactus.LockWalletResponse, error) {
	return &pactus.LockWalletResponse{}, nil
}

func (s *walletServer) UnlockWallet(_ context.Context,
	_ *pactus.UnlockWalletRequest,
) (*pactus.UnlockWalletResponse, error) {
	return &pactus.UnlockWalletResponse{}, nil
}
