package wallet

import (
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet/provider"
)

// TxOption defines a function type used to apply options to a txBuilder.
type TxOption func(builder *txBuilder) error

// OptionLockTime sets the lock time for the transaction.
func OptionLockTime(lockTime types.Height) func(builder *txBuilder) error {
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

// OptionDelegateOwner sets delegation owner address for bond/unbond transactions.
func OptionDelegateOwner(owner string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		if owner == "" {
			return nil
		}

		delegateOwner, err := crypto.AddressFromString(owner)
		if err != nil {
			return err
		}
		builder.delegateOwner = &delegateOwner

		return nil
	}
}

// OptionDelegateShare sets delegation owner reward share for delegated bond transactions.
func OptionDelegateShare(shareStr string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		if shareStr == "" {
			return nil
		}

		share, err := amount.FromString(shareStr)
		if err != nil {
			return err
		}
		builder.delegateShare = share

		return nil
	}
}

// OptionDelegateExpiry sets delegation expiry height for delegated bond transactions.
func OptionDelegateExpiry(expiry types.Height) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.delegateExpiry = expiry

		return nil
	}
}

// txBuilder helps build and configure a transaction before submitting it.
type txBuilder struct {
	provider       provider.IBlockchainProvider
	sender         *crypto.Address
	receiver       *crypto.Address
	pub            *bls.PublicKey
	delegateOwner  *crypto.Address
	delegateShare  amount.Amount
	delegateExpiry types.Height
	typ            payload.Type
	lockTime       types.Height
	amount         amount.Amount
	fee            amount.Amount
	memo           string
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

// setPublicKey sets the public key for the transaction, typically used in bond transactions.
func (m *txBuilder) setPublicKey(pubStr string) error {
	pub, err := bls.PublicKeyFromString(pubStr)
	if err != nil {
		return err
	}
	m.pub = pub

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
		val, _ := m.provider.GetValidator(m.receiver.String())
		if val != nil {
			// validator exists
			pub = nil
		}
		trx = tx.NewBondTx(m.lockTime, *m.sender, *m.receiver, pub, m.amount, m.fee, tx.WithMemo(m.memo))
		if m.delegateOwner != nil {
			bondPld := trx.Payload().(*payload.BondPayload)
			bondPld.DelegateOwner = *m.delegateOwner
			bondPld.DelegateShare = m.delegateShare
			bondPld.DelegateExpiry = m.delegateExpiry
		}

	case payload.TypeUnbond:
		trx = tx.NewUnbondTx(m.lockTime, *m.sender, tx.WithMemo(m.memo))
		if m.delegateOwner != nil {
			unbondPld := trx.Payload().(*payload.UnbondPayload)
			unbondPld.DelegateOwner = *m.delegateOwner
		}

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
		m.lockTime = height + 1
	}

	return nil
}
