package http

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type BlockchainResult struct {
	LastBlockHeight int32
	LastBlockHash   string
}
type BlockResult struct {
	Hash        string
	Time        time.Time
	Data        string
	Header      BlockHeaderResult
	Certificate BlockCertResult
	Txs         []*TransactionResult
}

type BlockHeaderResult struct {
	Version         uint8
	UnixTime        uint32
	PrevBlockHash   string
	StateRoot       string
	SortitionSeed   string
	ProposerAddress string
}

type BlockCertResult struct {
	Round      int16
	Committers []int32
	Absentees  []int32
	Signature  string
}

type TransactionResult struct {
	ID        string
	Data      string
	Version   uint8
	Stamp     string
	Sequence  int32
	Fee       int64
	Payload   string
	Memo      string
	PublicKey string
	Signature string
}

type TransactionSendPayloadResult struct {
	Sender   string
	Receiver string
	Amount   int64
}

type TransactionBondPayloadResult struct {
	Sender    string
	PublicKey string
	Stake     int64
}

type TransactionSortitionPayloadResult struct {
	Address string
	Proof   string
}

type SendTranscationResult struct {
	Status int
	ID     hash.Hash
}

type NetworkResult struct {
	SelfID peer.ID
	Peers  []*peerset.Peer
}

func txToResult(trx *tx.Tx) *TransactionResult {
	pldStr := ""
	switch trx.Payload().Type() {
	case payload.PayloadTypeBond:
		pld := new(TransactionBondPayloadResult)
		pld.PublicKey = trx.Payload().(*payload.BondPayload).PublicKey.String()
		pld.Sender = trx.Payload().(*payload.BondPayload).Sender.String()
		pld.Stake = trx.Payload().(*payload.BondPayload).Stake
		b, _ := json.Marshal(pld)
		pldStr = string(b)

	case payload.PayloadTypeSend:
		pld := new(TransactionSendPayloadResult)
		pld.Sender = trx.Payload().(*payload.SendPayload).Sender.String()
		pld.Receiver = trx.Payload().(*payload.SendPayload).Receiver.String()
		pld.Amount = trx.Payload().(*payload.SendPayload).Amount
		b, _ := json.Marshal(pld)
		pldStr = string(b)

	case payload.PayloadTypeSortition:
		pld := new(TransactionSortitionPayloadResult)
		pld.Address = trx.Payload().(*payload.SortitionPayload).Address.String()
		pld.Proof = hex.EncodeToString(trx.Payload().(*payload.SortitionPayload).Proof[:])
		b, _ := json.Marshal(pld)
		pldStr = string(b)
	}

	d, _ := trx.Bytes()
	out := new(TransactionResult)
	out.ID = trx.ID().String()
	out.Data = hex.EncodeToString(d)
	out.Version = trx.Version()
	out.Stamp = trx.Stamp().String()
	out.Sequence = trx.Sequence()
	out.Fee = trx.Fee()
	out.Memo = trx.Memo()
	out.Payload = pldStr
	out.PublicKey = trx.PublicKey().String()
	out.Signature = trx.Signature().String()
	return out
}
