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
	vars := mux.Vars(r)
	height, err := strconv.ParseInt(vars["height"], 10, 32)
	if err != nil {
		s.writeError(w, err)
		return
	}
	s.blockByHeight(w, uint32(height))
}

func (s *Server) GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHeight(s.ctx, func(p capnp.PactusServer_getBlockHeight_Params) error {
		vars := mux.Vars(r)
		blockHash, err := hash.FromString(vars["hash"])
		if err != nil {
			return err
		}
		return p.SetHash(blockHash.Bytes())
	})
	st, _ := res.Struct()
	blockHeight := st.Result()
	s.blockByHeight(w, blockHeight)
}

func (s *Server) blockByHeight(w http.ResponseWriter, blockHeight uint32) {
	res := s.capnp.GetBlock(s.ctx, func(p capnp.PactusServer_getBlock_Params) error {
		p.SetVerbosity(0)
		p.SetHeight(blockHeight)
		return nil
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
	tm.addRowInt("Height", int(blockHeight))
	tm.addRowBytes("Hash", b.Hash().Bytes())
	tm.addRowBytes("Data", d)
	tm.addRowString("--- Header", "---")
	tm.addRowInt("Version", int(b.Header().Version()))
	tm.addRowInt("UnixTime", int(b.Header().UnixTime()))
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
