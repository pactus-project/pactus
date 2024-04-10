// @generated
/// Request message for retrieving transaction details.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTransactionRequest {
    /// Transaction ID.
    #[prost(bytes="vec", tag="1")]
    pub id: ::prost::alloc::vec::Vec<u8>,
    /// Verbosity level for transaction details.
    #[prost(enumeration="TransactionVerbosity", tag="2")]
    pub verbosity: i32,
}
/// Response message containing details of a transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTransactionResponse {
    /// Height of the block containing the transaction.
    #[prost(uint32, tag="1")]
    pub block_height: u32,
    /// Time of the block containing the transaction.
    #[prost(uint32, tag="2")]
    pub block_time: u32,
    /// Information about the transaction.
    #[prost(message, optional, tag="3")]
    pub transaction: ::core::option::Option<TransactionInfo>,
}
/// Request message for calculating transaction fee.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CalculateFeeRequest {
    /// Transaction amount in NanoPAC.
    #[prost(int64, tag="1")]
    pub amount: i64,
    /// Type of transaction payload.
    #[prost(enumeration="PayloadType", tag="2")]
    pub payload_type: i32,
    /// Indicates that amount should be fixed and includes the fee.
    #[prost(bool, tag="3")]
    pub fixed_amount: bool,
}
/// Response message containing the calculated transaction fee.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CalculateFeeResponse {
    /// Calculated amount in NanoPAC.
    #[prost(int64, tag="1")]
    pub amount: i64,
    /// Calculated transaction fee in NanoPAC.
    #[prost(int64, tag="2")]
    pub fee: i64,
}
/// Request message for broadcasting a signed transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BroadcastTransactionRequest {
    /// Signed raw transaction data.
    #[prost(bytes="vec", tag="1")]
    pub signed_raw_transaction: ::prost::alloc::vec::Vec<u8>,
}
/// Response message containing the ID of the broadcasted transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BroadcastTransactionResponse {
    /// Transaction ID.
    #[prost(bytes="vec", tag="1")]
    pub id: ::prost::alloc::vec::Vec<u8>,
}
/// Request message for retrieving raw details of a transfer transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawTransferTransactionRequest {
    /// Lock time for the transaction.
    /// If not explicitly set, it sets to the last block height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// Sender's account address.
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    /// Receiver's account address.
    #[prost(string, tag="3")]
    pub receiver: ::prost::alloc::string::String,
    /// Transfer amount in NanoPAC.
    /// It should be greater than 0.
    #[prost(int64, tag="4")]
    pub amount: i64,
    /// Transaction fee in NanoPAC.
    /// If not explicitly set, it is calculated based on the amount.
    #[prost(int64, tag="5")]
    pub fee: i64,
    /// Transaction memo.
    #[prost(string, tag="6")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of a bond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawBondTransactionRequest {
    /// Lock time for the transaction.
    /// If not explicitly set, it sets to the last block height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// Sender's account address.
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    /// Receiver's validator address.
    #[prost(string, tag="3")]
    pub receiver: ::prost::alloc::string::String,
    /// Stake amount in NanoPAC.
    /// It should be greater than 0.
    #[prost(int64, tag="4")]
    pub stake: i64,
    /// Public key of the validator.
    #[prost(string, tag="5")]
    pub public_key: ::prost::alloc::string::String,
    /// Transaction fee in NanoPAC.
    /// If not explicitly set, it is calculated based on the stake.
    #[prost(int64, tag="6")]
    pub fee: i64,
    /// Transaction memo.
    #[prost(string, tag="7")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of an unbond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawUnbondTransactionRequest {
    /// Lock time for the transaction.
    /// If not explicitly set, it sets to the last block height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// Address of the validator to unbond from.
    #[prost(string, tag="3")]
    pub validator_address: ::prost::alloc::string::String,
    /// Transaction memo.
    #[prost(string, tag="4")]
    pub memo: ::prost::alloc::string::String,
}
/// Request message for retrieving raw details of a withdraw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawWithdrawTransactionRequest {
    /// Lock time for the transaction.
    /// If not explicitly set, it sets to the last block height.
    #[prost(uint32, tag="1")]
    pub lock_time: u32,
    /// Address of the validator to withdraw from.
    #[prost(string, tag="2")]
    pub validator_address: ::prost::alloc::string::String,
    /// Address of the account to withdraw to.
    #[prost(string, tag="3")]
    pub account_address: ::prost::alloc::string::String,
    /// Withdrawal amount in NanoPAC.
    /// It should be greater than 0.
    #[prost(int64, tag="4")]
    pub amount: i64,
    /// Transaction fee in NanoPAC.
    /// If not explicitly set, it is calculated based on the amount.
    #[prost(int64, tag="5")]
    pub fee: i64,
    /// Transaction memo.
    #[prost(string, tag="6")]
    pub memo: ::prost::alloc::string::String,
}
/// Response message containing raw transaction data.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetRawTransactionResponse {
    /// Raw transaction data.
    #[prost(bytes="vec", tag="1")]
    pub raw_transaction: ::prost::alloc::vec::Vec<u8>,
}
/// Payload for a transfer transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadTransfer {
    /// Sender's address.
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    /// Receiver's address.
    #[prost(string, tag="2")]
    pub receiver: ::prost::alloc::string::String,
    /// Transaction amount in NanoPAC.
    #[prost(int64, tag="3")]
    pub amount: i64,
}
/// Payload for a bond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadBond {
    /// Sender's address.
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    /// Receiver's address.
    #[prost(string, tag="2")]
    pub receiver: ::prost::alloc::string::String,
    /// Stake amount in NanoPAC.
    #[prost(int64, tag="3")]
    pub stake: i64,
}
/// Payload for a sortition transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadSortition {
    /// Address associated with the sortition.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    /// Proof for the sortition.
    #[prost(bytes="vec", tag="2")]
    pub proof: ::prost::alloc::vec::Vec<u8>,
}
/// Payload for an unbond transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadUnbond {
    /// Address of the validator to unbond from.
    #[prost(string, tag="1")]
    pub validator: ::prost::alloc::string::String,
}
/// Payload for a withdraw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PayloadWithdraw {
    /// Address to withdraw from.
    #[prost(string, tag="1")]
    pub from: ::prost::alloc::string::String,
    /// Address to withdraw to.
    #[prost(string, tag="2")]
    pub to: ::prost::alloc::string::String,
    /// Withdrawal amount in NanoPAC.
    #[prost(int64, tag="3")]
    pub amount: i64,
}
/// Information about a transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TransactionInfo {
    /// Transaction ID.
    #[prost(bytes="vec", tag="1")]
    pub id: ::prost::alloc::vec::Vec<u8>,
    /// Transaction data.
    #[prost(bytes="vec", tag="2")]
    pub data: ::prost::alloc::vec::Vec<u8>,
    /// Transaction version.
    #[prost(int32, tag="3")]
    pub version: i32,
    /// Lock time for the transaction.
    #[prost(uint32, tag="4")]
    pub lock_time: u32,
    /// Transaction value in NanoPAC.
    #[prost(int64, tag="5")]
    pub value: i64,
    /// Transaction fee in NanoPAC.
    #[prost(int64, tag="6")]
    pub fee: i64,
    /// Type of transaction payload.
    #[prost(enumeration="PayloadType", tag="7")]
    pub payload_type: i32,
    /// Transaction memo.
    #[prost(string, tag="8")]
    pub memo: ::prost::alloc::string::String,
    /// Public key associated with the transaction.
    #[prost(string, tag="9")]
    pub public_key: ::prost::alloc::string::String,
    /// Transaction signature.
    #[prost(bytes="vec", tag="10")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(oneof="transaction_info::Payload", tags="30, 31, 32, 33, 34")]
    pub payload: ::core::option::Option<transaction_info::Payload>,
}
/// Nested message and enum types in `TransactionInfo`.
pub mod transaction_info {
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
    pub enum Payload {
        /// Transfer payload.
        #[prost(message, tag="30")]
        Transfer(super::PayloadTransfer),
        /// Bond payload.
        #[prost(message, tag="31")]
        Bond(super::PayloadBond),
        /// Sortition payload.
        #[prost(message, tag="32")]
        Sortition(super::PayloadSortition),
        /// Unbond payload.
        #[prost(message, tag="33")]
        Unbond(super::PayloadUnbond),
        /// Withdraw payload.
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
/// Enumeration for verbosity level when requesting transaction details.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum TransactionVerbosity {
    /// Request transaction data only.
    TransactionData = 0,
    /// Request transaction details.
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
    /// Address of the account.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing the response with account information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAccountResponse {
    /// Account information.
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
    /// Address of the validator.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message to request validator information based on a validator number.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorByNumberRequest {
    /// Validator number.
    #[prost(int32, tag="1")]
    pub number: i32,
}
/// Message containing the response with validator information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorResponse {
    /// Validator information.
    #[prost(message, optional, tag="1")]
    pub validator: ::core::option::Option<ValidatorInfo>,
}
/// Message to request public key based on an address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetPublicKeyRequest {
    /// Address for which public key is requested.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing the response with the public key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetPublicKeyResponse {
    /// Public key of the account.
    #[prost(string, tag="1")]
    pub public_key: ::prost::alloc::string::String,
}
/// Message to request block information based on height and verbosity.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockRequest {
    /// Height of the block.
    #[prost(uint32, tag="1")]
    pub height: u32,
    /// Verbosity level for block information.
    #[prost(enumeration="BlockVerbosity", tag="2")]
    pub verbosity: i32,
}
/// Message containing the response with block information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockResponse {
    /// Height of the block.
    #[prost(uint32, tag="1")]
    pub height: u32,
    /// Hash of the block.
    #[prost(bytes="vec", tag="2")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
    /// Block data, only available if the verbosity level is set to BLOCK_DATA.
    #[prost(bytes="vec", tag="3")]
    pub data: ::prost::alloc::vec::Vec<u8>,
    /// Block timestamp.
    #[prost(uint32, tag="4")]
    pub block_time: u32,
    /// Block header information.
    #[prost(message, optional, tag="5")]
    pub header: ::core::option::Option<BlockHeaderInfo>,
    /// Certificate information of the previous block.
    #[prost(message, optional, tag="6")]
    pub prev_cert: ::core::option::Option<CertificateInfo>,
    /// List of transactions in the block.
    /// Transaction information is available when the verbosity level is set to BLOCK_TRANSACTIONS.
    #[prost(message, repeated, tag="7")]
    pub txs: ::prost::alloc::vec::Vec<TransactionInfo>,
}
/// Message to request block hash based on height.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHashRequest {
    /// Height of the block.
    #[prost(uint32, tag="1")]
    pub height: u32,
}
/// Message containing the response with the block hash.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHashResponse {
    /// Hash of the block.
    #[prost(bytes="vec", tag="1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
}
/// Message to request block height based on hash.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHeightRequest {
    /// Hash of the block.
    #[prost(bytes="vec", tag="1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
}
/// Message containing the response with the block height.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBlockHeightResponse {
    /// Height of the block.
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
    /// Height of the last block.
    #[prost(uint32, tag="1")]
    pub last_block_height: u32,
    /// Hash of the last block.
    #[prost(bytes="vec", tag="2")]
    pub last_block_hash: ::prost::alloc::vec::Vec<u8>,
    /// Total number of accounts.
    #[prost(int32, tag="3")]
    pub total_accounts: i32,
    /// Total number of validators.
    #[prost(int32, tag="4")]
    pub total_validators: i32,
    /// Total power in the blockchain.
    #[prost(int64, tag="5")]
    pub total_power: i64,
    /// Power of the committee.
    #[prost(int64, tag="6")]
    pub committee_power: i64,
    /// List of committee validators.
    #[prost(message, repeated, tag="7")]
    pub committee_validators: ::prost::alloc::vec::Vec<ValidatorInfo>,
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
    /// List of consensus instances.
    #[prost(message, repeated, tag="1")]
    pub instances: ::prost::alloc::vec::Vec<ConsensusInfo>,
}
/// Message containing information about a validator.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ValidatorInfo {
    /// Hash of the validator.
    #[prost(bytes="vec", tag="1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
    /// Validator data.
    #[prost(bytes="vec", tag="2")]
    pub data: ::prost::alloc::vec::Vec<u8>,
    /// Public key of the validator.
    #[prost(string, tag="3")]
    pub public_key: ::prost::alloc::string::String,
    /// Validator number.
    #[prost(int32, tag="4")]
    pub number: i32,
    /// Validator stake in NanoPAC.
    #[prost(int64, tag="5")]
    pub stake: i64,
    /// Last bonding height.
    #[prost(uint32, tag="6")]
    pub last_bonding_height: u32,
    /// Last sortition height.
    #[prost(uint32, tag="7")]
    pub last_sortition_height: u32,
    /// Unbonding height.
    #[prost(uint32, tag="8")]
    pub unbonding_height: u32,
    /// Address of the validator.
    #[prost(string, tag="9")]
    pub address: ::prost::alloc::string::String,
    /// Availability score of the validator.
    #[prost(double, tag="10")]
    pub availability_score: f64,
}
/// Message containing information about an account.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AccountInfo {
    /// Hash of the account.
    #[prost(bytes="vec", tag="1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
    /// Account data.
    #[prost(bytes="vec", tag="2")]
    pub data: ::prost::alloc::vec::Vec<u8>,
    /// Account number.
    #[prost(int32, tag="3")]
    pub number: i32,
    /// Account balance in NanoPAC.
    #[prost(int64, tag="4")]
    pub balance: i64,
    /// Address of the account.
    #[prost(string, tag="5")]
    pub address: ::prost::alloc::string::String,
}
/// Message containing information about the header of a block.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlockHeaderInfo {
    /// Block version.
    #[prost(int32, tag="1")]
    pub version: i32,
    /// Hash of the previous block.
    #[prost(bytes="vec", tag="2")]
    pub prev_block_hash: ::prost::alloc::vec::Vec<u8>,
    /// State root of the block.
    #[prost(bytes="vec", tag="3")]
    pub state_root: ::prost::alloc::vec::Vec<u8>,
    /// Sortition seed of the block.
    #[prost(bytes="vec", tag="4")]
    pub sortition_seed: ::prost::alloc::vec::Vec<u8>,
    /// Address of the proposer of the block.
    #[prost(string, tag="5")]
    pub proposer_address: ::prost::alloc::string::String,
}
/// Message containing information about a certificate.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CertificateInfo {
    /// Hash of the certificate.
    #[prost(bytes="vec", tag="1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
    /// Round of the certificate.
    #[prost(int32, tag="2")]
    pub round: i32,
    /// List of committers in the certificate.
    #[prost(int32, repeated, tag="3")]
    pub committers: ::prost::alloc::vec::Vec<i32>,
    /// List of absentees in the certificate.
    #[prost(int32, repeated, tag="4")]
    pub absentees: ::prost::alloc::vec::Vec<i32>,
    /// Certificate signature.
    #[prost(bytes="vec", tag="5")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// Message containing information about a vote.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct VoteInfo {
    /// Type of the vote.
    #[prost(enumeration="VoteType", tag="1")]
    pub r#type: i32,
    /// Voter's address.
    #[prost(string, tag="2")]
    pub voter: ::prost::alloc::string::String,
    /// Hash of the block being voted on.
    #[prost(bytes="vec", tag="3")]
    pub block_hash: ::prost::alloc::vec::Vec<u8>,
    /// Round of the vote.
    #[prost(int32, tag="4")]
    pub round: i32,
    /// Consensus round of the vote.
    #[prost(int32, tag="5")]
    pub cp_round: i32,
    /// Consensus value of the vote.
    #[prost(int32, tag="6")]
    pub cp_value: i32,
}
/// Message containing information about consensus.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConsensusInfo {
    /// Address of the consensus instance.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    /// Whether the consensus instance is active.
    #[prost(bool, tag="2")]
    pub active: bool,
    /// Height of the consensus instance.
    #[prost(uint32, tag="3")]
    pub height: u32,
    /// Round of the consensus instance.
    #[prost(int32, tag="4")]
    pub round: i32,
    /// List of votes in the consensus instance.
    #[prost(message, repeated, tag="5")]
    pub votes: ::prost::alloc::vec::Vec<VoteInfo>,
}
/// Enumeration for verbosity level when requesting block information.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum BlockVerbosity {
    /// Request block data only.
    BlockData = 0,
    /// Request block information and transaction IDs.
    BlockInfo = 1,
    /// Request block information and transaction details.
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
}
/// Response message containing information about the overall network.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNetworkInfoResponse {
    /// Name of the network.
    #[prost(string, tag="1")]
    pub network_name: ::prost::alloc::string::String,
    /// Total bytes sent across the network.
    #[prost(uint32, tag="2")]
    pub total_sent_bytes: u32,
    /// Total bytes received across the network.
    #[prost(uint32, tag="3")]
    pub total_received_bytes: u32,
    /// Number of connected peers.
    #[prost(uint32, tag="4")]
    pub connected_peers_count: u32,
    /// List of connected peers.
    #[prost(message, repeated, tag="5")]
    pub connected_peers: ::prost::alloc::vec::Vec<PeerInfo>,
    /// Bytes sent per peer ID.
    #[prost(map="uint32, uint64", tag="6")]
    pub sent_bytes: ::std::collections::HashMap<u32, u64>,
    /// Bytes received per peer ID.
    #[prost(map="uint32, uint64", tag="7")]
    pub received_bytes: ::std::collections::HashMap<u32, u64>,
}
/// Request message for retrieving information about a specific node in the
/// network.
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
    /// Agent information of the node.
    #[prost(string, tag="2")]
    pub agent: ::prost::alloc::string::String,
    /// Peer ID of the node.
    #[prost(bytes="vec", tag="3")]
    pub peer_id: ::prost::alloc::vec::Vec<u8>,
    /// Timestamp when the node started.
    #[prost(uint64, tag="4")]
    pub started_at: u64,
    /// Reachability status of the node.
    #[prost(string, tag="5")]
    pub reachability: ::prost::alloc::string::String,
    /// List of services provided by the node.
    #[prost(int32, repeated, tag="6")]
    pub services: ::prost::alloc::vec::Vec<i32>,
    /// Names of services provided by the node.
    #[prost(string, repeated, tag="7")]
    pub services_names: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// List of addresses associated with the node.
    #[prost(string, repeated, tag="8")]
    pub addrs: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// List of protocols supported by the node.
    #[prost(string, repeated, tag="9")]
    pub protocols: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// Information about a peer in the network.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PeerInfo {
    /// Status of the peer.
    #[prost(int32, tag="1")]
    pub status: i32,
    /// Moniker of the peer.
    #[prost(string, tag="2")]
    pub moniker: ::prost::alloc::string::String,
    /// Agent information of the peer.
    #[prost(string, tag="3")]
    pub agent: ::prost::alloc::string::String,
    /// Peer ID of the peer.
    #[prost(bytes="vec", tag="4")]
    pub peer_id: ::prost::alloc::vec::Vec<u8>,
    /// Consensus keys used by the peer.
    #[prost(string, repeated, tag="5")]
    pub consensus_keys: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Consensus address of the peer.
    #[prost(string, repeated, tag="6")]
    pub consensus_address: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Services provided by the peer.
    #[prost(uint32, tag="7")]
    pub services: u32,
    /// Hash of the last block the peer knows.
    #[prost(bytes="vec", tag="8")]
    pub last_block_hash: ::prost::alloc::vec::Vec<u8>,
    /// Height of the peer in the blockchain.
    #[prost(uint32, tag="9")]
    pub height: u32,
    /// Count of received messages.
    #[prost(int32, tag="10")]
    pub received_messages: i32,
    /// Count of invalid messages received.
    #[prost(int32, tag="11")]
    pub invalid_messages: i32,
    /// Timestamp of the last sent message.
    #[prost(int64, tag="12")]
    pub last_sent: i64,
    /// Timestamp of the last received message.
    #[prost(int64, tag="13")]
    pub last_received: i64,
    /// Bytes sent per message type.
    #[prost(map="int32, int64", tag="14")]
    pub sent_bytes: ::std::collections::HashMap<i32, i64>,
    /// Bytes received per message type.
    #[prost(map="int32, int64", tag="15")]
    pub received_bytes: ::std::collections::HashMap<i32, i64>,
    /// Network address of the peer.
    #[prost(string, tag="16")]
    pub address: ::prost::alloc::string::String,
    /// Direction of connection with the peer.
    #[prost(string, tag="17")]
    pub direction: ::prost::alloc::string::String,
    /// List of protocols supported by the peer.
    #[prost(string, repeated, tag="18")]
    pub protocols: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// Total sessions with the peer.
    #[prost(int32, tag="19")]
    pub total_sessions: i32,
    /// Completed sessions with the peer.
    #[prost(int32, tag="20")]
    pub completed_sessions: i32,
}
/// Message of address information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddressInfo {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub public_key: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub label: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub path: ::prost::alloc::string::String,
}
/// Message of address history information.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct HistoryInfo {
    /// Hash of transaction.
    #[prost(string, tag="1")]
    pub transaction_id: ::prost::alloc::string::String,
    /// transaction timestamp.
    #[prost(uint32, tag="2")]
    pub time: u32,
    /// payload type of transaction.
    #[prost(string, tag="3")]
    pub payload_type: ::prost::alloc::string::String,
    /// description of transaction.
    #[prost(string, tag="4")]
    pub description: ::prost::alloc::string::String,
    /// amount of transaction.
    #[prost(int64, tag="5")]
    pub amount: i64,
}
/// Request message to get an address transaction history.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressHistoryRequest {
    /// Name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Address to get the transaction history of it.
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
}
/// Response message to get an address transaction history.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAddressHistoryResponse {
    /// Array of address history and activities.
    #[prost(message, repeated, tag="1")]
    pub history_info: ::prost::alloc::vec::Vec<HistoryInfo>,
}
/// Request message for generating a new address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNewAddressRequest {
    /// Name of the wallet for which the new address is requested.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Address type for the new address.
    #[prost(enumeration="AddressType", tag="2")]
    pub address_type: i32,
    /// Label for the new address.
    #[prost(string, tag="3")]
    pub label: ::prost::alloc::string::String,
}
/// Response message containing the new address.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetNewAddressResponse {
    /// Name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Address information.
    #[prost(message, optional, tag="2")]
    pub address_info: ::core::option::Option<AddressInfo>,
}
/// Request message for creating a new wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateWalletRequest {
    /// Name of the new wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Mnemonic for wallet recovery.
    #[prost(string, tag="2")]
    pub mnemonic: ::prost::alloc::string::String,
    /// Language for the mnemonic.
    #[prost(string, tag="3")]
    pub language: ::prost::alloc::string::String,
    /// Password for securing the wallet.
    #[prost(string, tag="4")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the name of the created wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateWalletResponse {
    /// Name of the created wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for loading an existing wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LoadWalletRequest {
    /// Name of the wallet to load.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the name of the loaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LoadWalletResponse {
    /// Name of the loaded wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for unloading a currently loaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnloadWalletRequest {
    /// Name of the wallet to unload.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the name of the unloaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnloadWalletResponse {
    /// Name of the unloaded wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for locking a currently loaded wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LockWalletRequest {
    /// Name of the wallet to lock.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the name of the locked wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LockWalletResponse {
    /// Name of the locked wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for obtaining the validator address associated with a public
