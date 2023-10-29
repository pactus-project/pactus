package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// GetValidatorHandler returns a handler to get validator by address.
func (s *Server) GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, err := s.blockchain.GetValidator(s.ctx,
		&pactus.GetValidatorRequest{Address: vars["address"]})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := s.writeValidatorTable(res.Validator)
	s.writeHTML(w, tm.html())
}

// GetValidatorByNumberHandler returns a handler to get validator by number.
func (s *Server) GetValidatorByNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	num, err := strconv.ParseInt(vars["number"], 10, 32)
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, err := s.blockchain.GetValidatorByNumber(s.ctx, &pactus.GetValidatorByNumberRequest{
		Number: int32(num),
	})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := s.writeValidatorTable(res.Validator)
	s.writeHTML(w, tm.html())
}

func (s *Server) writeValidatorTable(val *pactus.ValidatorInfo) *tableMaker {
	tm := newTableMaker()
	tm.addRowString("Public Key", val.PublicKey)
	tm.addRowValAddress("Address", val.Address)
	tm.addRowInt("Number", int(val.Number))
	tm.addRowAmount("Stake", val.Stake)
	tm.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
	tm.addRowInt("LastSortitionHeight", int(val.LastSortitionHeight))
	tm.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
	tm.addRowBytes("Hash", val.Hash)

	return tm
}
