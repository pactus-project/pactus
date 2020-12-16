package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/www/capnp"
)

// GetAccountHandler returns a handler to get account by address
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	b := s.server.GetAccount(s.ctx, func(p capnp.ZarbServer_getAccount_Params) error {
		vars := mux.Vars(r)
		if err := p.SetAddress([]byte(vars["address"])); err != nil {
			return err
		}
		return nil
	})

	a, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := a.Result()
	d, _ := res.Data()
	acc := new(account.Account)
	err = acc.Decode(d)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, acc)
}
