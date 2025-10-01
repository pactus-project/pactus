# pactus-grpc

Python client for interacting with the [Pactus](https://pactus.org) blockchain via gRPC.

## Installation

```bash
pip install pactus-grpc
```

## Usage

```python
import asyncio
import grpc
from pactus_grpc import blockchain_pb2_grpc, blockchain_pb2, network_pb2_grpc, network_pb2


async def main():
    channel = grpc.aio.insecure_channel("127.0.0.1:50051")
    blockchain_stub = blockchain_pb2_grpc.BlockchainStub(channel)
    blockchain_request = blockchain_pb2.GetBlockchainInfoRequest()
    blockchain_response = await blockchain_stub.GetBlockchainInfo(blockchain_request)
    print(blockchain_response)

    await channel.close()


if __name__ == "__main__":
    asyncio.run(main())
```
