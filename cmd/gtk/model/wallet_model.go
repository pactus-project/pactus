//go:build gtk

package model

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/bech32m"
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

func NewWalletModel(
	ctx context.Context,
	walletClient pactus.WalletClient,
	transactionClient pactus.TransactionClient,
	blockchainClient pactus.BlockchainClient,
	walletName string,
) *WalletModel {
	return &WalletModel{
		ctx:               ctx,
		walletClient:      walletClient,
		transactionClient: transactionClient,
		blockchainClient:  blockchainClient,
		walletName:        walletName,
	}
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

func (model *WalletModel) WalletInfo() (*pactus.GetWalletInfoResponse, error) {
	return model.walletClient.GetWalletInfo(model.ctx, &pactus.GetWalletInfoRequest{
		WalletName: model.walletName,
	})
}

// MigrateWallet migrates a legacy JSON wallet to the SQLite format.
func (model *WalletModel) MigrateWallet() error {
	_, err := model.walletClient.MigrateWallet(model.ctx, &pactus.MigrateWalletRequest{
		WalletName: model.walletName,
	})

	return err
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

func (model *WalletModel) AddressInfo(addr string) *pactus.AddressInfo {
	res, err := model.walletClient.GetAddressInfo(model.ctx, &pactus.GetAddressInfoRequest{
		WalletName: model.walletName,
		Address:    addr,
	})
	if err != nil {
		return nil
	}

	return res.Addr
}

func (model *WalletModel) ListAccountAddresses() []*pactus.AddressInfo {
	lst, _ := model.ListAddresses(
		crypto.AddressTypeBLSAccount,
		crypto.AddressTypeEd25519Account,
		crypto.AddressTypeSecp256k1Account,
	)

	return lst
}

func (model *WalletModel) ListValidatorAddresses() []*pactus.AddressInfo {
	lst, _ := model.ListAddresses(crypto.AddressTypeValidator)

	return lst
}

func (model *WalletModel) ListAddresses(addressTypes ...crypto.AddressType) ([]*pactus.AddressInfo, error) {
	addressTypesPB := make([]pactus.AddressType, 0, len(addressTypes))
	for _, at := range addressTypes {
		addressTypesPB = append(addressTypesPB, pactus.AddressType(at))
	}

	res, err := model.walletClient.ListAddresses(model.ctx, &pactus.ListAddressesRequest{
		WalletName:     model.walletName,
		AddressTypes:   addressTypesPB,
		IncludeBalance: true,
		IncludeStake:   true,
	})
	if err != nil {
		return nil, err
	}

	return res.Addrs, nil
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

func (model *WalletModel) ValidatorInfo(addr string) *pactus.ValidatorInfo {
	res, err := model.blockchainClient.GetValidator(model.ctx, &pactus.GetValidatorRequest{
		Address: addr,
	})
	if err != nil {
		return nil
	}

	return res.Validator
}

func (model *WalletModel) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	res, err := model.walletClient.GetPrivateKey(model.ctx, &pactus.GetPrivateKeyRequest{
		WalletName: model.walletName,
		Password:   password,
		Address:    addr,
	})
	if err != nil {
		return nil, err
	}

	_, typ, _, err := bech32m.DecodeToBase256WithTypeNoLimit(res.PrivateKey)
	if err != nil {
		return nil, err
	}

	switch crypto.SignatureType(typ) {
	case crypto.SignatureTypeBLS:
		return bls.PrivateKeyFromString(res.PrivateKey)
	case crypto.SignatureTypeEd25519:
		return ed25519.PrivateKeyFromString(res.PrivateKey)
	case crypto.SignatureTypeSecp256k1:
		return secp256k1.PrivateKeyFromString(res.PrivateKey)
	default:
		return nil, crypto.InvalidSignatureTypeError(typ)
	}
}

func (model *WalletModel) Mnemonic(password string) (string, error) {
	res, err := model.walletClient.GetMnemonic(model.ctx, &pactus.GetMnemonicRequest{
		WalletName: model.walletName,
		Password:   password,
	})
	if err != nil {
		return "", err
	}

	return res.Mnemonic, nil
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
	_, err := model.walletClient.SetDefaultFee(model.ctx, &pactus.SetDefaultFeeRequest{
		WalletName: model.walletName,
		Amount:     int64(fee),
	})

	return err
}

func (model *WalletModel) NewAddress(
	addressType crypto.AddressType,
	label string,
	password string,
) (*types.AddressInfo, error) {
	res, err := model.walletClient.GetNewAddress(model.ctx, &pactus.GetNewAddressRequest{
		WalletName:  model.walletName,
		AddressType: pactus.AddressType(addressType),
		Label:       label,
		Password:    password,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddressInfo{
		Address:   res.Addr.Address,
		PublicKey: res.Addr.PublicKey,
		Label:     res.Addr.Label,
		Path:      res.Addr.Path,
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

	return res.Addr.Label
}

func (model *WalletModel) SetAddressLabel(addr, label string) error {
	_, err := model.walletClient.SetAddressLabel(model.ctx, &pactus.SetAddressLabelRequest{
		WalletName: model.walletName,
		Address:    addr,
		Label:      label,
	})

	return err
}

func (model *WalletModel) MakeTransferTx(
	sender, receiver string,
	amt amount.Amount,
	fee amount.Amount,
	memo string,
) (*tx.Tx, error) {
	res, err := model.transactionClient.GetRawTransferTransaction(model.ctx,
		&pactus.GetRawTransferTransactionRequest{
			Sender:   sender,
			Receiver: receiver,
			Amount:   int64(amt),
			Fee:      int64(fee),
			Memo:     memo,
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

// Delegation holds the PIP-49 fields of a delegated bond.
type Delegation struct {
	// Owner is the stake owner account. Only it can unbond and withdraw.
	Owner string
	// Share is the owner's part of each block reward, in nano PAC.
	Share amount.Amount
	// Expiry is the height at which the delegation ends.
	Expiry uint32
}

// MakeBondTx builds a bond transaction. A non-nil delegation makes it a
// delegated bond (PIP-49).
func (model *WalletModel) MakeBondTx(
	sender, receiver, publicKey string,
	amt amount.Amount,
	fee amount.Amount,
	memo string,
	delegation *Delegation,
) (*tx.Tx, error) {
	if publicKey == "" {
		valInfo := model.ValidatorInfo(receiver)
		if valInfo == nil {
			// Let's check if we can get public key from the wallet
			info := model.AddressInfo(receiver)
			if info != nil {
				publicKey = info.PublicKey
			}
		}
	}
	req := &pactus.GetRawBondTransactionRequest{
		Sender:    sender,
		Receiver:  receiver,
		PublicKey: publicKey,
		Stake:     int64(amt),
		Fee:       int64(fee),
		Memo:      memo,
	}
	if delegation != nil {
		req.DelegateOwner = delegation.Owner
		req.DelegateShare = int64(delegation.Share)
		req.DelegateExpiry = delegation.Expiry
	}
	res, err := model.transactionClient.GetRawBondTransaction(model.ctx, req)
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

// MakeUnbondTx builds an unbond transaction. For a delegated validator
// (PIP-49), delegateOwner is the stake owner that signs it; otherwise it is empty.
func (model *WalletModel) MakeUnbondTx(validatorAddr, delegateOwner, memo string) (*tx.Tx, error) {
	res, err := model.transactionClient.GetRawUnbondTransaction(model.ctx,
		&pactus.GetRawUnbondTransactionRequest{
			ValidatorAddress: validatorAddr,
			Memo:             memo,
			DelegateOwner:    delegateOwner,
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

// MakeWithdrawTx builds a withdraw transaction. For a delegated validator
// (PIP-49), delegateOwner is the stake owner that signs it; otherwise it is empty.
func (model *WalletModel) MakeWithdrawTx(
	sender, receiver string,
	amt amount.Amount,
	fee amount.Amount,
	memo string,
	delegateOwner string,
) (*tx.Tx, error) {
	res, err := model.transactionClient.GetRawWithdrawTransaction(model.ctx,
		&pactus.GetRawWithdrawTransactionRequest{
			ValidatorAddress: sender,
			AccountAddress:   receiver,
			Amount:           int64(amt),
			Fee:              int64(fee),
			Memo:             memo,
		})
	if err != nil {
		return nil, err
	}

	trx, err := tx.FromString(res.RawTransaction)
	if err != nil || delegateOwner == "" {
		return trx, err
	}

	return withWithdrawDelegateOwner(trx, delegateOwner)
}

// withWithdrawDelegateOwner sets the stake owner of an unsigned withdraw.
// GetRawWithdrawTransaction has no delegate owner field yet, so the wallet adds
// it here. The owner then signs the transaction, as PIP-49 requires.
func withWithdrawDelegateOwner(trx *tx.Tx, delegateOwner string) (*tx.Tx, error) {
	owner, err := crypto.AddressFromString(delegateOwner)
	if err != nil {
		return nil, err
	}
	pld, ok := trx.Payload().(*payload.WithdrawPayload)
	if !ok {
		return nil, errors.New("not a withdraw transaction")
	}
	pld.DelegateOwner = owner
	if err := pld.BasicCheck(); err != nil {
		return nil, err
	}

	// Encode and decode again, so the transaction has the delegation flag
	// and no cached ID of the transaction without the owner.
	raw, err := trx.Bytes()
	if err != nil {
		return nil, err
	}

	return tx.FromBytes(raw)
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

	// The transaction is replaced by the signed one. A transaction holds values that it
	// computes on first use, and they must not be copied: it is decoded again instead.
	signedData, err := signedTx.Bytes()
	if err != nil {
		return err
	}

	return trx.Decode(bytes.NewReader(signedData))
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

func (model *WalletModel) Transactions(count, skip int) ([]*pactus.WalletTransactionInfo, error) {
	res, err := model.walletClient.ListTransactions(model.ctx, &pactus.ListTransactionsRequest{
		WalletName: model.walletName,
		Count:      int32(count),
		Skip:       int32(skip),
	})
	if err != nil {
		return nil, err
	}

	return res.Txs, nil
}

// WalletNetwork returns the network name of the wallet addresses, as used by
// anchor links (PIP-57).
func (*WalletModel) WalletNetwork() string {
	if crypto.AddressHRP() == "pc" {
		return NetworkMainnet
	}

	return NetworkTestnet
}

// ChainStatus reads the last block height of the node and how far it is behind.
// It can wait while the node commits blocks, for example while it syncs, so
// never call it on the UI thread.
func (model *WalletModel) ChainStatus() (ChainStatus, error) {
	res, err := model.blockchainClient.GetBlockchainInfo(model.ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return ChainStatus{}, err
	}

	return ChainStatus{Height: res.LastBlockHeight, BlocksLeft: res.BlocksLeft}, nil
}

// Anchor returns the anchor of the account (PIP-50). found is false when the
// account has none.
func (model *WalletModel) Anchor(addr string) (anchor *pactus.AnchorInfo, found bool, err error) {
	res, err := model.blockchainClient.GetAnchor(model.ctx, &pactus.GetAnchorRequest{
		Address: addr,
	})
	if err != nil {
		return nil, false, err
	}

	return res.Anchor, res.Found, nil
}

// MakeAnchorSetTx builds a transaction that creates or updates the sender's
// anchor. The deposit is added to the amount already locked.
func (model *WalletModel) MakeAnchorSetTx(
	sender string,
	anchorType uint8,
	rootHash []byte,
	manifestURI string,
	deposit amount.Amount,
	fee amount.Amount,
	memo string,
) (*tx.Tx, error) {
	res, err := model.transactionClient.GetRawAnchorTransaction(model.ctx,
		&pactus.GetRawAnchorTransactionRequest{
			Sender:      sender,
			Action:      uint32(payload.AnchorActionSet),
			RootHash:    rootHash,
			ManifestUri: manifestURI,
			AnchorType:  uint32(anchorType),
			Deposit:     int64(deposit),
			Fee:         int64(fee),
			Memo:        memo,
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}

// MakeAnchorDeleteTx builds a transaction that deletes the sender's anchor.
// The locked deposit, minus the fee, returns to the balance.
func (model *WalletModel) MakeAnchorDeleteTx(sender string, fee amount.Amount, memo string) (*tx.Tx, error) {
	res, err := model.transactionClient.GetRawAnchorTransaction(model.ctx,
		&pactus.GetRawAnchorTransactionRequest{
			Sender: sender,
			Action: uint32(payload.AnchorActionDelete),
			Fee:    int64(fee),
			Memo:   memo,
		})
	if err != nil {
		return nil, err
	}

	return tx.FromString(res.RawTransaction)
}
