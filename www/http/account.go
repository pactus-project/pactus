package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/www/capnp"
)

// GetAccountHandler returns a handler to get account by address
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	b := s.capnp.GetAccount(s.ctx, func(p capnp.ZarbServer_getAccount_Params) error {
		vars := mux.Vars(r)
		return p.SetAddress([]byte(vars["address"]))
	})

	a, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := a.Result()
	d, _ := res.Data()
	acc, err := account.AccountFromBytes(d)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, acc)
}
