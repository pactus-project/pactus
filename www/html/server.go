package html

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	ret "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx         context.Context
	config      *Config
	listener    net.Listener
	server      *http.Server
	grpcConn    *grpc.ClientConn
	blockchain  pactus.BlockchainClient
	transaction pactus.TransactionClient
	network     pactus.NetworkClient
	router      *mux.Router
	enableAuth  bool
	logger      *logger.SubLogger
}

// init disables default pprof handlers registered by importing net/http/pprof.
// Your pprof is showing (https://mmcloughlin.com/posts/your-pprof-is-showing)
func init() {
	http.DefaultServeMux = http.NewServeMux()
}

func NewServer(ctx context.Context, conf *Config, enableAuth bool) *Server {
	return &Server{
		ctx:        ctx,
		config:     conf,
		enableAuth: enableAuth,
		logger:     logger.NewSubLogger("_html", nil),
	}
}

func (s *Server) StartServer(grpcServer string) error {
	if !s.config.Enable {
		return nil
	}

	dialOpts := make([]grpc.DialOption, 0)
	dialOpts = append(dialOpts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(ret.UnaryClientInterceptor()),
	)
	grpcConn, err := grpc.NewClient(
		grpcServer,
		dialOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcConn = grpcConn
	s.blockchain = pactus.NewBlockchainClient(grpcConn)
	s.transaction = pactus.NewTransactionClient(grpcConn)
	s.network = pactus.NewNetworkClient(grpcConn)

	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.RootHandler)
	s.router.HandleFunc("/blockchain/", s.BlockchainHandler)
	s.router.HandleFunc("/consensus", s.ConsensusHandler)
	s.router.HandleFunc("/network", s.NetworkHandler)
	s.router.HandleFunc("/node", s.NodeHandler)
	s.router.HandleFunc("/block/hash/{hash}", s.GetBlockByHashHandler)
	s.router.HandleFunc("/block/height/{height}", s.GetBlockByHeightHandler)
	s.router.HandleFunc("/transaction/id/{id}", s.GetTransactionHandler)
	s.router.HandleFunc("/txpool", s.GetTxPoolContentHandler)
	s.router.HandleFunc("/account/address/{address}", s.GetAccountHandler)
	s.router.HandleFunc("/validator/address/{address}", s.GetValidatorHandler)
	s.router.HandleFunc("/validator/number/{number}", s.GetValidatorByNumberHandler)
	s.router.HandleFunc("/metrics/prometheus", promhttp.Handler().ServeHTTP)

	if s.config.EnablePprof {
		http.HandleFunc("/debug/pprof/", pprof.Index)
		http.HandleFunc("/debug/pprof/profile", pprof.Profile)
		http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		http.HandleFunc("/debug/pprof/trace", pprof.Trace)
		s.router.HandleFunc("/debug/pprof", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/debug/pprof/", http.StatusPermanentRedirect)
		})
	}

	if s.enableAuth {
		http.Handle("/", handlers.RecoveryHandler()(basicAuth(s.router)))
	} else {
		http.Handle("/", handlers.RecoveryHandler()(s.router))
	}

	listener, err := util.NetworkListen(s.ctx, "tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.listener = listener
	s.server = &http.Server{
		Addr:              listener.Addr().String(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		s.logger.Info("HTML server start listening", "address", listener.Addr())
		if err := s.server.Serve(listener); err != nil {
			s.logger.Debug("error on HTML server", "error", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	if s.server != nil {
		_ = s.server.Shutdown(s.ctx)
		_ = s.server.Close()
		_ = s.listener.Close()
	}

	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	if s.enableAuth {
		if _, _, ok := r.BasicAuth(); !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}
	}

	buf := new(bytes.Buffer)
	buf.WriteString("<html><body><br>")

	err := s.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
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

func (*Server) writeHTML(w http.ResponseWriter, html string) int {
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

func (t *tableMaker) addRowBlockHash(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/block/hash/%s\">%s</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowAccAddress(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/account/address/%s\">%s</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowValAddress(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/validator/address/%s\">%s</a></td></tr>", key, val, val)
}

func (t *tableMaker) addRowTxID(key, val string) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td><a href=\"/transaction/id/%s\">%s</a></td></tr>", key, val, val)
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

func (t *tableMaker) addRowAmount(key string, amt amount.Amount) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%s</td></tr>",
		key, amt.String())
}

func (t *tableMaker) addRowPower(key string, power int64) {
	amt := amount.Amount(power)
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%s</td></tr>",
		key, amt.String())
}

func (t *tableMaker) addRowFloat64(key string, val float64) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%v</td></tr>", key, val)
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

func (t *tableMaker) addRowDouble(key string, val float64) {
	fmt.Fprintf(t.w, "<tr><td>%s</td><td>%f</td></tr>", key, val)
}

func (t *tableMaker) html() string {
	t.w.WriteString("</table>")

	return t.w.String()
}
