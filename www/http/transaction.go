package http

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	b := s.server.GetTransaction(s.ctx, func(p capnp.ZarbServer_getTransaction_Params) error {
		vars := mux.Vars(r)
		if err := p.SetHash([]byte(vars["hash"])); err != nil {
			return err
		}
		return nil
	})

	t, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := t.Result()
	trxData, _ := res.Data()
	trx := new(tx.Tx)
	err = trx.Decode(trxData)
	if err != nil {
		s.writeError(w, err)
		return
	}

	rec, _ := res.Receipt()
	recipetData, _ := rec.Data()
	recipet := new(tx.Receipt)
	err = recipet.Decode(recipetData)
	if err != nil {
		s.writeError(w, err)
		return
	}

	out := new(TransactionResult)
	out.Hash = trx.Hash()
	out.Tx = *trx
	out.Data = hex.EncodeToString(trxData)
	out.Receipt.Hash = recipet.Hash()
	out.Receipt.Data = hex.EncodeToString(recipetData)
	out.Receipt.Receipt = *recipet

	s.writeJSON(w, out)
}
