package types

import (
	"fmt"
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

// WalletInfo represents the information about the wallet.
type WalletInfo struct {
	Version    int
	Driver     string
	Path       string
	Network    genesis.ChainType
	DefaultFee amount.Amount
	UUID       string
	Encrypted  bool
	Neutered   bool
	CreatedAt  time.Time
}

// AddressInfo represents the information about a wallet address.
type AddressInfo struct {
	Address   string    `json:"address"`
	PublicKey string    `json:"public_key"`
	Label     string    `json:"label"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type TransactionStatus int

const (
	TransactionStatusFailed    = TransactionStatus(-1)
	TransactionStatusPending   = TransactionStatus(0)
	TransactionStatusConfirmed = TransactionStatus(1)
)

type TransactionInfo struct {
	ID          string
	Sender      string
	Receiver    string
	Amount      amount.Amount
	Fee         amount.Amount
	Memo        string
	Status      TransactionStatus
	BlockHeight uint32
	PayloadType payload.Type
	Data        []byte
	Comment     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MakeTransactionInfos builds TransactionInfo from trx.
// Returns might contains more than one entry per recipient, this covers the batch transfer transaction.s
// Data should be the serialized tx, otherwise it reurn error.
func MakeTransactionInfos(trx *tx.Tx) ([]*TransactionInfo, error) {
	data, err := trx.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize tx: %w", err)
	}

	receivers := make([]string, 0)
	switch pld := trx.Payload().(type) {
	case *payload.TransferPayload:
		receivers = append(receivers, pld.To.String())
	case *payload.BondPayload:
		receivers = append(receivers, pld.To.String())
	case *payload.UnbondPayload:
		receivers = append(receivers, pld.Validator.String())
	case *payload.WithdrawPayload:
		receivers = append(receivers, pld.To.String())
	case *payload.SortitionPayload:
		receivers = append(receivers, "")
	case *payload.BatchTransferPayload:
		for _, recipient := range pld.Recipients {
			receivers = append(receivers, recipient.To.String())
		}
	}

	infos := make([]*TransactionInfo, len(receivers))
	for i, receiver := range receivers {
		infos[i] = &TransactionInfo{
			ID:          trx.ID().String(),
			Sender:      trx.Payload().Signer().String(),
			Receiver:    receiver,
			Amount:      trx.Payload().Value(),
			Fee:         trx.Fee(),
			Memo:        trx.Memo(),
			Status:      TransactionStatusPending,
			BlockHeight: 0,
			PayloadType: trx.Payload().Type(),
			Data:        data,
		}
	}

	return infos, nil
}
