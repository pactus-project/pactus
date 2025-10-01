# pactus-jsonrpc

Rust client for interacting with the [Pactus](https://pactus.org) blockchain via JSON-RPC.

## Installation

```bash
cargo add pactus-jsonrpc
```

## Usage

```rust
use jsonrpc_client_http::HttpTransport;
use pactus_jsonrpc::pactus::PactusOpenRPC;

fn main() {
    let transport = HttpTransport::new().standalone().unwrap();
    let handle = transport.handle("http://127.0.0.1:8545").unwrap();
    let mut client = PactusOpenRPC::new(handle);

    let info = client.pactus_blockchain_get_blockchain_info().call().unwrap();
    println!("get_blockchain_info Response: {:?}", info);
}
```
