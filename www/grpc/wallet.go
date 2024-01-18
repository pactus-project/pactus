package grpc

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//
// TODO: default_wallet should be loaded on starting the node.
// TODO: We need to add a wallet service (or manaeger) to manage wallets.
// UnLocking wallet should happens inside the wallet manager (or no?)
//

type walletServer struct {
	*Server
	wallets   map[string]*wallet.Wallet
	chainType genesis.ChainType
}

func newWalletServer(server *Server, chainType genesis.ChainType) *walletServer {
	wallets := make(map[string]*wallet.Wallet)

	return &walletServer{
		Server:    server,
		wallets:   wallets,
		chainType: chainType,
	}
}

func (s *walletServer) walletPath(name string) string {
	return util.MakeAbs(filepath.Join(s.config.WalletsDir, name))
}

func (s *walletServer) CreateWallet(_ context.Context,
	req *pactus.CreateWalletRequest,
) (*pactus.CreateWalletResponse, error) {
	if req.WalletName == "" {
		return nil, fmt.Errorf("wallet name is required")
	}

	walletPath := s.walletPath(req.WalletName)
	w, err := wallet.Create(walletPath, req.Mnemonic, req.Language, s.chainType)
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

	return &pactus.CreateWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) LoadWallet(_ context.Context,
	req *pactus.LoadWalletRequest,
) (*pactus.LoadWalletResponse, error) {
	_, ok := s.wallets[req.WalletName]
	if !ok {
		// TODO: define special codes for errors
		return nil, status.Errorf(codes.AlreadyExists, "wallet already loaded")
	}

	walletPath := s.walletPath(req.WalletName)
	wlt, err := wallet.Open(walletPath, true)
	if err != nil {
		return nil, err
	}

	s.wallets[req.WalletName] = wlt

	return &pactus.LoadWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) UnloadWallet(_ context.Context,
	req *pactus.UnloadWalletRequest,
) (*pactus.UnloadWalletResponse, error) {
	_, ok := s.wallets[req.WalletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	delete(s.wallets, req.WalletName)

	return &pactus.UnloadWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) LockWallet(_ context.Context,
	_ *pactus.LockWalletRequest,
) (*pactus.LockWalletResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemeneted")
}

func (s *walletServer) UnlockWallet(_ context.Context,
	_ *pactus.UnlockWalletRequest,
) (*pactus.UnlockWalletResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemeneted")
}

func (s *walletServer) SignRawTransaction(_ context.Context,
	req *pactus.SignRawTransactionRequest,
) (*pactus.SignRawTransactionResponse, error) {
	wlt, ok := s.wallets[req.WalletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	trx, err := tx.FromBytes(req.RawTransaction)
	if err != nil {
		return nil, err
	}

	err = wlt.SignTransaction(req.Password, trx)
	if err != nil {
		return nil, err
	}

	data, err := trx.Bytes()
	if err != nil {
		return nil, err
	}

	return &pactus.SignRawTransactionResponse{
		TransactionId:        trx.ID().Bytes(),
		SignedRawTransaction: data,
	}, nil
}
