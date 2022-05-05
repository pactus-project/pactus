package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/capnp"
	"zombiezen.com/go/capnproto2/rpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	router   *mux.Router
	capnp    capnp.ZarbServer
	listener net.Listener
	logger   *logger.Logger
}

func NewServer(conf *Config) *Server {
	return &Server{
		ctx:    context.Background(),
		config: conf,
		logger: logger.NewLogger("_http", nil),
	}
}

func (s *Server) StartServer(capnpServer string) error {
	if !s.config.Enable {
		return nil
	}

	c, err := net.Dial("tcp", capnpServer)
	if err != nil {
		return err
	}

	conn := rpc.NewConn(rpc.StreamTransport(c))
	s.capnp = capnp.ZarbServer{Client: conn.Bootstrap(s.ctx)}
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.RootHandler)
	s.router.HandleFunc("/blockchain/", s.BlockchainHandler)
	s.router.HandleFunc("/block/hash/{hash}", s.GetBlockByHashHandler)
	s.router.HandleFunc("/block/height/{height}", s.GetBlockByHeightHandler)
	s.router.HandleFunc("/block_hash/height/{height}", s.GetBlockHashHandler)
	s.router.HandleFunc("/transaction/id/{id}", s.GetTransactionHandler)
	s.router.HandleFunc("/account/address/{address}", s.GetAccountHandler)
	s.router.HandleFunc("/validator/address/{address}", s.GetValidatorHandler)
	s.router.HandleFunc("/network", s.NetworkHandler)
	http.Handle("/", handlers.RecoveryHandler()(s.router))

	l, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.logger.Info("http started listening", "address", l.Addr().String())
	s.listener = l
	go func() {
		for {
			err := http.Serve(l, nil)
			if err != nil {
				s.logger.Error("error on a connection", "err", err)
			}
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.ctx.Done()

	if s.capnp.Client != nil {
		s.capnp.Client.Close()
	}
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
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
			buf.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a></br>", link, pathTemplate))
		}

		return nil
	})

	if err != nil {
		s.logger.Error("unable to walk through methods", "err", err)
		return
	}

	buf.WriteString("</body></html>")
	s.writeHTML(w, buf.String())
}

func (s *Server) writeError(w http.ResponseWriter, err error) int {
	s.logger.Error("an error occurred", "err", err)

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
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td><a href=\"/block/hash/%x\">%x</a></td></tr>", key, val, val))
}
func (t *tableMaker) addRowAccAddress(key, val string) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td><a href=\"/account/address/%s\">%s</a></td></tr>", key, val, val))
}
func (t *tableMaker) addRowValAddress(key, val string) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td><a href=\"/validator/address/%s\">%s</a></td></tr>", key, val, val))
}
func (t *tableMaker) addRowTxID(key string, val []byte) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td><a href=\"/transaction/id/%x\">%x</a></td></tr>", key, val, val))
}
func (t *tableMaker) addRowString(key, val string) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", key, val))
}
func (t *tableMaker) addRowAmount(key string, val int64) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d.%d</td></tr>",
		key, val/100000000, val%10000000))
}
func (t *tableMaker) addRowInt(key string, val int) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", key, val))
}
func (t *tableMaker) addRowInts(key string, vals []int32) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td>", key))
	for _, n := range vals {
		t.w.WriteString(fmt.Sprintf("%d, ", n))
	}
	t.w.WriteString("</td></tr>")
}
func (t *tableMaker) addRowBytes(key string, val []byte) {
	t.w.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%x</td></tr>", key, val))
}
func (t *tableMaker) html() string {
	t.w.WriteString("</table>")
	return t.w.String()
}
