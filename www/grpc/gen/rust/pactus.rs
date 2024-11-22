// @generated
/// Request message for retrieving transaction details.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTransactionRequest {
    /// The unique ID of the transaction to retrieve.
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    /// The verbosity level for transaction details.
    #[prost(enumeration="TransactionVerbosity", tag="2")]
    pub verbosity: i32,
}
/// Response message containing details of a transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTransactionResponse {
    /// The height of the block containing the transaction.
    #[prost(uint32, tag="1")]
    pub block_height: u32,
    /// The UNIX timestamp of the block containing the transaction.
    #[prost(uint32, tag="2")]
    pub block_time: u32,
    /// Detailed information about the transaction.
    #[prost(message, optional, tag="3")]
    pub transaction: ::core::option::Option<TransactionInfo>,
}
/// Request message for calculating transaction fee.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CalculateFeeRequest {
    /// The amount involved in the transaction, specified in NanoPAC.
    #[prost(int64, tag="1")]
    pub amount: i64,
    /// The type of transaction payload.
    #[prost(enumeration="PayloadType", tag="2")]
    pub payload_type: i32,
    /// Indicates if the amount should be fixed and include the fee.
    #[prost(bool, tag="3")]
    pub fixed_amount: bool,
}
/// Response message containing the calculated transaction fee.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CalculateFeeResponse {
    /// The calculated amount in NanoPAC.
    #[prost(int64, tag="1")]
    pub amount: i64,
    /// The calculated transaction fee in NanoPAC.
    #[prost(int64, tag="2")]
    pub fee: i64,
}
/// Request message for broadcasting a signed transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BroadcastTransactionRequest {
    /// The signed raw transaction data to be broadcasted.
    #[prost(string, tag="1")]
    pub signed_raw_transaction: ::prost::alloc::string::String,
}
/// Response message containing the ID of the broadcasted transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BroadcastTransactionResponse {
    /// The unique ID of the broadcasted transaction.
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of a transfer transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawTransferTransactionRequest {
    /// The lock time for the transaction. If not set, defaults to the last block
    /// height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// The sender's account address.
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    /// The receiver's account address.
    #[prost(string, tag="3")]
    pub receiver: ::prost::alloc::string::String,
    /// The amount to be transferred, specified in NanoPAC. Must be greater than 0.
    #[prost(int64, tag="4")]
    pub amount: i64,
    /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    #[prost(int64, tag="5")]
    pub fee: i64,
    /// A memo string for the transaction.
    #[prost(string, tag="6")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of a bond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawBondTransactionRequest {
    /// The lock time for the transaction. If not set, defaults to the last block
    /// height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// The sender's account address.
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    /// The receiver's validator address.
    #[prost(string, tag="3")]
    pub receiver: ::prost::alloc::string::String,
    /// The stake amount in NanoPAC. Must be greater than 0.
    #[prost(int64, tag="4")]
    pub stake: i64,
    /// The public key of the validator.
    #[prost(string, tag="5")]
    pub public_key: ::prost::alloc::string::String,
    /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    #[prost(int64, tag="6")]
    pub fee: i64,
    /// A memo string for the transaction.
    #[prost(string, tag="7")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of an unbond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawUnbondTransactionRequest {
    /// The lock time for the transaction. If not set, defaults to the last block
    /// height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// The address of the validator to unbond from.
    #[prost(string, tag="3")]
    pub validator_address: ::prost::alloc::string::String,
    /// A memo string for the transaction.
    #[prost(string, tag="4")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of a withdraw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawWithdrawTransactionRequest {
    /// The lock time for the transaction. If not set, defaults to the last block
    /// height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// The address of the validator to withdraw from.
    #[prost(string, tag="2")]
    pub validator_address: ::prost::alloc::string::String,
    /// The address of the account to withdraw to.
    #[prost(string, tag="3")]
    pub account_address: ::prost::alloc::string::String,
    /// The withdrawal amount in NanoPAC. Must be greater than 0.
    #[prost(int64, tag="4")]
    pub amount: i64,
    /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    #[prost(int64, tag="5")]
    pub fee: i64,
    /// A memo string for the transaction.
    #[prost(string, tag="6")]
    pub memo: ::prost::alloc::string::String,
}
/// Response message containing raw transaction data.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawTransactionResponse {
    /// The raw transaction data.
    #[prost(string, tag="1")]
    pub raw_transaction: ::prost::alloc::string::String,
    /// The unique ID of the transaction.
    #[prost(string, tag="2")]
    pub id: ::prost::alloc::string::String,
}
/// Payload for a transfer transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadTransfer {
    /// The sender's address.
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    /// The receiver's address.
    #[prost(string, tag="2")]
    pub receiver: ::prost::alloc::string::String,
    /// The amount to be transferred in NanoPAC.
    #[prost(int64, tag="3")]
    pub amount: i64,
}
/// Payload for a bond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadBond {
    /// The sender's address.
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    /// The receiver's address.
    #[prost(string, tag="2")]
    pub receiver: ::prost::alloc::string::String,
    /// The stake amount in NanoPAC.
    #[prost(int64, tag="3")]
    pub stake: i64,
    /// The public key of the validator.
    #[prost(string, tag="4")]
    pub public_key: ::prost::alloc::string::String,
}
/// Payload for a sortition transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadSortition {
    /// The validator address associated with the sortition proof.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    /// The proof for the sortition.
    #[prost(string, tag="2")]
    pub proof: ::prost::alloc::string::String,
}
/// Payload for an unbond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadUnbond {
    /// The address of the validator to unbond from.
    #[prost(string, tag="1")]
    pub validator: ::prost::alloc::string::String,
}
/// Payload for a withdraw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadWithdraw {
    /// The address of the validator to withdraw from.
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    /// The address of the account to withdraw to.
    #[prost(string, tag="2")]
    pub account_address: ::prost::alloc::string::String,
    /// The withdrawal amount in NanoPAC.
    #[prost(int64, tag="3")]
    pub amount: i64,
}
/// Information about a transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TransactionInfo {
    /// The unique ID of the transaction.
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    /// The raw transaction data.
    #[prost(string, tag="2")]
    pub data: ::prost::alloc::string::String,
    /// The version of the transaction.
    #[prost(int32, tag="3")]
    pub version: i32,
    /// The lock time for the transaction.
    #[prost(uint32, tag="4")]
    pub lock_time: u32,
    /// The value of the transaction in NanoPAC.
    #[prost(int64, tag="5")]
    pub value: i64,
    /// The fee for the transaction in NanoPAC.
    #[prost(int64, tag="6")]
    pub fee: i64,
    /// The type of transaction payload.
    #[prost(enumeration="PayloadType", tag="7")]
    pub payload_type: i32,
    /// A memo string for the transaction.
    #[prost(string, tag="8")]
    pub memo: ::prost::alloc::string::String,
    /// The public key associated with the transaction.
    #[prost(string, tag="9")]
    pub public_key: ::prost::alloc::string::String,
    /// The signature for the transaction.
    #[prost(string, tag="10")]
    pub signature: ::prost::alloc::string::String,
    #[prost(oneof="transaction_info::Payload", tags="30, 31, 32, 33, 34")]
    pub payload: ::core::option::Option<transaction_info::Payload>,
}
/// Nested message and enum types in `TransactionInfo`.
pub mod transaction_info {
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
    pub enum Payload {
        /// Transfer transaction payload.
        #[prost(message, tag="30")]
        Transfer(super::PayloadTransfer),
        /// Bond transaction payload.
        #[prost(message, tag="31")]
        Bond(super::PayloadBond),
        /// Sortition transaction payload.
        #[prost(message, tag="32")]
        Sortition(super::PayloadSortition),
        /// Unbond transaction payload.
        #[prost(message, tag="33")]
        Unbond(super::PayloadUnbond),
        /// Withdraw transaction payload.
        #[prost(message, tag="34")]
        Withdraw(super::PayloadWithdraw),
    }
}
/// Enumeration for different types of transaction payloads.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum PayloadType {
    /// Unknown payload type.
    Unknown = 0,
    /// Transfer payload type.
    TransferPayload = 1,
    /// Bond payload type.
    BondPayload = 2,
    /// Sortition payload type.
    SortitionPayload = 3,
    /// Unbond payload type.
    UnbondPayload = 4,
    /// Withdraw payload type.
    WithdrawPayload = 5,
}
impl PayloadType {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            PayloadType::Unknown => "UNKNOWN",
            PayloadType::TransferPayload => "TRANSFER_PAYLOAD",
            PayloadType::BondPayload => "BOND_PAYLOAD",
            PayloadType::SortitionPayload => "SORTITION_PAYLOAD",
            PayloadType::UnbondPayload => "UNBOND_PAYLOAD",
            PayloadType::WithdrawPayload => "WITHDRAW_PAYLOAD",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "UNKNOWN" => Some(Self::Unknown),
            "TRANSFER_PAYLOAD" => Some(Self::TransferPayload),
            "BOND_PAYLOAD" => Some(Self::BondPayload),
            "SORTITION_PAYLOAD" => Some(Self::SortitionPayload),
            "UNBOND_PAYLOAD" => Some(Self::UnbondPayload),
            "WITHDRAW_PAYLOAD" => Some(Self::WithdrawPayload),
            _ => None,
        }
    }
}
/// Enumeration for verbosity levels when requesting transaction details.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum TransactionVerbosity {
    /// Request transaction data only.
    TransactionData = 0,
    /// Request detailed transaction information.
    TransactionInfo = 1,
}
impl TransactionVerbosity {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            TransactionVerbosity::TransactionData => "TRANSACTION_DATA",
            TransactionVerbosity::TransactionInfo => "TRANSACTION_INFO",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "TRANSACTION_DATA" => Some(Self::TransactionData),
            "TRANSACTION_INFO" => Some(Self::TransactionInfo),
            _ => None,
        }
    }
}
/// Message to request account information based on an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAccountRequest {
    /// The address of the account to retrieve information for.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing the response with account information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAccountResponse {
    /// Detailed information about the account.
    #[prost(message, optional, tag="1")]
    pub account: ::core::option::Option<AccountInfo>,
}
/// Message to request validator addresses.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressesRequest {
}
/// Message containing the response with a list of validator addresses.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressesResponse {
    /// List of validator addresses.
    #[prost(string, repeated, tag="1")]
    pub addresses: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// Message to request validator information based on an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorRequest {
    /// The address of the validator to retrieve information for.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message to request validator information based on a validator number.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorByNumberRequest {
    /// The unique number of the validator to retrieve information for.
    #[prost(int32, tag="1")]
    pub number: i32,
}
/// Message containing the response with validator information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorResponse {
    /// Detailed information about the validator.
    #[prost(message, optional, tag="1")]
    pub validator: ::core::option::Option<ValidatorInfo>,
}
/// Message to request public key based on an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetPublicKeyRequest {
    /// The address for which to retrieve the public key.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing the response with the public key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetPublicKeyResponse {
    /// The public key associated with the provided address.
    #[prost(string, tag="1")]
    pub public_key: ::prost::alloc::string::String,
}
/// Message to request block information based on height and verbosity level.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockRequest {
    /// The height of the block to retrieve.
    #[prost(uint32, tag="1")]
    pub height: u32,
    /// The verbosity level for block information.
    #[prost(enumeration="BlockVerbosity", tag="2")]
    pub verbosity: i32,
}
/// Message containing the response with block information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockResponse {
    /// The height of the block.
    #[prost(uint32, tag="1")]
    pub height: u32,
    /// The hash of the block.
    #[prost(string, tag="2")]
    pub hash: ::prost::alloc::string::String,
    /// Block data, available only if verbosity level is set to BLOCK_DATA.
    #[prost(string, tag="3")]
    pub data: ::prost::alloc::string::String,
    /// The timestamp of the block.
    #[prost(uint32, tag="4")]
    pub block_time: u32,
    /// Header information of the block.
    #[prost(message, optional, tag="5")]
    pub header: ::core::option::Option<BlockHeaderInfo>,
    /// Certificate information of the previous block.
    #[prost(message, optional, tag="6")]
    pub prev_cert: ::core::option::Option<CertificateInfo>,
    /// List of transactions in the block, available when verbosity level is set to
    /// BLOCK_TRANSACTIONS.
    #[prost(message, repeated, tag="7")]
    pub txs: ::prost::alloc::vec::Vec<TransactionInfo>,
}
/// Message to request block hash based on height.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHashRequest {
    /// The height of the block to retrieve the hash for.
    #[prost(uint32, tag="1")]
    pub height: u32,
}
/// Message containing the response with the block hash.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHashResponse {
    /// The hash of the block.
    #[prost(string, tag="1")]
    pub hash: ::prost::alloc::string::String,
}
/// Message to request block height based on hash.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHeightRequest {
    /// The hash of the block to retrieve the height for.
    #[prost(string, tag="1")]
    pub hash: ::prost::alloc::string::String,
}
/// Message containing the response with the block height.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHeightResponse {
    /// The height of the block.
    #[prost(uint32, tag="1")]
    pub height: u32,
}
/// Message to request general information about the blockchain.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockchainInfoRequest {
}
/// Message containing the response with general blockchain information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockchainInfoResponse {
    /// The height of the last block in the blockchain.
    #[prost(uint32, tag="1")]
    pub last_block_height: u32,
    /// The hash of the last block in the blockchain.
    #[prost(string, tag="2")]
    pub last_block_hash: ::prost::alloc::string::String,
    /// The total number of accounts in the blockchain.
    #[prost(int32, tag="3")]
    pub total_accounts: i32,
    /// The total number of validators in the blockchain.
    #[prost(int32, tag="4")]
    pub total_validators: i32,
    /// The total power of the blockchain.
    #[prost(int64, tag="5")]
    pub total_power: i64,
    /// The power of the committee.
    #[prost(int64, tag="6")]
    pub committee_power: i64,
    /// List of committee validators.
    #[prost(message, repeated, tag="7")]
    pub committee_validators: ::prost::alloc::vec::Vec<ValidatorInfo>,
    /// If the blocks are subject to pruning.
    #[prost(bool, tag="8")]
    pub is_pruned: bool,
    /// Lowest-height block stored (only present if pruning is enabled)
    #[prost(uint32, tag="9")]
    pub pruning_height: u32,
    /// Timestamp of the last block in Unix format
    #[prost(int64, tag="10")]
    pub last_block_time: i64,
}
/// Message to request consensus information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetConsensusInfoRequest {
}
/// Message containing the response with consensus information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetConsensusInfoResponse {
    /// The proposal of the consensus info.
    #[prost(message, optional, tag="1")]
    pub proposal: ::core::option::Option<ProposalInfo>,
    /// List of consensus instances.
    #[prost(message, repeated, tag="2")]
    pub instances: ::prost::alloc::vec::Vec<ConsensusInfo>,
}
/// Request message to retrieve transactions in the transaction pool.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTxPoolContentRequest {
    /// The type of transactions to retrieve from the transaction pool. 0 means all
    /// types.
    #[prost(enumeration="PayloadType", tag="1")]
    pub payload_type: i32,
}
/// Response message containing transactions in the transaction pool.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTxPoolContentResponse {
    /// List of transactions currently in the pool.
    #[prost(message, repeated, tag="1")]
    pub txs: ::prost::alloc::vec::Vec<TransactionInfo>,
}
/// Message containing information about a validator.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ValidatorInfo {
    /// The hash of the validator.
    #[prost(string, tag="1")]
    pub hash: ::prost::alloc::string::String,
    /// The serialized data of the validator.
    #[prost(string, tag="2")]
    pub data: ::prost::alloc::string::String,
    /// The public key of the validator.
    #[prost(string, tag="3")]
    pub public_key: ::prost::alloc::string::String,
    /// The unique number assigned to the validator.
    #[prost(int32, tag="4")]
    pub number: i32,
    /// The stake of the validator in NanoPAC.
    #[prost(int64, tag="5")]
    pub stake: i64,
    /// The height at which the validator last bonded.
    #[prost(uint32, tag="6")]
    pub last_bonding_height: u32,
    /// The height at which the validator last participated in sortition.
    #[prost(uint32, tag="7")]
    pub last_sortition_height: u32,
    /// The height at which the validator will unbond.
    #[prost(uint32, tag="8")]
    pub unbonding_height: u32,
    /// The address of the validator.
    #[prost(string, tag="9")]
    pub address: ::prost::alloc::string::String,
    /// The availability score of the validator.
    #[prost(double, tag="10")]
    pub availability_score: f64,
}
/// Message containing information about an account.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AccountInfo {
    /// The hash of the account.
    #[prost(string, tag="1")]
    pub hash: ::prost::alloc::string::String,
    /// The serialized data of the account.
    #[prost(string, tag="2")]
    pub data: ::prost::alloc::string::String,
    /// The unique number assigned to the account.
    #[prost(int32, tag="3")]
    pub number: i32,
    /// The balance of the account in NanoPAC.
    #[prost(int64, tag="4")]
    pub balance: i64,
    /// The address of the account.
    #[prost(string, tag="5")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing information about the header of a block.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlockHeaderInfo {
    /// The version of the block.
    #[prost(int32, tag="1")]
    pub version: i32,
    /// The hash of the previous block.
    #[prost(string, tag="2")]
    pub prev_block_hash: ::prost::alloc::string::String,
    /// The state root hash of the blockchain.
    #[prost(string, tag="3")]
    pub state_root: ::prost::alloc::string::String,
    /// The sortition seed of the block.
    #[prost(string, tag="4")]
    pub sortition_seed: ::prost::alloc::string::String,
    /// The address of the proposer of the block.
    #[prost(string, tag="5")]
    pub proposer_address: ::prost::alloc::string::String,
}
/// Message containing information about a certificate.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CertificateInfo {
    /// The hash of the certificate.
    #[prost(string, tag="1")]
    pub hash: ::prost::alloc::string::String,
    /// The round of the certificate.
    #[prost(int32, tag="2")]
    pub round: i32,
    /// List of committers in the certificate.
    #[prost(int32, repeated, tag="3")]
    pub committers: ::prost::alloc::vec::Vec<i32>,
    /// List of absentees in the certificate.
    #[prost(int32, repeated, tag="4")]
    pub absentees: ::prost::alloc::vec::Vec<i32>,
    /// The signature of the certificate.
    #[prost(string, tag="5")]
    pub signature: ::prost::alloc::string::String,
}
/// Message containing information about a vote.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct VoteInfo {
    /// The type of the vote.
    #[prost(enumeration="VoteType", tag="1")]
    pub r#type: i32,
    /// The address of the voter.
    #[prost(string, tag="2")]
    pub voter: ::prost::alloc::string::String,
    /// The hash of the block being voted on.
    #[prost(string, tag="3")]
    pub block_hash: ::prost::alloc::string::String,
    /// The consensus round of the vote.
    #[prost(int32, tag="4")]
    pub round: i32,
    /// The change-proposer round of the vote.
    #[prost(int32, tag="5")]
    pub cp_round: i32,
    /// The change-proposer value of the vote.
    #[prost(int32, tag="6")]
    pub cp_value: i32,
}
/// Message containing information about a consensus instance.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConsensusInfo {
    /// The address of the consensus instance.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    /// Indicates whether the consensus instance is active and part of the
    /// committee.
    #[prost(bool, tag="2")]
    pub active: bool,
    /// The height of the consensus instance.
    #[prost(uint32, tag="3")]
    pub height: u32,
    /// The round of the consensus instance.
    #[prost(int32, tag="4")]
    pub round: i32,
    /// List of votes in the consensus instance.
    #[prost(message, repeated, tag="5")]
    pub votes: ::prost::alloc::vec::Vec<VoteInfo>,
}
/// Message containing information about a proposal.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ProposalInfo {
    /// The height of the proposal.
    #[prost(uint32, tag="1")]
    pub height: u32,
    /// The round of the proposal.
    #[prost(int32, tag="2")]
    pub round: i32,
    /// The block data of the proposal.
    #[prost(string, tag="3")]
    pub block_data: ::prost::alloc::string::String,
    /// The signature of the proposal, signed by the proposer.
    #[prost(string, tag="4")]
    pub signature: ::prost::alloc::string::String,
}
/// Enumeration for verbosity levels when requesting block information.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum BlockVerbosity {
    /// Request only block data.
    BlockData = 0,
    /// Request block information and transaction IDs.
    BlockInfo = 1,
    /// Request block information and detailed transaction data.
    BlockTransactions = 2,
}
impl BlockVerbosity {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            BlockVerbosity::BlockData => "BLOCK_DATA",
            BlockVerbosity::BlockInfo => "BLOCK_INFO",
            BlockVerbosity::BlockTransactions => "BLOCK_TRANSACTIONS",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "BLOCK_DATA" => Some(Self::BlockData),
            "BLOCK_INFO" => Some(Self::BlockInfo),
            "BLOCK_TRANSACTIONS" => Some(Self::BlockTransactions),
            _ => None,
        }
    }
}
/// Enumeration for types of votes.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum VoteType {
    /// Unknown vote type.
    VoteUnknown = 0,
    /// Prepare vote type.
    VotePrepare = 1,
    /// Precommit vote type.
    VotePrecommit = 2,
    /// Change proposer vote type.
    VoteChangeProposer = 3,
}
impl VoteType {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            VoteType::VoteUnknown => "VOTE_UNKNOWN",
            VoteType::VotePrepare => "VOTE_PREPARE",
            VoteType::VotePrecommit => "VOTE_PRECOMMIT",
            VoteType::VoteChangeProposer => "VOTE_CHANGE_PROPOSER",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "VOTE_UNKNOWN" => Some(Self::VoteUnknown),
            "VOTE_PREPARE" => Some(Self::VotePrepare),
            "VOTE_PRECOMMIT" => Some(Self::VotePrecommit),
            "VOTE_CHANGE_PROPOSER" => Some(Self::VoteChangeProposer),
            _ => None,
        }
    }
}
/// Request message for retrieving overall network information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNetworkInfoRequest {
    /// If true, returns only peers that are currently connected.
    #[prost(bool, tag="1")]
    pub only_connected: bool,
}
/// Response message containing information about the overall network.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNetworkInfoResponse {
    /// Name of the network.
    #[prost(string, tag="1")]
    pub network_name: ::prost::alloc::string::String,
    /// Number of connected peers.
    #[prost(uint32, tag="2")]
    pub connected_peers_count: u32,
    /// List of connected peers.
    #[prost(message, repeated, tag="3")]
    pub connected_peers: ::prost::alloc::vec::Vec<PeerInfo>,
    /// Metrics related to node activity.
    #[prost(message, optional, tag="4")]
    pub metric_info: ::core::option::Option<MetricInfo>,
}
/// Request message for retrieving information of the node.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNodeInfoRequest {
}
/// Response message containing information about a specific node in the network.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNodeInfoResponse {
    /// Moniker of the node.
    #[prost(string, tag="1")]
    pub moniker: ::prost::alloc::string::String,
    /// Version and agent details of the node.
    #[prost(string, tag="2")]
    pub agent: ::prost::alloc::string::String,
    /// Peer ID of the node.
    #[prost(string, tag="3")]
    pub peer_id: ::prost::alloc::string::String,
    /// Time the node was started (in epoch format).
    #[prost(uint64, tag="4")]
    pub started_at: u64,
    /// Reachability status of the node.
    #[prost(string, tag="5")]
    pub reachability: ::prost::alloc::string::String,
    /// Bitfield representing the services provided by the node.
    #[prost(int32, tag="6")]
    pub services: i32,
    /// Names of services provided by the node.
    #[prost(string, tag="7")]
    pub services_names: ::prost::alloc::string::String,
    /// List of addresses associated with the node.
    #[prost(string, repeated, tag="8")]
    pub local_addrs: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// List of protocols supported by the node.
    #[prost(string, repeated, tag="9")]
    pub protocols: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Offset between the node's clock and the network's clock (in seconds).
    #[prost(double, tag="13")]
    pub clock_offset: f64,
    /// Information about the node's connections.
    #[prost(message, optional, tag="14")]
    pub connection_info: ::core::option::Option<ConnectionInfo>,
}
/// Information about a peer in the network.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PeerInfo {
    /// Current status of the peer (e.g., connected, disconnected).
    #[prost(int32, tag="1")]
    pub status: i32,
    /// Moniker of the peer.
    #[prost(string, tag="2")]
    pub moniker: ::prost::alloc::string::String,
    /// Version and agent details of the peer.
    #[prost(string, tag="3")]
    pub agent: ::prost::alloc::string::String,
    /// Peer ID of the peer.
    #[prost(string, tag="4")]
    pub peer_id: ::prost::alloc::string::String,
    /// List of consensus keys used by the peer.
    #[prost(string, repeated, tag="5")]
    pub consensus_keys: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// List of consensus addresses used by the peer.
    #[prost(string, repeated, tag="6")]
    pub consensus_addresses: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Bitfield representing the services provided by the peer.
    #[prost(uint32, tag="7")]
    pub services: u32,
    /// Hash of the last block the peer knows.
    #[prost(string, tag="8")]
    pub last_block_hash: ::prost::alloc::string::String,
    /// Blockchain height of the peer.
    #[prost(uint32, tag="9")]
    pub height: u32,
    /// Time the last bundle sent to the peer (in epoch format).
    #[prost(int64, tag="10")]
    pub last_sent: i64,
    /// Time the last bundle received from the peer (in epoch format).
    #[prost(int64, tag="11")]
    pub last_received: i64,
    /// Network address of the peer.
    #[prost(string, tag="12")]
    pub address: ::prost::alloc::string::String,
    /// Connection direction (e.g., inbound, outbound).
    #[prost(string, tag="13")]
    pub direction: ::prost::alloc::string::String,
    /// List of protocols supported by the peer.
    #[prost(string, repeated, tag="14")]
    pub protocols: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Total download sessions with the peer.
    #[prost(int32, tag="15")]
    pub total_sessions: i32,
    /// Completed download sessions with the peer.
    #[prost(int32, tag="16")]
    pub completed_sessions: i32,
    /// Metrics related to peer activity.
    #[prost(message, optional, tag="17")]
    pub metric_info: ::core::option::Option<MetricInfo>,
}
/// ConnectionInfo contains information about the node's connections.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConnectionInfo {
    /// Total number of connections.
    #[prost(uint64, tag="1")]
    pub connections: u64,
    /// Number of inbound connections.
    #[prost(uint64, tag="2")]
    pub inbound_connections: u64,
    /// Number of outbound connections.
    #[prost(uint64, tag="3")]
    pub outbound_connections: u64,
}
/// MetricInfo contains data regarding network actvity.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MetricInfo {
    /// Total number of invalid bundles.
    #[prost(message, optional, tag="1")]
    pub total_invalid: ::core::option::Option<CounterInfo>,
    /// Total number of bundles sent.
    #[prost(message, optional, tag="2")]
    pub total_sent: ::core::option::Option<CounterInfo>,
    /// Total number of bundles received.
    #[prost(message, optional, tag="3")]
    pub total_received: ::core::option::Option<CounterInfo>,
    /// Number of sent bundles categorized by message type.
    #[prost(map="int32, message", tag="4")]
    pub message_sent: ::std::collections::HashMap<i32, CounterInfo>,
    /// Number of received bundles categorized by message type.
    #[prost(map="int32, message", tag="5")]
    pub message_received: ::std::collections::HashMap<i32, CounterInfo>,
}
/// CounterInfo holds data regarding byte and bundle counts.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CounterInfo {
    /// Total number of bytes.
    #[prost(uint64, tag="1")]
    pub bytes: u64,
    /// Total number of bundles.
    #[prost(uint64, tag="2")]
    pub bundles: u64,
}
/// Request message for sign message with private key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignMessageWithPrivateKeyRequest {
    /// The private key to sign the message.
    #[prost(string, tag="1")]
    pub private_key: ::prost::alloc::string::String,
    /// The message to sign.
    #[prost(string, tag="2")]
    pub message: ::prost::alloc::string::String,
}
/// Response message containing the generated signature.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignMessageWithPrivateKeyResponse {
    /// The signature of the message.
    #[prost(string, tag="1")]
    pub signature: ::prost::alloc::string::String,
}
/// Request message for verifying a message signature.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct VerifyMessageRequest {
    /// The signed message.
    #[prost(string, tag="1")]
    pub message: ::prost::alloc::string::String,
    /// The signature of the message.
    #[prost(string, tag="2")]
    pub signature: ::prost::alloc::string::String,
    /// The public key of the signer.
    #[prost(string, tag="3")]
    pub public_key: ::prost::alloc::string::String,
}
/// Response message containing the resualt of validation of signature and message.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct VerifyMessageResponse {
    /// Indicates if the signature is valid (true) or not (false).
    #[prost(bool, tag="1")]
    pub is_valid: bool,
}
/// Request message for aggregating BLS public keys.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlsPublicKeyAggregationRequest {
    /// The public keys to aggregate.
    #[prost(string, repeated, tag="1")]
    pub public_keys: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// Response message containing the aggregated public key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlsPublicKeyAggregationResponse {
    /// The aggregated public key.
    #[prost(string, tag="1")]
    pub public_key: ::prost::alloc::string::String,
    /// The aggregated public key account address.
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
}
/// Request message for aggregating BLS signatures.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlsSignatureAggregationRequest {
    /// The signatures to aggregate.
    #[prost(string, repeated, tag="1")]
    pub signatures: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// Response message containing the aggregated signature.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlsSignatureAggregationResponse {
    /// The aggregated signature.
    #[prost(string, tag="1")]
    pub signature: ::prost::alloc::string::String,
}
/// Message containing address information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddressInfo {
    /// The address string.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    /// The public key associated with the address.
    #[prost(string, tag="2")]
    pub public_key: ::prost::alloc::string::String,
    /// A label associated with the address.
    #[prost(string, tag="3")]
    pub label: ::prost::alloc::string::String,
    /// The Hierarchical Deterministic path of the address within the wallet.
    #[prost(string, tag="4")]
    pub path: ::prost::alloc::string::String,
}
/// Message containing transaction history information for an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct HistoryInfo {
    /// The transaction ID hash.
    #[prost(string, tag="1")]
    pub transaction_id: ::prost::alloc::string::String,
    /// The timestamp of the transaction.
    #[prost(uint32, tag="2")]
    pub time: u32,
    /// The payload type of the transaction.
    #[prost(string, tag="3")]
    pub payload_type: ::prost::alloc::string::String,
    /// A description of the transaction.
    #[prost(string, tag="4")]
    pub description: ::prost::alloc::string::String,
    /// The amount involved in the transaction.
    #[prost(int64, tag="5")]
    pub amount: i64,
}
/// Request message to get an address transaction history.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressHistoryRequest {
    /// The name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The address to retrieve the transaction history for.
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
}
/// Response message containing the address transaction history.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressHistoryResponse {
    /// Array of history information for the address.
    #[prost(message, repeated, tag="1")]
    pub history_info: ::prost::alloc::vec::Vec<HistoryInfo>,
}
/// Request message for generating a new address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNewAddressRequest {
    /// The name of the wallet to generate a new address.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The type of address to generate.
    #[prost(enumeration="AddressType", tag="2")]
    pub address_type: i32,
    /// A label for the new address.
    #[prost(string, tag="3")]
    pub label: ::prost::alloc::string::String,
    /// Password for the new address. It's required when address_type is ADDRESS_TYPE_ED25519_ACCOUNT.
    #[prost(string, tag="4")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the newly generated address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNewAddressResponse {
    /// The name of the wallet from which the address is generated.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Information about the newly generated address.
    #[prost(message, optional, tag="2")]
    pub address_info: ::core::option::Option<AddressInfo>,
}
/// Request message for restoring an existing wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RestoreWalletRequest {
    /// The name of the wallet to restore.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The mnemonic for wallet recovery.
    #[prost(string, tag="2")]
    pub mnemonic: ::prost::alloc::string::String,
    /// The password for securing the wallet.
    #[prost(string, tag="3")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the name of the restored wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RestoreWalletResponse {
    /// The name of the restored wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for creating a new wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateWalletRequest {
    /// The name of the new wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The password for securing the wallet.
    #[prost(string, tag="4")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the mnemonic for wallet recovery.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateWalletResponse {
    /// The mnemonic for wallet recovery.
    #[prost(string, tag="2")]
    pub mnemonic: ::prost::alloc::string::String,
}
/// Request message for loading an existing wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LoadWalletRequest {
    /// The name of the wallet to load.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the name of the loaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LoadWalletResponse {
    /// The name of the loaded wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for unloading a currently loaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnloadWalletRequest {
    /// The name of the wallet to unload.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the name of the unloaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnloadWalletResponse {
    /// The name of the unloaded wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for obtaining the validator address associated with a public
