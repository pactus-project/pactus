package tx

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func NewSubsidyTx(lockTime types.Height,
	recipients []payload.BatchRecipient, opts ...TxOption,
) *Tx {
	return NewBatchTransferTx(
		lockTime,
		crypto.TreasuryAddress,
		recipients,
		0,
		opts...,
	)
}

//revive:disable-next-line:argument-limit
func NewAnchorTx(lockTime types.Height, from crypto.Address, action uint8,
	rootHash []byte, manifestURI string, anchorType uint8,
	deposit, fee amount.Amount, opts ...TxOption,
) *Tx {
	pld := &payload.AnchorPayload{
		From:        from,
		Action:      action,
		RootHash:    append([]byte(nil), rootHash...),
		ManifestURI: manifestURI,
		AnchorType:  anchorType,
		Deposit:     deposit,
	}

	return newTx(lockTime, pld, fee, opts...)
}

func NewTransferTx(lockTime types.Height,
	sender, receiver crypto.Address,
	amt, fee amount.Amount, opts ...TxOption,
) *Tx {
	pld := &payload.TransferPayload{
		From:   sender,
		To:     receiver,
		Amount: amt,
	}

	return newTx(lockTime, pld, fee, opts...)
}

func NewBatchTransferTx(lockTime types.Height,
	sender crypto.Address, recipients []payload.BatchRecipient,
	fee amount.Amount, opts ...TxOption,
) *Tx {
	pld := &payload.BatchTransferPayload{
		From:       sender,
		Recipients: recipients,
	}

	return newTx(lockTime, pld, fee, opts...)
}

func NewBondTx(lockTime types.Height,
	sender, receiver crypto.Address,
	pubKey *bls.PublicKey,
	stake, fee amount.Amount, opts ...TxOption,
) *Tx {
	pld := &payload.BondPayload{
		From:      sender,
		To:        receiver,
		PublicKey: pubKey,
		Stake:     stake,
	}

	return newTx(lockTime, pld, fee, opts...)
}

func NewUnbondTx(lockTime types.Height,
	val crypto.Address,
	opts ...TxOption,
) *Tx {
	pld := &payload.UnbondPayload{
		Validator: val,
	}

	return newTx(lockTime, pld, 0, opts...)
}

func NewWithdrawTx(lockTime types.Height,
	val, acc crypto.Address,
	amt, fee amount.Amount,
	opts ...TxOption,
) *Tx {
	pld := &payload.WithdrawPayload{
		From:   val,
		To:     acc,
		Amount: amt,
	}

	return newTx(lockTime, pld, fee, opts...)
}

func NewSortitionTx(lockTime types.Height,
	addr crypto.Address,
	proof sortition.Proof,
) *Tx {
	pld := &payload.SortitionPayload{
		Address: addr,
		Proof:   proof,
	}

	return newTx(lockTime, pld, 0)
}
