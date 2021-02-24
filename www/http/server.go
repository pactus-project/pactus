package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/www/capnp"
	"zombiezen.com/go/capnproto2/rpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	router   *mux.Router
	server   capnp.ZarbServer
	listener net.Listener
	logger   *logger.Logger
}

func NewServer(conf *Config) (*Server, error) {
	return &Server{
		ctx:    context.Background(),
		config: conf,
		logger: logger.NewLogger("_http", nil),
	}, nil
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
	s.server = capnp.ZarbServer{Client: conn.Bootstrap(s.ctx)}
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.RootHandler)
	s.router.HandleFunc("/blockchain/", s.BlockchainHandler)
	s.router.HandleFunc("/block/height/{height}", s.GetBlockHandler)
	s.router.HandleFunc("/block_height/hash/{hash}", s.GetBlockHeightHandler)
	s.router.HandleFunc("/transaction/hash/{hash}", s.GetTransactionHandler)
	s.router.HandleFunc("/account/address/{address}", s.GetAccountHandler)
	s.router.HandleFunc("/validator/address/{address}", s.GetValidatorHandler)
	s.router.HandleFunc("/send_raw_transaction/{data}", s.SendRawTransactionHandler)
	s.router.HandleFunc("/network", s.NetworkHandler)
	http.Handle("/", handlers.RecoveryHandler()(s.router))

	l, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	s.logger.Info("Http started listening", "address", l.Addr().String())
	s.listener = l
	go func() {
		for {
			err := http.Serve(l, nil)
			if err != nil {
				s.logger.Error("Error on a connection", "err", err)
			}
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.ctx.Done()

	if s.server.Client != nil {
		s.server.Client.Close()
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
		s.logger.Error("Unable to walk through methods", "err", err)
		return
	}

	buf.WriteString("</body></html>")

	s.writeHTML(w, buf.String())
}

func (s *Server) writeJSON(w http.ResponseWriter, out interface{}) int {
	j, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		s.writeError(w, err)
		return 0
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	n, _ := io.WriteString(w, string(j))
	return n
}

func (s *Server) writePlainText(w http.ResponseWriter, out string) int {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	n, _ := io.WriteString(w, out)
	return n
}

func (s *Server) writeError(w http.ResponseWriter, err error) int {
	s.logger.Error("An error occurred", "err", err)

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
