//go:build gtk

package model

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type WalletModel struct {
	ctx               context.Context
	walletClient      pactus.WalletClient
	transactionClient pactus.TransactionClient
	blockchainClient  pactus.BlockchainClient
	walletName        string
}

// AddressRow is a UI-friendly but UI-agnostic representation of an address entry.
// Formatting (strings/markup) should be done by presenters/controllers, not here.
type AddressRow struct {
	No       int
	Address  string
	Label    string
	Path     string
	Imported bool
	Balance  amount.Amount
	Stake    amount.Amount
}

func NewWalletModel(
	ctx context.Context,
	walletClient pactus.WalletClient,
	transactionClient pactus.TransactionClient,
	blockchainClient pactus.BlockchainClient,
	walletName string,
) (*WalletModel, error) {
	return &WalletModel{
		ctx:               ctx,
		walletClient:      walletClient,
		transactionClient: transactionClient,
		blockchainClient:  blockchainClient,
		walletName:        walletName,
	}, nil
}

// WalletName returns the display name used in the UI.
func (model *WalletModel) WalletName() string {
	return model.walletName
}

func (model *WalletModel) IsEncrypted() bool {
	info, err := model.walletClient.GetWalletInfo(model.ctx, &pactus.GetWalletInfoRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return false
	}

	return info.Encrypted
}

func (model *WalletModel) WalletInfo() (*types.WalletInfo, error) {
	info, err := model.walletClient.GetWalletInfo(model.ctx, &pactus.GetWalletInfoRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return nil, err
	}

	chainType := genesis.Localnet
	switch info.Network {
	case "Mainnet":
		chainType = genesis.Mainnet
	case "Testnet":
		chainType = genesis.Testnet
	}

	return &types.WalletInfo{
		Path:       info.Path,
		Encrypted:  info.Encrypted,
		UUID:       info.Uuid,
		Network:    chainType,
		DefaultFee: amount.Amount(info.DefaultFee),
	}, nil
}

func (model *WalletModel) TotalBalance() (amount.Amount, error) {
	res, err := model.walletClient.GetTotalBalance(model.ctx, &pactus.GetTotalBalanceRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.TotalBalance), nil
}

func (model *WalletModel) TotalStake() (amount.Amount, error) {
	res, err := model.walletClient.GetTotalStake(model.ctx, &pactus.GetTotalStakeRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.TotalStake), nil
}

func (model *WalletModel) AddressInfo(addr string) *types.AddressInfo {
	res, err := model.walletClient.GetAddressInfo(model.ctx, &pactus.GetAddressInfoRequest{
		WalletName: model.walletName,
		Address:    addr,
	})
	if err != nil {
		return nil
	}

	return &types.AddressInfo{
		Address:   res.AddressInfo.Address,
		PublicKey: res.AddressInfo.PublicKey,
		Label:     res.AddressInfo.Label,
		Path:      res.AddressInfo.Path,
	}
}

func (model *WalletModel) ListAddresses(opts ...wallet.ListAddressOption) []types.AddressInfo {
	res, err := model.walletClient.ListAddresses(model.ctx, &pactus.ListAddressesRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return nil
	}

	infos := make([]types.AddressInfo, len(res.Data))
	for i, info := range res.Data {
		infos[i] = types.AddressInfo{
			Address:   info.Address,
			PublicKey: info.PublicKey,
			Label:     info.Label,
			Path:      info.Path,
		}
	}

	return infos
}

func (model *WalletModel) Balance(addr string) (amount.Amount, error) {
	res, err := model.blockchainClient.GetAccount(model.ctx, &pactus.GetAccountRequest{
		Address: addr,
	})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Account.Balance), nil
}

func (model *WalletModel) Stake(addr string) (amount.Amount, error) {
	res, err := model.blockchainClient.GetValidator(model.ctx, &pactus.GetValidatorRequest{
		Address: addr,
	})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Validator.Stake), nil
}

func (model *WalletModel) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return nil, errors.New("not implemented")
}

func (model *WalletModel) Mnemonic(password string) (string, error) {
	return "", errors.New("not implemented")
}

func (model *WalletModel) UpdatePassword(oldPassword, newPassword string) error {
	_, err := model.walletClient.UpdatePassword(model.ctx, &pactus.UpdatePasswordRequest{
		WalletName:  model.walletName,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})

	return err
}

func (model *WalletModel) SetDefaultFee(fee amount.Amount) error {
	return errors.New("not implemented")
}

