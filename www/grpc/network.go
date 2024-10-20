package grpc

import (
	"context"
	"encoding/hex"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/metric"
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

	clockOffset, err := s.sync.ClockOffset()
	if err != nil {
		s.logger.Warn("failed to get clock offset", "err", err)
	}

	return &pactus.GetNodeInfoResponse{
		Moniker:       s.sync.Moniker(),
		Agent:         version.NodeAgent.String(),
		PeerId:        hex.EncodeToString([]byte(s.sync.SelfID())),
		Reachability:  s.net.ReachabilityStatus(),
		LocalAddrs:    s.net.HostAddrs(),
		StartedAt:     uint64(ps.StartedAt().Unix()),
		Protocols:     s.net.Protocols(),
		Services:      int32(s.sync.Services()),
		ServicesNames: s.sync.Services().String(),
		ClockOffset:   clockOffset.Seconds(),
		ConnectionInfo: &pactus.ConnectionInfo{
			Connections:         uint64(s.net.NumConnectedPeers()),
			InboundConnections:  uint64(s.net.NumInbound()),
			OutboundConnections: uint64(s.net.NumOutbound()),
		},
		Fee: &pactus.FeeConfig{
			FixedFee:   s.txpoolConfig.Fee.FixedFee,
			DailyLimit: s.txpoolConfig.Fee.DailyLimit,
			UnitPrice:  s.txpoolConfig.Fee.UnitPrice,
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

		p.PeerId = hex.EncodeToString([]byte(peer.PeerID))
		p.Moniker = peer.Moniker
		p.Agent = peer.Agent
		p.Address = peer.Address
		p.Direction = peer.Direction
		p.Services = uint32(peer.Services)
		p.Height = peer.Height
		p.Protocols = peer.Protocols
		p.Status = int32(peer.Status)
		p.LastSent = peer.LastSent.Unix()
		p.LastReceived = peer.LastReceived.Unix()
		p.LastBlockHash = peer.LastBlockHash.String()
		p.TotalSessions = int32(peer.TotalSessions)
		p.CompletedSessions = int32(peer.CompletedSessions)
		p.MetricInfo = metricToProto(peer.Metric)

		for _, key := range peer.ConsensusKeys {
			p.ConsensusKeys = append(p.ConsensusKeys, key.String())
			p.ConsensusAddresses = append(p.ConsensusAddresses, key.ValidatorAddress().String())
		}

		return false
	})

	return &pactus.GetNetworkInfoResponse{
		NetworkName:         s.net.Name(),
		ConnectedPeersCount: uint32(len(peerInfos)),
		ConnectedPeers:      peerInfos,
		MetricInfo:          metricToProto(ps.Metric()),
	}, nil
}

func metricToProto(m metric.Metric) *pactus.MetricInfo {
	metricInfo := &pactus.MetricInfo{
		TotalInvalid: &pactus.CounterInfo{
			Bytes:   uint64(m.TotalInvalid.Bytes),
			Bundles: uint64(m.TotalInvalid.Bundles),
		},

		TotalSent: &pactus.CounterInfo{
			Bytes:   uint64(m.TotalSent.Bytes),
			Bundles: uint64(m.TotalSent.Bundles),
		},

		TotalReceived: &pactus.CounterInfo{
			Bytes:   uint64(m.TotalReceived.Bytes),
			Bundles: uint64(m.TotalReceived.Bundles),
		},
	}

	metricInfo.MessageSent = make(map[int32]*pactus.CounterInfo)
	for msgType, counter := range m.MessageSent {
		metricInfo.MessageSent[int32(msgType)] = &pactus.CounterInfo{
			Bytes:   uint64(counter.Bytes),
			Bundles: uint64(counter.Bundles),
		}
	}

	metricInfo.MessageReceived = make(map[int32]*pactus.CounterInfo)
	for msgType, counter := range m.MessageReceived {
		metricInfo.MessageReceived[int32(msgType)] = &pactus.CounterInfo{
			Bytes:   uint64(counter.Bytes),
			Bundles: uint64(counter.Bundles),
		}
	}

	return metricInfo
}
