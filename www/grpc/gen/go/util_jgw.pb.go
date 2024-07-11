// Code generated by protoc-gen-jrpc-gateway. DO NOT EDIT.
// source: util.proto

/*
Package pactus is a reverse proxy.

It translates gRPC into JSON-RPC 2.0
*/
package pactus

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type UtilJsonRPC struct {
	client UtilClient
}

type paramsAndHeadersUtil struct {
	Headers metadata.MD     `json:"headers,omitempty"`
	Params  json.RawMessage `json:"params"`
}

// RegisterUtilJsonRPC register the grpc client Util for json-rpc.
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterUtilJsonRPC(conn *grpc.ClientConn) *UtilJsonRPC {
	return &UtilJsonRPC{
		client: NewUtilClient(conn),
	}
}

func (s *UtilJsonRPC) Methods() map[string]func(ctx context.Context, message json.RawMessage) (any, error) {
	return map[string]func(ctx context.Context, params json.RawMessage) (any, error){

		"pactus.util.sign_message_with_private_key": func(ctx context.Context, data json.RawMessage) (any, error) {
			req := new(SignMessageWithPrivateKeyRequest)

			var jrpcData paramsAndHeadersUtil

			if err := json.Unmarshal(data, &jrpcData); err != nil {
				return nil, err
			}

			err := protojson.Unmarshal(jrpcData.Params, req)
			if err != nil {
				return nil, err
			}

			return s.client.SignMessageWithPrivateKey(metadata.NewOutgoingContext(ctx, jrpcData.Headers), req)
		},

		"pactus.util.verify_message": func(ctx context.Context, data json.RawMessage) (any, error) {
			req := new(VerifyMessageRequest)

			var jrpcData paramsAndHeadersUtil

			if err := json.Unmarshal(data, &jrpcData); err != nil {
				return nil, err
			}

			err := protojson.Unmarshal(jrpcData.Params, req)
			if err != nil {
				return nil, err
			}

			return s.client.VerifyMessage(metadata.NewOutgoingContext(ctx, jrpcData.Headers), req)
		},
	}
}
