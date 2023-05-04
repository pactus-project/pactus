package http

import (
	"net/http"

	"github.com/pactus-project/pactus/types/vote"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) ConsensusHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := s.blockchain.GetConsensusInfo(s.ctx,
		&pactus.GetConsensusInfoRequest{})
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	for i, cons := range res.Instances {
		tm.addRowInt("== Validator", i+1)
		tm.addRowValAddress("Address", cons.Address)
		tm.addRowBool("Active", cons.Active)
		tm.addRowInt("Height", int(cons.Height))
		tm.addRowInt("Round", int(cons.Round))
		tm.addRowString("Votes", "---")
		for i, v := range cons.Votes {
			tm.addRowInt("-- Vote #", i+1)
			tm.addRowString("Type", vote.Type(v.Type).String())
			tm.addRowString("Voter", v.Voter)
			tm.addRowInt("Round", int(v.Round))
			tm.addRowBlockHash("BlockHash", v.BlockHash)
		}
	}
	s.writeHTML(w, tm.html())
}
