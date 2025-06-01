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
	peerSet := s.sync.PeerSet()

	clockOffset, err := s.sync.ClockOffset()
	if err != nil {
		s.logger.Warn("failed to get clock offset", "err", err)
	}

	resp := &pactus.GetNodeInfoResponse{
		Moniker:       s.sync.Moniker(),
		Agent:         version.NodeAgent.String(),
		PeerId:        hex.EncodeToString([]byte(s.sync.SelfID())),
		Reachability:  s.net.ReachabilityStatus(),
		LocalAddrs:    s.net.HostAddrs(),
		StartedAt:     uint64(peerSet.StartedAt().Unix()),
		Protocols:     s.net.Protocols(),
		Services:      int32(s.sync.Services()),
		ServicesNames: s.sync.Services().String(),
		ClockOffset:   clockOffset.Seconds(),
		ConnectionInfo: &pactus.ConnectionInfo{
			Connections:         uint64(s.net.NumConnectedPeers()),
			InboundConnections:  uint64(s.net.NumInbound()),
			OutboundConnections: uint64(s.net.NumOutbound()),
		},
		ZmqPublishers: make([]*pactus.ZMQPublisherInfo, 0),
	}

	for _, publisher := range s.zmqPublishers {
		resp.ZmqPublishers = append(resp.ZmqPublishers, &pactus.ZMQPublisherInfo{
			Topic:   publisher.TopicName(),
			Address: publisher.Address(),
			Hwm:     int32(publisher.HWM()),
		})
	}

	return resp, nil
}

func (s *networkServer) GetNetworkInfo(_ context.Context,
	req *pactus.GetNetworkInfoRequest,
) (*pactus.GetNetworkInfoResponse, error) {
	peerSet := s.sync.PeerSet()
	peerInfos := make([]*pactus.PeerInfo, 0, peerSet.Len())

	peerSet.IteratePeers(func(peer *peer.Peer) bool {
		if req.OnlyConnected && !peer.Status.IsConnectedOrKnown() {
			return false
		}

		peerInfo := new(pactus.PeerInfo)
		peerInfos = append(peerInfos, peerInfo)

		data, err := cbor.Marshal(peer.Agent)
		if err != nil {
			s.logger.Error("couldn't marshal agent", "error", err)

			return false
		}
		peerInfo.Agent = string(data)

		peerInfo.PeerId = hex.EncodeToString([]byte(peer.PeerID))
		peerInfo.Moniker = peer.Moniker
		peerInfo.Agent = peer.Agent
		peerInfo.Address = peer.Address
		peerInfo.Direction = peer.Direction
		peerInfo.Services = uint32(peer.Services)
		peerInfo.Height = peer.Height
		peerInfo.Protocols = peer.Protocols
		peerInfo.Status = int32(peer.Status)
		peerInfo.LastSent = peer.LastSent.Unix()
		peerInfo.LastReceived = peer.LastReceived.Unix()
		peerInfo.LastBlockHash = peer.LastBlockHash.String()
		peerInfo.TotalSessions = int32(peer.TotalSessions)
		peerInfo.CompletedSessions = int32(peer.CompletedSessions)
		peerInfo.MetricInfo = metricToProto(peer.Metric)

		for _, key := range peer.ConsensusKeys {
			peerInfo.ConsensusKeys = append(peerInfo.ConsensusKeys, key.String())
			peerInfo.ConsensusAddresses = append(peerInfo.ConsensusAddresses, key.ValidatorAddress().String())
		}

		return false
	})

	return &pactus.GetNetworkInfoResponse{
		NetworkName:         s.net.Name(),
		ConnectedPeersCount: uint32(len(peerInfos)),
		ConnectedPeers:      peerInfos,
		MetricInfo:          metricToProto(peerSet.Metric()),
	}, nil
}

func metricToProto(metric metric.Metric) *pactus.MetricInfo {
	metricInfo := &pactus.MetricInfo{
		TotalInvalid: &pactus.CounterInfo{
			Bytes:   uint64(metric.TotalInvalid.Bytes),
			Bundles: uint64(metric.TotalInvalid.Bundles),
		},

		TotalSent: &pactus.CounterInfo{
			Bytes:   uint64(metric.TotalSent.Bytes),
			Bundles: uint64(metric.TotalSent.Bundles),
		},

		TotalReceived: &pactus.CounterInfo{
			Bytes:   uint64(metric.TotalReceived.Bytes),
			Bundles: uint64(metric.TotalReceived.Bundles),
		},
	}

	metricInfo.MessageSent = make(map[int32]*pactus.CounterInfo)
	for msgType, counter := range metric.MessageSent {
		metricInfo.MessageSent[int32(msgType)] = &pactus.CounterInfo{
			Bytes:   uint64(counter.Bytes),
			Bundles: uint64(counter.Bundles),
		}
	}

	metricInfo.MessageReceived = make(map[int32]*pactus.CounterInfo)
	for msgType, counter := range metric.MessageReceived {
		metricInfo.MessageReceived[int32(msgType)] = &pactus.CounterInfo{
			Bytes:   uint64(counter.Bytes),
			Bundles: uint64(counter.Bundles),
		}
	}

	return metricInfo
}
