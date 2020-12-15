package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	b := s.server.GetValidator(s.ctx, func(p capnp.ZarbServer_getValidator_Params) error {
		vars := mux.Vars(r)
		if err := p.SetAddress([]byte(vars["address"])); err != nil {
			return err
		}
		return nil
	})

	a, err := b.Struct()
	if err != nil {
		if _, err = io.WriteString(w, err.Error()); err != nil {
			s.logger.Error("Unable to write string", "err", err)
		}
		return
	}

	res, _ := a.Result()
	d, _ := res.Data()
	val := new(validator.Validator)
	err = val.Decode(d)
	if err != nil {
		s.logger.Error("Unable to decode account data", "err", err)
		return
	}

	j, _ := json.MarshalIndent(val, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, string(j)); err != nil {
		s.logger.Error("Unable to write string", "err", err)
	}
}
