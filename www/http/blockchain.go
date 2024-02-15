package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) BlockchainHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.blockchain.GetBlockchainInfo(r.Context(),
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := newTableMaker()
	tm.addRowBlockHash("Last Block Hash", res.LastBlockHash)
	tm.addRowInt("Last Block Height", int(res.LastBlockHeight))
	tm.addRowString("--- Committee", "---")
	tm.addRowAmount("Total Power", res.TotalPower)
	tm.addRowAmount("Committee Power", res.CommitteePower)
	for i, val := range res.CommitteeValidators {
		tm.addRowInt("--- Validator", i+1)
		tmVal := s.writeValidatorTable(val)
		tm.addRowString("", tmVal.html())
	}

	s.writeHTML(w, tm.html())
}

func (s *Server) GetBlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.ParseInt(vars["height"], 10, 32)
	if err != nil {
		s.writeError(w, err)

		return
	}
	s.blockByHeight(w, r, uint32(height))
}

func (s *Server) GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockHash, err := hash.FromString(vars["hash"])
	if err != nil {
		s.writeError(w, err)

		return
	}

	res, err := s.blockchain.GetBlockHeight(r.Context(),
		&pactus.GetBlockHeightRequest{Hash: blockHash.Bytes()})
	if err != nil {
		s.writeError(w, err)

		return
	}

	s.blockByHeight(w, r, res.Height)
}

func (s *Server) blockByHeight(w http.ResponseWriter, r *http.Request, blockHeight uint32) {
	res, err := s.blockchain.GetBlock(r.Context(),
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

// GetAccountHandler returns a handler to get account by address.
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := s.blockchain.GetAccount(r.Context(),
		&pactus.GetAccountRequest{Address: vars["address"]})
	if err != nil {
		s.writeError(w, err)

		return
	}

	acc := res.Account
	tm := newTableMaker()
	tm.addRowAccAddress("Address", acc.Address)
	tm.addRowInt("Number", int(acc.Number))
	tm.addRowAmount("Balance", acc.Balance)
	tm.addRowBytes("Hash", acc.Hash)

	s.writeHTML(w, tm.html())
}

// GetValidatorHandler returns a handler to get validator by address.
func (s *Server) GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := s.blockchain.GetValidator(r.Context(),
		&pactus.GetValidatorRequest{Address: vars["address"]})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := s.writeValidatorTable(res.Validator)
	s.writeHTML(w, tm.html())
}

// GetValidatorByNumberHandler returns a handler to get validator by number.
func (s *Server) GetValidatorByNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	num, err := strconv.ParseInt(vars["number"], 10, 32)
	if err != nil {
		s.writeError(w, err)

		return
	}

	res, err := s.blockchain.GetValidatorByNumber(r.Context(),
		&pactus.GetValidatorByNumberRequest{
			Number: int32(num),
		})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := s.writeValidatorTable(res.Validator)
	s.writeHTML(w, tm.html())
}

func (s *Server) writeValidatorTable(val *pactus.ValidatorInfo) *tableMaker {
	tm := newTableMaker()
	tm.addRowString("Public Key", val.PublicKey)
	tm.addRowValAddress("Address", val.Address)
	tm.addRowInt("Number", int(val.Number))
	tm.addRowAmount("Stake", val.Stake)
	tm.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
	tm.addRowInt("LastSortitionHeight", int(val.LastSortitionHeight))
	tm.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
	tm.addRowDouble("AvailabilityScore", val.AvailabilityScore)
	tm.addRowBytes("Hash", val.Hash)

	return tm
}

func (s *Server) ConsensusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.blockchain.GetConsensusInfo(r.Context(),
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
			tm.addRowInt("CPRound", int(v.CpRound))
			tm.addRowInt("CPValue", int(v.CpValue))
			tm.addRowBlockHash("BlockHash", v.BlockHash)
		}
	}
	s.writeHTML(w, tm.html())
}
