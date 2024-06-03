# JSON-RPC API Reference

Each node in the Pactus network can be configured to use the [JSON-RPC 2.0](https://www.jsonrpc.org/specification) protocol for communication.
Here you can find the list of all JSON-RPC methods and messages.

All the amounts and values in JSON-RPC endpoints are in NanoPAC units, which are atomic and the smallest unit in the Pactus blockchain.
Each PAC is equivalent to
 1,000,000,000 or 10<sup>9</sup> NanoPACs.


## Methods





- [pactus.transaction.get_transaction](#pactus.transaction.get_transaction)


- [pactus.transaction.calculate_fee](#pactus.transaction.calculate_fee)


- [pactus.transaction.broadcast_transaction](#pactus.transaction.broadcast_transaction)


- [pactus.transaction.get_raw_transfer_transaction](#pactus.transaction.get_raw_transfer_transaction)


- [pactus.transaction.get_raw_bond_transaction](#pactus.transaction.get_raw_bond_transaction)


- [pactus.transaction.get_raw_unbond_transaction](#pactus.transaction.get_raw_unbond_transaction)


- [pactus.transaction.get_raw_withdraw_transaction](#pactus.transaction.get_raw_withdraw_transaction)






- [pactus.blockchain.get_block](#pactus.blockchain.get_block)


- [pactus.blockchain.get_block_hash](#pactus.blockchain.get_block_hash)


- [pactus.blockchain.get_block_height](#pactus.blockchain.get_block_height)


- [pactus.blockchain.get_blockchain_info](#pactus.blockchain.get_blockchain_info)


- [pactus.blockchain.get_consensus_info](#pactus.blockchain.get_consensus_info)


- [pactus.blockchain.get_account](#pactus.blockchain.get_account)


- [pactus.blockchain.get_validator](#pactus.blockchain.get_validator)


- [pactus.blockchain.get_validator_by_number](#pactus.blockchain.get_validator_by_number)


- [pactus.blockchain.get_validator_addresses](#pactus.blockchain.get_validator_addresses)


- [pactus.blockchain.get_public_key](#pactus.blockchain.get_public_key)






- [pactus.network.get_network_info](#pactus.network.get_network_info)


- [pactus.network.get_node_info](#pactus.network.get_node_info)






- [pactus.wallet.create_wallet](#pactus.wallet.create_wallet)


- [pactus.wallet.restore_wallet](#pactus.wallet.restore_wallet)


- [pactus.wallet.load_wallet](#pactus.wallet.load_wallet)


- [pactus.wallet.unload_wallet](#pactus.wallet.unload_wallet)


- [pactus.wallet.lock_wallet](#pactus.wallet.lock_wallet)


- [pactus.wallet.unlock_wallet](#pactus.wallet.unlock_wallet)


- [pactus.wallet.get_total_balance](#pactus.wallet.get_total_balance)


- [pactus.wallet.sign_raw_transaction](#pactus.wallet.sign_raw_transaction)


- [pactus.wallet.get_validator_address](#pactus.wallet.get_validator_address)


- [pactus.wallet.get_new_address](#pactus.wallet.get_new_address)


- [pactus.wallet.get_address_history](#pactus.wallet.get_address_history)




---



<a id="pactus.transaction.get_transaction"></a>

## Method pactus.transaction.get_transaction

pactus.transaction.get_transaction retrieves transaction details based on the provided request
parameters.

### Parameters
```json
{
	"id": "str",	// (string) Transaction ID.
	"verbosity": "TRANSACTION_DATA or TRANSACTION_INFO"	// (string) Verbosity level for transaction details.
}
```

### Result
```json
{
	"block_height": n,	// (numeric) Height of the block containing the transaction.
	"block_time": n,	// (numeric) Time of the block containing the transaction.
	"transaction": {	// (json object) Information about the transaction.
		"bond": {	// (json object) Bond payload.
			"receiver": "str",	// (string) Receiver's address.
			"sender": "str",	// (string) Sender's address.
			"stake": n	// (numeric) Stake amount in NanoPAC.
		},
		"data": "str",	// (string) Transaction data.
		"fee": n,	// (numeric) Transaction fee in NanoPAC.
		"id": "str",	// (string) Transaction ID.
		"lock_time": n,	// (numeric) Lock time for the transaction.
		"memo": "str",	// (string) Transaction memo.
		"payload_type": "UNKNOWN or TRANSFER_PAYLOAD or BOND_PAYLOAD or SORTITION_PAYLOAD or UNBOND_PAYLOAD or WITHDRAW_PAYLOAD",	// (string) Type of transaction payload.
		"public_key": "str",	// (string) Public key associated with the transaction.
		"signature": "str",	// (string) Transaction signature.
		"sortition": {	// (json object) Sortition payload.
			"address": "str",	// (string) Address associated with the sortition.
			"proof": "str"	// (string) Proof for the sortition.
		},
		"transfer": {	// (json object) Transfer payload.
			"amount": n,	// (numeric) Transaction amount in NanoPAC.
			"receiver": "str",	// (string) Receiver's address.
			"sender": "str"	// (string) Sender's address.
		},
		"unbond": {	// (json object) Unbond payload.
			"validator": "str"	// (string) Address of the validator to unbond from.
		},
		"value": n,	// (numeric) Transaction value in NanoPAC.
		"version": n,	// (numeric) Transaction version.
		"withdraw": {	// (json object) Withdraw payload.
			"amount": n,	// (numeric) Withdrawal amount in NanoPAC.
			"from": "str",	// (string) Address to withdraw from.
			"to": "str"	// (string) Address to withdraw to.
		}
	}
}
```
---


<a id="pactus.transaction.calculate_fee"></a>

## Method pactus.transaction.calculate_fee

pactus.transaction.calculate_fee calculates the transaction fee based on the specified amount
and payload type.

### Parameters
```json
{
	"amount": n,	// (numeric) Transaction amount in NanoPAC.
	"fixed_amount": true|false,	// (boolean) Indicates that amount should be fixed and includes the fee.
	"payload_type": "UNKNOWN or TRANSFER_PAYLOAD or BOND_PAYLOAD or SORTITION_PAYLOAD or UNBOND_PAYLOAD or WITHDRAW_PAYLOAD"	// (string) Type of transaction payload.
}
```

### Result
```json
{
	"amount": n,	// (numeric) Calculated amount in NanoPAC.
	"fee": n	// (numeric) Calculated transaction fee in NanoPAC.
}
```
---


<a id="pactus.transaction.broadcast_transaction"></a>

## Method pactus.transaction.broadcast_transaction

pactus.transaction.broadcast_transaction broadcasts a signed transaction to the network.

### Parameters
```json
{
	"signed_raw_transaction": "str"	// (string) Signed raw transaction data.
}
```

### Result
```json
{
	"id": "str"	// (string) Transaction ID.
}
```
---


<a id="pactus.transaction.get_raw_transfer_transaction"></a>

## Method pactus.transaction.get_raw_transfer_transaction

pactus.transaction.get_raw_transfer_transaction retrieves raw details of a transfer transaction.

### Parameters
```json
{
	"amount": n,	// (numeric) Transfer amount in NanoPAC.\nIt should be greater than 0.
	"fee": n,	// (numeric) Transaction fee in NanoPAC.\nIf not explicitly set, it is calculated based on the amount.
	"lock_time": n,	// (numeric) Lock time for the transaction.\nIf not explicitly set, it sets to the last block height.
	"memo": "str",	// (string) Transaction memo.
	"receiver": "str",	// (string) Receiver's account address.
	"sender": "str"	// (string) Sender's account address.
}
```

### Result
```json
{
	"raw_transaction": "str"	// (string) Raw transaction data.
}
```
---


<a id="pactus.transaction.get_raw_bond_transaction"></a>

## Method pactus.transaction.get_raw_bond_transaction

pactus.transaction.get_raw_bond_transaction retrieves raw details of a bond transaction.

### Parameters
```json
{
	"fee": n,	// (numeric) Transaction fee in NanoPAC.\nIf not explicitly set, it is calculated based on the stake.
	"lock_time": n,	// (numeric) Lock time for the transaction.\nIf not explicitly set, it sets to the last block height.
	"memo": "str",	// (string) Transaction memo.
	"public_key": "str",	// (string) Public key of the validator.
	"receiver": "str",	// (string) Receiver's validator address.
	"sender": "str",	// (string) Sender's account address.
	"stake": n	// (numeric) Stake amount in NanoPAC.\nIt should be greater than 0.
}
```

### Result
```json
{
	"raw_transaction": "str"	// (string) Raw transaction data.
}
```
---


<a id="pactus.transaction.get_raw_unbond_transaction"></a>

## Method pactus.transaction.get_raw_unbond_transaction

pactus.transaction.get_raw_unbond_transaction retrieves raw details of an unbond transaction.

### Parameters
```json
{
	"lock_time": n,	// (numeric) Lock time for the transaction.\nIf not explicitly set, it sets to the last block height.
	"memo": "str",	// (string) Transaction memo.
	"validator_address": "str"	// (string) Address of the validator to unbond from.
}
```

### Result
```json
{
	"raw_transaction": "str"	// (string) Raw transaction data.
}
```
---


<a id="pactus.transaction.get_raw_withdraw_transaction"></a>

## Method pactus.transaction.get_raw_withdraw_transaction

pactus.transaction.get_raw_withdraw_transaction retrieves raw details of a withdraw transaction.

### Parameters
```json
{
	"account_address": "str",	// (string) Address of the account to withdraw to.
	"amount": n,	// (numeric) Withdrawal amount in NanoPAC.\nIt should be greater than 0.
	"fee": n,	// (numeric) Transaction fee in NanoPAC.\nIf not explicitly set, it is calculated based on the amount.
	"lock_time": n,	// (numeric) Lock time for the transaction.\nIf not explicitly set, it sets to the last block height.
	"memo": "str",	// (string) Transaction memo.
	"validator_address": "str"	// (string) Address of the validator to withdraw from.
}
```

### Result
```json
{
	"raw_transaction": "str"	// (string) Raw transaction data.
}
```
---






<a id="pactus.blockchain.get_block"></a>

## Method pactus.blockchain.get_block

pactus.blockchain.get_block retrieves information about a block based on the provided request
parameters.

### Parameters
```json
{
	"height": n,	// (numeric) Height of the block.
	"verbosity": "BLOCK_DATA or BLOCK_INFO or BLOCK_TRANSACTIONS"	// (string) Verbosity level for block information.
}
```

### Result
```json
{
	"block_time": n,	// (numeric) Block timestamp.
	"data": "str",	// (string) Block data, only available if the verbosity level is set to BLOCK_DATA.
	"hash": "str",	// (string) Hash of the block.
	"header": {	// (json object) Block header information.
		"prev_block_hash": "str",	// (string) Hash of the previous block.
		"proposer_address": "str",	// (string) Address of the proposer of the block.
		"sortition_seed": "str",	// (string) Sortition seed of the block.
		"state_root": "str",	// (string) State root of the block.
		"version": n	// (numeric) Block version.
	},
	"height": n,	// (numeric) Height of the block.
	"prev_cert": {	// (json object) Certificate information of the previous block.
		"absentees": [	// (json array) List of absentees in the certificate.
			n,
			...
		],
		"committers": [	// (json array) List of committers in the certificate.
			n,
			...
		],
		"hash": "str",	// (string) Hash of the certificate.
		"round": n,	// (numeric) Round of the certificate.
		"signature": "str"	// (string) Certificate signature.
	},
	"txs": [	// (json array) List of transactions in the block.\nTransaction information is available when the verbosity level is set to BLOCK_TRANSACTIONS.
		{
			"bond": {	// (json object) Bond payload.
				"receiver": "str",	// (string) Receiver's address.
				"sender": "str",	// (string) Sender's address.
				"stake": n	// (numeric) Stake amount in NanoPAC.
			},
			"data": "str",	// (string) Transaction data.
			"fee": n,	// (numeric) Transaction fee in NanoPAC.
			"id": "str",	// (string) Transaction ID.
			"lock_time": n,	// (numeric) Lock time for the transaction.
			"memo": "str",	// (string) Transaction memo.
			"payload_type": "UNKNOWN or TRANSFER_PAYLOAD or BOND_PAYLOAD or SORTITION_PAYLOAD or UNBOND_PAYLOAD or WITHDRAW_PAYLOAD",	// (string) Type of transaction payload.
			"public_key": "str",	// (string) Public key associated with the transaction.
			"signature": "str",	// (string) Transaction signature.
			"sortition": {	// (json object) Sortition payload.
				"address": "str",	// (string) Address associated with the sortition.
				"proof": "str"	// (string) Proof for the sortition.
			},
			"transfer": {	// (json object) Transfer payload.
				"amount": n,	// (numeric) Transaction amount in NanoPAC.
				"receiver": "str",	// (string) Receiver's address.
				"sender": "str"	// (string) Sender's address.
			},
			"unbond": {	// (json object) Unbond payload.
				"validator": "str"	// (string) Address of the validator to unbond from.
			},
			"value": n,	// (numeric) Transaction value in NanoPAC.
			"version": n,	// (numeric) Transaction version.
			"withdraw": {	// (json object) Withdraw payload.
				"amount": n,	// (numeric) Withdrawal amount in NanoPAC.
				"from": "str",	// (string) Address to withdraw from.
				"to": "str"	// (string) Address to withdraw to.
			}
		},
		...
	]
}
```
---


<a id="pactus.blockchain.get_block_hash"></a>

## Method pactus.blockchain.get_block_hash

pactus.blockchain.get_block_hash retrieves the hash of a block at the specified height.

### Parameters
```json
{
	"height": n	// (numeric) Height of the block.
}
```

### Result
```json
{
	"hash": "str"	// (string) Hash of the block.
}
```
---


<a id="pactus.blockchain.get_block_height"></a>

## Method pactus.blockchain.get_block_height

pactus.blockchain.get_block_height retrieves the height of a block with the specified hash.

### Parameters
```json
{
	"hash": "str"	// (string) Hash of the block.
}
```

### Result
```json
{
	"height": n	// (numeric) Height of the block.
}
```
---


<a id="pactus.blockchain.get_blockchain_info"></a>

## Method pactus.blockchain.get_blockchain_info

pactus.blockchain.get_blockchain_info retrieves general information about the blockchain.

### Parameters
```json
{}
```

### Result
```json
{
	"committee_power": n,	// (numeric) Power of the committee.
	"committee_validators": [	// (json array) List of committee validators.
		{
			"address": "str",	// (string) Address of the validator.
			"availability_score": n,	// (numeric) Availability score of the validator.
			"data": "str",	// (string) Validator data.
			"hash": "str",	// (string) Hash of the validator.
			"last_bonding_height": n,	// (numeric) Last bonding height.
			"last_sortition_height": n,	// (numeric) Last sortition height.
			"number": n,	// (numeric) Validator number.
			"public_key": "str",	// (string) Public key of the validator.
			"stake": n,	// (numeric) Validator stake in NanoPAC.
			"unbonding_height": n	// (numeric) Unbonding height.
		},
		...
	],
	"last_block_hash": "str",	// (string) Hash of the last block.
	"last_block_height": n,	// (numeric) Height of the last block.
	"total_accounts": n,	// (numeric) Total number of accounts.
	"total_power": n,	// (numeric) Total power in the blockchain.
	"total_validators": n	// (numeric) Total number of validators.
}
```
---


<a id="pactus.blockchain.get_consensus_info"></a>

## Method pactus.blockchain.get_consensus_info

pactus.blockchain.get_consensus_info retrieves information about the consensus instances.

### Parameters
```json
{}
```

### Result
```json
{
	"instances": [	// (json array) List of consensus instances.
		{
			"Active": true|false,	// (boolean) Whether the consensus instance is active.
			"address": "str",	// (string) Address of the consensus instance.
			"height": n,	// (numeric) Height of the consensus instance.
			"round": n,	// (numeric) Round of the consensus instance.
			"votes": [	// (json array) List of votes in the consensus instance.
				{
					"block_hash": "str",	// (string) Hash of the block being voted on.
					"cp_round": n,	// (numeric) Consensus round of the vote.
					"cp_value": n,	// (numeric) Consensus value of the vote.
					"round": n,	// (numeric) Round of the vote.
					"type": "VOTE_UNKNOWN or VOTE_PREPARE or VOTE_PRECOMMIT or VOTE_CHANGE_PROPOSER",	// (string) Type of the vote.
					"voter": "str"	// (string) Voter's address.
				},
				...
			]
		},
		...
	]
}
```
---


<a id="pactus.blockchain.get_account"></a>

## Method pactus.blockchain.get_account

pactus.blockchain.get_account retrieves information about an account based on the provided
address.

### Parameters
```json
{
	"address": "str"	// (string) Address of the account.
}
```

### Result
```json
{
	"account": {	// (json object) Account information.
		"address": "str",	// (string) Address of the account.
		"balance": n,	// (numeric) Account balance in NanoPAC.
		"data": "str",	// (string) Account data.
		"hash": "str",	// (string) Hash of the account.
		"number": n	// (numeric) Account number.
	}
}
```
---


<a id="pactus.blockchain.get_validator"></a>

## Method pactus.blockchain.get_validator

pactus.blockchain.get_validator retrieves information about a validator based on the provided
address.

### Parameters
```json
{
	"address": "str"	// (string) Address of the validator.
}
```

### Result
```json
{
	"validator": {	// (json object) Validator information.
		"address": "str",	// (string) Address of the validator.
		"availability_score": n,	// (numeric) Availability score of the validator.
		"data": "str",	// (string) Validator data.
		"hash": "str",	// (string) Hash of the validator.
		"last_bonding_height": n,	// (numeric) Last bonding height.
		"last_sortition_height": n,	// (numeric) Last sortition height.
		"number": n,	// (numeric) Validator number.
		"public_key": "str",	// (string) Public key of the validator.
		"stake": n,	// (numeric) Validator stake in NanoPAC.
		"unbonding_height": n	// (numeric) Unbonding height.
	}
}
```
---


<a id="pactus.blockchain.get_validator_by_number"></a>

## Method pactus.blockchain.get_validator_by_number

pactus.blockchain.get_validator_by_number retrieves information about a validator based on the
provided number.

### Parameters
```json
{
	"number": n	// (numeric) Validator number.
}
```

### Result
```json
{
	"validator": {	// (json object) Validator information.
		"address": "str",	// (string) Address of the validator.
		"availability_score": n,	// (numeric) Availability score of the validator.
		"data": "str",	// (string) Validator data.
		"hash": "str",	// (string) Hash of the validator.
		"last_bonding_height": n,	// (numeric) Last bonding height.
		"last_sortition_height": n,	// (numeric) Last sortition height.
		"number": n,	// (numeric) Validator number.
		"public_key": "str",	// (string) Public key of the validator.
		"stake": n,	// (numeric) Validator stake in NanoPAC.
		"unbonding_height": n	// (numeric) Unbonding height.
	}
}
```
---


<a id="pactus.blockchain.get_validator_addresses"></a>

## Method pactus.blockchain.get_validator_addresses

pactus.blockchain.get_validator_addresses retrieves a list of all validator addresses.

### Parameters
```json
{}
```

### Result
```json
{
	"addresses": [	// (json array) List of validator addresses.
		"str",
		...
	]
}
```
---


<a id="pactus.blockchain.get_public_key"></a>

## Method pactus.blockchain.get_public_key

pactus.blockchain.get_public_key retrieves the public key of an account based on the provided
address.

### Parameters
```json
{
	"address": "str"	// (string) Address for which public key is requested.
}
```

### Result
```json
{
	"public_key": "str"	// (string) Public key of the account.
}
```
---






<a id="pactus.network.get_network_info"></a>

## Method pactus.network.get_network_info

pactus.network.get_network_info retrieves information about the overall network.

### Parameters
```json
{
	"only_connected": true|false	// (boolean) Only returns the peers with connected status
}
```

### Result
```json
{
	"connected_peers": [	// (json array) List of connected peers.
		{
			"address": "str",	// (string) Network address of the peer.
			"agent": "str",	// (string) Agent information of the peer.
			"completed_sessions": n,	// (numeric) Completed sessions with the peer.
			"consensus_address": [	// (json array) Consensus address of the peer.
				"str",
				...
			],
			"consensus_keys": [	// (json array) Consensus keys used by the peer.
				"str",
				...
			],
			"direction": "str",	// (string) Direction of connection with the peer.
			"height": n,	// (numeric) Height of the peer in the blockchain.
			"invalid_bundles": n,	// (numeric) Count of invalid bundles received.
			"last_block_hash": "str",	// (string) Hash of the last block the peer knows.
			"last_received": n,	// (numeric) Timestamp of the last received bundle.
			"last_sent": n,	// (numeric) Timestamp of the last sent bundle.
			"moniker": "str",	// (string) Moniker of the peer.
			"peer_id": "str",	// (string) Peer ID of the peer.
			"protocols": [	// (json array) List of protocols supported by the peer.
				"str",
				...
			],
			"received_bundles": n,	// (numeric) Count of received bundles.
			"received_bytes": {	// (key:value json object) Bytes received per message type.
				...: ...,
				n: n
			},
			"sent_bytes": {	// (key:value json object) Bytes sent per message type.
				...: ...,
				n: n
			},
			"services": n,	// (numeric) Services provided by the peer.
			"status": n,	// (numeric) Status of the peer.
			"total_sessions": n	// (numeric) Total sessions with the peer.
		},
		...
	],
	"connected_peers_count": n,	// (numeric) Number of connected peers.
	"network_name": "str",	// (string) Name of the network.
	"received_bytes": {	// (key:value json object) Bytes received per peer ID.
		...: ...,
		n: n
	},
	"sent_bytes": {	// (key:value json object) Bytes sent per peer ID.
		...: ...,
		n: n
	},
	"total_received_bytes": n,	// (numeric) Total bytes received across the network.
	"total_sent_bytes": n	// (numeric) Total bytes sent across the network.
}
```
---


<a id="pactus.network.get_node_info"></a>

## Method pactus.network.get_node_info

pactus.network.get_node_info retrieves information about a specific node in the network.

### Parameters
```json
{}
```

### Result
```json
{
	"agent": "str",	// (string) Agent information of the node.
	"clock_offset": n,	// (numeric) Clock offset
	"connection_info": {	// (json object) Connection information
		"connections": n,	// (numeric) Total number of the connection.
		"inbound_connections": n,	// (numeric) Number of inbound connections.
		"outbound_connections": n	// (numeric) Number of outbound connections.
	},
	"local_addrs": [	// (json array) List of addresses associated with the node.
		"str",
		...
	],
	"moniker": "str",	// (string) Moniker of the node.
	"peer_id": "str",	// (string) Peer ID of the node.
	"protocols": [	// (json array) List of protocols supported by the node.
		"str",
		...
	],
	"reachability": "str",	// (string) Reachability status of the node.
	"services": [	// (json array) List of services provided by the node.
		n,
		...
	],
	"services_names": [	// (json array) Names of services provided by the node.
		"str",
		...
	],
	"started_at": n	// (numeric) Timestamp when the node started.
}
```
---






<a id="pactus.wallet.create_wallet"></a>

## Method pactus.wallet.create_wallet

pactus.wallet.create_wallet creates a new wallet with the specified parameters.

### Parameters
```json
{
	"password": "str",	// (string) Password for securing the wallet.
	"wallet_name": "str"	// (string) Name of the new wallet.
}
```

### Result
```json
{
	"mnemonic": "str"	// (string) Menomic for wallet recovery.
}
```
---


<a id="pactus.wallet.restore_wallet"></a>

## Method pactus.wallet.restore_wallet

pactus.wallet.restore_wallet restores an existing wallet with the given mnemonic.

### Parameters
```json
{
	"mnemonic": "str",	// (string) Menomic for wallet recovery.
	"password": "str",	// (string) Password for securing the wallet.
	"wallet_name": "str"	// (string) Name of the wallet to restore.
}
```

### Result
```json
{
	"wallet_name": "str"	// (string) Name of the restored wallet.
}
```
---


<a id="pactus.wallet.load_wallet"></a>

## Method pactus.wallet.load_wallet

pactus.wallet.load_wallet loads an existing wallet with the given name.

### Parameters
```json
{
	"wallet_name": "str"	// (string) Name of the wallet to load.
}
```

### Result
```json
{
	"wallet_name": "str"	// (string) Name of the loaded wallet.
}
```
---


<a id="pactus.wallet.unload_wallet"></a>

## Method pactus.wallet.unload_wallet

pactus.wallet.unload_wallet unloads a currently loaded wallet with the specified name.

### Parameters
```json
{
	"wallet_name": "str"	// (string) Name of the wallet to unload.
}
```

### Result
```json
{
	"wallet_name": "str"	// (string) Name of the unloaded wallet.
}
```
---


<a id="pactus.wallet.lock_wallet"></a>

## Method pactus.wallet.lock_wallet

pactus.wallet.lock_wallet locks a currently loaded wallet with the provided password and
timeout.

### Parameters
```json
{
	"wallet_name": "str"	// (string) Name of the wallet to lock.
}
```

### Result
```json
{
	"wallet_name": "str"	// (string) Name of the locked wallet.
}
```
---


<a id="pactus.wallet.unlock_wallet"></a>

## Method pactus.wallet.unlock_wallet

pactus.wallet.unlock_wallet unlocks a locked wallet with the provided password and
timeout.

### Parameters
```json
{
	"password": "str",	// (string) Password for unlocking the wallet.
	"timeout": n,	// (numeric) Timeout duration for the unlocked state.
	"wallet_name": "str"	// (string) Name of the wallet to unlock.
}
```

### Result
```json
{
	"wallet_name": "str"	// (string) Name of the unlocked wallet.
}
```
---


<a id="pactus.wallet.get_total_balance"></a>

## Method pactus.wallet.get_total_balance

pactus.wallet.get_total_balance returns the total available balance of the wallet.

### Parameters
```json
{
	"wallet_name": "str"	// (string) Name of the wallet.
}
```

### Result
```json
{
	"total_balance": n,	// (numeric) The total balance of the wallet in NanoPAC.
	"wallet_name": "str"	// (string) Name of the wallet.
}
```
---


<a id="pactus.wallet.sign_raw_transaction"></a>

## Method pactus.wallet.sign_raw_transaction

pactus.wallet.sign_raw_transaction signs a raw transaction for a specified wallet.

### Parameters
```json
{
	"password": "str",	// (string) Password for unlocking the wallet for signing.
	"raw_transaction": "str",	// (string) Raw transaction data to be signed.
	"wallet_name": "str"	// (string) Name of the wallet used for signing.
}
```

### Result
```json
{
	"signed_raw_transaction": "str",	// (string) Signed raw transaction data.
	"transaction_id": "str"	// (string) ID of the signed transaction.
}
```
---


<a id="pactus.wallet.get_validator_address"></a>

## Method pactus.wallet.get_validator_address

pactus.wallet.get_validator_address retrieves the validator address associated with a
public key.

### Parameters
```json
{
	"public_key": "str"	// (string) Public key for which the validator address is requested.
}
```

### Result
```json
{
	"address": "str"	// (string) Validator address associated with the public key.
}
```
---


<a id="pactus.wallet.get_new_address"></a>

## Method pactus.wallet.get_new_address

pactus.wallet.get_new_address generates a new address for the specified wallet.

### Parameters
```json
{
	"address_type": "ADDRESS_TYPE_TREASURY or ADDRESS_TYPE_VALIDATOR or ADDRESS_TYPE_BLS_ACCOUNT",	// (string) Address type for the new address.
	"label": "str",	// (string) Label for the new address.
	"wallet_name": "str"	// (string) Name of the wallet for which the new address is requested.
}
```

### Result
```json
{
	"address_info": {	// (json object) Address information.
		"address": "str",	// (string) 
		"label": "str",	// (string) 
		"path": "str",	// (string) 
		"public_key": "str"	// (string) 
	},
	"wallet_name": "str"	// (string) Name of the wallet.
}
```
---


<a id="pactus.wallet.get_address_history"></a>

## Method pactus.wallet.get_address_history

pactus.wallet.get_address_history retrieve transaction history of an address.

### Parameters
```json
{
	"address": "str",	// (string) Address to get the transaction history of it.
	"wallet_name": "str"	// (string) Name of the wallet.
}
```

### Result
```json
{
	"history_info": [	// (json array) Array of address history and activities.
		{
			"amount": n,	// (numeric) amount of transaction.
			"description": "str",	// (string) description of transaction.
			"payload_type": "str",	// (string) payload type of transaction.
			"time": n,	// (numeric) transaction timestamp.
			"transaction_id": "str"	// (string) Hash of transaction.
		},
		...
	]
}
```
---




