package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/www/capnp"
)

func (s *Server) GetBlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHash(s.ctx, func(p capnp.PactusServer_getBlockHash_Params) error {
		vars := mux.Vars(r)
		height, _ := strconv.ParseInt(vars["height"], 10, 32)
		p.SetHeight(uint32(height))
		return nil
	})
	st, _ := res.Struct()
	data, _ := st.Result()
	h, _ := hash.FromBytes(data)
	s.blockByHash(w, h)
}

func (s *Server) GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h, _ := hash.FromString(vars["hash"])
	s.blockByHash(w, h)
}

func (s *Server) blockByHash(w http.ResponseWriter, blockHash hash.Hash) {
	res := s.capnp.GetBlock(s.ctx, func(p capnp.PactusServer_getBlock_Params) error {
		p.SetVerbosity(0)
		return p.SetHash(blockHash.Bytes())
	}).Result()

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	d, _ := st.Data()
	b, err := block.FromBytes(d)
	if err != nil {
		s.writeError(w, err)
		return
	}

	seed := b.Header().SortitionSeed()

	tm := newTableMaker()
	tm.addRowString("Time", b.Header().Time().String())
	tm.addRowBytes("Hash", b.Hash().Bytes())
	tm.addRowBytes("Data", d)
	tm.addRowString("--- Header", "---")
	tm.addRowInt("Version", int(b.Header().Version()))
	tm.addRowInt("UnixTime", int(b.Header().Time().Unix()))
	tm.addRowBlockHash("PrevBlockHash", b.Header().PrevBlockHash().Bytes())
	tm.addRowBytes("StateRoot", b.Header().StateRoot().Bytes())
	tm.addRowBytes("SortitionSeed", seed[:])
	tm.addRowValAddress("ProposerAddress", b.Header().ProposerAddress().String())
	if b.PrevCertificate() != nil {
		tm.addRowString("--- PrevCertificate", "---")
		tm.addRowBytes("Hash", b.PrevCertificate().Hash().Bytes())
		tm.addRowInt("Round", int(b.PrevCertificate().Round()))
		tm.addRowInts("Committers", b.PrevCertificate().Committers())
		tm.addRowInts("Absentees", b.PrevCertificate().Absentees())
		tm.addRowBytes("Signature", b.PrevCertificate().Signature().Bytes())
	}
	tm.addRowString("--- Transactions", "---")
	for i, trx := range b.Transactions() {
		tm.addRowInt("Transaction #", i+1)
		txToTable(trx, tm)
	}

	s.writeHTML(w, tm.html())
}

func (s *Server) GetBlockHashHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHash(s.ctx, func(p capnp.PactusServer_getBlockHash_Params) error {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 32)
		if err != nil {
			return err
		}
		p.SetHeight(uint32(height))
		return nil
	})

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	data, _ := st.Result()
	hash, err := hash.FromBytes(data)
	if err != nil {
		s.writeError(w, err)
		return
	}
	tm := newTableMaker()
	tm.addRowBytes("Hash", hash.Bytes())
	s.writeHTML(w, tm.html())
}
