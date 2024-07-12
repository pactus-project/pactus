package grpc

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

//
// TODO: default_wallet should be loaded on starting the node.

type walletServer struct {
	*Server
	walletManager *wallet.Manager
}

func newWalletServer(server *Server, manager *wallet.Manager) *walletServer {
	return &walletServer{
		Server:        server,
		walletManager: manager,
	}
}

func (*walletServer) mapHistoryInfo(hi []wallet.HistoryInfo) []*pactus.HistoryInfo {
	historyInfo := make([]*pactus.HistoryInfo, 0)
	for _, hi := range hi {
		historyInfo = append(historyInfo, &pactus.HistoryInfo{
			TransactionId: hi.TxID,
			// Time:          uint32(hi.Time.Unix()),  // TODO: Fix me
			PayloadType: hi.PayloadType,
			Description: hi.Desc,
			Amount:      hi.Amount.ToNanoPAC(),
		})
	}

	return historyInfo
}

func (s *walletServer) GetValidatorAddress(_ context.Context,
	req *pactus.GetValidatorAddressRequest,
) (*pactus.GetValidatorAddressResponse, error) {
	adr, err := s.walletManager.GetValidatorAddress(req.PublicKey)
	if err != nil {
		return nil, err
	}

	return &pactus.GetValidatorAddressResponse{
		Address: adr,
	}, nil
}

func (s *walletServer) CreateWallet(_ context.Context,
	req *pactus.CreateWalletRequest,
) (*pactus.CreateWalletResponse, error) {
	if req.WalletName == "" {
		return nil, fmt.Errorf("wallet name is required")
	}

	mnemonic, err := s.walletManager.CreateWallet(
		req.WalletName, req.Password,
	)
	if err != nil {
		return nil, err
	}

	return &pactus.CreateWalletResponse{
		Mnemonic: mnemonic,
	}, nil
}

func (s *walletServer) RestoreWallet(_ context.Context,
	req *pactus.RestoreWalletRequest,
) (*pactus.RestoreWalletResponse, error) {
	if req.WalletName == "" {
		return nil, fmt.Errorf("wallet name is required")
	}
	if req.Mnemonic == "" {
		return nil, fmt.Errorf("mnemonic is required")
	}

	if err := s.walletManager.RestoreWallet(
		req.WalletName, req.Mnemonic, req.Password,
	); err != nil {
		return nil, err
	}

	return &pactus.RestoreWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) LoadWallet(_ context.Context,
	req *pactus.LoadWalletRequest,
) (*pactus.LoadWalletResponse, error) {
	if err := s.walletManager.LoadWallet(req.WalletName, s.Address()); err != nil {
		return nil, err
	}

	return &pactus.LoadWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) UnloadWallet(_ context.Context,
	req *pactus.UnloadWalletRequest,
) (*pactus.UnloadWalletResponse, error) {
	if err := s.walletManager.UnloadWallet(req.WalletName); err != nil {
		return nil, err
	}

	return &pactus.UnloadWalletResponse{
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) GetTotalBalance(_ context.Context,
	req *pactus.GetTotalBalanceRequest,
) (*pactus.GetTotalBalanceResponse, error) {
	balance, err := s.walletManager.TotalBalance(req.WalletName)
	if err != nil {
		return nil, err
	}

	return &pactus.GetTotalBalanceResponse{
		WalletName:   req.WalletName,
		TotalBalance: balance.ToNanoPAC(),
	}, nil
}

func (s *walletServer) SignRawTransaction(_ context.Context,
	req *pactus.SignRawTransactionRequest,
) (*pactus.SignRawTransactionResponse, error) {
	rawBytes, err := hex.DecodeString(req.RawTransaction)
	if err != nil {
		return nil, err
	}

	id, data, err := s.walletManager.SignRawTransaction(
		req.WalletName, req.Password, rawBytes,
	)
	if err != nil {
		return nil, err
	}

	return &pactus.SignRawTransactionResponse{
		TransactionId:        hex.EncodeToString(id),
		SignedRawTransaction: hex.EncodeToString(data),
	}, nil
}

func (s *walletServer) GetNewAddress(_ context.Context,
	req *pactus.GetNewAddressRequest,
) (*pactus.GetNewAddressResponse, error) {
	data, err := s.walletManager.GetNewAddress(
		req.WalletName,
		req.Label,
		crypto.AddressType(req.AddressType),
	)
	if err != nil {
		return nil, err
	}

	return &pactus.GetNewAddressResponse{
		WalletName: req.WalletName,
		AddressInfo: &pactus.AddressInfo{
			Address:   data.Address,
			Label:     data.Label,
			PublicKey: data.PublicKey,
			Path:      data.Path,
		},
	}, nil
}

func (s *walletServer) GetAddressHistory(_ context.Context,
	req *pactus.GetAddressHistoryRequest,
) (*pactus.GetAddressHistoryResponse, error) {
	data, err := s.walletManager.AddressHistory(req.WalletName, req.Address)
	if err != nil {
		return nil, err
	}

	return &pactus.GetAddressHistoryResponse{
		HistoryInfo: s.mapHistoryInfo(data),
	}, nil
}

func (s *walletServer) SignMessage(_ context.Context,
	req *pactus.SignMessageRequest,
) (*pactus.SignMessageResponse, error) {
	sig, err := s.walletManager.SignMessage(req.Message, req.Password, req.Address, req.WalletName)
	if err != nil {
		return nil, err
	}

	return &pactus.SignMessageResponse{
		Signature: sig,
	}, nil
}
