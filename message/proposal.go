package message

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

type ProposalPayload struct {
	Proposal vote.Proposal `cbor:"1,keyasint"`
	Txs      []tx.Tx       `cbor:"2,keyasint"`
}

func NewProposalMessage(proposal vote.Proposal, txs []tx.Tx) Message {
	return Message{
		Type: PayloadTypeProposal,
		Payload: &ProposalPayload{
			Proposal: proposal,
			Txs:      txs,
		},
	}
}

func (p *ProposalPayload) SanityCheck() error {
	if err := p.Proposal.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}
	// TODO: compare tx.hash() with tx_hashes insied block
	if p.Proposal.Block().TxHashes().Count() != len(p.Txs) {
		return errors.Errorf(errors.ErrInvalidMessage, "Not enough transactions")
	}
	return nil
}

func (p *ProposalPayload) Type() PayloadType {
	return PayloadTypeProposal
}

func (p *ProposalPayload) Fingerprint() string {
	return p.Proposal.Fingerprint()
}
