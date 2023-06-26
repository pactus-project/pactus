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

	val := res.Validator
	tm := newTableMaker()
	tm.addRowString("Public Key", val.PublicKey)
	tm.addRowValAddress("Address", val.Address)
	tm.addRowInt("Number", int(val.Number))
	tm.addRowInt("Sequence", int(val.Sequence))
	tm.addRowAmount("Stake", val.Stake)
	tm.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
	tm.addRowInt("LastJoinedHeight", int(val.LastJoinedHeight))
	tm.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
	tm.addRowBytes("Hash", val.Hash)

	s.writeHTML(w, tm.html())
}

// GetValidatorByNumberHandler returns a handler to get validator by number.
func (s *Server) GetValidatorByNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	num, err := strconv.Atoi(vars["number"])
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

	val := res.Validator
	tm := newTableMaker()
	tm.addRowString("Public Key", val.PublicKey)
	tm.addRowValAddress("Address", val.Address)
	tm.addRowInt("Number", int(val.Number))
	tm.addRowInt("Sequence", int(val.Sequence))
	tm.addRowAmount("Stake", val.Stake)
	tm.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
	tm.addRowInt("LastJoinedHeight", int(val.LastJoinedHeight))
	tm.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
	tm.addRowBytes("Hash", val.Hash)

	s.writeHTML(w, tm.html())
}
