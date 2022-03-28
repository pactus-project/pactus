package tx

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx/payload"
)

func NewMintbaseTx(stamp hash.Stamp, seq int32, receiver crypto.Address, amount int64, memo string) *Tx {
	return NewSendTx(
		stamp,
		seq,
		crypto.TreasuryAddress,
		receiver,
		amount,
		0,
		memo)
}

func NewSendTx(stamp hash.Stamp,
	seq int32,
	sender, receiver crypto.Address,
	amount, fee int64, memo string) *Tx {
	pld := &payload.SendPayload{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
	return NewTx(stamp, seq, pld, fee, memo)
}

func NewBondTx(stamp hash.Stamp,
	seq int32,
	sender crypto.Address,
	val *bls.PublicKey,
	stake, fee int64, memo string) *Tx {
	pld := &payload.BondPayload{
		Sender:    sender,
		PublicKey: val,
		Stake:     stake,
	}
	return NewTx(stamp, seq, pld, fee, memo)
}

func NewUnbondTx(stamp hash.Stamp,
	seq int32,
	val crypto.Address,
	memo string) *Tx {
	pld := &payload.UnbondPayload{
		Validator: val,
	}
	return NewTx(stamp, seq, pld, 0, memo)
}

func NewWithdrawTx(stamp hash.Stamp,
	seq int32,
	val crypto.Address,
	acc crypto.Address,
	amount, fee int64,
	memo string) *Tx {
	pld := &payload.WithdrawPayload{
		From:   val,
		To:     acc,
		Amount: amount,
	}
	return NewTx(stamp, seq, pld, fee, memo)
}

func NewSortitionTx(stamp hash.Stamp,
	seq int32,
	addr crypto.Address,
	proof sortition.Proof) *Tx {
	pld := &payload.SortitionPayload{
		Address: addr,
		Proof:   proof,
	}
	return NewTx(stamp, seq, pld, 0, "")
}
