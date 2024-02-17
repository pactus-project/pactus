# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [transaction.proto](#transaction-proto)
    - [BroadcastTransactionRequest](#pactus-BroadcastTransactionRequest)
    - [BroadcastTransactionResponse](#pactus-BroadcastTransactionResponse)
    - [CalculateFeeRequest](#pactus-CalculateFeeRequest)
    - [CalculateFeeResponse](#pactus-CalculateFeeResponse)
    - [GetRawBondTransactionRequest](#pactus-GetRawBondTransactionRequest)
    - [GetRawTransactionResponse](#pactus-GetRawTransactionResponse)
    - [GetRawTransferTransactionRequest](#pactus-GetRawTransferTransactionRequest)
    - [GetRawUnBondTransactionRequest](#pactus-GetRawUnBondTransactionRequest)
    - [GetRawWithdrawTransactionRequest](#pactus-GetRawWithdrawTransactionRequest)
    - [GetTransactionRequest](#pactus-GetTransactionRequest)
    - [GetTransactionResponse](#pactus-GetTransactionResponse)
    - [PayloadBond](#pactus-PayloadBond)
    - [PayloadSortition](#pactus-PayloadSortition)
    - [PayloadTransfer](#pactus-PayloadTransfer)
    - [PayloadUnbond](#pactus-PayloadUnbond)
    - [PayloadWithdraw](#pactus-PayloadWithdraw)
    - [TransactionInfo](#pactus-TransactionInfo)
  
    - [PayloadType](#pactus-PayloadType)
    - [TransactionVerbosity](#pactus-TransactionVerbosity)
  
    - [Transaction](#pactus-Transaction)
  
- [blockchain.proto](#blockchain-proto)
    - [AccountInfo](#pactus-AccountInfo)
    - [BlockHeaderInfo](#pactus-BlockHeaderInfo)
    - [CertificateInfo](#pactus-CertificateInfo)
    - [ConsensusInfo](#pactus-ConsensusInfo)
    - [GetAccountRequest](#pactus-GetAccountRequest)
    - [GetAccountResponse](#pactus-GetAccountResponse)
    - [GetBlockHashRequest](#pactus-GetBlockHashRequest)
    - [GetBlockHashResponse](#pactus-GetBlockHashResponse)
    - [GetBlockHeightRequest](#pactus-GetBlockHeightRequest)
    - [GetBlockHeightResponse](#pactus-GetBlockHeightResponse)
    - [GetBlockRequest](#pactus-GetBlockRequest)
    - [GetBlockResponse](#pactus-GetBlockResponse)
    - [GetBlockchainInfoRequest](#pactus-GetBlockchainInfoRequest)
    - [GetBlockchainInfoResponse](#pactus-GetBlockchainInfoResponse)
    - [GetConsensusInfoRequest](#pactus-GetConsensusInfoRequest)
    - [GetConsensusInfoResponse](#pactus-GetConsensusInfoResponse)
    - [GetPublicKeyRequest](#pactus-GetPublicKeyRequest)
    - [GetPublicKeyResponse](#pactus-GetPublicKeyResponse)
    - [GetValidatorAddressesRequest](#pactus-GetValidatorAddressesRequest)
    - [GetValidatorAddressesResponse](#pactus-GetValidatorAddressesResponse)
    - [GetValidatorByNumberRequest](#pactus-GetValidatorByNumberRequest)
    - [GetValidatorRequest](#pactus-GetValidatorRequest)
    - [GetValidatorResponse](#pactus-GetValidatorResponse)
    - [ValidatorInfo](#pactus-ValidatorInfo)
    - [VoteInfo](#pactus-VoteInfo)
  
    - [BlockVerbosity](#pactus-BlockVerbosity)
    - [VoteType](#pactus-VoteType)
  
    - [Blockchain](#pactus-Blockchain)
  
- [network.proto](#network-proto)
    - [GetNetworkInfoRequest](#pactus-GetNetworkInfoRequest)
    - [GetNetworkInfoResponse](#pactus-GetNetworkInfoResponse)
    - [GetNetworkInfoResponse.ReceivedBytesEntry](#pactus-GetNetworkInfoResponse-ReceivedBytesEntry)
    - [GetNetworkInfoResponse.SentBytesEntry](#pactus-GetNetworkInfoResponse-SentBytesEntry)
    - [GetNodeInfoRequest](#pactus-GetNodeInfoRequest)
    - [GetNodeInfoResponse](#pactus-GetNodeInfoResponse)
    - [PeerInfo](#pactus-PeerInfo)
    - [PeerInfo.ReceivedBytesEntry](#pactus-PeerInfo-ReceivedBytesEntry)
    - [PeerInfo.SentBytesEntry](#pactus-PeerInfo-SentBytesEntry)
  
    - [Network](#pactus-Network)
  
- [wallet.proto](#wallet-proto)
    - [CreateWalletRequest](#pactus-CreateWalletRequest)
    - [CreateWalletResponse](#pactus-CreateWalletResponse)
    - [GetValidatorAddressRequest](#pactus-GetValidatorAddressRequest)
    - [GetValidatorAddressResponse](#pactus-GetValidatorAddressResponse)
    - [LoadWalletRequest](#pactus-LoadWalletRequest)
    - [LoadWalletResponse](#pactus-LoadWalletResponse)
    - [LockWalletRequest](#pactus-LockWalletRequest)
    - [LockWalletResponse](#pactus-LockWalletResponse)
    - [SignRawTransactionRequest](#pactus-SignRawTransactionRequest)
    - [SignRawTransactionResponse](#pactus-SignRawTransactionResponse)
    - [UnloadWalletRequest](#pactus-UnloadWalletRequest)
    - [UnloadWalletResponse](#pactus-UnloadWalletResponse)
    - [UnlockWalletRequest](#pactus-UnlockWalletRequest)
    - [UnlockWalletResponse](#pactus-UnlockWalletResponse)
  
    - [Wallet](#pactus-Wallet)
  
- [Scalar Value Types](#scalar-value-types)



<a name="transaction-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## transaction.proto



<a name="pactus-BroadcastTransactionRequest"></a>

### BroadcastTransactionRequest
Request message for broadcasting a signed transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signed_raw_transaction | [bytes](#bytes) |  | Signed raw transaction data. |






<a name="pactus-BroadcastTransactionResponse"></a>

### BroadcastTransactionResponse
Response message containing the ID of the broadcasted transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [bytes](#bytes) |  | Transaction ID. |






<a name="pactus-CalculateFeeRequest"></a>

### CalculateFeeRequest
Request message for calculating transaction fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [int64](#int64) |  | Transaction amount. |
| payload_type | [PayloadType](#pactus-PayloadType) |  | Type of transaction payload. |






<a name="pactus-CalculateFeeResponse"></a>

### CalculateFeeResponse
Response message containing the calculated transaction fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fee | [int64](#int64) |  | Calculated transaction fee. |






<a name="pactus-GetRawBondTransactionRequest"></a>

### GetRawBondTransactionRequest
Request message for retrieving raw details of a bond transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lock_time | [uint32](#uint32) |  | Lock time for the transaction. |
| sender | [string](#string) |  | Sender&#39;s address. |
| receiver | [string](#string) |  | Receiver&#39;s address. |
| stake | [int64](#int64) |  | Stake amount. |
| public_key | [string](#string) |  | Public key of the validator. |
| fee | [int64](#int64) |  | Transaction fee. |
| memo | [string](#string) |  | Transaction memo. |






<a name="pactus-GetRawTransactionResponse"></a>

### GetRawTransactionResponse
Response message containing raw transaction data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| raw_transaction | [bytes](#bytes) |  | Raw transaction data. |






<a name="pactus-GetRawTransferTransactionRequest"></a>

### GetRawTransferTransactionRequest
Request message for retrieving raw details of a transfer transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lock_time | [uint32](#uint32) |  | Lock time for the transaction. |
| sender | [string](#string) |  | Sender&#39;s address. |
| receiver | [string](#string) |  | Receiver&#39;s address. |
| amount | [int64](#int64) |  | Transaction amount. |
| fee | [int64](#int64) |  | Transaction fee. |
| memo | [string](#string) |  | Transaction memo. |






<a name="pactus-GetRawUnBondTransactionRequest"></a>

### GetRawUnBondTransactionRequest
Request message for retrieving raw details of an unbond transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lock_time | [uint32](#uint32) |  | Lock time for the transaction. |
| validator_address | [string](#string) |  | Address of the validator to unbond from. |
| memo | [string](#string) |  | Transaction memo. |






<a name="pactus-GetRawWithdrawTransactionRequest"></a>

### GetRawWithdrawTransactionRequest
Request message for retrieving raw details of a withdraw transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lock_time | [uint32](#uint32) |  | Lock time for the transaction. |
| validator_address | [string](#string) |  | Address of the validator to withdraw from. |
| account_address | [string](#string) |  | Address of the account to withdraw to. |
| fee | [int64](#int64) |  | Transaction fee. |
| amount | [int64](#int64) |  | Withdrawal amount. |
| memo | [string](#string) |  | Transaction memo. |






<a name="pactus-GetTransactionRequest"></a>

### GetTransactionRequest
Request message for retrieving transaction details.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [bytes](#bytes) |  | Transaction ID. |
| verbosity | [TransactionVerbosity](#pactus-TransactionVerbosity) |  | Verbosity level for transaction details. |






<a name="pactus-GetTransactionResponse"></a>

### GetTransactionResponse
Response message containing details of a transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_height | [uint32](#uint32) |  | Height of the block containing the transaction. |
| block_time | [uint32](#uint32) |  | Time of the block containing the transaction. |
| transaction | [TransactionInfo](#pactus-TransactionInfo) |  | Information about the transaction. |






<a name="pactus-PayloadBond"></a>

### PayloadBond
Payload for a bond transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender&#39;s address. |
| receiver | [string](#string) |  | Receiver&#39;s address. |
| stake | [int64](#int64) |  | Stake amount. |






<a name="pactus-PayloadSortition"></a>

### PayloadSortition
Payload for a sortition transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Address associated with the sortition. |
| proof | [bytes](#bytes) |  | Proof for the sortition. |






<a name="pactus-PayloadTransfer"></a>

### PayloadTransfer
Payload for a transfer transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender&#39;s address. |
| receiver | [string](#string) |  | Receiver&#39;s address. |
| amount | [int64](#int64) |  | Transaction amount. |






<a name="pactus-PayloadUnbond"></a>

### PayloadUnbond
Payload for an unbond transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| validator | [string](#string) |  | Address of the validator to unbond from. |






<a name="pactus-PayloadWithdraw"></a>

### PayloadWithdraw
Payload for a withdraw transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from | [string](#string) |  | Address to withdraw from. |
| to | [string](#string) |  | Address to withdraw to. |
| amount | [int64](#int64) |  | Withdrawal amount. |






<a name="pactus-TransactionInfo"></a>

### TransactionInfo
Information about a transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [bytes](#bytes) |  | Transaction ID. |
| data | [bytes](#bytes) |  | Transaction data. |
| version | [int32](#int32) |  | Transaction version. |
| lock_time | [uint32](#uint32) |  | Lock time for the transaction. |
| value | [int64](#int64) |  | Transaction value. |
| fee | [int64](#int64) |  | Transaction fee. |
| payload_type | [PayloadType](#pactus-PayloadType) |  | Type of transaction payload. |
| transfer | [PayloadTransfer](#pactus-PayloadTransfer) |  | Transfer payload. |
| bond | [PayloadBond](#pactus-PayloadBond) |  | Bond payload. |
| sortition | [PayloadSortition](#pactus-PayloadSortition) |  | Sortition payload. |
| unbond | [PayloadUnbond](#pactus-PayloadUnbond) |  | Unbond payload. |
| withdraw | [PayloadWithdraw](#pactus-PayloadWithdraw) |  | Withdraw payload. |
| memo | [string](#string) |  | Transaction memo. |
| public_key | [string](#string) |  | Public key associated with the transaction. |
| signature | [bytes](#bytes) |  | Transaction signature. |





 


<a name="pactus-PayloadType"></a>

### PayloadType
Enumeration for different types of transaction payloads.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | Unknown payload type. |
| TRANSFER_PAYLOAD | 1 | Transfer payload type. |
| BOND_PAYLOAD | 2 | Bond payload type. |
| SORTITION_PAYLOAD | 3 | Sortition payload type. |
| UNBOND_PAYLOAD | 4 | Unbond payload type. |
| WITHDRAW_PAYLOAD | 5 | Withdraw payload type. |



<a name="pactus-TransactionVerbosity"></a>

### TransactionVerbosity
Enumeration for verbosity level when requesting transaction details.

| Name | Number | Description |
| ---- | ------ | ----------- |
| TRANSACTION_DATA | 0 | Request only transaction data. |
| TRANSACTION_INFO | 1 | Request detailed transaction information. |


 

 


<a name="pactus-Transaction"></a>

### Transaction
Transaction service defines various RPC methods for interacting with
transactions.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetTransaction | [GetTransactionRequest](#pactus-GetTransactionRequest) | [GetTransactionResponse](#pactus-GetTransactionResponse) | GetTransaction retrieves transaction details based on the provided request parameters. |
| CalculateFee | [CalculateFeeRequest](#pactus-CalculateFeeRequest) | [CalculateFeeResponse](#pactus-CalculateFeeResponse) | CalculateFee calculates the transaction fee based on the specified amount and payload type. |
| BroadcastTransaction | [BroadcastTransactionRequest](#pactus-BroadcastTransactionRequest) | [BroadcastTransactionResponse](#pactus-BroadcastTransactionResponse) | BroadcastTransaction broadcasts a signed transaction to the network. |
| GetRawTransferTransaction | [GetRawTransferTransactionRequest](#pactus-GetRawTransferTransactionRequest) | [GetRawTransactionResponse](#pactus-GetRawTransactionResponse) | GetRawTransferTransaction retrieves raw details of a transfer transaction. |
| GetRawBondTransaction | [GetRawBondTransactionRequest](#pactus-GetRawBondTransactionRequest) | [GetRawTransactionResponse](#pactus-GetRawTransactionResponse) | GetRawBondTransaction retrieves raw details of a bond transaction. |
| GetRawUnBondTransaction | [GetRawUnBondTransactionRequest](#pactus-GetRawUnBondTransactionRequest) | [GetRawTransactionResponse](#pactus-GetRawTransactionResponse) | GetRawUnBondTransaction retrieves raw details of an unbond transaction. |
| GetRawWithdrawTransaction | [GetRawWithdrawTransactionRequest](#pactus-GetRawWithdrawTransactionRequest) | [GetRawTransactionResponse](#pactus-GetRawTransactionResponse) | GetRawWithdrawTransaction retrieves raw details of a withdraw transaction. |

 



<a name="blockchain-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## blockchain.proto



<a name="pactus-AccountInfo"></a>

### AccountInfo
Message containing information about an account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [bytes](#bytes) |  | Hash of the account. |
| data | [bytes](#bytes) |  | Account data. |
| number | [int32](#int32) |  | Account number. |
| balance | [int64](#int64) |  | Account balance. |
| address | [string](#string) |  | Address of the account. |






<a name="pactus-BlockHeaderInfo"></a>

### BlockHeaderInfo
Message containing information about the header of a block.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [int32](#int32) |  | Block version. |
| prev_block_hash | [bytes](#bytes) |  | Hash of the previous block. |
| state_root | [bytes](#bytes) |  | State root of the block. |
| sortition_seed | [bytes](#bytes) |  | Sortition seed of the block. |
| proposer_address | [string](#string) |  | Address of the proposer of the block. |






<a name="pactus-CertificateInfo"></a>

### CertificateInfo
Message containing information about a certificate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [bytes](#bytes) |  | Hash of the certificate. |
| round | [int32](#int32) |  | Round of the certificate. |
| committers | [int32](#int32) | repeated | List of committers in the certificate. |
| absentees | [int32](#int32) | repeated | List of absentees in the certificate. |
| signature | [bytes](#bytes) |  | Certificate signature. |






<a name="pactus-ConsensusInfo"></a>

### ConsensusInfo
Message containing information about consensus.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Address of the consensus instance. |
| Active | [bool](#bool) |  | Whether the consensus instance is active. |
| height | [uint32](#uint32) |  | Height of the consensus instance. |
| round | [int32](#int32) |  | Round of the consensus instance. |
| votes | [VoteInfo](#pactus-VoteInfo) | repeated | List of votes in the consensus instance. |






<a name="pactus-GetAccountRequest"></a>

### GetAccountRequest
Message to request account information based on an address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Address of the account. |






<a name="pactus-GetAccountResponse"></a>

### GetAccountResponse
Message containing the response with account information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [AccountInfo](#pactus-AccountInfo) |  | Account information. |






<a name="pactus-GetBlockHashRequest"></a>

### GetBlockHashRequest
Message to request block hash based on height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint32](#uint32) |  | Height of the block. |






<a name="pactus-GetBlockHashResponse"></a>

### GetBlockHashResponse
Message containing the response with the block hash.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [bytes](#bytes) |  | Hash of the block. |






<a name="pactus-GetBlockHeightRequest"></a>

### GetBlockHeightRequest
Message to request block height based on hash.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [bytes](#bytes) |  | Hash of the block. |






<a name="pactus-GetBlockHeightResponse"></a>

### GetBlockHeightResponse
Message containing the response with the block height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint32](#uint32) |  | Height of the block. |






<a name="pactus-GetBlockRequest"></a>

### GetBlockRequest
Message to request block information based on height and verbosity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint32](#uint32) |  | Height of the block. |
| verbosity | [BlockVerbosity](#pactus-BlockVerbosity) |  | Verbosity level for block information. |






<a name="pactus-GetBlockResponse"></a>

### GetBlockResponse
Message containing the response with block information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint32](#uint32) |  | Height of the block. |
| hash | [bytes](#bytes) |  | Hash of the block. |
| data | [bytes](#bytes) |  | Block data. |
| block_time | [uint32](#uint32) |  | Block timestamp. |
| header | [BlockHeaderInfo](#pactus-BlockHeaderInfo) |  | Block header information. |
| prev_cert | [CertificateInfo](#pactus-CertificateInfo) |  | Certificate information of the previous block. |
| txs | [TransactionInfo](#pactus-TransactionInfo) | repeated | List of transactions in the block. |






<a name="pactus-GetBlockchainInfoRequest"></a>

### GetBlockchainInfoRequest
Message to request general information about the blockchain.






<a name="pactus-GetBlockchainInfoResponse"></a>

### GetBlockchainInfoResponse
Message containing the response with general blockchain information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| last_block_height | [uint32](#uint32) |  | Height of the last block. |
| last_block_hash | [bytes](#bytes) |  | Hash of the last block. |
| total_accounts | [int32](#int32) |  | Total number of accounts. |
| total_validators | [int32](#int32) |  | Total number of validators. |
| total_power | [int64](#int64) |  | Total power in the blockchain. |
| committee_power | [int64](#int64) |  | Power of the committee. |
| committee_validators | [ValidatorInfo](#pactus-ValidatorInfo) | repeated | List of committee validators. |






<a name="pactus-GetConsensusInfoRequest"></a>

### GetConsensusInfoRequest
Message to request consensus information.






<a name="pactus-GetConsensusInfoResponse"></a>

### GetConsensusInfoResponse
Message containing the response with consensus information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instances | [ConsensusInfo](#pactus-ConsensusInfo) | repeated | List of consensus instances. |






<a name="pactus-GetPublicKeyRequest"></a>

### GetPublicKeyRequest
Message to request public key based on an address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Address for which public key is requested. |






<a name="pactus-GetPublicKeyResponse"></a>

### GetPublicKeyResponse
Message containing the response with the public key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_key | [string](#string) |  | Public key of the account. |






<a name="pactus-GetValidatorAddressesRequest"></a>

### GetValidatorAddressesRequest
Message to request validator addresses.






<a name="pactus-GetValidatorAddressesResponse"></a>

### GetValidatorAddressesResponse
Message containing the response with a list of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| addresses | [string](#string) | repeated | List of validator addresses. |






<a name="pactus-GetValidatorByNumberRequest"></a>

### GetValidatorByNumberRequest
Message to request validator information based on a validator number.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| number | [int32](#int32) |  | Validator number. |






<a name="pactus-GetValidatorRequest"></a>

### GetValidatorRequest
Message to request validator information based on an address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Address of the validator. |






<a name="pactus-GetValidatorResponse"></a>

### GetValidatorResponse
Message containing the response with validator information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| validator | [ValidatorInfo](#pactus-ValidatorInfo) |  | Validator information. |






<a name="pactus-ValidatorInfo"></a>

### ValidatorInfo
Message containing information about a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [bytes](#bytes) |  | Hash of the validator. |
| data | [bytes](#bytes) |  | Validator data. |
| public_key | [string](#string) |  | Public key of the validator. |
| number | [int32](#int32) |  | Validator number. |
| stake | [int64](#int64) |  | Validator stake. |
| last_bonding_height | [uint32](#uint32) |  | Last bonding height. |
| last_sortition_height | [uint32](#uint32) |  | Last sortition height. |
| unbonding_height | [uint32](#uint32) |  | Unbonding height. |
| address | [string](#string) |  | Address of the validator. |
| availability_score | [double](#double) |  | Availability score of the validator. |






<a name="pactus-VoteInfo"></a>

### VoteInfo
Message containing information about a vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [VoteType](#pactus-VoteType) |  | Type of the vote. |
| voter | [string](#string) |  | Voter&#39;s address. |
| block_hash | [bytes](#bytes) |  | Hash of the block being voted on. |
| round | [int32](#int32) |  | Round of the vote. |
| cp_round | [int32](#int32) |  | Consensus round of the vote. |
| cp_value | [int32](#int32) |  | Consensus value of the vote. |





 


<a name="pactus-BlockVerbosity"></a>

### BlockVerbosity
Enumeration for verbosity level when requesting block information.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BLOCK_DATA | 0 | Request block data only. |
| BLOCK_INFO | 1 | Request block information only. |
| BLOCK_TRANSACTIONS | 2 | Request block transactions only. |



<a name="pactus-VoteType"></a>

### VoteType
Enumeration for types of votes.

| Name | Number | Description |
| ---- | ------ | ----------- |
| VOTE_UNKNOWN | 0 | Unknown vote type. |
| VOTE_PREPARE | 1 | Prepare vote type. |
| VOTE_PRECOMMIT | 2 | Precommit vote type. |
| VOTE_CHANGE_PROPOSER | 3 | Change proposer vote type. |


 

 


<a name="pactus-Blockchain"></a>

### Blockchain
Blockchain service defines RPC methods for interacting with the blockchain.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBlock | [GetBlockRequest](#pactus-GetBlockRequest) | [GetBlockResponse](#pactus-GetBlockResponse) | GetBlock retrieves information about a block based on the provided request parameters. |
| GetBlockHash | [GetBlockHashRequest](#pactus-GetBlockHashRequest) | [GetBlockHashResponse](#pactus-GetBlockHashResponse) | GetBlockHash retrieves the hash of a block at the specified height. |
| GetBlockHeight | [GetBlockHeightRequest](#pactus-GetBlockHeightRequest) | [GetBlockHeightResponse](#pactus-GetBlockHeightResponse) | GetBlockHeight retrieves the height of a block with the specified hash. |
| GetBlockchainInfo | [GetBlockchainInfoRequest](#pactus-GetBlockchainInfoRequest) | [GetBlockchainInfoResponse](#pactus-GetBlockchainInfoResponse) | GetBlockchainInfo retrieves general information about the blockchain. |
| GetConsensusInfo | [GetConsensusInfoRequest](#pactus-GetConsensusInfoRequest) | [GetConsensusInfoResponse](#pactus-GetConsensusInfoResponse) | GetConsensusInfo retrieves information about the consensus instances. |
| GetAccount | [GetAccountRequest](#pactus-GetAccountRequest) | [GetAccountResponse](#pactus-GetAccountResponse) | GetAccount retrieves information about an account based on the provided address. |
| GetValidator | [GetValidatorRequest](#pactus-GetValidatorRequest) | [GetValidatorResponse](#pactus-GetValidatorResponse) | GetValidator retrieves information about a validator based on the provided address. |
| GetValidatorByNumber | [GetValidatorByNumberRequest](#pactus-GetValidatorByNumberRequest) | [GetValidatorResponse](#pactus-GetValidatorResponse) | GetValidatorByNumber retrieves information about a validator based on the provided number. |
| GetValidatorAddresses | [GetValidatorAddressesRequest](#pactus-GetValidatorAddressesRequest) | [GetValidatorAddressesResponse](#pactus-GetValidatorAddressesResponse) | GetValidatorAddresses retrieves a list of all validator addresses. |
| GetPublicKey | [GetPublicKeyRequest](#pactus-GetPublicKeyRequest) | [GetPublicKeyResponse](#pactus-GetPublicKeyResponse) | GetPublicKey retrieves the public key of an account based on the provided address. |

 



<a name="network-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## network.proto



<a name="pactus-GetNetworkInfoRequest"></a>

### GetNetworkInfoRequest
Request message for retrieving overall network information.






<a name="pactus-GetNetworkInfoResponse"></a>

### GetNetworkInfoResponse
Response message containing information about the overall network.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network_name | [string](#string) |  | Name of the network. |
| total_sent_bytes | [uint32](#uint32) |  | Total bytes sent across the network. |
| total_received_bytes | [uint32](#uint32) |  | Total bytes received across the network. |
| connected_peers_count | [uint32](#uint32) |  | Number of connected peers. |
| connected_peers | [PeerInfo](#pactus-PeerInfo) | repeated | List of connected peers. |
| sent_bytes | [GetNetworkInfoResponse.SentBytesEntry](#pactus-GetNetworkInfoResponse-SentBytesEntry) | repeated | Bytes sent per peer ID. |
| received_bytes | [GetNetworkInfoResponse.ReceivedBytesEntry](#pactus-GetNetworkInfoResponse-ReceivedBytesEntry) | repeated | Bytes received per peer ID. |






<a name="pactus-GetNetworkInfoResponse-ReceivedBytesEntry"></a>

### GetNetworkInfoResponse.ReceivedBytesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [uint64](#uint64) |  |  |






<a name="pactus-GetNetworkInfoResponse-SentBytesEntry"></a>

### GetNetworkInfoResponse.SentBytesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [uint32](#uint32) |  |  |
| value | [uint64](#uint64) |  |  |






<a name="pactus-GetNodeInfoRequest"></a>

### GetNodeInfoRequest
Request message for retrieving information about a specific node in the
network.






<a name="pactus-GetNodeInfoResponse"></a>

### GetNodeInfoResponse
Response message containing information about a specific node in the network.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| moniker | [string](#string) |  | Moniker of the node. |
| agent | [string](#string) |  | Agent information of the node. |
| peer_id | [bytes](#bytes) |  | Peer ID of the node. |
| started_at | [uint64](#uint64) |  | Timestamp when the node started. |
| reachability | [string](#string) |  | Reachability status of the node. |
| services | [int32](#int32) | repeated | List of services provided by the node. |
| services_names | [string](#string) | repeated | Names of services provided by the node. |
| addrs | [string](#string) | repeated | List of addresses associated with the node. |
| protocols | [string](#string) | repeated | List of protocols supported by the node. |






<a name="pactus-PeerInfo"></a>

### PeerInfo
Information about a peer in the network.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [int32](#int32) |  | Status of the peer. |
| moniker | [string](#string) |  | Moniker of the peer. |
| agent | [string](#string) |  | Agent information of the peer. |
| peer_id | [bytes](#bytes) |  | Peer ID of the peer. |
| consensus_keys | [string](#string) | repeated | Consensus keys used by the peer. |
| consensus_address | [string](#string) | repeated | Consensus address of the peer. |
| services | [uint32](#uint32) |  | Services provided by the peer. |
| last_block_hash | [bytes](#bytes) |  | Hash of the last block the peer knows. |
| height | [uint32](#uint32) |  | Height of the peer in the blockchain. |
| received_messages | [int32](#int32) |  | Count of received messages. |
| invalid_messages | [int32](#int32) |  | Count of invalid messages received. |
| last_sent | [int64](#int64) |  | Timestamp of the last sent message. |
| last_received | [int64](#int64) |  | Timestamp of the last received message. |
| sent_bytes | [PeerInfo.SentBytesEntry](#pactus-PeerInfo-SentBytesEntry) | repeated | Bytes sent per message type. |
| received_bytes | [PeerInfo.ReceivedBytesEntry](#pactus-PeerInfo-ReceivedBytesEntry) | repeated | Bytes received per message type. |
| address | [string](#string) |  | Network address of the peer. |
| direction | [string](#string) |  | Direction of connection with the peer. |
| protocols | [string](#string) | repeated | List of protocols supported by the peer. |
| total_sessions | [int32](#int32) |  | Total sessions with the peer. |
| completed_sessions | [int32](#int32) |  | Completed sessions with the peer. |






<a name="pactus-PeerInfo-ReceivedBytesEntry"></a>

### PeerInfo.ReceivedBytesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [int32](#int32) |  |  |
| value | [int64](#int64) |  |  |






<a name="pactus-PeerInfo-SentBytesEntry"></a>

### PeerInfo.SentBytesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [int32](#int32) |  |  |
| value | [int64](#int64) |  |  |





 

 

 


<a name="pactus-Network"></a>

### Network
Network service provides RPCs for retrieving information about the network.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetNetworkInfo | [GetNetworkInfoRequest](#pactus-GetNetworkInfoRequest) | [GetNetworkInfoResponse](#pactus-GetNetworkInfoResponse) | GetNetworkInfo retrieves information about the overall network. |
| GetNodeInfo | [GetNodeInfoRequest](#pactus-GetNodeInfoRequest) | [GetNodeInfoResponse](#pactus-GetNodeInfoResponse) | GetNodeInfo retrieves information about a specific node in the network. |

 



<a name="wallet-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## wallet.proto



<a name="pactus-CreateWalletRequest"></a>

### CreateWalletRequest
Request message for creating a new wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the new wallet. |
| mnemonic | [string](#string) |  | Mnemonic for wallet recovery. |
| language | [string](#string) |  | Language for the mnemonic. |
| password | [string](#string) |  | Password for securing the wallet. |






<a name="pactus-CreateWalletResponse"></a>

### CreateWalletResponse
Response message containing the name of the created wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the created wallet. |






<a name="pactus-GetValidatorAddressRequest"></a>

### GetValidatorAddressRequest
Request message for obtaining the validator address associated with a public
key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_key | [string](#string) |  | Public key for which the validator address is requested. |






<a name="pactus-GetValidatorAddressResponse"></a>

### GetValidatorAddressResponse
Response message containing the validator address corresponding to a public
key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Validator address associated with the public key. |






<a name="pactus-LoadWalletRequest"></a>

### LoadWalletRequest
Request message for loading an existing wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the wallet to load. |






<a name="pactus-LoadWalletResponse"></a>

### LoadWalletResponse
Response message containing the name of the loaded wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the loaded wallet. |






<a name="pactus-LockWalletRequest"></a>

### LockWalletRequest
Request message for locking a currently loaded wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the wallet to lock. |






<a name="pactus-LockWalletResponse"></a>

### LockWalletResponse
Response message containing the name of the locked wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the locked wallet. |






<a name="pactus-SignRawTransactionRequest"></a>

### SignRawTransactionRequest
Request message for signing a raw transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the wallet used for signing. |
| raw_transaction | [bytes](#bytes) |  | Raw transaction data to be signed. |
| password | [string](#string) |  | Password for unlocking the wallet for signing. |






<a name="pactus-SignRawTransactionResponse"></a>

### SignRawTransactionResponse
Response message containing the transaction ID and signed raw transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| transaction_id | [bytes](#bytes) |  | ID of the signed transaction. |
| signed_raw_transaction | [bytes](#bytes) |  | Signed raw transaction data. |






<a name="pactus-UnloadWalletRequest"></a>

### UnloadWalletRequest
Request message for unloading a currently loaded wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the wallet to unload. |






<a name="pactus-UnloadWalletResponse"></a>

### UnloadWalletResponse
Response message containing the name of the unloaded wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the unloaded wallet. |






<a name="pactus-UnlockWalletRequest"></a>

### UnlockWalletRequest
Request message for unlocking a wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the wallet to unlock. |
| password | [string](#string) |  | Password for unlocking the wallet. |
| timeout | [int32](#int32) |  | Timeout duration for the unlocked state. |






<a name="pactus-UnlockWalletResponse"></a>

### UnlockWalletResponse
Response message containing the name of the unlocked wallet.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wallet_name | [string](#string) |  | Name of the unlocked wallet. |





 

 

 


<a name="pactus-Wallet"></a>

### Wallet
Define the Wallet service with various RPC methods for wallet management.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateWallet | [CreateWalletRequest](#pactus-CreateWalletRequest) | [CreateWalletResponse](#pactus-CreateWalletResponse) | CreateWallet creates a new wallet with the specified parameters. |
| LoadWallet | [LoadWalletRequest](#pactus-LoadWalletRequest) | [LoadWalletResponse](#pactus-LoadWalletResponse) | LoadWallet loads an existing wallet with the given name. |
| UnloadWallet | [UnloadWalletRequest](#pactus-UnloadWalletRequest) | [UnloadWalletResponse](#pactus-UnloadWalletResponse) | UnloadWallet unloads a currently loaded wallet with the specified name. |
| LockWallet | [LockWalletRequest](#pactus-LockWalletRequest) | [LockWalletResponse](#pactus-LockWalletResponse) | LockWallet locks a currently loaded wallet with the provided password and timeout. |
| UnlockWallet | [UnlockWalletRequest](#pactus-UnlockWalletRequest) | [UnlockWalletResponse](#pactus-UnlockWalletResponse) | UnlockWallet unlocks a locked wallet with the provided password and timeout. |
| SignRawTransaction | [SignRawTransactionRequest](#pactus-SignRawTransactionRequest) | [SignRawTransactionResponse](#pactus-SignRawTransactionResponse) | SignRawTransaction signs a raw transaction for a specified wallet. |
| GetValidatorAddress | [GetValidatorAddressRequest](#pactus-GetValidatorAddressRequest) | [GetValidatorAddressResponse](#pactus-GetValidatorAddressResponse) | GetValidatorAddress retrieves the validator address associated with a public key. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

