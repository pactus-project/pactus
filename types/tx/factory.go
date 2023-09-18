package tx

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func NewSubsidyTx(stamp hash.Stamp, lockTime uint32,
	receiver crypto.Address, amount int64, memo string,
) *Tx {
	return NewTransferTx(
		stamp,
		lockTime,
		crypto.TreasuryAddress,
		receiver,
		amount,
		0,
		memo)
}

func NewTransferTx(stamp hash.Stamp, lockTime uint32,
	sender, receiver crypto.Address,
	amount, fee int64, memo string,
) *Tx {
	pld := &payload.TransferPayload{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
	return NewTx(stamp, lockTime, pld, fee, memo)
}

func NewBondTx(stamp hash.Stamp, lockTime uint32,
	sender, receiver crypto.Address,
	pubKey *bls.PublicKey,
	stake, fee int64, memo string,
) *Tx {
	pld := &payload.BondPayload{
		Sender:    sender,
		Receiver:  receiver,
		PublicKey: pubKey,
		Stake:     stake,
	}
	return NewTx(stamp, lockTime, pld, fee, memo)
}

func NewUnbondTx(stamp hash.Stamp, lockTime uint32,
	val crypto.Address,
	memo string,
) *Tx {
	pld := &payload.UnbondPayload{
		Validator: val,
	}
	return NewTx(stamp, lockTime, pld, 0, memo)
}

func NewWithdrawTx(stamp hash.Stamp, lockTime uint32,
	val crypto.Address,
	acc crypto.Address,
	amount, fee int64,
	memo string,
) *Tx {
	pld := &payload.WithdrawPayload{
		From:   val,
		To:     acc,
		Amount: amount,
	}
	return NewTx(stamp, lockTime, pld, fee, memo)
}

func NewSortitionTx(stamp hash.Stamp, lockTime uint32,
	addr crypto.Address,
	proof sortition.Proof,
) *Tx {
	pld := &payload.SortitionPayload{
		Address: addr,
		Proof:   proof,
	}
	return NewTx(stamp, lockTime, pld, 0, "")
}
