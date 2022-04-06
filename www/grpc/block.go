package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (zs *zarbServer) GetBlockHash(ctx context.Context, request *zarb.BlockHashRequest) (*zarb.BlockHashResponse, error) {
	height := request.GetHeight()
	hash := zs.state.BlockHash(height)
	if hash.IsUndef() {
		return nil, status.Errorf(codes.NotFound, "block hash not found with this height")
	}
	return &zarb.BlockHashResponse{
		Hash: hash.Bytes(),
	}, nil
}

func (zs *zarbServer) GetBlock(ctx context.Context, request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	hash, err := hash.FromBytes(request.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "provided hash is not Valid")
	}
	block := zs.state.Block(hash)
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "block not found")
	}
	timestamp := timestamppb.New(block.Header().Time())
	header := &zarb.BlockHeaderInfo{}
	var prevCert *zarb.CertificateInfo

	//populate BLOCK_DATA
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

	//TODO: Cache for better performance
	//populate BLOCK_TRANSACTIONS
	tranactions := make([]*zarb.TransactionInfo, 0, block.Transactions().Len())
	if request.Verbosity.Number() > 1 {
		for _, trx := range block.Transactions() {
			tranactions = append(tranactions, transactionToProto(trx))
		}
	}

	res := &zarb.BlockResponse{
		// Height: , // TODO: fix me
		Hash:      hash.Bytes(),
		BlockTime: timestamp,
		Header:    header,
		Txs:       tranactions,
		PrevCert:  prevCert,
	}

	return res, nil
}
