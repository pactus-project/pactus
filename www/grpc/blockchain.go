package grpc

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func (zs *zarbServer) GetBlockchainInfo(ctx context.Context,
	request *zarb.BlockchainInfoRequest) (*zarb.BlockchainInfoResponse, error) {
	height := zs.state.LastBlockHeight()

	return &zarb.BlockchainInfoResponse{
		LastBlockHeight: height,
		LastBlockHash:   zs.state.LastBlockHash().Bytes(),
	}, nil
}
func (zs *zarbServer) GetNetworkInfo(ctx context.Context,
	request *zarb.NetworkInfoRequest) (*zarb.NetworkInfoResponse, error) {
	// Create response peers
	rps := make([]*zarb.PeerInfo, int32(len(zs.sync.Peers())))

	// cast to response peers from synced peers
	for i, peer := range zs.sync.Peers() {
		rps[i] = new(zarb.PeerInfo)
		p := rps[i]

		bs, err := cbor.Marshal(peer.Agent)
		if err != nil {
			zs.logger.Error("couldn't marshal agent", "err", err)
			continue
		}
		p.Agent = string(bs)

		p.PeerId = []byte(peer.PeerID)
		p.Moniker = peer.Moniker
		p.Agent = peer.Agent
		p.PublicKey = peer.PublicKey.Bytes()
		p.Flags = int32(peer.Flags)
		p.Height = peer.Height
		p.ReceivedMessages = int32(peer.ReceivedBundles)
		p.InvalidMessages = int32(peer.InvalidBundles)
		p.ReceivedBytes = int32(peer.ReceivedBytes)
	}

	return &zarb.NetworkInfoResponse{
		SelfId: []byte(zs.sync.SelfID()),
		Peers:  rps,
	}, nil
}
