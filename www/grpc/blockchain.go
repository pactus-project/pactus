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
	state   state.Facade
	consMgr consensus.ManagerReader
	logger  *logger.Logger
}

func (s *blockchainServer) GetBlockchainInfo(_ context.Context,
	_ *pactus.GetBlockchainInfoRequest) (*pactus.GetBlockchainInfoResponse, error) {
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

func (s *blockchainServer) GetConsensusInfo(_ context.Context,
	_ *pactus.GetConsensusInfoRequest) (*pactus.GetConsensusInfoResponse, error) {
	instances := make([]*pactus.ConsensusInfo, 0)
	for _, cons := range s.consMgr.Instances() {
		height, round := cons.HeightRound()
		votes := cons.AllVotes()
		voteInfos := make([]*pactus.VoteInfo, 0, len(votes))
		for _, v := range votes {
			voteInfos = append(voteInfos, voteToProto(v))
		}

		instances = append(instances,
			&pactus.ConsensusInfo{
				Address: cons.SignerKey().Address().String(),
				Active:  cons.IsActive(),
				Height:  height,
				Round:   int32(round),
				Votes:   voteInfos,
			})
	}

	return &pactus.GetConsensusInfoResponse{Instances: instances}, nil
}

func (s *blockchainServer) GetBlockHash(_ context.Context,
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

func (s *blockchainServer) GetBlockHeight(_ context.Context,
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

func (s *blockchainServer) GetBlock(_ context.Context,
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

func (s *blockchainServer) GetAccount(_ context.Context,
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

func (s *blockchainServer) GetAccountByNumber(_ context.Context,
	req *pactus.GetAccountByNumberRequest) (*pactus.GetAccountResponse, error) {
	acc := s.state.AccountByNumber(req.Number)
	if acc == nil {
		return nil, status.Errorf(codes.InvalidArgument, "account not found")
	}

	return &pactus.GetAccountResponse{
		Account: accountToProto(acc),
	}, nil
}

func (s *blockchainServer) GetValidatorByNumber(_ context.Context,
	req *pactus.GetValidatorByNumberRequest) (*pactus.GetValidatorResponse, error) {
	val := s.state.ValidatorByNumber(req.Number)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	return &pactus.GetValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidator(_ context.Context,
	req *pactus.GetValidatorRequest) (*pactus.GetValidatorResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %v", err.Error())
	}
	val := s.state.ValidatorByAddress(addr)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	return &pactus.GetValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidatorAddresses(_ context.Context,
	_ *pactus.GetValidatorAddressesRequest) (*pactus.GetValidatorAddressesResponse, error) {
	addresses := s.state.ValidatorAddresses()
	addressesPB := make([]string, 0, len(addresses))
	for _, address := range addresses {
		addressesPB = append(addressesPB, address.String())
	}
	return &pactus.GetValidatorAddressesResponse{Addresses: addressesPB}, nil
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
