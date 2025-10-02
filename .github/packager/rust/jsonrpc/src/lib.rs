//! # Pactus JSON-RPC Client
//!
//! A Rust client library for interacting with the Pactus blockchain via JSON-RPC.
//!
//! ## Example
//!
//! ```rust
//! use jsonrpsee::http_client::HttpClient;
//! use pactus_jsonrpc::pactus::PactusOpenRPC;
//!
//! #[tokio::main]
//! async fn main() {
//!     let client = HttpClient::builder().build("http://127.0.0.1:8545").unwrap();
//!     let rpc: PactusOpenRPC<HttpClient> = PactusOpenRPC::new(client);
//!
//!     let info = rpc.pactus_blockchain_get_blockchain_info().await.unwrap();
//!     println!("get_blockchain_info Response: {:?}", info);
//! }
//! ```

pub mod pactus;
