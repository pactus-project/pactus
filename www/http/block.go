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

func (s *Server) GetBlockHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlock(s.ctx, func(p capnp.ZarbServer_getBlock_Params) error {
		vars := mux.Vars(r)
		h, _ := hash.FromString(vars["hash"])
		p.SetVerbosity(0)
		return p.SetHash(h.Bytes())
	}).Result()

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	d, _ := st.Data()
	h, _ := st.Hash()
	b, err := block.FromBytes(d)
	if err != nil {
		s.writeError(w, err)
		return
	}

	out := new(BlockResult)
	out.Block = b
	out.Hash, _ = hash.FromBytes(h)
	out.Data = hex.EncodeToString(d)
	out.Time = b.Header().Time()

	s.writeJSON(w, out)
}

func (s *Server) GetBlockHeightHandler(w http.ResponseWriter, r *http.Request) {
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
