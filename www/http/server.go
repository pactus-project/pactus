package http

import (
	"bytes"
	"context"
	"fmt"
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

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	c, err := net.Dial("tcp", s.config.CapnpServer)
	if err != nil {
		return err
	}

	conn := rpc.NewConn(rpc.StreamTransport(c))
	s.server = capnp.ZarbServer{Client: conn.Bootstrap(s.ctx)}
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.RootHandler)
	s.router.HandleFunc("/block/height/{height}", s.BlockByHeightHandler)
	s.router.HandleFunc("/block/hash/{hash}", s.BlockByHashHandler)
	s.router.HandleFunc("/tx/hash/{hash}", s.TxHandler)
	s.router.HandleFunc("/account/number/{number}", s.AccountNumberHandler)
	http.Handle("/", handlers.RecoveryHandler()(s.router))

	l, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	s.logger.Info("Http started listening", "address", l.Addr().String())
	s.listener = l
	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("Recovered from a panic", r)
				}
			}()

			err := http.Serve(l, nil)
			if err != nil {
				s.logger.Error("Error on a connection", "error", err)
			}
		}
	}()

	return nil
}

func (s *Server) StopServer() error {
	if s.server.Client != nil {
		s.server.Client.Close()
	}

	return nil
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.WriteString("<html><body><br>")

	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

	buf.WriteString("</body></html>")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write(buf.Bytes())

}
