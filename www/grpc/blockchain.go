package grpc

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetBlockchainInfo(ctx context.Context, request *zarb.BlockchainInfoRequest) (*zarb.BlockchainInfoResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented yet!")
}

func (zs *zarbServer) GetNetworkInfo(ctx context.Context, request *zarb.NetworkInfoRequest) (*zarb.NetworkInfoResponse, error) {
	// Create response peers
	rps := make([]*zarb.Peer, int32(len(zs.sync.Peers())))

	// cast to response peers from synced peers
	for i, peer := range zs.sync.Peers() {
		rps[i] = new(zarb.Peer)
		p := rps[i]

		bs, err := cbor.Marshal(peer.NodeVersion())
		if err != nil {
			zs.logger.Error("Couldn't marshal peer version", "err", err)
			continue
		}
		p.NodeVersion = bs

		p.PeerId = peer.PeerID().String()
		p.Moniker = peer.Moniker()
		p.PublicKey = peer.PublicKey().String()
		p.InitialBlockDownload = peer.InitialBlockDownload()
		p.Height = int32(peer.Height())
		p.ReceivedMessages = int32(peer.ReceivedMessages())
		p.InvalidMessages = int32(peer.InvalidMessages())
		p.ReceivedBytes = int32(peer.ReceivedBytes())
	}

	return &zarb.NetworkInfoResponse{
		PeerId: zs.sync.PeerID().String(),
		Peers:  rps,
	}, nil
}
