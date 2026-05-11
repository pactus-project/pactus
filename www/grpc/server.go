package grpc

import (
	"context"
	"net"

	consmgr "github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/zmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	ctx           context.Context
	config        *Config
	listener      net.Listener
	server        *grpc.Server
	address       string
	state         state.Facade
	net           network.Network
	sync          sync.Synchronizer
	consMgr       consmgr.ManagerReader
	walletMgr     wltmgr.IManager
	zmqPublishers []zmq.Publisher
	logger        *logger.SubLogger
}

func NewServer(ctx context.Context, conf *Config, state state.Facade, sync sync.Synchronizer,
	network network.Network, consMgr consmgr.ManagerReader,
	walletMgr wltmgr.IManager,
	zmqPublishers []zmq.Publisher,
) *Server {
	return &Server{
		ctx:           ctx,
		config:        conf,
		state:         state,
		sync:          sync,
		net:           network,
		consMgr:       consMgr,
		walletMgr:     walletMgr,
		zmqPublishers: zmqPublishers,
		logger:        logger.NewSubLogger("_grpc", nil),
	}
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	listener, err := util.NetworkListen(s.ctx, "tcp", s.config.Listen)
	if err != nil {
		return err
	}

	return s.startListening(listener)
}

func (s *Server) startListening(listener net.Listener) error {
	unaryOpts := make([]grpc.UnaryServerInterceptor, 0)
	streamOpts := make([]grpc.StreamServerInterceptor, 0)

	if s.config.BasicAuth != "" {
		unaryOpts = append(unaryOpts, BasicAuth(s.config.BasicAuth))
		streamOpts = append(streamOpts, BasicAuthStream(s.config.BasicAuth))
	}

	unaryOpts = append(unaryOpts, s.Recovery())
	streamOpts = append(streamOpts, s.RecoveryStream())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryOpts...),
		grpc.ChainStreamInterceptor(streamOpts...),
	)

	blockchainServer := newBlockchainServer(s)
	transactionServer := newTransactionServer(s)
	networkServer := newNetworkServer(s)
	utilServer := newUtilsServer(s)

	pactus.RegisterBlockchainServer(grpcServer, blockchainServer)
	pactus.RegisterTransactionServer(grpcServer, transactionServer)
	pactus.RegisterNetworkServer(grpcServer, networkServer)
	pactus.RegisterUtilsServer(grpcServer, utilServer)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pactus.Blockchain_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pactus.Transaction_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pactus.Network_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pactus.Utils_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	if s.config.EnableWallet {
		walletServer := newWalletServer(s, s.walletMgr)

		pactus.RegisterWalletServer(grpcServer, walletServer)
		healthServer.SetServingStatus(pactus.Wallet_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	}

	s.listener = listener
	s.address = listener.Addr().String()
	s.server = grpcServer

	go func() {
		s.logger.Info("gRPC server start listening", "address", listener.Addr())
		if err := s.server.Serve(listener); err != nil {
			s.logger.Debug("error on gRPC server", "error", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	if s.server != nil {
		s.server.Stop()
		_ = s.listener.Close()
	}
}
