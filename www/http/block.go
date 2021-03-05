package http

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetBlockHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlock(s.ctx, func(p capnp.ZarbServer_getBlock_Params) error {
		vars := mux.Vars(r)
		height, err := strconv.Atoi(vars["height"])
		if err != nil {
			return err
		}
		p.SetHeight(uint64(height))
		p.SetVerbosity(0)
		return nil
	}).Result()

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	d, _ := st.Data()
	h, _ := st.Hash()
	b := new(block.Block)
	if err = b.Decode(d); err != nil {
		s.writeError(w, err)
		return
	}

	out := new(BlockResult)
	out.Block = *b
	out.Hash, _ = crypto.HashFromRawBytes(h)
	out.Data = hex.EncodeToString(d)
	out.Time = b.Header().Time()

	s.writeJSON(w, out)
}

func (s *Server) GetBlockHeightHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetBlockHeight(s.ctx, func(p capnp.ZarbServer_getBlockHeight_Params) error {
		vars := mux.Vars(r)
		if err := p.SetHash([]byte(vars["hash"])); err != nil {
			return err
		}
		return nil
	})

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}

	height := st.Result()

	s.writePlainText(w, strconv.FormatUint(height, 10))
}
