package http

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	b := s.capnp.GetTransaction(s.ctx, func(p capnp.ZarbServer_getTransaction_Params) error {
		vars := mux.Vars(r)
		return p.SetId([]byte(vars["id"]))
	})

	t, err := b.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := t.Result()
	trxData, _ := res.Data()
	fmt.Printf("%x\n", trxData)
	trx := new(tx.Tx)
	err = trx.Decode(trxData)
	if err != nil {
		s.writeError(w, err)
		return
	}

	out := new(TransactionResult)
	out.ID = trx.ID()
	out.Tx = *trx
	out.Data = hex.EncodeToString(trxData)

	s.writeJSON(w, out)
}

func (s *Server) SendRawTransactionHandler(w http.ResponseWriter, r *http.Request) {
	txRes := s.capnp.SendRawTransaction(s.ctx, func(p capnp.ZarbServer_sendRawTransaction_Params) error {
		vars := mux.Vars(r)
		d, err := hex.DecodeString(vars["data"])
		if err != nil {
			return err
		}
		return p.SetRawTx(d)
	})

	t, err := txRes.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, _ := t.Result()

	out := new(SendTranscationResult)
	txID, err := res.Id()
	if err != nil {
		s.writeError(w, err)
		return
	}
	out.ID, err = hash.FromRawBytes(txID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	out.Status = int(res.Status())

	s.writeJSON(w, out)
}
