package types

import (
	"fmt"
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
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

func (ts TransactionStatus) String() string {
	switch ts {
	case TransactionStatusFailed:
		return "failed"
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusConfirmed:
		return "confirmed"
	default:
		return "unknown"
	}
}

// TxDirection indicates whether to include incoming or outgoing transactions.
type TxDirection int

const (
	// TxDirectionAny includes both incoming and outgoing transactions.
	TxDirectionAny TxDirection = 0
	// TxDirectionIncoming includes only incoming transactions where the wallet receives funds.
	TxDirectionIncoming = 1
	// TxDirectionOutgoing includes only outgoing transactions where the wallet sends funds.
	TxDirectionOutgoing = 2
)

type TransactionInfo struct {
	ID          string
	Sender      string
	Receiver    string
	Direction   TxDirection
	Amount      amount.Amount
	Fee         amount.Amount
	Memo        string
	Status      TransactionStatus
	BlockHeight block.Height
	PayloadType payload.Type
	Data        []byte
	Comment     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MakeTransactionInfos builds TransactionInfo from trx.
// Returns might contains more than one entry per recipient, this covers the batch transfer transaction.s
// Data should be the serialized tx, otherwise it return error.
func MakeTransactionInfos(trx *tx.Tx, status TransactionStatus, blockHeight block.Height) ([]*TransactionInfo, error) {
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
			Direction:   TxDirectionAny,
			Amount:      trx.Payload().Value(),
			Fee:         trx.Fee(),
			Memo:        trx.Memo(),
			Status:      status,
			BlockHeight: blockHeight,
			PayloadType: trx.Payload().Type(),
			Data:        data,
		}
	}

	return infos, nil
}
