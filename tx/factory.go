package tx

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx/payload"
)

func NewMintbaseTx(stamp crypto.Hash, receiver crypto.Address, amount int64, memo string) *Tx {
	return NewSendTx(
		stamp,
		crypto.MintbaseAddress,
		receiver,
		amount,
		0,
		memo,
		nil,
		nil)
}

func NewSendTx(stamp crypto.Hash,
	sender, receiver crypto.Address,
	amount, fee int64, memo string,
	publicKey *crypto.PublicKey, signature *crypto.Signature) *Tx {
	return &Tx{
		data: txData{
			Stamp:   stamp,
			Version: 1,
			Type:    payload.PayloadTypeSend,
			Payload: &payload.SendPayload{
				Sender:   sender,
				Receiver: receiver,
				Amount:   amount,
			},
			Fee:       fee,
			Memo:      memo,
			PublicKey: publicKey,
			Signature: signature,
		},
	}
}


