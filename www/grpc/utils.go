package grpc

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/crypto/bls"
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
	prvKey, err := bls.PrivateKeyFromString(req.PrivateKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid private key")
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
		return nil, status.Error(codes.InvalidArgument, "signature is invalid")
	}

	pub, err := bls.PublicKeyFromString(req.PublicKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "public key is invalid")
	}

	if err := pub.Verify([]byte(req.Message), sig); err == nil {
		return &pactus.VerifyMessageResponse{
			IsValid: true,
		}, nil
	}

	return &pactus.VerifyMessageResponse{
		IsValid: false,
	}, nil
}

func (*utilServer) BLSPublicKeyAggregation(_ context.Context,
	req *pactus.BLSPublicKeyAggregationRequest,
) (*pactus.BLSPublicKeyAggregationResponse, error) {
	if len(req.PublicKeys) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no public keys provided")
	}
	pubs := make([]*bls.PublicKey, len(req.PublicKeys))
	for i, pubKey := range req.PublicKeys {
		p, err := bls.PublicKeyFromString(pubKey)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid public key %s", pubKey))
		}
		pubs[i] = p
	}

	return &pactus.BLSPublicKeyAggregationResponse{
		PublicKey: bls.PublicKeyAggregate(pubs...).String(),
	}, nil
}

func (*utilServer) BLSSignatureAggregation(_ context.Context,
	req *pactus.BLSSignatureAggregationRequest,
) (*pactus.BLSSignatureAggregationResponse, error) {
	if len(req.Signatures) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no signatures provided")
	}
	sigs := make([]*bls.Signature, len(req.Signatures))
	for i, sig := range req.Signatures {
		s, err := bls.SignatureFromString(sig)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid signature %s", sig))
		}
		sigs[i] = s
	}

	return &pactus.BLSSignatureAggregationResponse{
		Signature: bls.SignatureAggregate(sigs...).String(),
	}, nil
}
