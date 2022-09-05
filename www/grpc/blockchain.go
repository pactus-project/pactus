package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type blockchainServer struct {
	state  state.Facade
	logger *logger.Logger
}

func (s *blockchainServer) GetBlockchainInfo(ctx context.Context,
	request *zarb.BlockchainInfoRequest) (*zarb.BlockchainInfoResponse, error) {
	height := s.state.LastBlockHeight()

	return &zarb.BlockchainInfoResponse{
		LastBlockHeight: height,
		LastBlockHash:   s.state.LastBlockHash().Bytes(),
	}, nil
}

func (s *blockchainServer) GetBlockHash(ctx context.Context,
	request *zarb.BlockHashRequest) (*zarb.BlockHashResponse, error) {
	height := request.GetHeight()
	hash := s.state.BlockHash(height)
	if hash.IsUndef() {
		return nil, status.Errorf(codes.NotFound, "block hash not found with this height")
	}
	return &zarb.BlockHashResponse{
		Hash: hash.Bytes(),
	}, nil
}

func (s *blockchainServer) GetBlock(ctx context.Context,
	request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	hash, err := hash.FromBytes(request.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "provided hash is not Valid")
	}
	block := s.state.Block(hash)
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "block not found")
	}
	timestamp := timestamppb.New(block.Header().Time())
	header := &zarb.BlockHeaderInfo{}
	var prevCert *zarb.CertificateInfo

	if request.Verbosity.Number() > 0 {
		seed := block.Header().SortitionSeed()

		cert := block.PrevCertificate()
		if cert != nil {
			committers := make([]int32, len(block.PrevCertificate().Committers()))
			for i, n := range block.PrevCertificate().Committers() {
				committers[i] = n
			}
			absentees := make([]int32, len(block.PrevCertificate().Absentees()))
			for i, n := range block.PrevCertificate().Absentees() {
				absentees[i] = n
			}
			prevCert = &zarb.CertificateInfo{
				Round:      int32(block.PrevCertificate().Round()),
				Committers: committers,
				Absentees:  absentees,
				Signature:  block.PrevCertificate().Signature().Bytes(),
			}
		}
		header = &zarb.BlockHeaderInfo{
			Version:         int32(block.Header().Version()),
			PrevBlockHash:   block.Header().PrevBlockHash().Bytes(),
			StateRoot:       block.Header().StateRoot().Bytes(),
			SortitionSeed:   seed[:],
			ProposerAddress: block.Header().ProposerAddress().String(),
		}
	}

	// TODO: Cache for better performance
	trxs := make([]*zarb.TransactionInfo, 0, block.Transactions().Len())
	if request.Verbosity.Number() > 1 {
		for _, trx := range block.Transactions() {
			trxs = append(trxs, transactionToProto(trx))
		}
	}

	res := &zarb.BlockResponse{
		// Height: , // TODO: fix me
		Hash:      hash.Bytes(),
		BlockTime: timestamp,
		Header:    header,
		Txs:       trxs,
		PrevCert:  prevCert,
	}

	return res, nil
}

func (s *blockchainServer) GetAccount(ctx context.Context,
	request *zarb.AccountRequest) (*zarb.AccountResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %v", err)
	}
	acc := s.state.AccountByAddress(addr)
	if acc == nil {
		return nil, status.Errorf(codes.InvalidArgument, "account not found")
	}
	res := &zarb.AccountResponse{
		Account: accountToProto(acc),
	}

	return res, nil
}

func (s *blockchainServer) GetValidatorByNumber(ctx context.Context,
	request *zarb.ValidatorByNumberRequest) (*zarb.ValidatorResponse, error) {
	val := s.state.ValidatorByNumber(request.Number)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &zarb.ValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidator(ctx context.Context,
	request *zarb.ValidatorRequest) (*zarb.ValidatorResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %v", err.Error())
	}
	val := s.state.ValidatorByAddress(addr)
	if val == nil {
		return nil, status.Errorf(codes.NotFound, "validator not found")
	}

	// TODO: make a function
	// proto validator from native validator
	return &zarb.ValidatorResponse{
		Validator: validatorToProto(val),
	}, nil
}

func (s *blockchainServer) GetValidators(ctx context.Context,
	request *zarb.ValidatorsRequest) (*zarb.ValidatorsResponse, error) {
	validators := s.state.CommitteeValidators()
	validatorsResp := make([]*zarb.ValidatorInfo, 0)
	for _, val := range validators {
		validatorsResp = append(validatorsResp, validatorToProto(val))
	}
	return &zarb.ValidatorsResponse{Validators: validatorsResp}, nil
}

func validatorToProto(val *validator.Validator) *zarb.ValidatorInfo {
	return &zarb.ValidatorInfo{
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

func accountToProto(acc *account.Account) *zarb.AccountInfo {
	return &zarb.AccountInfo{
		Address:  acc.Address().String(),
		Number:   acc.Number(),
		Sequence: acc.Sequence(),
		Balance:  acc.Balance(),
	}
}
