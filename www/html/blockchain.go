package html

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/vote"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) BlockchainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.blockchain.GetBlockchainInfo(ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tmk := newTableMaker()
	tmk.addRowBlockHash("Last Block Hash", res.LastBlockHash)
	tmk.addRowInt("Last Block Height", int(res.LastBlockHeight))
	tmk.addRowString("Last Block Time", time.Unix(res.LastBlockTime, 0).String())
	tmk.addRowInt("Total Accounts", int(res.TotalAccounts))
	tmk.addRowInt("Total Validators", int(res.TotalValidators))
	tmk.addRowInt("Active Validators", int(res.ActiveValidators))
	tmk.addRowPower("Total Power", res.TotalPower)
	tmk.addRowBool("Is Pruned", res.IsPruned)
	tmk.addRowInt("Pruning Height", int(res.PruningHeight))

	s.writeHTML(w, tmk.html())
}

func (s *Server) CommitteeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.blockchain.GetCommitteeInfo(ctx,
		&pactus.GetCommitteeInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tmk := newTableMaker()
	tmk.addRowPower("Committee Power", res.CommitteePower)
	tmk.addRowInt("Validators", len(res.Validators))
	tmk.addRowString("--- Protocol Versions", "---")
	for ver, percentage := range res.ProtocolVersions {
		tmk.addRowDouble(fmt.Sprintf("Version %d", ver), percentage)
	}
	tmk.addRowString("--- Validators", "---")
	for i, val := range res.Validators {
		tmk.addRowInt("--- Validator", i+1)
		tmVal := s.writeValidatorTable(val)
		tmk.addRowString("", tmVal.html())
	}

	s.writeHTML(w, tmk.html())
}

func (s *Server) GetBlockByHeightHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	height, err := strconv.ParseInt(vars["height"], 10, 32)
	if err != nil {
		s.writeError(w, err)

		return
	}
	s.blockByHeight(ctx, w, uint32(height))
}

func (s *Server) GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	blockHash, err := hash.FromString(vars["hash"])
	if err != nil {
		s.writeError(w, err)

		return
	}

	res, err := s.blockchain.GetBlockHeight(ctx,
		&pactus.GetBlockHeightRequest{Hash: blockHash.String()})
	if err != nil {
		s.writeError(w, err)

		return
	}

	s.blockByHeight(ctx, w, res.Height)
}

func (s *Server) blockByHeight(ctx context.Context, w http.ResponseWriter, blockHeight uint32) {
	res, err := s.blockchain.GetBlock(ctx,
		&pactus.GetBlockRequest{
			Height:    blockHeight,
			Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS,
		},
	)
	if err != nil {
		s.writeError(w, err)

		return
	}

	tmk := newTableMaker()
	tmk.addRowString("Time", time.Unix(int64(res.BlockTime), 0).String())
	tmk.addRowInt("Height", int(res.Height))
	tmk.addRowString("Hash", res.Hash)
	tmk.addRowString("Data", res.Data)
	if res.Header != nil {
		tmk.addRowString("--- Header", "---")
		tmk.addRowInt("Version", int(res.Header.Version))
		tmk.addRowInt("UnixTime", int(res.BlockTime))
		tmk.addRowBlockHash("PrevBlockHash", res.Header.PrevBlockHash)
		tmk.addRowString("StateRoot", res.Header.StateRoot)
		tmk.addRowString("SortitionSeed", res.Header.SortitionSeed)
		tmk.addRowValAddress("ProposerAddress", res.Header.ProposerAddress)
	}
	if res.PrevCert != nil {
		tmk.addRowString("--- PrevCertificate", "---")
		tmk.addRowString("Hash", res.PrevCert.Hash)
		tmk.addRowInt("Round", int(res.PrevCert.Round))
		tmk.addRowInts("Committers", res.PrevCert.Committers)
		tmk.addRowInts("Absentees", res.PrevCert.Absentees)
		tmk.addRowString("Signature", res.PrevCert.Signature)
	}
	tmk.addRowString("--- Transactions", "---")
	for i, trx := range res.Txs {
		tmk.addRowInt("Transaction #", i+1)
		txToTable(tmk, trx)
	}

	s.writeHTML(w, tmk.html())
}

