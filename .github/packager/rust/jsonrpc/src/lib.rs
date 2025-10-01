//! # Pactus JSON-RPC Client
//!
//! A Rust client library for interacting with the Pactus blockchain via JSON-RPC.
//!
//! ## Example
//!
//! ```rust
//! use jsonrpc_client_http::HttpTransport;
//! use pactus_jsonrpc::pactus::PactusOpenRPC;
//!
//! fn main(){
//!     let transport = HttpTransport::new().standalone().unwrap();
//!     let handle = transport.handle("http://localhost:8545").unwrap();
//!     let mut client = PactusOpenRPC::new(handle);
//!
//!     let info = client.pactus_blockchain_get_blockchain_info().call().unwrap();
//!     println!("get_blockchain_info Response: {:?}", info);
//! }
//! ```

pub mod pactus;
