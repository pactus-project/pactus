package grpc

import (
	"context"
	"encoding/hex"

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
		LastBlockHash:       s.state.LastBlockHash().String(),
		TotalAccounts:       s.state.TotalAccounts(),
		TotalValidators:     s.state.TotalValidators(),
		TotalPower:          s.state.TotalPower(),
		CommitteePower:      s.state.CommitteePower(),
		IsPruned:            s.state.IsPruned(),
		PruningHeight:       s.state.PruningHeight(),
		LastBlockTime:       s.state.LastBlockTime().Unix(),
		CommitteeValidators: cv,
	}, nil
}

func (s *blockchainServer) GetConsensusInfo(_ context.Context,
	_ *pactus.GetConsensusInfoRequest,
) (*pactus.GetConsensusInfoResponse, error) {
	instances := make([]*pactus.ConsensusInfo, 0)
	var proposal *pactus.Proposal

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

		if proposal == nil && cons.Proposal() != nil {
			// All instances have the same proposal, so we just need to fill proposal once.
			p := cons.Proposal()

			var blockData string
			if p.Block() != nil {
				data, err := p.Block().Bytes()
				if err != nil {
					return nil, status.Errorf(codes.Internal, err.Error())
				}

				blockData = string(data)
			}

			proposal = &pactus.Proposal{
				Height:        p.Height(),
				Round:         int32(p.Round()),
				BlockData:     blockData,
				SignatureData: string(p.Signature().Bytes()),
			}
		}
	}

	return &pactus.GetConsensusInfoResponse{Instances: instances, Proposal: proposal}, nil
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
		Hash: h.String(),
	}, nil
}

func (s *blockchainServer) GetBlockHeight(_ context.Context,
	req *pactus.GetBlockHeightRequest,
) (*pactus.GetBlockHeightResponse, error) {
	h, err := hash.FromString(req.GetHash())
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
	cBlk := s.state.CommittedBlock(height)
	if cBlk == nil {
		return nil, status.Errorf(codes.NotFound, "block not found")
	}
	res := &pactus.GetBlockResponse{
		Height: cBlk.Height,
		Hash:   cBlk.BlockHash.String(),
	}

	switch req.Verbosity {
	case pactus.BlockVerbosity_BLOCK_DATA:
		res.Data = hex.EncodeToString(cBlk.Data)

	case pactus.BlockVerbosity_BLOCK_INFO,
		pactus.BlockVerbosity_BLOCK_TRANSACTIONS:
		block, err := cBlk.ToBlock()
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
				Hash:       cert.Hash().String(),
				Round:      int32(cert.Round()),
				Committers: committers,
				Absentees:  absentees,
				Signature:  cert.Signature().String(),
			}
		}
		header := &pactus.BlockHeaderInfo{
			Version:         int32(block.Header().Version()),
			PrevBlockHash:   block.Header().PrevBlockHash().String(),
			StateRoot:       block.Header().StateRoot().String(),
			SortitionSeed:   hex.EncodeToString(seed[:]),
			ProposerAddress: block.Header().ProposerAddress().String(),
		}

		trxs := make([]*pactus.TransactionInfo, 0, block.Transactions().Len())
		for _, trx := range block.Transactions() {
			if req.Verbosity == pactus.BlockVerbosity_BLOCK_INFO {
				data, _ := trx.Bytes()
				trxs = append(trxs, &pactus.TransactionInfo{
					Id:   trx.ID().String(),
					Data: hex.EncodeToString(data),
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

func (s *blockchainServer) GetTxPoolContent(_ context.Context,
	_ *pactus.GetTxPoolContentRequest,
) (*pactus.GetTxPoolContentResponse, error) {
	result := make([]*pactus.TransactionInfo, 0)
	for _, t := range s.state.AllPendingTxs() {
		result = append(result, transactionToProto(t))
	}

	return &pactus.GetTxPoolContentResponse{
		Txs: result,
	}, nil
}

func (s *blockchainServer) validatorToProto(val *validator.Validator) *pactus.ValidatorInfo {
	data, _ := val.Bytes()

	return &pactus.ValidatorInfo{
		Hash:                val.Hash().String(),
		Data:                hex.EncodeToString(data),
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

func (*blockchainServer) accountToProto(addr crypto.Address, acc *account.Account) *pactus.AccountInfo {
	data, _ := acc.Bytes()

	return &pactus.AccountInfo{
		Hash:    acc.Hash().String(),
		Data:    hex.EncodeToString(data),
		Number:  acc.Number(),
		Balance: acc.Balance().ToNanoPAC(),
		Address: addr.String(),
	}
}

func (*blockchainServer) voteToProto(v *vote.Vote) *pactus.VoteInfo {
	cpRound := int32(0)
	cpValue := int32(0)
	if v.IsCPVote() {
		cpRound = int32(v.CPRound())
		cpValue = int32(v.CPValue())
	}

	return &pactus.VoteInfo{
		Type:      pactus.VoteType(v.Type()),
		Voter:     v.Signer().String(),
		BlockHash: v.BlockHash().String(),
		Round:     int32(v.Round()),
		CpRound:   cpRound,
		CpValue:   cpValue,
	}
}
