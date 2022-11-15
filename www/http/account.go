package http

import (
	"net/http"

	"github.com/gorilla/mux"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// GetAccountHandler returns a handler to get account by address.
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := s.blockchain.GetAccount(s.ctx,
		&pactus.GetAccountRequest{Address: vars["address"]})
	if err != nil {
		s.writeError(w, err)
		return
	}

	acc := res.Account
	tm := newTableMaker()
	tm.addRowAccAddress("Address", acc.Address)
	tm.addRowInt("Number", int(acc.Number))
	tm.addRowInt("Sequence", int(acc.Sequence))
	tm.addRowAmount("Balance", acc.Balance)
	tm.addRowBytes("Hash", acc.Hash)

	s.writeHTML(w, tm.html())
}
