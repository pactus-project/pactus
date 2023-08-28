package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsensusInfo(t *testing.T) {
	td := setup(t)

	v1, _ := td.GenerateTestPrepareVote(1000, 1)
	v2, _ := td.GenerateTestPrecommitVote(1000, 2)
	td.mockConsMgr.AddVote(v1)
	td.mockConsMgr.AddVote(v2)
	td.mockConsMgr.MoveToNewHeight()
	td.mockConsMgr.MoveToNewHeight()

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.ConsensusHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<td>2</td>")
	assert.Contains(t, w.Body.String(), v2.Signer().String())
}
