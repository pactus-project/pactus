# gRPC

This directory contains required files for [gRPC](https://github.com/grpc-ecosystem/grpc-gateway) service.

In order to compile [pactus.proto](./proto/pactus.proto) file, run this command:

```bash
make proto
```

## gRPC-gateway

Pactus utilizes the [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to transform gRPC services
into RESTful APIs, accompanied by a [Swagger](https://swagger.io/) UI.

