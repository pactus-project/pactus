package grpc

import (
	"context"
	"unsafe"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type networkServer struct {
	*Server
}

func (s *networkServer) GetPeersInfo(_ context.Context,
	_ *pactus.GetPeersInfoRequest,
) (*pactus.GetPeersInfoResponse, error) {
	ps := s.sync.PeerSet()
	peerInfos := make([]*pactus.PeerInfo, 0, ps.Len())

	ps.IteratePeers(func(peer *peerset.Peer) bool {
		p := new(pactus.PeerInfo)
		peerInfos = append(peerInfos, p)

		bs, err := cbor.Marshal(peer.Agent)
		if err != nil {
			s.logger.Error("couldn't marshal agent", "error", err)
			return false
		}
		p.Agent = string(bs)

		p.PeerId = []byte(peer.PeerID)
		p.Moniker = peer.Moniker
		p.Agent = peer.Agent
		p.Address = peer.Address
		p.Direction = peer.Direction
		p.Services = uint32(peer.Services)
		p.Height = peer.Height
		p.Protocols = peer.Protocols
		p.ReceivedMessages = int32(peer.ReceivedBundles)
		p.InvalidMessages = int32(peer.InvalidBundles)
		p.ReceivedBytes = *(*map[int32]int64)(unsafe.Pointer(&peer.ReceivedBytes))
		p.SentBytes = *(*map[int32]int64)(unsafe.Pointer(&peer.SentBytes))
		p.Status = int32(peer.Status)
		p.LastSent = peer.LastSent.Unix()
		p.LastReceived = peer.LastReceived.Unix()
		p.LastBlockHash = peer.LastBlockHash.Bytes()
		p.TotalSessions = int32(peer.TotalSessions)
		p.CompletedSessions = int32(peer.CompletedSessions)

		for _, key := range peer.ConsensusKeys {
			p.ConsensusKeys = append(p.ConsensusKeys, key.String())
		}

		return false
	})

	sentBytes := ps.SentBytes()
	receivedBytes := ps.ReceivedBytes()

	return &pactus.GetPeersInfoResponse{
		TotalSentBytes:     int32(ps.TotalSentBytes()),
		TotalReceivedBytes: int32(ps.TotalReceivedBytes()),
		SentBytes:          *(*map[int32]int64)(unsafe.Pointer(&sentBytes)),
		ReceivedBytes:      *(*map[int32]int64)(unsafe.Pointer(&receivedBytes)),
		StartedAt:          ps.StartedAt().Unix(),
		Peers:              peerInfos,
	}, nil
}

func (s *networkServer) GetNodeInfo(_ context.Context,
	_ *pactus.GetNodeInfoRequest,
) (*pactus.GetNodeInfoResponse, error) {
	return &pactus.GetNodeInfoResponse{
		Moniker:      s.sync.Moniker(),
		Agent:        version.Agent(),
		PeerId:       []byte(s.sync.SelfID()),
		Reachability: s.net.ReachabilityStatus(),
		Addrs:        s.net.HostAddrs(),
	}, nil
}

func (s *networkServer) GetNetworkInfo(_ context.Context,
	_ *pactus.GetNetworkInfoRequest,
) (*pactus.GetNetworkInfoResponse, error) {
	return &pactus.GetNetworkInfoResponse{
		ProtocolVersion: int32(version.ProtocolVersion),
		ConnectedPeers:  int32(s.net.NumConnectedPeers()),
		NetworkName:     s.net.Name(),
		Protocols:       s.net.Protocols(),
		LocalAddress:    s.net.HostAddrs(),
	}, nil
}
