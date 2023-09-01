package wallet

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxOption func(builder *txBuilder) error

func OptionStamp(stamp string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		if stamp != "" {
			stamp, err := hash.StampFromString(stamp)
			if err != nil {
				return err
			}
			builder.stamp = &stamp
		}
		return nil
	}
}

func OptionSequence(seq int32) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.seq = seq
		return nil
	}
}

func OptionFee(fee int64) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.fee = fee
		return nil
	}
}

func OptionMemo(memo string) func(builder *txBuilder) error {
	return func(builder *txBuilder) error {
		builder.memo = memo
		return nil
	}
}

type txBuilder struct {
	client *grpcClient
	stamp  *hash.Stamp
	from   *crypto.Address
	to     *crypto.Address
	pub    *bls.PublicKey
	typ    payload.Type
	seq    int32
	amount int64
	fee    int64
	memo   string
}

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

func (m *txBuilder) setFromAddr(addr string) error {
	from, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.from = &from
	return nil
}

func (m *txBuilder) setToAddress(addr string) error {
	to, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.to = &to
	return nil
}

func (m *txBuilder) build() (*tx.Tx, error) {
	err := m.setStamp()
	if err != nil {
		return nil, err
	}

	err = m.setSequence()
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
		trx = tx.NewTransferTx(*m.stamp, m.seq, *m.from, *m.to, m.amount, m.fee, m.memo)
	case payload.TypeBond:
		{
			pub := m.pub
			val, _ := m.client.getValidator(*m.to)
			if val != nil {
				// validator exists
				pub = nil
			}
			trx = tx.NewBondTx(*m.stamp, m.seq, *m.from, *m.to, pub, m.amount, m.fee, m.memo)
		}
	case payload.TypeUnbond:
		trx = tx.NewUnbondTx(*m.stamp, m.seq, *m.from, m.memo)
	case payload.TypeWithdraw:
		trx = tx.NewWithdrawTx(*m.stamp, m.seq, *m.from, *m.to, m.amount, m.fee, m.memo)
	}

	return trx, nil
}

func (m *txBuilder) setStamp() error {
	if m.stamp == nil {
		if m.client == nil {
			return ErrOffline
		}

		stamp, err := m.client.getStamp()
		if err != nil {
			return err
		}
		m.stamp = &stamp
	}

	return nil
}

func (m *txBuilder) setSequence() error {
	if m.seq == 0 {
		if m.client == nil {
			return ErrOffline
		}

		switch m.typ {
		case payload.TypeTransfer,
			payload.TypeBond:
			{
				acc, err := m.client.getAccount(*m.from)
				if err != nil {
					return err
				}
				m.seq = acc.Sequence + 1
			}

		case payload.TypeUnbond,
			payload.TypeWithdraw:
			{
				val, err := m.client.getValidator(*m.from)
				if err != nil {
					return err
				}
				m.seq = val.Sequence + 1
			}
		}
	}
	return nil
}

func (m *txBuilder) setFee() error {
	if m.fee == 0 {
		if m.client == nil {
			return ErrOffline
		}
		fee, err := m.client.getFee(m.amount, m.typ)
		if err != nil {
			return err
		}
		m.fee = fee
	}
	return nil
}
