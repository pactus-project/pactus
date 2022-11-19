package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto/hash"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
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
	vars := mux.Vars(r)
	blockHash, err := hash.FromString(vars["hash"])
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, err := s.blockchain.GetBlockHeight(s.ctx,
		&pactus.GetBlockHeightRequest{Hash: blockHash.Bytes()})
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.blockByHeight(w, res.Height)
}

func (s *Server) blockByHeight(w http.ResponseWriter, blockHeight uint32) {
	res, err := s.blockchain.GetBlock(s.ctx,
		&pactus.GetBlockRequest{
			Height:    blockHeight,
			Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS,
		},
	)
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	tm.addRowString("Time", time.Unix(int64(res.BlockTime), 0).String())
	tm.addRowInt("Height", int(res.Height))
	tm.addRowBytes("Hash", res.Hash)
	tm.addRowBytes("Data", res.Data)
	if res.Header != nil {
		tm.addRowString("--- Header", "---")
		tm.addRowInt("Version", int(res.Header.Version))
		tm.addRowInt("UnixTime", int(res.BlockTime))
		tm.addRowBlockHash("PrevBlockHash", res.Header.PrevBlockHash)
		tm.addRowBytes("StateRoot", res.Header.StateRoot)
		tm.addRowBytes("SortitionSeed", res.Header.SortitionSeed)
		tm.addRowValAddress("ProposerAddress", res.Header.ProposerAddress)
	}
	if res.PrevCert != nil {
		tm.addRowString("--- PrevCertificate", "---")
		tm.addRowBytes("Hash", res.PrevCert.Hash)
		tm.addRowInt("Round", int(res.PrevCert.Round))
		tm.addRowInts("Committers", res.PrevCert.Committers)
		tm.addRowInts("Absentees", res.PrevCert.Absentees)
		tm.addRowBytes("Signature", res.PrevCert.Signature)
	}
	tm.addRowString("--- Transactions", "---")
	for i, trx := range res.Txs {
		tm.addRowInt("Transaction #", i+1)
		txToTable(trx, tm)
	}

	s.writeHTML(w, tm.html())
}
