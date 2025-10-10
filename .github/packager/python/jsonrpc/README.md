# pactus-jsonrpc

Python client for interacting with the [Pactus](https://pactus.org) blockchain via JSON-RPC.

## Installation

```bash
pip install pactus-jsonrpc
```

## Usage

```python
import asyncio
from pactus_jsonrpc.client import PactusOpenRPCClient


async def main():
    client = PactusOpenRPCClient(
        headers={},
        client_url="http://127.0.0.1:8545"
    )

    blockchain_info = await client.pactus.blockchain.get_blockchain_info()
    print(blockchain_info)


if __name__ == "__main__":
    asyncio.run(main())
```
