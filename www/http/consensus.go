package http

import (
	"net/http"

	"github.com/pactus-project/pactus/types/vote"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) ConsensusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.blockchain.GetConsensusInfo(s.ctx,
		&pactus.GetConsensusInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowInt("Height", int(res.Height))
	tm.addRowInt("Round", int(res.Round))
	tm.addRowString("Votes", "---")
	for i, v := range res.Votes {
		tm.addRowInt("-- Vote #", i+1)
		tm.addRowString("Type", vote.Type(v.Type).String())
		tm.addRowString("Voter", v.Voter)
		tm.addRowInt("Round", int(v.Round))
		tm.addRowBlockHash("BlockHash", v.BlockHash)
	}
	s.writeHTML(w, tm.html())
}
