package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (zs *zarbServer) GetBlockHeight(ctx context.Context, request *zarb.BlockHeightRequest) (*zarb.BlockHeightResponse, error) {
	h, err := hash.FromString(request.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Hash provided is not Valid")
	}
	height := zs.state.BlockHeight(h)
	if height == -1 {
		return nil, status.Errorf(codes.NotFound, "No block found with the Hash provided")
	}
	return &zarb.BlockHeightResponse{
		Height: int64(height),
	}, nil
}

func (zs *zarbServer) GetBlock(ctx context.Context, request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	height := request.GetHeight()
	block := zs.state.Block(int(height))
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Block not found")
	}
	hash := block.Hash().String()
	timestamp := timestamppb.New(block.Header().Time())
	header := &zarb.BlockHeaderInfo{}
	tranactions := make([]*zarb.TransactionInfo, 0)
	var prevCert *zarb.CertificateInfo

	//populate BLOCK_DATA
	if request.Verbosity.Number() > 0 {
		seed := block.Header().SortitionSeed()

		sortitionSeed, err := seed.MarshalText()
		if err != nil {
			zs.logger.Error("couldn't marshal sortition seed: %v", err)
		}
		cert := block.PrevCertificate()
		if cert != nil {
			committers := make([]int32, len(block.PrevCertificate().Committers()))
			for c := range block.PrevCertificate().Committers() {
				committers = append(committers, int32(c))
			}
			absentees := make([]int32, len(block.PrevCertificate().Absentees()))
			for c := range block.PrevCertificate().Absentees() {
				absentees = append(absentees, int32(c))
			}
			prevCert = &zarb.CertificateInfo{
				Round:      int64(block.PrevCertificate().Round()),
				Committers: committers,
				Absentees:  absentees,
				Signature:  block.PrevCertificate().Signature().String(),
			}

		}
		header = &zarb.BlockHeaderInfo{
			Version:             int32(block.Header().Version()),
			PrevBlockHash:       block.Header().PrevBlockHash().String(),
			StateHash:           block.Header().StateHash().String(),
			TxIdsHash:           block.TxIDs().Hash().String(),
			PrevCertificateHash: block.Header().PrevCertificateHash().String(),
			SortitionSeed:       sortitionSeed,
			ProposerAddress:     block.Header().ProposerAddress().String(),
		}

	}

	//TODO: Cache for better performance
	//populate BLOCK_TRANSACTIONS
	if request.Verbosity.Number() > 1 {
		for _, id := range block.TxIDs().IDs() {
			t := zs.state.Transaction(id)
			tranactions = append(tranactions, zs.encodeTransaction(t))
		}
	}

	res := &zarb.BlockResponse{
		Hash:                hash,
		BlockTime:           timestamp,
		Header:              header,
		Tranactions:         tranactions,
		PreviousCertificate: prevCert,
	}

	return res, nil
}
