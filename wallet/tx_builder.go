package wallet

import (
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet/provider"
)

// TxOption defines a function type used to apply options to a txBuilder.
type TxOption func(builder *txBuilder) error

// OptionLockTime sets the lock time for the transaction.
func OptionLockTime(lockTime uint32) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.lockTime = lockTime

		return nil
	}
}

// OptionFee sets the transaction fee using a string input.
func OptionFee(feeStr string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		if feeStr == "" {
			return nil
		}

		fee, err := amount.FromString(feeStr)
		if err != nil {
			return err
		}
		builder.fee = fee

		return nil
	}
}

// OptionMemo sets a memo or note for the transaction.
func OptionMemo(memo string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.memo = memo

		return nil
	}
}

// txBuilder helps build and configure a transaction before submitting it.
type txBuilder struct {
	provider provider.IBlockchainProvider
	sender   *crypto.Address
	receiver *crypto.Address
	pub      *bls.PublicKey
	typ      payload.Type
	lockTime uint32
	amount   amount.Amount
	fee      amount.Amount
	memo     string
}

// setSenderAddr sets the sender's address for the transaction.
func (m *txBuilder) setSenderAddr(addr string) error {
	sender, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.sender = &sender

	return nil
}

// setReceiverAddress sets the recipient's address for the transaction.
func (m *txBuilder) setReceiverAddress(addr string) error {
	receiver, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.receiver = &receiver

	return nil
}

// build constructs and finalizes the transaction, selecting the appropriate type based on the builder's configuration.
func (m *txBuilder) build() (*tx.Tx, error) {
	err := m.setLockTime()
	if err != nil {
		return nil, err
	}
	var trx *tx.Tx
	switch m.typ {
	case payload.TypeTransfer:
		trx = tx.NewTransferTx(m.lockTime, *m.sender, *m.receiver, m.amount, m.fee, tx.WithMemo(m.memo))
	case payload.TypeBond:
		pub := m.pub
		_, err := m.provider.GetValidator(m.receiver.String())
		if err != nil {
			// validator exists
			pub = nil
		}
		trx = tx.NewBondTx(m.lockTime, *m.sender, *m.receiver, pub, m.amount, m.fee, tx.WithMemo(m.memo))

	case payload.TypeUnbond:
		trx = tx.NewUnbondTx(m.lockTime, *m.sender, tx.WithMemo(m.memo))

	case payload.TypeWithdraw:
		trx = tx.NewWithdrawTx(m.lockTime, *m.sender, *m.receiver, m.amount, m.fee, tx.WithMemo(m.memo))

	case payload.TypeBatchTransfer:
		return nil, errors.New("BatchTransfer is not implemented yet")

	case payload.TypeSortition:
		return nil, errors.New("unable to build sortition transactions")
	}

	return trx, nil
}

// setLockTime assigns a lock time to the transaction.
// If not provided, it retrieves the last block height and increments it.
func (m *txBuilder) setLockTime() error {
	if m.lockTime == 0 {
		height, err := m.provider.LastBlockHeight()
		if err != nil {
			return err
		}
		m.lockTime = uint32(height + 1)
	}

	return nil
}