/// key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressRequest {
    /// Public key for which the validator address is requested.
    #[prost(string, tag="1")]
    pub public_key: ::prost::alloc::string::String,
}
/// Response message containing the validator address corresponding to a public
/// key.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetValidatorAddressResponse {
    /// Validator address associated with the public key.
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
/// Request message for unlocking a wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnlockWalletRequest {
    /// Name of the wallet to unlock.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Password for unlocking the wallet.
    #[prost(string, tag="2")]
    pub password: ::prost::alloc::string::String,
    /// Timeout duration for the unlocked state.
    #[prost(int32, tag="3")]
    pub timeout: i32,
}
/// Response message containing the name of the unlocked wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnlockWalletResponse {
    /// Name of the unlocked wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Request message for signing a raw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignRawTransactionRequest {
    /// Name of the wallet used for signing.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// Raw transaction data to be signed.
    #[prost(bytes="vec", tag="2")]
    pub raw_transaction: ::prost::alloc::vec::Vec<u8>,
    /// Password for unlocking the wallet for signing.
    #[prost(string, tag="3")]
    pub password: ::prost::alloc::string::String,
}
/// Response message containing the transaction ID and signed raw transaction.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignRawTransactionResponse {
    /// ID of the signed transaction.
    #[prost(bytes="vec", tag="1")]
    pub transaction_id: ::prost::alloc::vec::Vec<u8>,
    /// Signed raw transaction data.
    #[prost(bytes="vec", tag="2")]
    pub signed_raw_transaction: ::prost::alloc::vec::Vec<u8>,
}
/// Request message for obtaining the available balance of a wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalBalanceRequest {
    /// Name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
}
/// Response message containing the available balance of the wallet.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTotalBalanceResponse {
    /// Name of the wallet.
    #[prost(string, tag="1")]
    pub wallet_name: ::prost::alloc::string::String,
    /// The total balance of the wallet in NanoPAC.
    #[prost(int64, tag="2")]
    pub total_balance: i64,
}
/// Enum for the address type.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum AddressType {
    Treasury = 0,
    Validator = 1,
    BlsAccount = 2,
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
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "ADDRESS_TYPE_TREASURY" => Some(Self::Treasury),
            "ADDRESS_TYPE_VALIDATOR" => Some(Self::Validator),
            "ADDRESS_TYPE_BLS_ACCOUNT" => Some(Self::BlsAccount),
            _ => None,
        }
    }
}
include!("pactus.serde.rs");
include!("pactus.tonic.rs");
// @@protoc_insertion_point(module)