package grpc

import (
	"context"

	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type utilityServer struct {
	*Server
}

func newUtilityServer(server *Server) *utilityServer {
	return &utilityServer{
		Server: server,
	}
}

func (s *utilityServer) CalculateFee(_ context.Context,
	req *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	amount := req.Amount
	fee := s.state.CalculateFee(amount, payload.Type(req.PayloadType))

	if req.FixedAmount {
		amount -= fee
	}

	return &pactus.CalculateFeeResponse{
		Amount: amount,
		Fee:    fee,
	}, nil
}
