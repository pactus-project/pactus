package grpc

import (
	"context"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blockchainServer struct {
	state     state.Facade
	consensus consensus.Reader
	logger    *logger.Logger
}

func (s *blockchainServer) GetBlockchainInfo(ctx context.Context,
	req *pactus.GetBlockchainInfoRequest) (*pactus.GetBlockchainInfoResponse, error) {
	vals := s.state.CommitteeValidators()
	cv := make([]*pactus.ValidatorInfo, 0, len(vals))
	for _, v := range vals {
		cv = append(cv, validatorToProto(v))
	}

	return &pactus.GetBlockchainInfoResponse{
		LastBlockHeight:     s.state.LastBlockHeight(),
		LastBlockHash:       s.state.LastBlockHash().Bytes(),
		TotalAccounts:       s.state.TotalAccounts(),
		TotalValidators:     s.state.TotalValidators(),
		TotalPower:          s.state.TotalPower(),
		CommitteePower:      s.state.CommitteePower(),
		CommitteeValidators: cv,
	}, nil
}

func (s *blockchainServer) GetConsensusInfo(ctx context.Context,
	req *pactus.GetConsensusInfoRequest) (*pactus.GetConsensusInfoResponse, error) {
	height, round := s.consensus.HeightRound()
	votes := s.consensus.AllVotes()
	vinfo := make([]*pactus.VoteInfo, 0, len(votes))
	for _, v := range votes {
		vinfo = append(vinfo, voteToProto(v))
	}

	return &pactus.GetConsensusInfoResponse{
		Height: height,
		Round:  int32(round),
		Votes:  vinfo,
	}, nil
}

func (s *blockchainServer) GetBlockHash(ctx context.Context,
	req *pactus.GetBlockHashRequest) (*pactus.GetBlockHashResponse, error) {
	height := req.GetHeight()
	hash := s.state.BlockHash(height)
	if hash.IsUndef() {
		return nil, status.Errorf(codes.NotFound, "block not found with this height")
	}
	return &pactus.GetBlockHashResponse{
		Hash: hash.Bytes(),
	}, nil
}

func (s *blockchainServer) GetBlockHeight(ctx context.Context,
	req *pactus.GetBlockHeightRequest) (*pactus.GetBlockHeightResponse, error) {
	hash, err := hash.FromBytes(req.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hash: %v", err)
	}
	height := s.state.BlockHeight(hash)
	if height == 0 {
		return nil, status.Errorf(codes.NotFound, "block not found with this hash")
	}
	return &pactus.GetBlockHeightResponse{
		Height: height,
	}, nil
}

func (s *blockchainServer) GetBlock(ctx context.Context,
	req *pactus.GetBlockRequest) (*pactus.GetBlockResponse, error) {
	height := req.GetHeight()
	storedBlock := s.state.StoredBlock(height)
	if storedBlock == nil {
		return nil, status.Errorf(codes.InvalidArgument, "block not found")
	}
	res := &pactus.GetBlockResponse{
		Height: storedBlock.Height,
		Hash:   storedBlock.BlockHash.Bytes(),
		Data:   storedBlock.Data,
	}

	if req.Verbosity > pactus.BlockVerbosity_BLOCK_DATA {
		block := storedBlock.ToBlock()
		blockTime := block.Header().UnixTime()
		seed := block.Header().SortitionSeed()
		cert := block.PrevCertificate()
		var prevCert *pactus.CertificateInfo

		if cert != nil {
			committers := make([]int32, len(cert.Committers()))
			for i, n := range cert.Committers() {
				committers[i] = n
			}
			absentees := make([]int32, len(cert.Absentees()))
			for i, n := range cert.Absentees() {
				absentees[i] = n
			}
			prevCert = &pactus.CertificateInfo{
				Hash:       cert.Hash().Bytes(),
				Round:      int32(cert.Round()),
				Committers: committers,
				Absentees:  absentees,
				Signature:  cert.Signature().Bytes(),
			}
		}
		header := &pactus.BlockHeaderInfo{
			Version:         int32(block.Header().Version()),
			PrevBlockHash:   block.Header().PrevBlockHash().Bytes(),
			StateRoot:       block.Header().StateRoot().Bytes(),
			SortitionSeed:   seed[:],
			ProposerAddress: block.Header().ProposerAddress().String(),
		}

		trxs := make([]*pactus.TransactionInfo, 0, block.Transactions().Len())
		for _, trx := range block.Transactions() {
			if req.Verbosity == pactus.BlockVerbosity_BLOCK_INFO {
				trxs = append(trxs, &pactus.TransactionInfo{Id: trx.ID().Bytes()})
			} else {
				trxs = append(trxs, transactionToProto(trx))
			}
		}

		res.BlockTime = blockTime
		res.Header = header
		res.Txs = trxs
		res.PrevCert = prevCert
	}

	return res, nil
}

func (s *blockchainServer) GetAccount(ctx context.Context,
	req *pactus.GetAccountRequest) (*pactus.GetAccountResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %v", err)
	}
	acc := s.state.AccountByAddress(addr)
	if acc == nil {
		return nil, status.Errorf(codes.InvalidArgument, "account not found")
	}
	res := &pactus.GetAccountResponse{
		Account: accountToProto(acc),
	}

	return res, nil
}

func (s *blockchainServer) GetValidatorByNumber(ctx context.Context,
	req *pactus.GetValidatorByNumberRequest) (*pactus.GetValidatorResponse, error) {
	val := s.state.ValidatorByNumber(req.Number)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &pactus.GetValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidator(ctx context.Context,
	req *pactus.GetValidatorRequest) (*pactus.GetValidatorResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %v", err.Error())
	}
	val := s.state.ValidatorByAddress(addr)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &pactus.GetValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidators(ctx context.Context,
	req *pactus.GetValidatorsRequest) (*pactus.GetValidatorsResponse, error) {
	validators := s.state.CommitteeValidators()
	validatorsResp := make([]*pactus.ValidatorInfo, 0)
	for _, val := range validators {
		validatorsResp = append(validatorsResp, validatorToProto(val))
	}
	return &pactus.GetValidatorsResponse{Validators: validatorsResp}, nil
}

func validatorToProto(val *validator.Validator) *pactus.ValidatorInfo {
	data, _ := val.Bytes()
	return &pactus.ValidatorInfo{
		Hash:              val.Hash().Bytes(),
		Data:              data,
		PublicKey:         val.PublicKey().String(),
		Address:           val.Address().String(),
		Number:            val.Number(),
		Sequence:          val.Sequence(),
		Stake:             val.Stake(),
		LastBondingHeight: val.LastBondingHeight(),
		LastJoinedHeight:  val.LastJoinedHeight(),
		UnbondingHeight:   val.UnbondingHeight(),
	}
}

func accountToProto(acc *account.Account) *pactus.AccountInfo {
	data, _ := acc.Bytes()
	return &pactus.AccountInfo{
		Hash:     acc.Hash().Bytes(),
		Data:     data,
		Address:  acc.Address().String(),
		Number:   acc.Number(),
		Sequence: acc.Sequence(),
		Balance:  acc.Balance(),
	}
}

func voteToProto(v *vote.Vote) *pactus.VoteInfo {
	return &pactus.VoteInfo{
		Type:      pactus.VoteType(v.Type()),
		Voter:     v.Signer().String(),
		BlockHash: v.BlockHash().Bytes(),
		Round:     int32(v.Round()),
	}
}
