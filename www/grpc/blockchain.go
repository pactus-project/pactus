package grpc

import (
	"context"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blockchainServer struct {
	*Server
}

func newBlockchainServer(server *Server) *blockchainServer {
	return &blockchainServer{
		Server: server,
	}
}

func (s *blockchainServer) GetBlockchainInfo(_ context.Context,
	_ *pactus.GetBlockchainInfoRequest,
) (*pactus.GetBlockchainInfoResponse, error) {
	vals := s.state.CommitteeValidators()
	cv := make([]*pactus.ValidatorInfo, 0, len(vals))
	for _, v := range vals {
		cv = append(cv, s.validatorToProto(v))
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
	_ *pactus.GetConsensusInfoRequest,
) (*pactus.GetConsensusInfoResponse, error) {
	instances := make([]*pactus.ConsensusInfo, 0)
	for _, cons := range s.consMgr.Instances() {
		height, round := cons.HeightRound()
		votes := cons.AllVotes()
		voteInfos := make([]*pactus.VoteInfo, 0, len(votes))
		for _, v := range votes {
			voteInfos = append(voteInfos, s.voteToProto(v))
		}

		instances = append(instances,
			&pactus.ConsensusInfo{
				Address: cons.ConsensusKey().ValidatorAddress().String(),
				Active:  cons.IsActive(),
				Height:  height,
				Round:   int32(round),
				Votes:   voteInfos,
			})
	}

	return &pactus.GetConsensusInfoResponse{Instances: instances}, nil
}

func (s *blockchainServer) GetBlockHash(_ context.Context,
	req *pactus.GetBlockHashRequest,
) (*pactus.GetBlockHashResponse, error) {
	height := req.GetHeight()
	h := s.state.BlockHash(height)
	if h.IsUndef() {
		return nil, status.Errorf(codes.NotFound, "block not found with this height")
	}

	return &pactus.GetBlockHashResponse{
		Hash: h.Bytes(),
	}, nil
}

func (s *blockchainServer) GetBlockHeight(_ context.Context,
	req *pactus.GetBlockHeightRequest,
) (*pactus.GetBlockHeightResponse, error) {
	h, err := hash.FromBytes(req.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hash: %v", err)
	}
	height := s.state.BlockHeight(h)
	if height == 0 {
		return nil, status.Errorf(codes.NotFound, "block not found with this hash")
	}

	return &pactus.GetBlockHeightResponse{
		Height: height,
	}, nil
}

func (s *blockchainServer) GetBlock(_ context.Context,
	req *pactus.GetBlockRequest,
) (*pactus.GetBlockResponse, error) {
	height := req.GetHeight()
	committedBlock := s.state.CommittedBlock(height)
	if committedBlock == nil {
		return nil, status.Errorf(codes.NotFound, "block not found")
	}
	res := &pactus.GetBlockResponse{
		Height: committedBlock.Height,
		Hash:   committedBlock.BlockHash.Bytes(),
	}

	switch req.Verbosity {
	case pactus.BlockVerbosity_BLOCK_DATA:
		res.Data = committedBlock.Data

	case pactus.BlockVerbosity_BLOCK_INFO,
		pactus.BlockVerbosity_BLOCK_TRANSACTIONS:
		block, err := committedBlock.ToBlock()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
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
				data, _ := trx.Bytes()
				trxs = append(trxs, &pactus.TransactionInfo{
					Id:   trx.ID().Bytes(),
					Data: data,
				})
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
	req *pactus.GetAccountRequest,
) (*pactus.GetAccountResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %v", err)
	}
	acc := s.state.AccountByAddress(addr)
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}
	res := &pactus.GetAccountResponse{
		Account: s.accountToProto(addr, acc),
	}

	return res, nil
}

func (s *blockchainServer) GetValidatorByNumber(_ context.Context,
	req *pactus.GetValidatorByNumberRequest,
) (*pactus.GetValidatorResponse, error) {
	val := s.state.ValidatorByNumber(req.Number)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	return &pactus.GetValidatorResponse{
		Validator: s.validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidator(_ context.Context,
	req *pactus.GetValidatorRequest,
) (*pactus.GetValidatorResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %v", err.Error())
	}
	val := s.state.ValidatorByAddress(addr)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	return &pactus.GetValidatorResponse{
		Validator: s.validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidatorAddresses(_ context.Context,
	_ *pactus.GetValidatorAddressesRequest,
) (*pactus.GetValidatorAddressesResponse, error) {
	addresses := s.state.ValidatorAddresses()
	addressesPB := make([]string, 0, len(addresses))
	for _, address := range addresses {
		addressesPB = append(addressesPB, address.String())
	}

	return &pactus.GetValidatorAddressesResponse{Addresses: addressesPB}, nil
}

func (s *blockchainServer) GetPublicKey(_ context.Context,
	req *pactus.GetPublicKeyRequest,
) (*pactus.GetPublicKeyResponse, error) {
	addr, err := crypto.AddressFromString(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account address: %v", err.Error())
	}

	publicKey, err := s.state.PublicKey(addr)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "public key not found")
	}

	return &pactus.GetPublicKeyResponse{PublicKey: publicKey.String()}, nil
}

func (s *blockchainServer) validatorToProto(val *validator.Validator) *pactus.ValidatorInfo {
	data, _ := val.Bytes()

	return &pactus.ValidatorInfo{
		Hash:                val.Hash().Bytes(),
		Data:                data,
		PublicKey:           val.PublicKey().String(),
		Address:             val.Address().String(),
		Number:              val.Number(),
		Stake:               val.Stake().ToNanoPAC(),
		LastBondingHeight:   val.LastBondingHeight(),
		LastSortitionHeight: val.LastSortitionHeight(),
		UnbondingHeight:     val.UnbondingHeight(),
		AvailabilityScore:   s.state.AvailabilityScore(val.Number()),
	}
}

func (s *blockchainServer) accountToProto(addr crypto.Address, acc *account.Account) *pactus.AccountInfo {
	data, _ := acc.Bytes()

	return &pactus.AccountInfo{
		Hash:    acc.Hash().Bytes(),
		Data:    data,
		Number:  acc.Number(),
		Balance: acc.Balance().ToNanoPAC(),
		Address: addr.String(),
	}
}

func (s *blockchainServer) voteToProto(v *vote.Vote) *pactus.VoteInfo {
	cpRound := int32(0)
	cpValue := int32(0)
	if v.IsCPVote() {
		cpRound = int32(v.CPRound())
		cpValue = int32(v.CPValue())
	}

	return &pactus.VoteInfo{
		Type:      pactus.VoteType(v.Type()),
		Voter:     v.Signer().String(),
		BlockHash: v.BlockHash().Bytes(),
		Round:     int32(v.Round()),
		CpRound:   cpRound,
		CpValue:   cpValue,
	}
}
