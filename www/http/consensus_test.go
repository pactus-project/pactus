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

	v1, _ := vote.GenerateTestChangeProposerVote(1000, 0)
	v2, _ := vote.GenerateTestPrepareVote(1000, 1)
	v3, _ := vote.GenerateTestPrecommitVote(1000, 2)
	tMockConsMgr.AddVote(v1)
	tMockConsMgr.AddVote(v2)
	tMockConsMgr.AddVote(v3)
	tMockConsMgr.MoveToNewHeight()
	tMockConsMgr.MoveToNewHeight()

	w := httptest.NewRecorder()
	r := new(http.Request)

	tHTTPServer.ConsensusHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<td>2</td>")
	assert.Contains(t, w.Body.String(), v2.Signer().String())
}
