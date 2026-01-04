package grpc

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type walletServer struct {
	*Server
	walletManager wltmgr.IManager
}

func newWalletServer(server *Server, manager wltmgr.IManager) *walletServer {
	return &walletServer{
		Server:        server,
		walletManager: manager,
	}
}

func (*walletServer) addressInfoToProto(ai *types.AddressInfo) *pactus.AddressInfo {
	return &pactus.AddressInfo{
		Address:   ai.Address,
		Label:     ai.Label,
		PublicKey: ai.PublicKey,
		Path:      ai.Path,
	}
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
		WalletName: req.WalletName,
		Mnemonic:   mnemonic,
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
	if err := s.walletManager.LoadWallet(req.WalletName); err != nil {
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

func (s *walletServer) ListWallets(_ context.Context,
	req *pactus.ListWalletsRequest,
) (*pactus.ListWalletsResponse, error) {
	wallets, err := s.walletManager.ListWallets(req.IncludeUnloaded)
	if err != nil {
		return nil, err
	}

	return &pactus.ListWalletsResponse{
		Wallets: wallets,
	}, nil
}

func (s *walletServer) IsWalletLoaded(_ context.Context,
	req *pactus.IsWalletLoadedRequest,
) (*pactus.IsWalletLoadedResponse, error) {
	loaded := s.walletManager.IsWalletLoaded(req.WalletName)

	return &pactus.IsWalletLoadedResponse{
		WalletName: req.WalletName,
		Loaded:     loaded,
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

func (s *walletServer) SignMessage(_ context.Context,
	req *pactus.SignMessageRequest,
) (*pactus.SignMessageResponse, error) {
	sig, err := s.walletManager.SignMessage(req.WalletName, req.Password, req.Address, req.Message)
	if err != nil {
		return nil, err
	}

	return &pactus.SignMessageResponse{
		Signature: sig,
	}, nil
}

func (s *walletServer) GetNewAddress(_ context.Context,
	req *pactus.GetNewAddressRequest,
) (*pactus.GetNewAddressResponse, error) {
	info, err := s.walletManager.NewAddress(
		req.WalletName,
		crypto.AddressType(req.AddressType),
		req.Label,
		wallet.WithPassword(req.Password),
	)
	if err != nil {
		return nil, err
	}

	return &pactus.GetNewAddressResponse{
		WalletName:  req.WalletName,
		AddressInfo: s.addressInfoToProto(info),
	}, nil
}

func (s *walletServer) GetAddressInfo(_ context.Context,
	req *pactus.GetAddressInfoRequest,
) (*pactus.GetAddressInfoResponse, error) {
	info, err := s.walletManager.AddressInfo(req.WalletName, req.Address)
	if err != nil {
		return nil, err
	}

	return &pactus.GetAddressInfoResponse{
		WalletName:  req.WalletName,
		AddressInfo: s.addressInfoToProto(info),
	}, nil
}

func (s *walletServer) SetAddressLabel(_ context.Context,
	req *pactus.SetAddressLabelRequest,
) (*pactus.SetAddressLabelResponse, error) {
	err := s.walletManager.SetAddressLabel(req.WalletName, req.Address, req.Label)
	if err != nil {
		return nil, err
	}

	return &pactus.SetAddressLabelResponse{
		WalletName: req.WalletName,
		Address:    req.Address,
		Label:      req.Label,
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
		WalletName: req.WalletName,
		Version:    int32(info.Version),
		Network:    info.Network.String(),
		Encrypted:  info.Encrypted,
		Uuid:       info.UUID,
		CreatedAt:  info.CreatedAt.Unix(),
		DefaultFee: info.DefaultFee.ToNanoPAC(),
		Driver:     info.Driver,
		Path:       info.Path,
	}, nil
}

func (s *walletServer) ListAddresses(_ context.Context,
	req *pactus.ListAddressesRequest,
) (*pactus.ListAddressesResponse, error) {
	addressTypes := make([]crypto.AddressType, 0)
	for _, addrType := range req.AddressTypes {
		addressTypes = append(addressTypes, crypto.AddressType(addrType))
	}

	addrs, err := s.walletManager.ListAddresses(req.WalletName, wallet.WithAddressTypes(addressTypes))
	if err != nil {
		return nil, err
	}

	addrsPB := make([]*pactus.AddressInfo, 0, len(addrs))
	for _, info := range addrs {
		addrsPB = append(addrsPB, s.addressInfoToProto(&info))
	}

	return &pactus.ListAddressesResponse{
		WalletName: req.WalletName,
		Data:       addrsPB,
	}, nil
}

func (s *walletServer) ListTransactions(_ context.Context,
	req *pactus.ListTransactionsRequest,
) (*pactus.ListTransactionsResponse, error) {
	infos, err := s.walletManager.ListTransactions(req.WalletName,
		wallet.WithDirection(wallet.TxDirection(req.Direction)),
		wallet.WithAddress(req.Address),
		wallet.WithCount(int(req.Count)),
		wallet.WithSkip(int(req.Skip)))
	if err != nil {
		return nil, err
	}

	lastBlockHeight := s.state.LastBlockHeight()
	trxs := make([]*pactus.TransactionInfo, 0, len(infos))
	for _, info := range infos {
		trx, _ := tx.FromBytes(info.Data)

		confirmations := 0
		if info.BlockHeight > 0 {
			confirmations = int(lastBlockHeight) - int(info.BlockHeight)
		}
		trxs = append(trxs, transactionToProto(trx, uint32(info.BlockHeight), confirmations))
	}

	return &pactus.ListTransactionsResponse{
		WalletName: req.WalletName,
		Txs:        trxs,
	}, nil
}

func (s *walletServer) UpdatePassword(_ context.Context,
	req *pactus.UpdatePasswordRequest,
) (*pactus.UpdatePasswordResponse, error) {
	err := s.walletManager.UpdatePassword(req.WalletName, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &pactus.UpdatePasswordResponse{
		WalletName: req.WalletName,
	}, nil
}
