package grpc

import (
	"context"
	"unsafe"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type networkServer struct {
	sync   sync.Synchronizer
	logger *logger.SubLogger
}

func (s *networkServer) GetNetworkInfo(_ context.Context,
	_ *pactus.GetNetworkInfoRequest,
) (*pactus.GetNetworkInfoResponse, error) {
	peerset := s.sync.PeerSet()
	peers := make([]*pactus.PeerInfo, peerset.Len())

	for i, peer := range peerset.GetPeerList() {
		peers[i] = new(pactus.PeerInfo)
		p := peers[i]

		bs, err := cbor.Marshal(peer.Agent)
		if err != nil {
			s.logger.Error("couldn't marshal agent", "err", err)
			continue
		}
		p.Agent = string(bs)

		p.PeerId = []byte(peer.PeerID)
		p.Moniker = peer.Moniker
		p.Agent = peer.Agent
		p.Services = uint32(peer.Services)
		p.Height = peer.Height
		p.ReceivedMessages = int32(peer.ReceivedBundles)
		p.InvalidMessages = int32(peer.InvalidBundles)
		p.ReceivedBytes = *(*map[int32]int64)(unsafe.Pointer(&peer.ReceivedBytes))
		p.SentBytes = *(*map[int32]int64)(unsafe.Pointer(&peer.SentBytes))
		p.Status = int32(peer.Status)
		p.LastReceived = peer.LastReceived.Unix()
		p.LastBlockHash = peer.LastBlockHash.Bytes()

		for _, key := range peer.ConsensusKeys {
			p.ConsensusKeys = append(p.ConsensusKeys, key.String())
		}
	}

	sentBytes := peerset.SentBytes()
	receivedBytes := peerset.ReceivedBytes()

	return &pactus.GetNetworkInfoResponse{
		TotalSentBytes:     int32(peerset.TotalSentBytes()),
		TotalReceivedBytes: int32(peerset.TotalReceivedBytes()),
		SentBytes:          *(*map[int32]int64)(unsafe.Pointer(&sentBytes)),
		ReceivedBytes:      *(*map[int32]int64)(unsafe.Pointer(&receivedBytes)),
		StartedAt:          peerset.StartedAt().Unix(),
		Peers:              peers,
	}, nil
}

func (s *networkServer) GetNodeInfo(_ context.Context,
	_ *pactus.GetNodeInfoRequest,
) (*pactus.GetNodeInfoResponse, error) {
	return &pactus.GetNodeInfoResponse{
		Moniker: s.sync.Moniker(),
		Agent:   version.Agent(),
		PeerId:  []byte(s.sync.SelfID()),
	}, nil
}
