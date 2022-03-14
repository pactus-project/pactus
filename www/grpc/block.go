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
	hash := zs.state.BlockHash(int(height))
	if hash.IsUndef() {
		return nil, status.Errorf(codes.NotFound, "No block found with this height")
	}
	return &zarb.BlockHashResponse{
		Hash: hash.RawBytes(),
	}, nil
}

func (zs *zarbServer) GetBlock(ctx context.Context, request *zarb.BlockRequest) (*zarb.BlockResponse, error) {
	hash, err := hash.FromRawBytes(request.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "provided hash is not Valid")
	}
	block := zs.state.Block(hash)
	if block == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Block not found")
	}
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
			Version:         int32(block.Header().Version()),
			PrevBlockHash:   block.Header().PrevBlockHash().RawBytes(),
			StateRoot:       block.Header().StateRoot().RawBytes(),
			TxsRoot:         block.Header().TxsRoot().RawBytes(),
			PrevCertHash:    block.Header().PrevCertificateHash().RawBytes(),
			SortitionSeed:   sortitionSeed,
			ProposerAddress: block.Header().ProposerAddress().String(),
		}

	}

	//TODO: Cache for better performance
	//populate BLOCK_TRANSACTIONS
	if request.Verbosity.Number() > 1 {
		for _, trx := range block.Transactions() {
			tranactions = append(tranactions, zs.encodeTransaction(trx))
		}
	}

	res := &zarb.BlockResponse{
		Hash:      hash.RawBytes(),
		BlockTime: timestamp,
		Header:    header,
		Txs:       tranactions,
		PrevCert:  prevCert,
	}

	return res, nil
}
