package http

import (
	"net/http"

	"github.com/zarbchain/zarb-go/types/vote"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func (s *Server) ConsensusHandler(w http.ResponseWriter, r *http.Request) {
	res := s.capnp.GetConsensusInfo(s.ctx, func(p capnp.ZarbServer_getConsensusInfo_Params) error {
		return nil
	}).Result()
	st, err := res.Struct()
	if err != nil {
		s.writeError(w, err)
		return
	}
	tm := newTableMaker()

	tm.addRowInt("Height", int(st.Height()))
	tm.addRowInt("Round", int(st.Round()))
	tm.addRowString("Votes", "---")
	voteList, _ := st.Votes()
	for i := 0; i < voteList.Len(); i++ {
		v := voteList.At(i)
		tm.addRowInt("-- Vote #", i+1)
		tm.addRowString("Type", vote.Type(v.Type()).String())
		voter, _ := v.Voter()
		tm.addRowString("Voter", voter)
		tm.addRowInt("Round", int(v.Round()))
		hash, _ := v.BlockHash()
		tm.addRowBlockHash("BlockHash", hash)
	}

	s.writeHTML(w, tm.html())
}
