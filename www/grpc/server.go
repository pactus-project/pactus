package grpc

import (
	"context"
	"net"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/zmq"
	"google.golang.org/grpc"
)

type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    *Config
	listener  net.Listener
	address   string
	grpc      *grpc.Server
	state     state.Facade
	net       network.Network
	sync      sync.Synchronizer
	consMgr   consensus.ManagerReader
	walletMgr *wallet.Manager
	zmq       *zmq.Server
	logger    *logger.SubLogger
}

func NewServer(conf *Config, state state.Facade, sync sync.Synchronizer,
	network network.Network, consMgr consensus.ManagerReader,
	walletMgr *wallet.Manager,
	zmq *zmq.Server,
) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		config:    conf,
		state:     state,
		sync:      sync,
		net:       network,
		consMgr:   consMgr,
		walletMgr: walletMgr,
		zmq:       zmq,
		logger:    logger.NewSubLogger("_grpc", nil),
	}
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	return s.startListening(listener)
}

func (s *Server) startListening(listener net.Listener) error {
	opts := make([]grpc.UnaryServerInterceptor, 0)

	if s.config.BasicAuth != "" {
		opts = append(opts, BasicAuth(s.config.BasicAuth))
	}

	opts = append(opts, s.Recovery())

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))

	blockchainServer := newBlockchainServer(s)
	transactionServer := newTransactionServer(s)
	networkServer := newNetworkServer(s)
	utilServer := newUtilsServer(s)

	pactus.RegisterBlockchainServiceServer(grpcServer, blockchainServer)
	pactus.RegisterTransactionServiceServer(grpcServer, transactionServer)
	pactus.RegisterNetworkServiceServer(grpcServer, networkServer)
	pactus.RegisterUtilsServiceServer(grpcServer, utilServer)

	if s.config.EnableWallet {
		walletServer := newWalletServer(s, s.walletMgr)

		pactus.RegisterWalletServiceServer(grpcServer, walletServer)
	}

	s.listener = listener
	s.address = listener.Addr().String()
	s.grpc = grpcServer

	s.logger.Info("grpc started listening", "address", listener.Addr().String())
	go func() {
		s.logger.Info("grpc server started", "addr", listener.Addr())
		if err := s.grpc.Serve(listener); err != nil {
			s.logger.Error("error on grpc serve", "error", err)
		}
	}()

	return s.startGateway(s.address)
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.grpc != nil {
		s.grpc.Stop()
		_ = s.listener.Close()
	}
}
