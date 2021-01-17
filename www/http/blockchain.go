package http

import (
	"net/http"

	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	res := s.server.GetBlockchainInfo(s.ctx, func(p capnp.ZarbServer_getBlockchainInfo_Params) error {
		return nil
	}).Result()

	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	out := new(BlockchainResult)
	out.Height = int(st.Height())
	s.writeJSON(w, out)
}
