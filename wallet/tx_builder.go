package wallet

import (
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
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

// OptionFeeFromString sets the transaction fee using a string input.
func OptionFeeFromString(feeStr string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		if feeStr == "" {
			return nil
		}

		fee, err := amount.FromString(feeStr)
		if err != nil {
			return err
		}
		builder.fee = &fee

		return nil
	}
}

// OptionFee sets the transaction fee using an Amount input.
func OptionFee(fee amount.Amount) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.fee = &fee

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
	client   *grpcClient
	sender   *crypto.Address
	receiver *crypto.Address
	pub      *bls.PublicKey
	typ      payload.Type
	lockTime uint32
	amount   amount.Amount
	fee      *amount.Amount
	memo     string
}

// newTxBuilder initializes a txBuilder with provided options, allowing for flexible configuration of the transaction.
func newTxBuilder(client *grpcClient, options ...TxOption) (*txBuilder, error) {
	builder := &txBuilder{
		client: client,
	}
	for _, op := range options {
		err := op(builder)
		if err != nil {
			return nil, err
		}
	}

	return builder, nil
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

	err = m.setFee()
	if err != nil {
		return nil, err
	}

	var trx *tx.Tx
	switch m.typ {
	case payload.TypeTransfer:
		trx = tx.NewTransferTx(m.lockTime, *m.sender, *m.receiver, m.amount, *m.fee, tx.WithMemo(m.memo))
	case payload.TypeBond:
		pub := m.pub
		val, _ := m.client.getValidator(m.receiver.String())
		if val != nil {
			// validator exists
			pub = nil
		}
		trx = tx.NewBondTx(m.lockTime, *m.sender, *m.receiver, pub, m.amount, *m.fee, tx.WithMemo(m.memo))

	case payload.TypeUnbond:
		trx = tx.NewUnbondTx(m.lockTime, *m.sender, tx.WithMemo(m.memo))

	case payload.TypeWithdraw:
		trx = tx.NewWithdrawTx(m.lockTime, *m.sender, *m.receiver, m.amount, *m.fee, tx.WithMemo(m.memo))

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
		if m.client == nil {
			return ErrOffline
		}

		info, err := m.client.getBlockchainInfo()
		if err != nil {
			return err
		}
		m.lockTime = info.LastBlockHeight + 1
	}

	return nil
}

// setFee determines the fee for the transaction.
// If not set, it retrieves the fee from the client based on amount and transaction type.
func (m *txBuilder) setFee() error {
	if m.fee == nil {
		if m.client == nil {
			return ErrOffline
		}
		fee, err := m.client.getFee(m.amount, m.typ)
		if err != nil {
			return err
		}
		m.fee = &fee
	}

	return nil
}
