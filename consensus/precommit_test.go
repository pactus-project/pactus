package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(0)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	prop := td.makeProposal(t, h, r)
	propBlockHash := prop.Block().Hash()

	extraSignBytes := vote.VoteTypeCPMainVote.Bytes()
	extraSignBytes = append(extraSignBytes, util.Int16ToSlice(0)...)
	extraSignBytes = append(extraSignBytes, byte(vote.CPValueNo))
	signBytes := certificate.BlockCertificateSignBytes(propBlockHash, h, r)
	signBytes = append(signBytes, extraSignBytes...)
	sigX := td.consX.valKey.Sign(signBytes)
	sigY := td.consY.valKey.Sign(signBytes)
	sigM := td.consM.valKey.Sign(signBytes)
	sig := bls.SignatureAggregate(sigX, sigY, sigM)
	cert := certificate.NewCertificate(h, r, []int32{0, 1, 2, 3, 4, 5}, []int32{2, 3, 5}, sig)
	just := &vote.JustDecided{
		QCert: cert,
	}
	decideVote := vote.NewCPDecidedVote(propBlockHash, h, r, 0, vote.CPValueNo, just, td.consX.valKey.Address())
	td.HelperSignVote(td.consX.valKey, decideVote)

	td.consP.AddVote(decideVote)
	assert.Equal(t, "precommit", td.consP.currentState.name())

	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexX)
	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexY)
	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexM)
	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexN)

	td.shouldPublishQueryProposal(t, td.consP, h)
}
