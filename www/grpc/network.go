package grpc

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type networkServer struct {
	*Server
}

func newNetworkServer(server *Server) *networkServer {
	return &networkServer{
		Server: server,
	}
}

func (s *networkServer) GetNodeInfo(_ context.Context,
	_ *pactus.GetNodeInfoRequest,
) (*pactus.GetNodeInfoResponse, error) {
	ps := s.sync.PeerSet()

	services := []int32{}
	servicesNames := []string{}

	if s.sync.Services().IsNetwork() {
		services = append(services, int32(service.Network))
		servicesNames = append(servicesNames, "NETWORK")
	}

	clockOffset, err := s.sync.ClockOffset()
	if err != nil {
		s.logger.Warn("failed to get clock offset", "err", err)
	}

	return &pactus.GetNodeInfoResponse{
		Moniker:       s.sync.Moniker(),
		Agent:         version.NodeAgent.String(),
		PeerId:        []byte(s.sync.SelfID()),
		Reachability:  s.net.ReachabilityStatus(),
		LocalAddrs:    s.net.HostAddrs(),
		StartedAt:     uint64(ps.StartedAt().Unix()),
		Protocols:     s.net.Protocols(),
		Services:      services,
		ServicesNames: servicesNames,
		ClockOffset:   clockOffset.Seconds(),
		ConnectionInfo: &pactus.ConnectionInfo{
			Connections:         uint64(s.net.NumConnectedPeers()),
			InboundConnections:  uint64(s.net.NumInbound()),
			OutboundConnections: uint64(s.net.NumOutbound()),
		},
	}, nil
}

func (s *networkServer) GetNetworkInfo(_ context.Context,
	req *pactus.GetNetworkInfoRequest,
) (*pactus.GetNetworkInfoResponse, error) {
	ps := s.sync.PeerSet()
	peerInfos := make([]*pactus.PeerInfo, 0, ps.Len())

	ps.IteratePeers(func(peer *peer.Peer) bool {
		if req.OnlyConnected && !peer.Status.IsConnectedOrKnown() {
			return false
		}

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
		p.ReceivedBundles = int32(peer.ReceivedBundles)
		p.InvalidBundles = int32(peer.InvalidBundles)
		p.Status = int32(peer.Status)
		p.LastSent = peer.LastSent.Unix()
		p.LastReceived = peer.LastReceived.Unix()
		p.LastBlockHash = peer.LastBlockHash.Bytes()
		p.TotalSessions = int32(peer.TotalSessions)
		p.CompletedSessions = int32(peer.CompletedSessions)

		p.ReceivedBytes = make(map[int32]int64)
		for msgType, bytes := range peer.ReceivedBytes {
			p.ReceivedBytes[int32(msgType)] = bytes
		}

		p.SentBytes = make(map[int32]int64)
		for msgType, bytes := range peer.SentBytes {
			p.SentBytes[int32(msgType)] = bytes
		}

		for _, key := range peer.ConsensusKeys {
			p.ConsensusKeys = append(p.ConsensusKeys, key.String())
			p.ConsensusAddress = append(p.ConsensusAddress, key.ValidatorAddress().String())
		}

		return false
	})

	sentBytes := make(map[int32]int64)
	for msgType, bytes := range ps.SentBytes() {
		sentBytes[int32(msgType)] = bytes
	}

	receivedBytes := make(map[int32]int64)
	for msgType, bytes := range ps.ReceivedBytes() {
		receivedBytes[int32(msgType)] = bytes
	}

	return &pactus.GetNetworkInfoResponse{
		TotalSentBytes:      ps.TotalSentBytes(),
		TotalReceivedBytes:  ps.TotalReceivedBytes(),
		NetworkName:         s.net.Name(),
		ConnectedPeersCount: uint32(len(peerInfos)),
		ConnectedPeers:      peerInfos,
		SentBytes:           sentBytes,
		ReceivedBytes:       receivedBytes,
	}, nil
}
