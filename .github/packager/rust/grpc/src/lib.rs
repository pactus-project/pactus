//! # Pactus gRPC Client
//!
//! A Rust client library for interacting with the Pactus blockchain via gRPC.
//!
//! ## Example
//!
//! ```rust
//! use pactus_grpc::{blockchain_client::BlockchainClient, GetBlockchainInfoRequest};
//! use tonic::transport::Channel;
//!
//! #[tokio::main]
//! async fn main() -> Result<(), Box<dyn std::error::Error>> {
//!     let channel = Channel::from_static("http://127.0.0.1:50051")
//!         .connect()
//!         .await?;
//!
//!     let mut client = BlockchainClient::new(channel);
//!
//!     let request = tonic::Request::new(GetBlockchainInfoRequest {});
//!     let response = client.get_blockchain_info(request).await?;
//!     let info = response.into_inner();
//!
//!     println!("get_blockchain_info Response: {:?}", info);
//!
//!     Ok(())
//! }
//! ```

pub mod pactus;

// Re-export the main message types
pub use pactus::*;
