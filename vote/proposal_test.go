package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/block"
)

func TestMarshaling(t *testing.T) {
	b1, _, pv := block.GenerateTestBlock()
	p1 := NewProposal(100, 5, b1)
	sig := pv.Sign(p1.SignBytes())
	p1.SetSignature(sig)

	bz1, err := p1.MarshalCBOR()
	assert.NoError(t, err)
	var p2 Proposal
	err = p2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	bz2, _ := p2.MarshalCBOR()

	assert.Equal(t, bz1, bz2)
	assert.Equal(t, p1.Hash(), p2.Hash())
}
