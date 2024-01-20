package tx

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func NewSubsidyTx(lockTime uint32,
	receiver crypto.Address, amount int64, memo string,
) *Tx {
	return NewTransferTx(
		lockTime,
		crypto.TreasuryAddress,
		receiver,
		amount,
		0,
		memo)
}

func NewTransferTx(lockTime uint32,
	sender, receiver crypto.Address,
	amount, fee int64, memo string,
) *Tx {
	pld := &payload.TransferPayload{
		From:   sender,
		To:     receiver,
		Amount: amount,
	}

	return newTx(lockTime, pld, fee, memo)
}

func NewBondTx(lockTime uint32,
	sender, receiver crypto.Address,
	pubKey *bls.PublicKey,
	stake, fee int64, memo string,
) *Tx {
	pld := &payload.BondPayload{
		From:      sender,
		To:        receiver,
		PublicKey: pubKey,
		Stake:     stake,
	}

	return newTx(lockTime, pld, fee, memo)
}

func NewUnbondTx(lockTime uint32,
	val crypto.Address,
	memo string,
) *Tx {
	pld := &payload.UnbondPayload{
		Validator: val,
	}

	return newTx(lockTime, pld, 0, memo)
}

func NewWithdrawTx(lockTime uint32,
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

	return newTx(lockTime, pld, fee, memo)
}

func NewSortitionTx(lockTime uint32,
	addr crypto.Address,
	proof sortition.Proof,
) *Tx {
	pld := &payload.SortitionPayload{
		Validator: addr,
		Proof:     proof,
	}

	return newTx(lockTime, pld, 0, "")
}
