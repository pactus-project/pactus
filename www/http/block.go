package http

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetBlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHash(s.ctx, func(p capnp.ZarbServer_getBlockHash_Params) error {
		vars := mux.Vars(r)
		height, _ := strconv.ParseInt(vars["height"], 10, 32)
		p.SetHeight(int32(height))
		return nil
	})
	st, _ := res.Struct()
	data, _ := st.Result()
	h, _ := hash.FromBytes(data)
	s.blockByHash(w, r, h)
}

func (s *Server) GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h, _ := hash.FromString(vars["hash"])
	s.blockByHash(w, r, h)
}

func (s *Server) blockByHash(w http.ResponseWriter, r *http.Request, blockHash hash.Hash) {
	res := s.capnp.GetBlock(s.ctx, func(p capnp.ZarbServer_getBlock_Params) error {
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
	out := new(BlockResult)
	out.Header = BlockHeaderResult{
		Version:         b.Header().Version(),
		UnixTime:        uint32(b.Header().Time().Unix()),
		PrevBlockHash:   b.Header().PrevBlockHash().String(),
		StateRoot:       b.Header().StateRoot().String(),
		SortitionSeed:   hex.EncodeToString(seed[:]),
		ProposerAddress: b.Header().ProposerAddress().String(),
	}
	for _, trx := range b.Transactions() {
		out.Txs = append(out.Txs, txToResult(trx))
	}
	out.Hash = b.Hash().String()
	out.Data = hex.EncodeToString(d)
	out.Time = b.Header().Time()

	s.writeJSON(w, out)
}

func (s *Server) GetBlockHashHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHash(s.ctx, func(p capnp.ZarbServer_getBlockHash_Params) error {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 32)
		if err != nil {
			return err
		}
		p.SetHeight(int32(height))
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
	s.writePlainText(w, hash.String())
}
