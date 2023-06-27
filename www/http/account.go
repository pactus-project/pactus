package http

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

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
	tm.addRowInt("Number", int(acc.Number))
	tm.addRowInt("Sequence", int(acc.Sequence))
	tm.addRowAmount("Balance", acc.Balance)
	tm.addRowBytes("Hash", acc.Hash)

	s.writeHTML(w, tm.html())
}

// GetAccountByNumberHandler returns a handler to get account by number.
func (s *Server) GetAccountByNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	num, err := strconv.ParseInt(vars["number"], 10, 32)
	if err != nil {
		s.writeError(w, err)
		return
	}

	if num > math.MaxInt32 || num < math.MinInt32 {
		s.writeError(w, fmt.Errorf("integer overflow detected"))
		return
	}

	res, err := s.blockchain.GetAccountByNumber(s.ctx, &pactus.GetAccountByNumberRequest{
		Number: int32(num),
	})
	if err != nil {
		s.writeError(w, err)
		return
	}

	acc := res.Account
	tm := newTableMaker()
	tm.addRowInt("Number", int(acc.Number))
	tm.addRowInt("Sequence", int(acc.Sequence))
	tm.addRowAmount("Balance", acc.Balance)
	tm.addRowBytes("Hash", acc.Hash)

	s.writeHTML(w, tm.html())
}