/// key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressRequest {
    /// The public key for which the validator address is requested.
    #[prost(string, tag="1")]
    pub public_key: ::prost::alloc::string::String,
}
/// Response message containing the validator address corresponding to a public
/// key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressResponse {
    /// The validator address associated with the public key.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Request message for signing a raw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignRawTransactionRequest {
    /// The name of the wallet used for signing.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The raw transaction data to be signed.
    #[prost(string, tag="2")]
    pub raw_transaction: ::prost::alloc::string::String,
    /// The password for unlocking the wallet for signing.
    #[prost(string, tag="3")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the transaction ID and signed raw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignRawTransactionResponse {
    /// The ID of the signed transaction.
    #[prost(string, tag="1")]
    pub transaction_id: ::prost::alloc::string::String,
    /// The signed raw transaction data.
    #[prost(string, tag="2")]
    pub signed_raw_transaction: ::prost::alloc::string::String,
}
/// Request message for obtaining the available balance of a wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalBalanceRequest {
    /// The name of the wallet to get the total balance.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the available balance of the wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalBalanceResponse {
    /// The name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The total balance of the wallet in NanoPAC.
    #[prost(int64, tag="2")]
    pub total_balance: i64,
}
/// Request message to sign an arbitrary message.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignMessageRequest {
    /// The name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The password for unlocking the wallet for signing.
    #[prost(string, tag="2")]
    pub password: ::prost::alloc::string::String,
    /// The account address associated with the private key.
    #[prost(string, tag="3")]
    pub address: ::prost::alloc::string::String,
    /// The arbitrary message to be signed.
    #[prost(string, tag="4")]
    pub message: ::prost::alloc::string::String,
}
/// Response message containing the available balance of the wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignMessageResponse {
    /// Signature of the message.
    #[prost(string, tag="1")]
    pub signature: ::prost::alloc::string::String,
}
/// Request message for get total of stake.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalStakeRequest {
    /// The name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response return total stake in wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalStakeResponse {
    #[prost(int64, tag="1")]
    pub total_stake: i64,
    #[prost(string, tag="2")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request get address information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressInfoRequest {
    /// The name of the wallet to generate a new address.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
}
/// Response return address information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressInfoResponse {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub label: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub public_key: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub path: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request update label an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SetLabelRequest {
    /// The name of the wallet to generate a new address.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The password for unlocking the wallet for signing.
    #[prost(string, tag="3")]
    pub password: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub label: ::prost::alloc::string::String,
}
/// Response return empty response.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SetLabelResponse {
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListWalletRequest {
}
/// Response return list wallet name.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListWalletResponse {
    #[prost(string, repeated, tag="1")]
    pub wallets: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// Request get wallet information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetWalletInfoRequest {
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response return wallet information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetWalletInfoResponse {
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    #[prost(int64, tag="2")]
    pub version: i64,
    #[prost(string, tag="3")]
    pub network: ::prost::alloc::string::String,
    #[prost(bool, tag="4")]
    pub encrypted: bool,
    #[prost(string, tag="5")]
    pub uuid: ::prost::alloc::string::String,
    #[prost(int64, tag="6")]
    pub created_at: i64,
}
/// Request get list of addresses in wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListAddressRequest {
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response return list addresses.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListAddressResponse {
    #[prost(message, repeated, tag="1")]
    pub data: ::prost::alloc::vec::Vec<AddressInfo>,
}
/// Enum for the address type.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum AddressType {
    /// Treasury address type.
    /// Should not be used to generate new addresses.
    Treasury = 0,
    /// Validator address type.
    Validator = 1,
    /// Account address type with BLS signature scheme.
    BlsAccount = 2,
    /// Account address type with Ed25519 signature scheme.
    /// Note: Generating a new Ed25519 address requires the wallet password.
    Ed25519Account = 3,
}
impl AddressType {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            AddressType::Treasury => "ADDRESS_TYPE_TREASURY",
            AddressType::Validator => "ADDRESS_TYPE_VALIDATOR",
            AddressType::BlsAccount => "ADDRESS_TYPE_BLS_ACCOUNT",
            AddressType::Ed25519Account => "ADDRESS_TYPE_ED25519_ACCOUNT",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "ADDRESS_TYPE_TREASURY" => Some(Self::Treasury),
            "ADDRESS_TYPE_VALIDATOR" => Some(Self::Validator),
            "ADDRESS_TYPE_BLS_ACCOUNT" => Some(Self::BlsAccount),
            "ADDRESS_TYPE_ED25519_ACCOUNT" => Some(Self::Ed25519Account),
            _ => None,
        }
    }
}
include!("pactus.serde.rs");
include!("pactus.tonic.rs");
// @@protoc_insertion_point(module)