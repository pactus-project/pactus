// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: transaction.proto

/*
Package pactus is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package pactus

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Transaction_GetTransaction_0(ctx context.Context, marshaler runtime.Marshaler, client TransactionClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetTransactionRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.Bytes(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	val, ok = pathParams["verbosity"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "verbosity")
	}

	e, err = runtime.Enum(val, TransactionVerbosity_value)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "verbosity", err)
	}

	protoReq.Verbosity = TransactionVerbosity(e)

	msg, err := client.GetTransaction(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Transaction_GetTransaction_0(ctx context.Context, marshaler runtime.Marshaler, server TransactionServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetTransactionRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.Bytes(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	val, ok = pathParams["verbosity"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "verbosity")
	}

	e, err = runtime.Enum(val, TransactionVerbosity_value)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "verbosity", err)
	}

	protoReq.Verbosity = TransactionVerbosity(e)

	msg, err := server.GetTransaction(ctx, &protoReq)
	return msg, metadata, err

}

func request_Transaction_CalculateFee_0(ctx context.Context, marshaler runtime.Marshaler, client TransactionClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CalculateFeeRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["amount"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "amount")
	}

	protoReq.Amount, err = runtime.Int64(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "amount", err)
	}

	val, ok = pathParams["payloadType"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "payloadType")
	}

	e, err = runtime.Enum(val, PayloadType_value)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "payloadType", err)
	}

	protoReq.PayloadType = PayloadType(e)

	msg, err := client.CalculateFee(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Transaction_CalculateFee_0(ctx context.Context, marshaler runtime.Marshaler, server TransactionServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CalculateFeeRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["amount"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "amount")
	}

	protoReq.Amount, err = runtime.Int64(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "amount", err)
	}

	val, ok = pathParams["payloadType"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "payloadType")
	}

	e, err = runtime.Enum(val, PayloadType_value)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "payloadType", err)
	}

	protoReq.PayloadType = PayloadType(e)

	msg, err := server.CalculateFee(ctx, &protoReq)
	return msg, metadata, err

}

func request_Transaction_SendRawTransaction_0(ctx context.Context, marshaler runtime.Marshaler, client TransactionClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SendRawTransactionRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.Bytes(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := client.SendRawTransaction(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Transaction_SendRawTransaction_0(ctx context.Context, marshaler runtime.Marshaler, server TransactionServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SendRawTransactionRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.Bytes(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := server.SendRawTransaction(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterTransactionHandlerServer registers the http handlers for service Transaction to "mux".
// UnaryRPC     :call TransactionServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterTransactionHandlerFromEndpoint instead.
func RegisterTransactionHandlerServer(ctx context.Context, mux *runtime.ServeMux, server TransactionServer) error {

	mux.Handle("GET", pattern_Transaction_GetTransaction_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/pactus.Transaction/GetTransaction", runtime.WithHTTPPathPattern("/v1/transaction/get/id/{id}/verbosity/{verbosity}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Transaction_GetTransaction_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_GetTransaction_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Transaction_CalculateFee_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/pactus.Transaction/CalculateFee", runtime.WithHTTPPathPattern("/v1/transaction/calculate_fee/amount/{amount}/payload_type/{payloadType}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Transaction_CalculateFee_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_CalculateFee_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("PUT", pattern_Transaction_SendRawTransaction_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/pactus.Transaction/SendRawTransaction", runtime.WithHTTPPathPattern("/v1/transaction/send/data/{data}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Transaction_SendRawTransaction_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_SendRawTransaction_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterTransactionHandlerFromEndpoint is same as RegisterTransactionHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterTransactionHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterTransactionHandler(ctx, mux, conn)
}

// RegisterTransactionHandler registers the http handlers for service Transaction to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterTransactionHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterTransactionHandlerClient(ctx, mux, NewTransactionClient(conn))
}

// RegisterTransactionHandlerClient registers the http handlers for service Transaction
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "TransactionClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "TransactionClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "TransactionClient" to call the correct interceptors.
func RegisterTransactionHandlerClient(ctx context.Context, mux *runtime.ServeMux, client TransactionClient) error {

	mux.Handle("GET", pattern_Transaction_GetTransaction_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/pactus.Transaction/GetTransaction", runtime.WithHTTPPathPattern("/v1/transaction/get/id/{id}/verbosity/{verbosity}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Transaction_GetTransaction_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_GetTransaction_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Transaction_CalculateFee_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/pactus.Transaction/CalculateFee", runtime.WithHTTPPathPattern("/v1/transaction/calculate_fee/amount/{amount}/payload_type/{payloadType}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Transaction_CalculateFee_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_CalculateFee_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("PUT", pattern_Transaction_SendRawTransaction_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/pactus.Transaction/SendRawTransaction", runtime.WithHTTPPathPattern("/v1/transaction/send/data/{data}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Transaction_SendRawTransaction_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Transaction_SendRawTransaction_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_Transaction_GetTransaction_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3, 2, 4, 1, 0, 4, 1, 5, 4}, []string{"v1", "transaction", "get", "id", "verbosity"}, ""))

	pattern_Transaction_CalculateFee_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3, 2, 4, 1, 0, 4, 1, 5, 5}, []string{"v1", "transaction", "calculate_fee", "amount", "payload_type", "payloadType"}, ""))

	pattern_Transaction_SendRawTransaction_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3}, []string{"v1", "transaction", "send", "data"}, ""))
)

var (
	forward_Transaction_GetTransaction_0 = runtime.ForwardResponseMessage

	forward_Transaction_CalculateFee_0 = runtime.ForwardResponseMessage

	forward_Transaction_SendRawTransaction_0 = runtime.ForwardResponseMessage
)