// GetAccountHandler returns a handler to get account by address.
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	res, err := s.blockchain.GetAccount(ctx,
		&pactus.GetAccountRequest{Address: vars["address"]})
	if err != nil {
		s.writeError(w, err)

		return
	}

	acc := res.Account
	tmk := newTableMaker()
	tmk.addRowAccAddress("Address", acc.Address)
	tmk.addRowInt("Number", int(acc.Number))
	tmk.addRowAmount("Balance", amount.Amount(acc.Balance))
	tmk.addRowString("Hash", acc.Hash)

	s.writeHTML(w, tmk.html())
}

// GetValidatorHandler returns a handler to get validator by address.
func (s *Server) GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	res, err := s.blockchain.GetValidator(ctx,
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
	ctx := r.Context()

	vars := mux.Vars(r)

	num, err := strconv.ParseInt(vars["number"], 10, 32)
	if err != nil {
		s.writeError(w, err)

		return
	}

	res, err := s.blockchain.GetValidatorByNumber(ctx,
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

func (s *Server) GetTxPoolContentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.blockchain.GetTxPoolContent(ctx, &pactus.GetTxPoolContentRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tm := newTableMaker()
	for i, trx := range res.Txs {
		tm.addRowString("\n-------------- ", fmt.Sprintf("%d --------------\n", i))
		txToTable(tm, trx)
	}
	s.writeHTML(w, tm.html())
}

func (*Server) writeValidatorTable(val *pactus.ValidatorInfo) *tableMaker {
	tmk := newTableMaker()
	tmk.addRowString("Public Key", val.PublicKey)
	tmk.addRowValAddress("Address", val.Address)
	tmk.addRowInt("Number", int(val.Number))
	tmk.addRowAmount("Stake", amount.Amount(val.Stake))
	tmk.addRowInt("LastBondingHeight", int(val.LastBondingHeight))
	tmk.addRowInt("LastSortitionHeight", int(val.LastSortitionHeight))
	tmk.addRowInt("UnbondingHeight", int(val.UnbondingHeight))
	tmk.addRowDouble("AvailabilityScore", val.AvailabilityScore)
	tmk.addRowInt("ProtocolVersion", int(val.ProtocolVersion))
	tmk.addRowString("Hash", val.Hash)

	return tmk
}

func (s *Server) ConsensusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := s.blockchain.GetConsensusInfo(ctx,
		&pactus.GetConsensusInfoRequest{})
	if err != nil {
		s.writeError(w, err)

		return
	}

	tmk := newTableMaker()

	tmk.addRowString("== Proposal", "")
	if res.Proposal != nil {
		tmk.addRowInt("Height", int(res.Proposal.Height))
		tmk.addRowInt("Round", int(res.Proposal.Round))
		tmk.addRowString("BlockData", res.Proposal.BlockData)
		tmk.addRowString("Signature", res.Proposal.Signature)
	}

	for i, cons := range res.Instances {
		tmk.addRowInt("== Validator", i+1)
		tmk.addRowValAddress("Address", cons.Address)
		tmk.addRowBool("Active", cons.Active)
		tmk.addRowInt("Height", int(cons.Height))
		tmk.addRowInt("Round", int(cons.Round))
		tmk.addRowString("Votes", "---")
		for index, vte := range cons.Votes {
			tmk.addRowInt("-- Vote #", index+1)
			tmk.addRowBlockHash("BlockHash", vte.BlockHash)
			tmk.addRowString("Type", vote.Type(vte.Type).String())
			tmk.addRowString("Voter", vte.Voter)
			tmk.addRowInt("Round", int(vte.Round))
			tmk.addRowInt("CPRound", int(vte.CpRound))
			tmk.addRowInt("CPValue", int(vte.CpValue))
		}
	}

	s.writeHTML(w, tmk.html())
}
