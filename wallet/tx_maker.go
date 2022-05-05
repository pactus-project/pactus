package wallet

import (
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
)

type TxOption func(maker *txMaker) error

func OptionStamp(stamp string) func(maker *txMaker) error {
	return func(maker *txMaker) error {
		if stamp != "" {
			stamp, err := hash.StampFromString(stamp)
			if err != nil {
				return err
			}
			maker.stamp = &stamp
		}
		return nil
	}
}
func OptionSequence(seq int32) func(maker *txMaker) error {
	return func(maker *txMaker) error {
		maker.seq = seq
		return nil
	}
}

func OptionFee(fee int64) func(maker *txMaker) error {
	return func(maker *txMaker) error {
		maker.fee = fee
		return nil
	}
}

func OptionMemo(memo string) func(maker *txMaker) error {
	return func(maker *txMaker) error {
		maker.memo = memo
		return nil
	}
}

type txMaker struct {
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

func newTxMaker(client *grpcClient, options ...TxOption) (*txMaker, error) {
	maker := &txMaker{
		client: client,
	}
	for _, op := range options {
		err := op(maker)
		if err != nil {
			return nil, err
		}
	}
	return maker, nil
}

func (m *txMaker) setFromAddr(addr string) error {
	from, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.from = &from
	return nil
}
func (m *txMaker) setToAddress(addr string) error {
	to, err := crypto.AddressFromString(addr)
	if err != nil {
		return err
	}
	m.to = &to
	return nil
}

func (m *txMaker) build() (*tx.Tx, error) {
	err := m.checkStamp()
	if err != nil {
		return nil, err
	}

	err = m.checkSequence()
	if err != nil {
		return nil, err
	}

	m.checkFee()

	var trx *tx.Tx
	switch m.typ {
	case payload.PayloadTypeSend:
		{
			trx = tx.NewSendTx(*m.stamp, m.seq, *m.from, *m.to, m.amount, m.fee,
				m.memo)
		}
	case payload.PayloadTypeBond:
		{
			trx = tx.NewBondTx(*m.stamp, m.seq, *m.from, *m.to, m.pub, m.amount,
				m.fee, m.memo)
		}
	case payload.PayloadTypeUnbond:
		{
			trx = tx.NewUnbondTx(*m.stamp, m.seq, *m.from, m.memo)
		}
	case payload.PayloadTypeWithdraw:
		{
			trx = tx.NewWithdrawTx(*m.stamp, m.seq, *m.from, *m.to, m.amount,
				m.fee, m.memo)
		}
	}

	return trx, nil
}

func (m *txMaker) checkStamp() error {
	if m.stamp == nil {
		stamp, err := m.client.getStamp()
		if err != nil {
			return err
		}
		m.stamp = &stamp
	}

	return nil
}

func (m *txMaker) checkSequence() error {
	if m.seq == 0 {
		switch m.typ {
		case payload.PayloadTypeSend,
			payload.PayloadTypeBond:
			{
				seq, err := m.client.getAccountSequence(*m.from)
				if err != nil {
					return err
				}
				m.seq = seq
			}

		case payload.PayloadTypeUnbond,
			payload.PayloadTypeWithdraw:
			{
				seq, err := m.client.GetValidatorSequence(*m.from)
				if err != nil {
					return err
				}
				m.seq = seq
			}
		}
	}
	return nil
}

func (m *txMaker) checkFee() {
	if m.fee == 0 {
		switch m.typ {
		case payload.PayloadTypeSend,
			payload.PayloadTypeBond,
			payload.PayloadTypeWithdraw:
			{
				fee := m.amount / 10000
				if fee < 10000 {
					fee = 10000
				}
				m.fee = fee
			}

		case payload.PayloadTypeUnbond:
			{
				m.fee = 0
			}
		}
	}
}
