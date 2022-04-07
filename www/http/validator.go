package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	b := s.capnp.GetValidator(s.ctx, func(p capnp.ZarbServer_getValidator_Params) error {
		vars := mux.Vars(r)
		return p.SetAddress(vars["address"])
	})

	a, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := a.Result()
	d, _ := res.Data()
	val, err := validator.FromBytes(d)
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowString("Public Key", val.PublicKey().String())
	tm.addRowValAddress("Address", val.Address().String())
	tm.addRowInt("Number", int(val.Number()))
	tm.addRowInt("Sequence", int(val.Sequence()))
	tm.addRowInt("Stake", int(val.Stake()))
	tm.addRowBytes("Hash", val.Hash().Bytes())
	tm.addRowBytes("Data", d)

	s.writeHTML(w, tm.html())
}
