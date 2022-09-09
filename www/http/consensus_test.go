package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestConsensusInfo(t *testing.T) {
	setup(t)

	tMockConsensus.Round = 2
	tMockState.TestStore.LastHeight = 1000
	v1, _ := vote.GenerateTestChangeProposerVote(1000, 0)
	v2, _ := vote.GenerateTestPrepareVote(1000, 1)
	v3, _ := vote.GenerateTestPrecommitVote(1000, 2)
	tMockConsensus.AddVote(v1)
	tMockConsensus.AddVote(v2)
	tMockConsensus.AddVote(v3)

	w := httptest.NewRecorder()
	r := new(http.Request)

	tHTTPServer.ConsensusHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "1001")
	assert.Contains(t, w.Body.String(), v1.Signer().String())
}
