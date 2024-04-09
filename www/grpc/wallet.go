package grpc

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//
// TODO: default_wallet should be loaded on starting the node.
// TODO: We need to add a wallet service (or manager) to manage wallets.
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

func (s *walletServer) GetValidatorAddress(_ context.Context,
	req *pactus.GetValidatorAddressRequest,
) (*pactus.GetValidatorAddressResponse, error) {
	pubKey, err := bls.PublicKeyFromString(req.PublicKey)
	if err != nil {
		return nil, err
	}

	return &pactus.GetValidatorAddressResponse{
		Address: pubKey.ValidatorAddress().String(),
	}, nil
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
	if ok {
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

func (s *walletServer) GetTotalBalance(_ context.Context,
	req *pactus.GetTotalBalanceRequest,
) (*pactus.GetTotalBalanceResponse, error) {
	wlt, ok := s.wallets[req.WalletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	totalBalance := wlt.TotalBalance()

	return &pactus.GetTotalBalanceResponse{
		WalletName:   req.WalletName,
		TotalBalance: totalBalance.ToNanoPAC(),
	}, nil
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

func (s *walletServer) GetNewAddress(_ context.Context,
	req *pactus.GetNewAddressRequest,
) (*pactus.GetNewAddressResponse, error) {
	wlt, ok := s.wallets[req.WalletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	var addressInfo *vault.AddressInfo
	switch req.AddressType {
	case pactus.AddressType(crypto.AddressTypeBLSAccount):
		info, err := wlt.NewBLSAccountAddress(req.Label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case pactus.AddressType(crypto.AddressTypeValidator):
		info, err := wlt.NewValidatorAddress(req.Label)
		if err != nil {
			return nil, err
		}
		addressInfo = info

	case pactus.AddressType(crypto.AddressTypeTreasury):
		return nil, status.Errorf(codes.InvalidArgument, "invalid address type")

	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid address type")
	}

	return &pactus.GetNewAddressResponse{
		WalletName: req.WalletName,
		AddressInfo: &pactus.AddressInfo{
			Address:   addressInfo.Address,
			PublicKey: addressInfo.PublicKey,
			Label:     addressInfo.Label,
			Path:      addressInfo.Path,
		},
	}, nil
}

func (s *walletServer) GetAddressHistory(_ context.Context,
	req *pactus.GetAddressHistoryRequest,
) (*pactus.GetAddressHistoryResponse, error) {
	wlt, ok := s.wallets[req.WalletName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "wallet is not loaded")
	}

	historyInfo := wlt.GetHistory(req.Address)

	return &pactus.GetAddressHistoryResponse{
		HistoryInfo: s.mapHistoryInfo(historyInfo),
	}, nil
}

func (s *walletServer) mapHistoryInfo(hi []wallet.HistoryInfo) []*pactus.HistoryInfo {
	historyInfo := make([]*pactus.HistoryInfo, 0)
	for _, hi := range hi {
		historyInfo = append(historyInfo, &pactus.HistoryInfo{
			TransactionId: hi.TxID,
			Time:          uint32(hi.Time.Unix()),
			PayloadType:   hi.PayloadType,
			Description:   hi.Desc,
			Amount:        hi.Amount.ToNanoPAC(),
		})
	}

	return historyInfo
}
