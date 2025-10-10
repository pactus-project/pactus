# pactus-jsonrpc

Rust client for interacting with the [Pactus](https://pactus.org) blockchain via JSON-RPC.

## Installation

```bash
cargo add pactus-jsonrpc
```

## Usage

```rust
use jsonrpsee::http_client::HttpClient;
use pactus_jsonrpc::pactus::PactusOpenRPC;

#[tokio::main]
async fn main() {
    let client = HttpClient::builder().build("http://127.0.0.1:8545").unwrap();
    let rpc: PactusOpenRPC<HttpClient> = PactusOpenRPC::new(client);

    let info = rpc.pactus_blockchain_get_blockchain_info().await.unwrap();
    println!("get_blockchain_info Response: {:?}", info);
}
```
