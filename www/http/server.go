package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx         context.Context
	cancel      context.CancelFunc
	config      *Config
	router      *mux.Router
	grpcClient  *grpc.ClientConn
	httpServer  *http.Server
	blockchain  pactus.BlockchainClient
	transaction pactus.TransactionClient
	network     pactus.NetworkClient
	listener    net.Listener
	logger      *logger.SubLogger
}

func NewServer(conf *Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,
		config: conf,
		logger: logger.NewSubLogger("_http", nil),
	}
}

func (s *Server) StartServer(grpcServer string) error {
	if !s.config.Enable {
		return nil
	}

	conn, err := grpc.DialContext(
		s.ctx,
		grpcServer,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcClient = conn
	s.blockchain = pactus.NewBlockchainClient(conn)
	s.transaction = pactus.NewTransactionClient(conn)
	s.network = pactus.NewNetworkClient(conn)

	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.RootHandler)
	s.router.HandleFunc("/blockchain/", s.BlockchainHandler)
	s.router.HandleFunc("/consensus", s.ConsensusHandler)
	s.router.HandleFunc("/network", s.NetworkHandler)
	s.router.HandleFunc("/node", s.NodeHandler)
	s.router.HandleFunc("/block/hash/{hash}", s.GetBlockByHashHandler)
	s.router.HandleFunc("/block/height/{height}", s.GetBlockByHeightHandler)
	s.router.HandleFunc("/transaction/id/{id}", s.GetTransactionHandler)
	s.router.HandleFunc("/account/address/{address}", s.GetAccountHandler)
	s.router.HandleFunc("/validator/address/{address}", s.GetValidatorHandler)
	s.router.HandleFunc("/validator/number/{number}", s.GetValidatorByNumberHandler)
	s.router.HandleFunc("/metrics/prometheus", promhttp.Handler().ServeHTTP)
	http.Handle("/", handlers.RecoveryHandler()(s.router))

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.logger.Info("http started listening", "address", listener.Addr().String())
	s.listener = listener

	s.httpServer = &http.Server{
		Addr:              listener.Addr().String(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return

			default:
				err := s.httpServer.Serve(listener)
				if err != nil {
					s.logger.Error("error on a connection", "error", err)
				}
			}
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.httpServer != nil {
		_ = s.httpServer.Shutdown(s.ctx)
		s.httpServer.Close()
		s.listener.Close()
	}

	if s.grpcClient != nil {
		s.grpcClient.Close()
	}
}

func (s *Server) RootHandler(w http.ResponseWriter, _ *http.Request) {
	buf := new(bytes.Buffer)
	buf.WriteString("<html><body><br>")

	err := s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			link := pathTemplate
			i := strings.Index(link, "{")
			if i != -1 {
				link = link[0:i]
			}
			fmt.Fprintf(buf, "<a href=\"%s\">%s</a></br>", link, pathTemplate)
		}

		return nil
	})
	if err != nil {
		s.logger.Error("unable to walk through methods", "error", err)

		return
	}

	buf.WriteString("</body></html>")
	s.writeHTML(w, buf.String())
}

func (s *Server) writeError(w http.ResponseWriter, err error) int {
	s.logger.Error("an error occurred", "error", err)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	n, _ := io.WriteString(w, err.Error())

	return n
}

func (s *Server) writeHTML(w http.ResponseWriter, html string) int {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	n, _ := io.WriteString(w, html)

	return n
}

type tableMaker struct {
	w *bytes.Buffer
}

func newTableMaker() *tableMaker {
	t := &tableMaker{
		w: bytes.NewBufferString("<table>"),
	}

	return t
}

func (t *tableMaker) addRowBlockHash(key string, val []byte) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/block/hash/%x\">%x</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowAccAddress(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/account/address/%s\">%s</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowValAddress(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/validator/address/%s\">%s</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowTxID(key string, val []byte) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/transaction/id/%x\">%x</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowString(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%s</td></tr>", key, val)
}

func (t *tableMaker) addRowStrings(key string, val []string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%v</td></tr>", key, strings.Join(val, ","))
}

func (t *tableMaker) addRowTime(key string, sec int64) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%s</td></tr>", key, time.Unix(sec, 0).String())
}

func (t *tableMaker) addRowAmount(key string, change int64) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%s</td></tr>",
		key, util.ChangeToString(change))
}

func (t *tableMaker) addRowInt(key string, val int) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%d</td></tr>", key, val)
}

func (t *tableMaker) addRowBool(key string, val bool) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%v</td></tr>", key, val)
}

func (t *tableMaker) addRowInts(key string, vals []int32) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>", key)
	for _, n := range vals {
		fmt.Fprintf(t.w, "%d, ", n)
	}
	t.w.WriteString("</td></tr>")
}

func (t *tableMaker) addRowBytes(key string, val []byte) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%x</td></tr>", key, val)
}

func (t *tableMaker) addRowDouble(key string, val float64) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%f</td></tr>", key, val)
}

func (t *tableMaker) html() string {
	t.w.WriteString("</table>")

	return t.w.String()
}
