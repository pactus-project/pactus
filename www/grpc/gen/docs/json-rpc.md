---
title: JSON-RPC API Reference
weight: 2
---

Each node in the Pactus network can be configured to use the 
[JSON-RPC](https://www.jsonrpc.org/specification) protocol for communication.
Here, you can find the list of all JSON-RPC methods and messages.

All the amounts and values in JSON-RPC endpoints are in NanoPAC units,
which are atomic and the smallest unit in the Pactus blockchain.
Each PAC is equivalent to 1,000,000,000 or 10<sup>9</sup> NanoPACs.

## Example

To call JSON-RPC methods, you need to create the JSON-RPC request:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "pactus.network.get_node_info",
  "params": {}
}
```

> Make sure you always add the `params` field, even if no parameters are needed, and ensure you use curly braces.

Then you use the `curl` command to send the request to the node:

```bash
curl --location 'http://localhost:8545/' \
--header 'Content-Type: application/json' \
--data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "pactus.network.get_node_info",
    "params": {}
}'
```

> Before sending the request, you need to enable the JSON-RPC service inside the
> [configuration](/get-started/configuration/).

### Using Basic Auth

If you have enabled the [gRPC Basic Authentication](/tutorials/grpc-sign-transactions/),
then you need to set the `Authorization` header.

```bash
curl --location 'http://localhost:8545/' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic <BASIC-AUTH-TOKEN>' \
--data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "pactus.blockchain.get_account",
    "params": {
        "address": "pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z"
    }
}'
```

<h2>JSON-RPC Methods</h2>

<div id="toc-container">
  <ul class="">
  <li> Transaction Service
      <ul>
        <li>
          <a href="#pactus.transaction.get_transaction">
          <span class="rpc-badge"></span> pactus.transaction.get_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.calculate_fee">
          <span class="rpc-badge"></span> pactus.transaction.calculate_fee</a>
        </li>
        <li>
          <a href="#pactus.transaction.broadcast_transaction">
          <span class="rpc-badge"></span> pactus.transaction.broadcast_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.get_raw_transfer_transaction">
          <span class="rpc-badge"></span> pactus.transaction.get_raw_transfer_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.get_raw_bond_transaction">
          <span class="rpc-badge"></span> pactus.transaction.get_raw_bond_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.get_raw_unbond_transaction">
          <span class="rpc-badge"></span> pactus.transaction.get_raw_unbond_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.get_raw_withdraw_transaction">
          <span class="rpc-badge"></span> pactus.transaction.get_raw_withdraw_transaction</a>
        </li>
        <li>
          <a href="#pactus.transaction.decode_raw_transaction">
          <span class="rpc-badge"></span> pactus.transaction.decode_raw_transaction</a>
        </li>
        </ul>
    </li>
    <li> Blockchain Service
      <ul>
        <li>
          <a href="#pactus.blockchain.get_block">
          <span class="rpc-badge"></span> pactus.blockchain.get_block</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_block_hash">
          <span class="rpc-badge"></span> pactus.blockchain.get_block_hash</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_block_height">
          <span class="rpc-badge"></span> pactus.blockchain.get_block_height</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_blockchain_info">
          <span class="rpc-badge"></span> pactus.blockchain.get_blockchain_info</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_consensus_info">
          <span class="rpc-badge"></span> pactus.blockchain.get_consensus_info</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_account">
          <span class="rpc-badge"></span> pactus.blockchain.get_account</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_validator">
          <span class="rpc-badge"></span> pactus.blockchain.get_validator</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_validator_by_number">
          <span class="rpc-badge"></span> pactus.blockchain.get_validator_by_number</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_validator_addresses">
          <span class="rpc-badge"></span> pactus.blockchain.get_validator_addresses</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_public_key">
          <span class="rpc-badge"></span> pactus.blockchain.get_public_key</a>
        </li>
        <li>
          <a href="#pactus.blockchain.get_tx_pool_content">
          <span class="rpc-badge"></span> pactus.blockchain.get_tx_pool_content</a>
        </li>
        </ul>
    </li>
    <li> Network Service
      <ul>
        <li>
          <a href="#pactus.network.get_network_info">
          <span class="rpc-badge"></span> pactus.network.get_network_info</a>
        </li>
        <li>
          <a href="#pactus.network.get_node_info">
          <span class="rpc-badge"></span> pactus.network.get_node_info</a>
        </li>
        </ul>
    </li>
    <li> Utils Service
      <ul>
        <li>
          <a href="#pactus.utils.sign_message_with_private_key">
          <span class="rpc-badge"></span> pactus.utils.sign_message_with_private_key</a>
        </li>
        <li>
          <a href="#pactus.utils.verify_message">
          <span class="rpc-badge"></span> pactus.utils.verify_message</a>
        </li>
        <li>
          <a href="#pactus.utils.b_l_s_public_key_aggregation">
          <span class="rpc-badge"></span> pactus.utils.b_l_s_public_key_aggregation</a>
        </li>
        <li>
          <a href="#pactus.utils.b_l_s_signature_aggregation">
          <span class="rpc-badge"></span> pactus.utils.b_l_s_signature_aggregation</a>
        </li>
        </ul>
    </li>
    <li> Wallet Service
      <ul>
        <li>
          <a href="#pactus.wallet.create_wallet">
          <span class="rpc-badge"></span> pactus.wallet.create_wallet</a>
        </li>
        <li>
          <a href="#pactus.wallet.restore_wallet">
          <span class="rpc-badge"></span> pactus.wallet.restore_wallet</a>
        </li>
        <li>
          <a href="#pactus.wallet.load_wallet">
          <span class="rpc-badge"></span> pactus.wallet.load_wallet</a>
        </li>
        <li>
          <a href="#pactus.wallet.unload_wallet">
          <span class="rpc-badge"></span> pactus.wallet.unload_wallet</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_total_balance">
          <span class="rpc-badge"></span> pactus.wallet.get_total_balance</a>
        </li>
        <li>
          <a href="#pactus.wallet.sign_raw_transaction">
          <span class="rpc-badge"></span> pactus.wallet.sign_raw_transaction</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_validator_address">
          <span class="rpc-badge"></span> pactus.wallet.get_validator_address</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_new_address">
          <span class="rpc-badge"></span> pactus.wallet.get_new_address</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_address_history">
          <span class="rpc-badge"></span> pactus.wallet.get_address_history</a>
        </li>
        <li>
          <a href="#pactus.wallet.sign_message">
          <span class="rpc-badge"></span> pactus.wallet.sign_message</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_total_stake">
          <span class="rpc-badge"></span> pactus.wallet.get_total_stake</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_address_info">
          <span class="rpc-badge"></span> pactus.wallet.get_address_info</a>
        </li>
        <li>
          <a href="#pactus.wallet.set_address_label">
          <span class="rpc-badge"></span> pactus.wallet.set_address_label</a>
        </li>
        <li>
          <a href="#pactus.wallet.list_wallet">
          <span class="rpc-badge"></span> pactus.wallet.list_wallet</a>
        </li>
        <li>
          <a href="#pactus.wallet.get_wallet_info">
          <span class="rpc-badge"></span> pactus.wallet.get_wallet_info</a>
        </li>
        <li>
          <a href="#pactus.wallet.list_address">
          <span class="rpc-badge"></span> pactus.wallet.list_address</a>
        </li>
        </ul>
    </li>
    </ul>
</div>

<div class="api-doc">

## Transaction Service

<p>Transaction service defines various RPC methods for interacting with
transactions.</p>

### pactus.transaction.get_transaction <span id="pactus.transaction.get_transaction" class="rpc-badge"></span>

<p>GetTransaction retrieves transaction details based on the provided request
parameters.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the transaction to retrieve.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">verbosity</td>
    <td> numeric</td>
    <td>
    (Enum)The verbosity level for transaction details.
    <br>Available values:<ul>
      <li>TRANSACTION_DATA = 0 (Request transaction data only.)</li>
      <li>TRANSACTION_INFO = 1 (Request detailed transaction information.)</li>
      </ul>
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">block_height</td>
    <td> numeric</td>
    <td>
    The height of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">block_time</td>
    <td> numeric</td>
    <td>
    The UNIX timestamp of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">transaction</td>
    <td> object</td>
    <td>
    Detailed information about the transaction.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">transaction.id</td>
        <td> string</td>
        <td>
        The unique ID of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.data</td>
        <td> string</td>
        <td>
        The raw transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.version</td>
        <td> numeric</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.lock_time</td>
        <td> numeric</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.value</td>
        <td> numeric</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.fee</td>
        <td> numeric</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.payload_type</td>
        <td> numeric</td>
        <td>
        (Enum)The type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = 0 (Unknown payload type.)</li>
          <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
          <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
          <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
          <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
          <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.transfer</td>
        <td> object</td>
        <td>
        (OneOf)Transfer transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.transfer.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.amount</td>
            <td> numeric</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.bond</td>
        <td> object</td>
        <td>
        (OneOf)Bond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.bond.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.stake</td>
            <td> numeric</td>
            <td>
            The stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.public_key</td>
            <td> string</td>
            <td>
            The public key of the validator.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.sortition</td>
        <td> object</td>
        <td>
        (OneOf)Sortition transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.sortition.address</td>
            <td> string</td>
            <td>
            The validator address associated with the sortition proof.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.sortition.proof</td>
            <td> string</td>
            <td>
            The proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.unbond</td>
        <td> object</td>
        <td>
        (OneOf)Unbond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.unbond.validator</td>
            <td> string</td>
            <td>
            The address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.withdraw</td>
        <td> object</td>
        <td>
        (OneOf)Withdraw transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.withdraw.validator_address</td>
            <td> string</td>
            <td>
            The address of the validator to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.account_address</td>
            <td> string</td>
            <td>
            The address of the account to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.amount</td>
            <td> numeric</td>
            <td>
            The withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.memo</td>
        <td> string</td>
        <td>
        A memo string for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.public_key</td>
        <td> string</td>
        <td>
        The public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.signature</td>
        <td> string</td>
        <td>
        The signature for the transaction.
        </td>
      </tr>
         </tbody>
</table>

### pactus.transaction.calculate_fee <span id="pactus.transaction.calculate_fee" class="rpc-badge"></span>

<p>CalculateFee calculates the transaction fee based on the specified amount
and payload type.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    The amount involved in the transaction, specified in NanoPAC.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> numeric</td>
    <td>
    (Enum)The type of transaction payload.
    <br>Available values:<ul>
      <li>UNKNOWN = 0 (Unknown payload type.)</li>
      <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
      <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
      <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
      <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
      <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
      </ul>
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fixed_amount</td>
    <td> boolean</td>
    <td>
    Indicates if the amount should be fixed and include the fee.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    The calculated amount in NanoPAC.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    The calculated transaction fee in NanoPAC.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.broadcast_transaction <span id="pactus.transaction.broadcast_transaction" class="rpc-badge"></span>

<p>BroadcastTransaction broadcasts a signed transaction to the network.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">signed_raw_transaction</td>
    <td> string</td>
    <td>
    The signed raw transaction data to be broadcasted.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the broadcasted transaction.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.get_raw_transfer_transaction <span id="pactus.transaction.get_raw_transfer_transaction" class="rpc-badge"></span>

<p>GetRawTransferTransaction retrieves raw details of a transfer transaction.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> numeric</td>
    <td>
    The lock time for the transaction. If not set, defaults to the last block
height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">sender</td>
    <td> string</td>
    <td>
    The sender's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">receiver</td>
    <td> string</td>
    <td>
    The receiver's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    The amount to be transferred, specified in NanoPAC. Must be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    A memo string for the transaction.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the transaction.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.get_raw_bond_transaction <span id="pactus.transaction.get_raw_bond_transaction" class="rpc-badge"></span>

<p>GetRawBondTransaction retrieves raw details of a bond transaction.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> numeric</td>
    <td>
    The lock time for the transaction. If not set, defaults to the last block
height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">sender</td>
    <td> string</td>
    <td>
    The sender's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">receiver</td>
    <td> string</td>
    <td>
    The receiver's validator address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">stake</td>
    <td> numeric</td>
    <td>
    The stake amount in NanoPAC. Must be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    The public key of the validator.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    A memo string for the transaction.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the transaction.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.get_raw_unbond_transaction <span id="pactus.transaction.get_raw_unbond_transaction" class="rpc-badge"></span>

<p>GetRawUnbondTransaction retrieves raw details of an unbond transaction.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> numeric</td>
    <td>
    The lock time for the transaction. If not set, defaults to the last block
height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">validator_address</td>
    <td> string</td>
    <td>
    The address of the validator to unbond from.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    A memo string for the transaction.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the transaction.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.get_raw_withdraw_transaction <span id="pactus.transaction.get_raw_withdraw_transaction" class="rpc-badge"></span>

<p>GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> numeric</td>
    <td>
    The lock time for the transaction. If not set, defaults to the last block
height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">validator_address</td>
    <td> string</td>
    <td>
    The address of the validator to withdraw from.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">account_address</td>
    <td> string</td>
    <td>
    The address of the account to withdraw to.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    The withdrawal amount in NanoPAC. Must be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    A memo string for the transaction.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">id</td>
    <td> string</td>
    <td>
    The unique ID of the transaction.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.decode_raw_transaction <span id="pactus.transaction.decode_raw_transaction" class="rpc-badge"></span>

<p>DecodeRawTransaction accepts raw transaction and returnes decoded transaction.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">transaction</td>
    <td> object</td>
    <td>
    The decoded transaction.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">transaction.id</td>
        <td> string</td>
        <td>
        The unique ID of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.data</td>
        <td> string</td>
        <td>
        The raw transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.version</td>
        <td> numeric</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.lock_time</td>
        <td> numeric</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.value</td>
        <td> numeric</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.fee</td>
        <td> numeric</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.payload_type</td>
        <td> numeric</td>
        <td>
        (Enum)The type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = 0 (Unknown payload type.)</li>
          <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
          <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
          <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
          <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
          <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.transfer</td>
        <td> object</td>
        <td>
        (OneOf)Transfer transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.transfer.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.amount</td>
            <td> numeric</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.bond</td>
        <td> object</td>
        <td>
        (OneOf)Bond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.bond.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.stake</td>
            <td> numeric</td>
            <td>
            The stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.public_key</td>
            <td> string</td>
            <td>
            The public key of the validator.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.sortition</td>
        <td> object</td>
        <td>
        (OneOf)Sortition transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.sortition.address</td>
            <td> string</td>
            <td>
            The validator address associated with the sortition proof.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.sortition.proof</td>
            <td> string</td>
            <td>
            The proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.unbond</td>
        <td> object</td>
        <td>
        (OneOf)Unbond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.unbond.validator</td>
            <td> string</td>
            <td>
            The address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.withdraw</td>
        <td> object</td>
        <td>
        (OneOf)Withdraw transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.withdraw.validator_address</td>
            <td> string</td>
            <td>
            The address of the validator to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.account_address</td>
            <td> string</td>
            <td>
            The address of the account to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.amount</td>
            <td> numeric</td>
            <td>
            The withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.memo</td>
        <td> string</td>
        <td>
        A memo string for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.public_key</td>
        <td> string</td>
        <td>
        The public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.signature</td>
        <td> string</td>
        <td>
        The signature for the transaction.
        </td>
      </tr>
         </tbody>
</table>

## Blockchain Service

<p>Blockchain service defines RPC methods for interacting with the blockchain.</p>

### pactus.blockchain.get_block <span id="pactus.blockchain.get_block" class="rpc-badge"></span>

<p>GetBlock retrieves information about a block based on the provided request
parameters.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> numeric</td>
    <td>
    The height of the block to retrieve.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">verbosity</td>
    <td> numeric</td>
    <td>
    (Enum)The verbosity level for block information.
    <br>Available values:<ul>
      <li>BLOCK_DATA = 0 (Request only block data.)</li>
      <li>BLOCK_INFO = 1 (Request block information and transaction IDs.)</li>
      <li>BLOCK_TRANSACTIONS = 2 (Request block information and detailed transaction data.)</li>
      </ul>
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> numeric</td>
    <td>
    The height of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">hash</td>
    <td> string</td>
    <td>
    The hash of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">data</td>
    <td> string</td>
    <td>
    Block data, available only if verbosity level is set to BLOCK_DATA.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">block_time</td>
    <td> numeric</td>
    <td>
    The timestamp of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">header</td>
    <td> object</td>
    <td>
    Header information of the block.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">header.version</td>
        <td> numeric</td>
        <td>
        The version of the block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.prev_block_hash</td>
        <td> string</td>
        <td>
        The hash of the previous block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.state_root</td>
        <td> string</td>
        <td>
        The state root hash of the blockchain.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.sortition_seed</td>
        <td> string</td>
        <td>
        The sortition seed of the block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.proposer_address</td>
        <td> string</td>
        <td>
        The address of the proposer of the block.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">prev_cert</td>
    <td> object</td>
    <td>
    Certificate information of the previous block.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">prev_cert.hash</td>
        <td> string</td>
        <td>
        The hash of the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.round</td>
        <td> numeric</td>
        <td>
        The round of the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.committers</td>
        <td>repeated numeric</td>
        <td>
        List of committers in the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.absentees</td>
        <td>repeated numeric</td>
        <td>
        List of absentees in the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.signature</td>
        <td> string</td>
        <td>
        The signature of the certificate.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">txs</td>
    <td>repeated object</td>
    <td>
    List of transactions in the block, available when verbosity level is set to
BLOCK_TRANSACTIONS.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">txs[].id</td>
        <td> string</td>
        <td>
        The unique ID of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].data</td>
        <td> string</td>
        <td>
        The raw transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].version</td>
        <td> numeric</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> numeric</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> numeric</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> numeric</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> numeric</td>
        <td>
        (Enum)The type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = 0 (Unknown payload type.)</li>
          <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
          <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
          <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
          <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
          <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].transfer</td>
        <td> object</td>
        <td>
        (OneOf)Transfer transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].transfer.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.amount</td>
            <td> numeric</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> object</td>
        <td>
        (OneOf)Bond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].bond.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.stake</td>
            <td> numeric</td>
            <td>
            The stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.public_key</td>
            <td> string</td>
            <td>
            The public key of the validator.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].sortition</td>
        <td> object</td>
        <td>
        (OneOf)Sortition transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].sortition.address</td>
            <td> string</td>
            <td>
            The validator address associated with the sortition proof.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].sortition.proof</td>
            <td> string</td>
            <td>
            The proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].unbond</td>
        <td> object</td>
        <td>
        (OneOf)Unbond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].unbond.validator</td>
            <td> string</td>
            <td>
            The address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].withdraw</td>
        <td> object</td>
        <td>
        (OneOf)Withdraw transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].withdraw.validator_address</td>
            <td> string</td>
            <td>
            The address of the validator to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.account_address</td>
            <td> string</td>
            <td>
            The address of the account to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.amount</td>
            <td> numeric</td>
            <td>
            The withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].memo</td>
        <td> string</td>
        <td>
        A memo string for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].public_key</td>
        <td> string</td>
        <td>
        The public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].signature</td>
        <td> string</td>
        <td>
        The signature for the transaction.
        </td>
      </tr>
         </tbody>
</table>

### pactus.blockchain.get_block_hash <span id="pactus.blockchain.get_block_hash" class="rpc-badge"></span>

<p>GetBlockHash retrieves the hash of a block at the specified height.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> numeric</td>
    <td>
    The height of the block to retrieve the hash for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">hash</td>
    <td> string</td>
    <td>
    The hash of the block.
    </td>
  </tr>
     </tbody>
</table>

### pactus.blockchain.get_block_height <span id="pactus.blockchain.get_block_height" class="rpc-badge"></span>

<p>GetBlockHeight retrieves the height of a block with the specified hash.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">hash</td>
    <td> string</td>
    <td>
    The hash of the block to retrieve the height for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> numeric</td>
    <td>
    The height of the block.
    </td>
  </tr>
     </tbody>
</table>

### pactus.blockchain.get_blockchain_info <span id="pactus.blockchain.get_blockchain_info" class="rpc-badge"></span>

<p>GetBlockchainInfo retrieves general information about the blockchain.</p>

<h4>Parameters</h4>

Parameters has no fields.
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">last_block_height</td>
    <td> numeric</td>
    <td>
    The height of the last block in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">last_block_hash</td>
    <td> string</td>
    <td>
    The hash of the last block in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_accounts</td>
    <td> numeric</td>
    <td>
    The total number of accounts in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_validators</td>
    <td> numeric</td>
    <td>
    The total number of validators in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_power</td>
    <td> numeric</td>
    <td>
    The total power of the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">committee_power</td>
    <td> numeric</td>
    <td>
    The power of the committee.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">committee_validators</td>
    <td>repeated object</td>
    <td>
    List of committee validators.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">committee_validators[].hash</td>
        <td> string</td>
        <td>
        The hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].data</td>
        <td> string</td>
        <td>
        The serialized data of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].public_key</td>
        <td> string</td>
        <td>
        The public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].number</td>
        <td> numeric</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].stake</td>
        <td> numeric</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_bonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_sortition_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].unbonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator will unbond.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].address</td>
        <td> string</td>
        <td>
        The address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].availability_score</td>
        <td> numeric</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">is_pruned</td>
    <td> boolean</td>
    <td>
    If the blocks are subject to pruning.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">pruning_height</td>
    <td> numeric</td>
    <td>
    Lowest-height block stored (only present if pruning is enabled)
    </td>
  </tr>
     <tr>
    <td class="fw-bold">last_block_time</td>
    <td> numeric</td>
    <td>
    Timestamp of the last block in Unix format
    </td>
  </tr>
     </tbody>
</table>

### pactus.blockchain.get_consensus_info <span id="pactus.blockchain.get_consensus_info" class="rpc-badge"></span>

<p>GetConsensusInfo retrieves information about the consensus instances.</p>

<h4>Parameters</h4>

Parameters has no fields.
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">proposal</td>
    <td> object</td>
    <td>
    The proposal of the consensus info.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">proposal.height</td>
        <td> numeric</td>
        <td>
        The height of the proposal.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">proposal.round</td>
        <td> numeric</td>
        <td>
        The round of the proposal.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">proposal.block_data</td>
        <td> string</td>
        <td>
        The block data of the proposal.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">proposal.signature</td>
        <td> string</td>
        <td>
        The signature of the proposal, signed by the proposer.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">instances</td>
    <td>repeated object</td>
    <td>
    List of consensus instances.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">instances[].address</td>
        <td> string</td>
        <td>
        The address of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].active</td>
        <td> boolean</td>
        <td>
        Indicates whether the consensus instance is active and part of the
committee.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].height</td>
        <td> numeric</td>
        <td>
        The height of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].round</td>
        <td> numeric</td>
        <td>
        The round of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].votes</td>
        <td>repeated object</td>
        <td>
        List of votes in the consensus instance.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">instances[].votes[].type</td>
            <td> numeric</td>
            <td>
            (Enum)The type of the vote.
            <br>Available values:<ul>
              <li>VOTE_UNKNOWN = 0 (Unknown vote type.)</li>
              <li>VOTE_PREPARE = 1 (Prepare vote type.)</li>
              <li>VOTE_PRECOMMIT = 2 (Precommit vote type.)</li>
              <li>VOTE_CP_PRE_VOTE = 3 (Change-proposer:pre-vote vote type.)</li>
              <li>VOTE_CP_MAIN_VOTE = 4 (change-proposer:main-vote vote type.)</li>
              <li>VOTE_CP_DECIDED = 5 (change-proposer:decided vote type.)</li>
              </ul>
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].voter</td>
            <td> string</td>
            <td>
            The address of the voter.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].block_hash</td>
            <td> string</td>
            <td>
            The hash of the block being voted on.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].round</td>
            <td> numeric</td>
            <td>
            The consensus round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_round</td>
            <td> numeric</td>
            <td>
            The change-proposer round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_value</td>
            <td> numeric</td>
            <td>
            The change-proposer value of the vote.
            </td>
          </tr>
          </tbody>
</table>

### pactus.blockchain.get_account <span id="pactus.blockchain.get_account" class="rpc-badge"></span>

<p>GetAccount retrieves information about an account based on the provided
address.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The address of the account to retrieve information for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">account</td>
    <td> object</td>
    <td>
    Detailed information about the account.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">account.hash</td>
        <td> string</td>
        <td>
        The hash of the account.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.data</td>
        <td> string</td>
        <td>
        The serialized data of the account.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.number</td>
        <td> numeric</td>
        <td>
        The unique number assigned to the account.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.balance</td>
        <td> numeric</td>
        <td>
        The balance of the account in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.address</td>
        <td> string</td>
        <td>
        The address of the account.
        </td>
      </tr>
         </tbody>
</table>

### pactus.blockchain.get_validator <span id="pactus.blockchain.get_validator" class="rpc-badge"></span>

<p>GetValidator retrieves information about a validator based on the provided
address.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The address of the validator to retrieve information for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">validator</td>
    <td> object</td>
    <td>
    Detailed information about the validator.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">validator.hash</td>
        <td> string</td>
        <td>
        The hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.data</td>
        <td> string</td>
        <td>
        The serialized data of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.public_key</td>
        <td> string</td>
        <td>
        The public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.number</td>
        <td> numeric</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> numeric</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator will unbond.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.address</td>
        <td> string</td>
        <td>
        The address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.availability_score</td>
        <td> numeric</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         </tbody>
</table>

### pactus.blockchain.get_validator_by_number <span id="pactus.blockchain.get_validator_by_number" class="rpc-badge"></span>

<p>GetValidatorByNumber retrieves information about a validator based on the
provided number.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">number</td>
    <td> numeric</td>
    <td>
    The unique number of the validator to retrieve information for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">validator</td>
    <td> object</td>
    <td>
    Detailed information about the validator.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">validator.hash</td>
        <td> string</td>
        <td>
        The hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.data</td>
        <td> string</td>
        <td>
        The serialized data of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.public_key</td>
        <td> string</td>
        <td>
        The public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.number</td>
        <td> numeric</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> numeric</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> numeric</td>
        <td>
        The height at which the validator will unbond.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.address</td>
        <td> string</td>
        <td>
        The address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.availability_score</td>
        <td> numeric</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         </tbody>
</table>

### pactus.blockchain.get_validator_addresses <span id="pactus.blockchain.get_validator_addresses" class="rpc-badge"></span>

<p>GetValidatorAddresses retrieves a list of all validator addresses.</p>

<h4>Parameters</h4>

Parameters has no fields.
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">addresses</td>
    <td>repeated string</td>
    <td>
    List of validator addresses.
    </td>
  </tr>
     </tbody>
</table>

### pactus.blockchain.get_public_key <span id="pactus.blockchain.get_public_key" class="rpc-badge"></span>

<p>GetPublicKey retrieves the public key of an account based on the provided
address.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The address for which to retrieve the public key.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    The public key associated with the provided address.
    </td>
  </tr>
     </tbody>
</table>

### pactus.blockchain.get_tx_pool_content <span id="pactus.blockchain.get_tx_pool_content" class="rpc-badge"></span>

<p>GetTxPoolContent retrieves current transactions in the transaction pool.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> numeric</td>
    <td>
    (Enum)The type of transactions to retrieve from the transaction pool. 0 means all
types.
    <br>Available values:<ul>
      <li>UNKNOWN = 0 (Unknown payload type.)</li>
      <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
      <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
      <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
      <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
      <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
      </ul>
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">txs</td>
    <td>repeated object</td>
    <td>
    List of transactions currently in the pool.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">txs[].id</td>
        <td> string</td>
        <td>
        The unique ID of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].data</td>
        <td> string</td>
        <td>
        The raw transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].version</td>
        <td> numeric</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> numeric</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> numeric</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> numeric</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> numeric</td>
        <td>
        (Enum)The type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = 0 (Unknown payload type.)</li>
          <li>TRANSFER_PAYLOAD = 1 (Transfer payload type.)</li>
          <li>BOND_PAYLOAD = 2 (Bond payload type.)</li>
          <li>SORTITION_PAYLOAD = 3 (Sortition payload type.)</li>
          <li>UNBOND_PAYLOAD = 4 (Unbond payload type.)</li>
          <li>WITHDRAW_PAYLOAD = 5 (Withdraw payload type.)</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].transfer</td>
        <td> object</td>
        <td>
        (OneOf)Transfer transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].transfer.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.amount</td>
            <td> numeric</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> object</td>
        <td>
        (OneOf)Bond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].bond.sender</td>
            <td> string</td>
            <td>
            The sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.receiver</td>
            <td> string</td>
            <td>
            The receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.stake</td>
            <td> numeric</td>
            <td>
            The stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.public_key</td>
            <td> string</td>
            <td>
            The public key of the validator.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].sortition</td>
        <td> object</td>
        <td>
        (OneOf)Sortition transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].sortition.address</td>
            <td> string</td>
            <td>
            The validator address associated with the sortition proof.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].sortition.proof</td>
            <td> string</td>
            <td>
            The proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].unbond</td>
        <td> object</td>
        <td>
        (OneOf)Unbond transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].unbond.validator</td>
            <td> string</td>
            <td>
            The address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].withdraw</td>
        <td> object</td>
        <td>
        (OneOf)Withdraw transaction payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].withdraw.validator_address</td>
            <td> string</td>
            <td>
            The address of the validator to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.account_address</td>
            <td> string</td>
            <td>
            The address of the account to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.amount</td>
            <td> numeric</td>
            <td>
            The withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].memo</td>
        <td> string</td>
        <td>
        A memo string for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].public_key</td>
        <td> string</td>
        <td>
        The public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].signature</td>
        <td> string</td>
        <td>
        The signature for the transaction.
        </td>
      </tr>
         </tbody>
</table>

## Network Service

<p>Network service provides RPCs for retrieving information about the network.</p>

### pactus.network.get_network_info <span id="pactus.network.get_network_info" class="rpc-badge"></span>

<p>GetNetworkInfo retrieves information about the overall network.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">only_connected</td>
    <td> boolean</td>
    <td>
    If true, returns only peers that are currently connected.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">network_name</td>
    <td> string</td>
    <td>
    Name of the network.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connected_peers_count</td>
    <td> numeric</td>
    <td>
    Number of connected peers.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connected_peers</td>
    <td>repeated object</td>
    <td>
    List of connected peers.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">connected_peers[].status</td>
        <td> numeric</td>
        <td>
        Current status of the peer (e.g., connected, disconnected).
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].moniker</td>
        <td> string</td>
        <td>
        Moniker of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].agent</td>
        <td> string</td>
        <td>
        Version and agent details of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].peer_id</td>
        <td> string</td>
        <td>
        Peer ID of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].consensus_keys</td>
        <td>repeated string</td>
        <td>
        List of consensus keys used by the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].consensus_addresses</td>
        <td>repeated string</td>
        <td>
        List of consensus addresses used by the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].services</td>
        <td> numeric</td>
        <td>
        Bitfield representing the services provided by the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_block_hash</td>
        <td> string</td>
        <td>
        Hash of the last block the peer knows.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].height</td>
        <td> numeric</td>
        <td>
        Blockchain height of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_sent</td>
        <td> numeric</td>
        <td>
        Time the last bundle sent to the peer (in epoch format).
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_received</td>
        <td> numeric</td>
        <td>
        Time the last bundle received from the peer (in epoch format).
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].address</td>
        <td> string</td>
        <td>
        Network address of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].direction</td>
        <td> string</td>
        <td>
        Connection direction (e.g., inbound, outbound).
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].protocols</td>
        <td>repeated string</td>
        <td>
        List of protocols supported by the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].total_sessions</td>
        <td> numeric</td>
        <td>
        Total download sessions with the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].completed_sessions</td>
        <td> numeric</td>
        <td>
        Completed download sessions with the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].metric_info</td>
        <td> object</td>
        <td>
        Metrics related to peer activity.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">connected_peers[].metric_info.TotalInvalid</td>
            <td> object</td>
            <td>
            Total number of invalid bundles.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">connected_peers[].metric_info.TotalSent</td>
            <td> object</td>
            <td>
            Total number of bundles sent.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">connected_peers[].metric_info.TotalReceived</td>
            <td> object</td>
            <td>
            Total number of bundles received.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">connected_peers[].metric_info.MessageSent</td>
            <td> object</td>
            <td>
            Number of sent bundles categorized by message type.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">connected_peers[].metric_info.MessageReceived</td>
            <td> object</td>
            <td>
            Number of received bundles categorized by message type.
            </td>
          </tr>
          <tr>
    <td class="fw-bold">metric_info</td>
    <td> object</td>
    <td>
    Metrics related to node activity.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">metric_info.TotalInvalid</td>
        <td> object</td>
        <td>
        Total number of invalid bundles.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">metric_info.TotalInvalid.Bytes</td>
            <td> numeric</td>
            <td>
            Total number of bytes.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">metric_info.TotalInvalid.Bundles</td>
            <td> numeric</td>
            <td>
            Total number of bundles.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">metric_info.TotalSent</td>
        <td> object</td>
        <td>
        Total number of bundles sent.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">metric_info.TotalSent.Bytes</td>
            <td> numeric</td>
            <td>
            Total number of bytes.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">metric_info.TotalSent.Bundles</td>
            <td> numeric</td>
            <td>
            Total number of bundles.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">metric_info.TotalReceived</td>
        <td> object</td>
        <td>
        Total number of bundles received.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">metric_info.TotalReceived.Bytes</td>
            <td> numeric</td>
            <td>
            Total number of bytes.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">metric_info.TotalReceived.Bundles</td>
            <td> numeric</td>
            <td>
            Total number of bundles.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">metric_info.MessageSent</td>
        <td> object</td>
        <td>
        Number of sent bundles categorized by message type.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">metric_info.MessageReceived</td>
        <td> object</td>
        <td>
        Number of received bundles categorized by message type.
        </td>
      </tr>
         </tbody>
</table>

### pactus.network.get_node_info <span id="pactus.network.get_node_info" class="rpc-badge"></span>

<p>GetNodeInfo retrieves information about a specific node in the network.</p>

<h4>Parameters</h4>

Parameters has no fields.
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">moniker</td>
    <td> string</td>
    <td>
    Moniker of the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">agent</td>
    <td> string</td>
    <td>
    Version and agent details of the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">peer_id</td>
    <td> string</td>
    <td>
    Peer ID of the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">started_at</td>
    <td> numeric</td>
    <td>
    Time the node was started (in epoch format).
    </td>
  </tr>
     <tr>
    <td class="fw-bold">reachability</td>
    <td> string</td>
    <td>
    Reachability status of the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">services</td>
    <td> numeric</td>
    <td>
    Bitfield representing the services provided by the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">services_names</td>
    <td> string</td>
    <td>
    Names of services provided by the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">local_addrs</td>
    <td>repeated string</td>
    <td>
    List of addresses associated with the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">protocols</td>
    <td>repeated string</td>
    <td>
    List of protocols supported by the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">clock_offset</td>
    <td> numeric</td>
    <td>
    Offset between the node's clock and the network's clock (in seconds).
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connection_info</td>
    <td> object</td>
    <td>
    Information about the node's connections.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">connection_info.connections</td>
        <td> numeric</td>
        <td>
        Total number of connections.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connection_info.inbound_connections</td>
        <td> numeric</td>
        <td>
        Number of inbound connections.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connection_info.outbound_connections</td>
        <td> numeric</td>
        <td>
        Number of outbound connections.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">zmq_publishers</td>
    <td>repeated object</td>
    <td>
    List of active ZeroMQ publishers.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">zmq_publishers[].topic</td>
        <td> string</td>
        <td>
        The topic associated with the publisher.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">zmq_publishers[].address</td>
        <td> string</td>
        <td>
        The address of the publisher.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">zmq_publishers[].hwm</td>
        <td> numeric</td>
        <td>
        The high-water mark (HWM) for the publisher, indicating the
maximum number of messages to queue before dropping older ones.
        </td>
      </tr>
         </tbody>
</table>

## Utils Service

<p>Utils service defines RPC methods for utility functions such as message
signing and verification.</p>

### pactus.utils.sign_message_with_private_key <span id="pactus.utils.sign_message_with_private_key" class="rpc-badge"></span>

<p>SignMessageWithPrivateKey signs message with provided private key.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">private_key</td>
    <td> string</td>
    <td>
    The private key to sign the message.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">message</td>
    <td> string</td>
    <td>
    The message to sign.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">signature</td>
    <td> string</td>
    <td>
    The signature of the message.
    </td>
  </tr>
     </tbody>
</table>

### pactus.utils.verify_message <span id="pactus.utils.verify_message" class="rpc-badge"></span>

<p>VerifyMessage verifies signature with public key and message.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">message</td>
    <td> string</td>
    <td>
    The signed message.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">signature</td>
    <td> string</td>
    <td>
    The signature of the message.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    The public key of the signer.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">is_valid</td>
    <td> boolean</td>
    <td>
    Indicates if the signature is valid (true) or not (false).
    </td>
  </tr>
     </tbody>
</table>

### pactus.utils.b_l_s_public_key_aggregation <span id="pactus.utils.b_l_s_public_key_aggregation" class="rpc-badge"></span>

<p>BLSPublicKeyAggregation aggregates bls public keys.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">public_keys</td>
    <td>repeated string</td>
    <td>
    The public keys to aggregate.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    The aggregated public key.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The aggregated public key account address.
    </td>
  </tr>
     </tbody>
</table>

### pactus.utils.b_l_s_signature_aggregation <span id="pactus.utils.b_l_s_signature_aggregation" class="rpc-badge"></span>

<p>BLSSignatureAggregation aggregates bls signatures.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">signatures</td>
    <td>repeated string</td>
    <td>
    The signatures to aggregate.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">signature</td>
    <td> string</td>
    <td>
    The aggregated signature.
    </td>
  </tr>
     </tbody>
</table>

## Wallet Service

<p>Define the Wallet service with various RPC methods for wallet management.</p>

### pactus.wallet.create_wallet <span id="pactus.wallet.create_wallet" class="rpc-badge"></span>

<p>CreateWallet creates a new wallet with the specified parameters.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the new wallet.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    The password for securing the wallet.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">mnemonic</td>
    <td> string</td>
    <td>
    The mnemonic for wallet recovery.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.restore_wallet <span id="pactus.wallet.restore_wallet" class="rpc-badge"></span>

<p>RestoreWallet restores an existing wallet with the given mnemonic.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to restore.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">mnemonic</td>
    <td> string</td>
    <td>
    The mnemonic for wallet recovery.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    The password for securing the wallet.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the restored wallet.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.load_wallet <span id="pactus.wallet.load_wallet" class="rpc-badge"></span>

<p>LoadWallet loads an existing wallet with the given name.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to load.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the loaded wallet.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.unload_wallet <span id="pactus.wallet.unload_wallet" class="rpc-badge"></span>

<p>UnloadWallet unloads a currently loaded wallet with the specified name.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to unload.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the unloaded wallet.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_total_balance <span id="pactus.wallet.get_total_balance" class="rpc-badge"></span>

<p>GetTotalBalance returns the total available balance of the wallet.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to get the total balance.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_balance</td>
    <td> numeric</td>
    <td>
    The total balance of the wallet in NanoPAC.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.sign_raw_transaction <span id="pactus.wallet.sign_raw_transaction" class="rpc-badge"></span>

<p>SignRawTransaction signs a raw transaction for a specified wallet.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet used for signing.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    The raw transaction data to be signed.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    The password for unlocking the wallet for signing.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">transaction_id</td>
    <td> string</td>
    <td>
    The ID of the signed transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">signed_raw_transaction</td>
    <td> string</td>
    <td>
    The signed raw transaction data.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_validator_address <span id="pactus.wallet.get_validator_address" class="rpc-badge"></span>

<p>GetValidatorAddress retrieves the validator address associated with a
public key.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    The public key for which the validator address is requested.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The validator address associated with the public key.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_new_address <span id="pactus.wallet.get_new_address" class="rpc-badge"></span>

<p>GetNewAddress generates a new address for the specified wallet.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to generate a new address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address_type</td>
    <td> numeric</td>
    <td>
    (Enum)The type of address to generate.
    <br>Available values:<ul>
      <li>ADDRESS_TYPE_TREASURY = 0 (Treasury address type.
Should not be used to generate new addresses.)</li>
      <li>ADDRESS_TYPE_VALIDATOR = 1 (Validator address type.)</li>
      <li>ADDRESS_TYPE_BLS_ACCOUNT = 2 (Account address type with BLS signature scheme.)</li>
      <li>ADDRESS_TYPE_ED25519_ACCOUNT = 3 (Account address type with Ed25519 signature scheme.
Note: Generating a new Ed25519 address requires the wallet password.)</li>
      </ul>
    </td>
  </tr>
  <tr>
    <td class="fw-bold">label</td>
    <td> string</td>
    <td>
    A label for the new address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    Password for the new address. It's required when address_type is ADDRESS_TYPE_ED25519_ACCOUNT.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet from which the address is generated.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">address_info</td>
    <td> object</td>
    <td>
    Information about the newly generated address.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">address_info.address</td>
        <td> string</td>
        <td>
        The address string.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">address_info.public_key</td>
        <td> string</td>
        <td>
        The public key associated with the address.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">address_info.label</td>
        <td> string</td>
        <td>
        A label associated with the address.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">address_info.path</td>
        <td> string</td>
        <td>
        The Hierarchical Deterministic path of the address within the wallet.
        </td>
      </tr>
         </tbody>
</table>

### pactus.wallet.get_address_history <span id="pactus.wallet.get_address_history" class="rpc-badge"></span>

<p>GetAddressHistory retrieves the transaction history of an address.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The address to retrieve the transaction history for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">history_info</td>
    <td>repeated object</td>
    <td>
    Array of history information for the address.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">history_info[].transaction_id</td>
        <td> string</td>
        <td>
        The transaction ID hash.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].time</td>
        <td> numeric</td>
        <td>
        The timestamp of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].payload_type</td>
        <td> string</td>
        <td>
        The payload type of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].description</td>
        <td> string</td>
        <td>
        A description of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].amount</td>
        <td> numeric</td>
        <td>
        The amount involved in the transaction.
        </td>
      </tr>
         </tbody>
</table>

### pactus.wallet.sign_message <span id="pactus.wallet.sign_message" class="rpc-badge"></span>

<p>SignMessage signs an arbitrary message.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    The password for unlocking the wallet for signing.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    The account address associated with the private key.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">message</td>
    <td> string</td>
    <td>
    The arbitrary message to be signed.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">signature</td>
    <td> string</td>
    <td>
    Signature of the message.
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_total_stake <span id="pactus.wallet.get_total_stake" class="rpc-badge"></span>

<p>GetTotalStake return total stake of wallet.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet.
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">total_stake</td>
    <td> numeric</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_address_info <span id="pactus.wallet.get_address_info" class="rpc-badge"></span>

<p>GetAddressInfo return address information.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to generate a new address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">label</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">path</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.set_address_label <span id="pactus.wallet.set_address_label" class="rpc-badge"></span>

<p>SetAddressLabel set label for given address.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    The name of the wallet to generate a new address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    The password for unlocking the wallet for signing.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
  <tr>
    <td class="fw-bold">label</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  </tbody>
</table>

### pactus.wallet.list_wallet <span id="pactus.wallet.list_wallet" class="rpc-badge"></span>

<p>ListWallet return list wallet name.</p>

<h4>Parameters</h4>

Parameters has no fields.
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallets</td>
    <td>repeated string</td>
    <td>
    
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.get_wallet_info <span id="pactus.wallet.get_wallet_info" class="rpc-badge"></span>

<p>GetWalletInfo return wallet information.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">version</td>
    <td> numeric</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">network</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">encrypted</td>
    <td> boolean</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">uuid</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">created_at</td>
    <td> numeric</td>
    <td>
    
    </td>
  </tr>
     </tbody>
</table>

### pactus.wallet.list_address <span id="pactus.wallet.list_address" class="rpc-badge"></span>

<p>ListAddress return list address in wallet.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">wallet_name</td>
    <td> string</td>
    <td>
    
    </td>
  </tr>
  </tbody>
</table>
  <h4>Result</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">data</td>
    <td>repeated object</td>
    <td>
    
    </td>
  </tr>
     <tr>
        <td class="fw-bold">data[].address</td>
        <td> string</td>
        <td>
        The address string.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">data[].public_key</td>
        <td> string</td>
        <td>
        The public key associated with the address.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">data[].label</td>
        <td> string</td>
        <td>
        A label associated with the address.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">data[].path</td>
        <td> string</td>
        <td>
        The Hierarchical Deterministic path of the address within the wallet.
        </td>
      </tr>
         </tbody>
</table>
