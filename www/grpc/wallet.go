package grpc

import (
	"context"
	"encoding/hex"
	"errors"

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

func (*walletServer) mapHistoryInfo(his []wallet.HistoryInfo) []*pactus.HistoryInfo {
	historyInfo := make([]*pactus.HistoryInfo, 0)
	for _, info := range his {
		historyInfo = append(historyInfo, &pactus.HistoryInfo{
			TransactionId: info.TxID,
			// Time:          uint32(hi.Time.Unix()),  // TODO: Fix me
			PayloadType: info.PayloadType,
			Description: info.Desc,
			Amount:      info.Amount.ToNanoPAC(),
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
		return nil, errors.New("wallet name is required")
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
		return nil, errors.New("wallet name is required")
	}
	if req.Mnemonic == "" {
		return nil, errors.New("mnemonic is required")
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

	txID, data, err := s.walletManager.SignRawTransaction(
		req.WalletName, req.Password, rawBytes,
	)
	if err != nil {
		return nil, err
	}

	return &pactus.SignRawTransactionResponse{
		TransactionId:        hex.EncodeToString(txID),
		SignedRawTransaction: hex.EncodeToString(data),
	}, nil
}

func (s *walletServer) GetNewAddress(_ context.Context,
	req *pactus.GetNewAddressRequest,
) (*pactus.GetNewAddressResponse, error) {
	data, err := s.walletManager.GetNewAddress(
		req.WalletName,
		req.Label,
		req.Password,
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

func (s *walletServer) GetTotalStake(_ context.Context,
	req *pactus.GetTotalStakeRequest,
) (*pactus.GetTotalStakeResponse, error) {
	stake, err := s.walletManager.TotalStake(req.WalletName)
	if err != nil {
		return nil, err
	}

	return &pactus.GetTotalStakeResponse{
		TotalStake: stake.ToNanoPAC(),
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) GetAddressInfo(_ context.Context,
	req *pactus.GetAddressInfoRequest,
) (*pactus.GetAddressInfoResponse, error) {
	info, err := s.walletManager.GetAddressInfo(req.WalletName, req.Address)
	if err != nil {
		return nil, err
	}

	return &pactus.GetAddressInfoResponse{
		Address:    info.Address,
		Path:       info.Path,
		PublicKey:  info.PublicKey,
		Label:      info.Label,
		WalletName: req.WalletName,
	}, nil
}

func (s *walletServer) SetAddressLabel(_ context.Context,
	req *pactus.SetAddressLabelRequest,
) (*pactus.SetAddressLabelResponse, error) {
	return &pactus.SetAddressLabelResponse{}, s.walletMgr.SetAddressLabel(req.WalletName, req.Address, req.Label)
}

func (s *walletServer) ListWallet(_ context.Context,
	_ *pactus.ListWalletRequest,
) (*pactus.ListWalletResponse, error) {
	wallets, err := s.walletManager.ListWallet()
	if err != nil {
		return nil, err
	}

	return &pactus.ListWalletResponse{
		Wallets: wallets,
	}, nil
}

func (s *walletServer) GetWalletInfo(_ context.Context,
	req *pactus.GetWalletInfoRequest,
) (*pactus.GetWalletInfoResponse, error) {
	info, err := s.walletManager.WalletInfo(req.WalletName)
	if err != nil {
		return nil, err
	}

	return &pactus.GetWalletInfoResponse{
		WalletName: info.WalletName,
		Version:    int32(info.Version),
		Network:    info.Network,
		Encrypted:  info.Encrypted,
		Uuid:       info.UUID,
		CreatedAt:  info.CreatedAt.Unix(),
		DefaultFee: info.DefaultFee.ToNanoPAC(),
	}, nil
}

func (s *walletServer) ListAddress(_ context.Context,
	req *pactus.ListAddressRequest,
) (*pactus.ListAddressResponse, error) {
	addrs, err := s.walletManager.ListAddress(req.WalletName)
	if err != nil {
		return nil, err
	}

	addrsPB := make([]*pactus.AddressInfo, 0, len(addrs))
	for _, addr := range addrs {
		addrsPB = append(addrsPB, &pactus.AddressInfo{
			Address:   addr.Address,
			Label:     addr.Label,
			PublicKey: addr.PublicKey,
			Path:      addr.Path,
		})
	}

	return &pactus.ListAddressResponse{
		Data: addrsPB,
	}, nil
}
