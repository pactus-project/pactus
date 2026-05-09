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
	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)

	if s.config.BasicAuth != "" {
		unaryInterceptors = append(unaryInterceptors, BasicAuth(s.config.BasicAuth))
		streamInterceptors = append(streamInterceptors, BasicAuthStream(s.config.BasicAuth))
	}

	unaryInterceptors = append(unaryInterceptors, s.Recovery())
	streamInterceptors = append(streamInterceptors, s.StreamRecovery())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	blockchainServer := newBlockchainServer(s)
	transactionServer := newTransactionServer(s)
	networkServer := newNetworkServer(s)
	utilServer := newUtilsServer(s)

	pactus.RegisterBlockchainServer(grpcServer, blockchainServer)
	pactus.RegisterTransactionServer(grpcServer, transactionServer)
	pactus.RegisterNetworkServer(grpcServer, networkServer)
	pactus.RegisterUtilsServer(grpcServer, utilServer)

	if s.config.EnableWallet {
		walletServer := newWalletServer(s, s.walletMgr)

		pactus.RegisterWalletServer(grpcServer, walletServer)
	}

	s.registerHealthServer(grpcServer)

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

func (s *Server) registerHealthServer(grpcServer *grpc.Server) {
	healthServer := health.NewServer()

	serviceNames := []string{
		"",
		pactus.Blockchain_ServiceDesc.ServiceName,
		pactus.Transaction_ServiceDesc.ServiceName,
		pactus.Network_ServiceDesc.ServiceName,
		pactus.Utils_ServiceDesc.ServiceName,
	}
	if s.config.EnableWallet {
		serviceNames = append(serviceNames, pactus.Wallet_ServiceDesc.ServiceName)
	}

	for _, serviceName := range serviceNames {
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	}

	healthpb.RegisterHealthServer(grpcServer, healthServer)
}

func (s *Server) StopServer() {
	if s.server != nil {
		s.server.Stop()
		_ = s.listener.Close()
	}
}
