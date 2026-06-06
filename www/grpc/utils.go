package grpc

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/wallet/vault"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	prvKey, err := vault.PrivateKeyFromString(req.PrivateKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sig := prvKey.Sign([]byte(req.Message)).String()

	return &pactus.SignMessageWithPrivateKeyResponse{
		Signature: sig,
	}, nil
}

func (*utilServer) VerifyMessage(_ context.Context,
	req *pactus.VerifyMessageRequest,
) (*pactus.VerifyMessageResponse, error) {
	pubKey, err := vault.PublicKeyFromString(req.PublicKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sig, err := vault.SignatureFromString(req.Signature, pubKey.Type())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = pubKey.Verify([]byte(req.Message), sig); err == nil {
		return &pactus.VerifyMessageResponse{
			IsValid: true,
		}, nil
	}

	return &pactus.VerifyMessageResponse{
		IsValid: false,
	}, nil
}

func (*utilServer) PublicKeyAggregation(_ context.Context,
	req *pactus.PublicKeyAggregationRequest,
) (*pactus.PublicKeyAggregationResponse, error) {
	pubs := make([]*bls.PublicKey, len(req.PublicKeys))
	for i, pubKey := range req.PublicKeys {
		p, err := bls.PublicKeyFromString(pubKey)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid public key %s", pubKey))
		}
		pubs[i] = p
	}

	aggPub, err := bls.PublicKeyAggregate(pubs...)
	if err != nil {
		return nil, err
	}

	return &pactus.PublicKeyAggregationResponse{
		PublicKey: aggPub.String(),
		Address:   aggPub.AccountAddress().String(),
	}, nil
}

func (*utilServer) SignatureAggregation(_ context.Context,
	req *pactus.SignatureAggregationRequest,
) (*pactus.SignatureAggregationResponse, error) {
	sigs := make([]*bls.Signature, len(req.Signatures))
	for i, sig := range req.Signatures {
		s, err := bls.SignatureFromString(sig)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid signature %s", sig))
		}
		sigs[i] = s
	}

	aggSig, err := bls.SignatureAggregate(sigs...)
	if err != nil {
		return nil, err
	}

	return &pactus.SignatureAggregationResponse{
		Signature: aggSig.String(),
	}, nil
}
