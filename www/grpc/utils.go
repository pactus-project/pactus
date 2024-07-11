package grpc

import (
	"context"

	"github.com/pactus-project/pactus/crypto/bls"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type utilServer struct {
	*Server
}

func newUtilsServer(server *Server) *utilServer {
	return &utilServer{
		Server: server,
	}
}

func (*utilServer) SignMessageWithPrivateKey(_ context.Context,
	req *pactus.SignMessageWithPrivateKeyRequest,
) (*pactus.SignMessageWithPrivateKeyResponse, error) {
	prvKey, err := bls.PrivateKeyFromString(req.PrivateKey)
	if err != nil {
		return nil, err
	}

	sig := prvKey.Sign([]byte(req.Message)).String()

	return &pactus.SignMessageWithPrivateKeyResponse{
		Signature: sig,
	}, nil
}

func (*utilServer) VerifyMessage(_ context.Context,
	req *pactus.VerifyMessageRequest,
) (*pactus.VerifyMessageResponse, error) {
	sig, err := bls.SignatureFromString(req.Signature)
	if err != nil {
		return &pactus.VerifyMessageResponse{
			IsValid: false,
		}, err
	}

	pub, err := bls.PublicKeyFromString(req.PublicKey)
	if err != nil {
		return &pactus.VerifyMessageResponse{
			IsValid: false,
		}, err
	}

	if err := pub.Verify([]byte(req.Message), sig); err != nil {
		return &pactus.VerifyMessageResponse{
			IsValid: false,
		}, err
	}

	return &pactus.VerifyMessageResponse{
		IsValid: true,
	}, nil
}
