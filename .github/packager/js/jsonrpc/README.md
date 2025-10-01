# pactus-jsonrpc

JavaScript client for interacting with the [Pactus](https://pactus.org) blockchain via JSON-RPC.

## Installation

```bash
npm install pactus-jsonrpc
```

## Usage

```javascript
import PactusOpenRPC from "pactus-jsonrpc";

const jsonrpcClient = new PactusOpenRPC({
  transport: {
    type: "http",
    host: "127.0.0.1",
    port: 8545
  },
});

const blockchainInfo = await jsonrpcClient.pactusBlockchainGetBlockchainInfo();
console.log(JSON.stringify(blockchainInfo, null, 2));
```
