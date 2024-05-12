package ntp

import (
	"math"
	"time"

	"github.com/beevik/ntp"
	"github.com/pactus-project/pactus/util/logger"
)

type Server struct {
	threshold  time.Duration
	logger     *logger.SubLogger
	serverList []string
}

func NewNtpServer() *Server {
	server := &Server{
		threshold: 1 * time.Second,
		serverList: []string{
			"0.pool.ntp.org",
			"1.pool.ntp.org",
			"2.pool.ntp.org",
			"3.pool.ntp.org",
		},
	}

	server.logger = logger.NewSubLogger("_ntp", server)

	return server
}

func (s *Server) ClockOffset() time.Duration {
	clockOffset := time.Duration(math.MinInt64)

	for _, server := range s.serverList {
		remoteTime, err := ntp.Time(server)
		if err != nil {
			s.logger.Debug("ntp error", "server", server, "error", err)

			continue
		}

		clockOffset = time.Since(remoteTime)

		break
	}

	return clockOffset
}

func (s *Server) String() string {
	return "ntp"
}

func (s *Server) GetThreshold() time.Duration {
	return s.threshold
}

func (s *Server) OutOfSync(offset time.Duration) bool {
	return math.Abs(float64(offset)) > float64(s.threshold)
}