func (model *WalletModel) NewAddress(
	addressType crypto.AddressType,
	label string,
	opts ...wallet.NewAddressOption,
) (*types.AddressInfo, error) {
	// TODO: handle opts?
	res, err := model.walletClient.GetNewAddress(model.ctx, &pactus.GetNewAddressRequest{
		WalletName:  model.walletName,
		AddressType: pactus.AddressType(addressType),
		Label:       label,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddressInfo{
		Address:   res.AddressInfo.Address,
		PublicKey: res.AddressInfo.PublicKey,
		Label:     res.AddressInfo.Label,
		Path:      res.AddressInfo.Path,
	}, nil
}

func (model *WalletModel) AddressLabel(addr string) string {
	res, err := model.walletClient.GetAddressInfo(model.ctx, &pactus.GetAddressInfoRequest{
		WalletName: model.walletName,
		Address:    addr,
	})
	if err != nil {
		return ""
	}

	return res.AddressInfo.Label
}

func (model *WalletModel) SetAddressLabel(addr, label string) error {
	_, err := model.walletClient.SetAddressLabel(model.ctx, &pactus.SetAddressLabelRequest{
		WalletName: model.walletName,
		Address:    addr,
		Label:      label,
	})

	return err
}

// AddressRows returns typed address rows with domain data only.
func (model *WalletModel) AddressRows() []AddressRow {
	rows := make([]AddressRow, 0)
	res, err := model.walletClient.ListAddresses(model.ctx, &pactus.ListAddressesRequest{
		WalletName: model.walletName,
	})
	if err != nil {
		return rows
	}
	for no, info := range res.Data {
		balance, _ := model.Balance(info.Address)
		stake, _ := model.Stake(info.Address)

		rows = append(rows, AddressRow{
			No:       no + 1,
			Address:  info.Address,
			Label:    info.Label,
			Path:     info.Path,
			Imported: info.Path == "",
			Balance:  balance,
			Stake:    stake,
		})
	}

	return rows
}

func (model *WalletModel) MakeTransferTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	// TODO: handle opts?
	res, err := model.transactionClient.GetRawTransferTransaction(model.ctx,
		&pactus.GetRawTransferTransactionRequest{
			Sender:   sender,
			Receiver: receiver,
			Amount:   int64(amt),
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

func (model *WalletModel) MakeBondTx(
	sender, receiver, publicKey string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	// TODO: handle opts?
	res, err := model.transactionClient.GetRawBondTransaction(model.ctx,
		&pactus.GetRawBondTransactionRequest{
			Sender:    sender,
			Receiver:  receiver,
			PublicKey: publicKey,
			Stake:     int64(amt),
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

func (model *WalletModel) MakeUnbondTx(validator string, opts ...wallet.TxOption) (*tx.Tx, error) {
	// TODO: handle opts?
	res, err := model.transactionClient.GetRawUnbondTransaction(model.ctx,
		&pactus.GetRawUnbondTransactionRequest{
			ValidatorAddress: validator,
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

func (model *WalletModel) MakeWithdrawTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	// TODO: handle opts?
	res, err := model.transactionClient.GetRawWithdrawTransaction(model.ctx,
		&pactus.GetRawWithdrawTransactionRequest{
			ValidatorAddress: sender,
			AccountAddress:   receiver,
			Amount:           int64(amt),
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

func (model *WalletModel) SignTransaction(password string, trx *tx.Tx) error {
	raw, err := trx.Bytes()
	if err != nil {
		return err
	}
	res, err := model.walletClient.SignRawTransaction(model.ctx, &pactus.SignRawTransactionRequest{
		WalletName:     model.walletName,
		RawTransaction: hex.EncodeToString(raw),
		Password:       password,
	})
	if err != nil {
		return err
	}

	signedTx, err := tx.FromString(res.SignedRawTransaction)
	if err != nil {
		return err
	}

	*trx = *signedTx

	return nil
}

func (model *WalletModel) BroadcastTransaction(trx *tx.Tx) (string, error) {
	raw, err := trx.Bytes()
	if err != nil {
		return "", err
	}
	res, err := model.transactionClient.BroadcastTransaction(model.ctx, &pactus.BroadcastTransactionRequest{
		SignedRawTransaction: hex.EncodeToString(raw),
	})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (model *WalletModel) Transactions(count, skip int) []*pactus.TransactionInfo {
	res, err := model.walletClient.ListTransactions(model.ctx, &pactus.ListTransactionsRequest{
		WalletName: model.walletName,
		Count:      int32(count),
		Skip:       int32(skip),
	})
	if err != nil {
		return nil
	}

	return res.Txs
}
