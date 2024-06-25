---
title: JSON-RPC API Reference
weight: 2
---

Each node in the Pactus network can be configured to use the [gRPC](https://grpc.io/) protocol for communication.
Here, you can find the list of all gRPC methods and messages.

All the amounts and values in gRPC endpoints are in NanoPAC units,
which are atomic and the smallest unit in the Pactus blockchain.
Each PAC is equivalent to 1,000,000,000 or 10<sup>9</sup> NanoPACs.

All binary data, including hash and transaction data,
are encoded using the [base64](https://en.wikipedia.org/wiki/Base64) decoder.

<h2>JSON-RPC Services</h2>

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
          <a href="#pactus.transaction.get_transaction_pool">
          <span class="rpc-badge"></span> pactus.transaction.get_transaction_pool</a>
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
    Transaction ID.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">verbosity</td>
    <td> string</td>
    <td>
    (Enum) Verbosity level for transaction details.
    <br>Available values:<ul>
      <li>TRANSACTION_DATA = Request transaction data only.</li>
      <li>TRANSACTION_INFO = Request transaction details.</li>
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
    Height of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">block_time</td>
    <td> numeric</td>
    <td>
    Time of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">transaction</td>
    <td> object</td>
    <td>
    Information about the transaction.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">transaction.id</td>
        <td> string</td>
        <td>
        Transaction ID.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.data</td>
        <td> string</td>
        <td>
        Transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.version</td>
        <td> numeric</td>
        <td>
        Transaction version.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.lock_time</td>
        <td> numeric</td>
        <td>
        Lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.value</td>
        <td> numeric</td>
        <td>
        Transaction value in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.fee</td>
        <td> numeric</td>
        <td>
        Transaction fee in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.payload_type</td>
        <td> string</td>
        <td>
        (Enum) Type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = Unknown payload type.</li>
          <li>TRANSFER_PAYLOAD = Transfer payload type.</li>
          <li>BOND_PAYLOAD = Bond payload type.</li>
          <li>SORTITION_PAYLOAD = Sortition payload type.</li>
          <li>UNBOND_PAYLOAD = Unbond payload type.</li>
          <li>WITHDRAW_PAYLOAD = Withdraw payload type.</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.transfer</td>
        <td> object</td>
        <td>
        (OneOf) Transfer payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.transfer.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.transfer.amount</td>
            <td> numeric</td>
            <td>
            Transaction amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.bond</td>
        <td> object</td>
        <td>
        (OneOf) Bond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.bond.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.bond.stake</td>
            <td> numeric</td>
            <td>
            Stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.sortition</td>
        <td> object</td>
        <td>
        (OneOf) Sortition payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.sortition.address</td>
            <td> string</td>
            <td>
            Address associated with the sortition.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.sortition.proof</td>
            <td> string</td>
            <td>
            Proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.unbond</td>
        <td> object</td>
        <td>
        (OneOf) Unbond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.unbond.validator</td>
            <td> string</td>
            <td>
            Address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.withdraw</td>
        <td> object</td>
        <td>
        (OneOf) Withdraw payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">transaction.withdraw.from</td>
            <td> string</td>
            <td>
            Address to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.to</td>
            <td> string</td>
            <td>
            Address to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">transaction.withdraw.amount</td>
            <td> numeric</td>
            <td>
            Withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.memo</td>
        <td> string</td>
        <td>
        Transaction memo.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.public_key</td>
        <td> string</td>
        <td>
        Public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.signature</td>
        <td> string</td>
        <td>
        Transaction signature.
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
    Transaction amount in NanoPAC.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> string</td>
    <td>
    (Enum) Type of transaction payload.
    <br>Available values:<ul>
      <li>UNKNOWN = Unknown payload type.</li>
      <li>TRANSFER_PAYLOAD = Transfer payload type.</li>
      <li>BOND_PAYLOAD = Bond payload type.</li>
      <li>SORTITION_PAYLOAD = Sortition payload type.</li>
      <li>UNBOND_PAYLOAD = Unbond payload type.</li>
      <li>WITHDRAW_PAYLOAD = Withdraw payload type.</li>
      </ul>
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fixed_amount</td>
    <td> boolean</td>
    <td>
    Indicates that amount should be fixed and includes the fee.
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
    Calculated amount in NanoPAC.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    Calculated transaction fee in NanoPAC.
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
    Signed raw transaction data.
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
    Transaction ID.
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
    Lock time for the transaction.
If not explicitly set, it sets to the last block height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">sender</td>
    <td> string</td>
    <td>
    Sender's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">receiver</td>
    <td> string</td>
    <td>
    Receiver's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    Transfer amount in NanoPAC.
It should be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    Transaction fee in NanoPAC.
If not explicitly set, it is calculated based on the amount.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    Transaction memo.
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
    Raw transaction data.
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
    Lock time for the transaction.
If not explicitly set, it sets to the last block height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">sender</td>
    <td> string</td>
    <td>
    Sender's account address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">receiver</td>
    <td> string</td>
    <td>
    Receiver's validator address.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">stake</td>
    <td> numeric</td>
    <td>
    Stake amount in NanoPAC.
It should be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">public_key</td>
    <td> string</td>
    <td>
    Public key of the validator.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    Transaction fee in NanoPAC.
If not explicitly set, it is calculated based on the stake.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    Transaction memo.
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
    Raw transaction data.
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
    Lock time for the transaction.
If not explicitly set, it sets to the last block height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">validator_address</td>
    <td> string</td>
    <td>
    Address of the validator to unbond from.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    Transaction memo.
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
    Raw transaction data.
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
    Lock time for the transaction.
If not explicitly set, it sets to the last block height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">validator_address</td>
    <td> string</td>
    <td>
    Address of the validator to withdraw from.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">account_address</td>
    <td> string</td>
    <td>
    Address of the account to withdraw to.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">amount</td>
    <td> numeric</td>
    <td>
    Withdrawal amount in NanoPAC.
It should be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> numeric</td>
    <td>
    Transaction fee in NanoPAC.
If not explicitly set, it is calculated based on the amount.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    Transaction memo.
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
    Raw transaction data.
    </td>
  </tr>
     </tbody>
</table>

### pactus.transaction.get_transaction_pool <span id="pactus.transaction.get_transaction_pool" class="rpc-badge"></span>

<p>GetTransactionPool retrieves current transactions on the TXPool.</p>

<h4>Parameters</h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> string</td>
    <td>
    (Enum) Payload type of tranactions in the tx pool, 0 is all types.
    <br>Available values:<ul>
      <li>UNKNOWN = Unknown payload type.</li>
      <li>TRANSFER_PAYLOAD = Transfer payload type.</li>
      <li>BOND_PAYLOAD = Bond payload type.</li>
      <li>SORTITION_PAYLOAD = Sortition payload type.</li>
      <li>UNBOND_PAYLOAD = Unbond payload type.</li>
      <li>WITHDRAW_PAYLOAD = Withdraw payload type.</li>
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
    List of the transaction in the pool.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">txs[].id</td>
        <td> string</td>
        <td>
        Transaction ID.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].data</td>
        <td> string</td>
        <td>
        Transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].version</td>
        <td> numeric</td>
        <td>
        Transaction version.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> numeric</td>
        <td>
        Lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> numeric</td>
        <td>
        Transaction value in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> numeric</td>
        <td>
        Transaction fee in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> string</td>
        <td>
        (Enum) Type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = Unknown payload type.</li>
          <li>TRANSFER_PAYLOAD = Transfer payload type.</li>
          <li>BOND_PAYLOAD = Bond payload type.</li>
          <li>SORTITION_PAYLOAD = Sortition payload type.</li>
          <li>UNBOND_PAYLOAD = Unbond payload type.</li>
          <li>WITHDRAW_PAYLOAD = Withdraw payload type.</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].transfer</td>
        <td> object</td>
        <td>
        (OneOf) Transfer payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].transfer.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.amount</td>
            <td> numeric</td>
            <td>
            Transaction amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> object</td>
        <td>
        (OneOf) Bond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].bond.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.stake</td>
            <td> numeric</td>
            <td>
            Stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].sortition</td>
        <td> object</td>
        <td>
        (OneOf) Sortition payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].sortition.address</td>
            <td> string</td>
            <td>
            Address associated with the sortition.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].sortition.proof</td>
            <td> string</td>
            <td>
            Proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].unbond</td>
        <td> object</td>
        <td>
        (OneOf) Unbond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].unbond.validator</td>
            <td> string</td>
            <td>
            Address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].withdraw</td>
        <td> object</td>
        <td>
        (OneOf) Withdraw payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].withdraw.from</td>
            <td> string</td>
            <td>
            Address to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.to</td>
            <td> string</td>
            <td>
            Address to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.amount</td>
            <td> numeric</td>
            <td>
            Withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].memo</td>
        <td> string</td>
        <td>
        Transaction memo.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].public_key</td>
        <td> string</td>
        <td>
        Public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].signature</td>
        <td> string</td>
        <td>
        Transaction signature.
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
    Height of the block.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">verbosity</td>
    <td> string</td>
    <td>
    (Enum) Verbosity level for block information.
    <br>Available values:<ul>
      <li>BLOCK_DATA = Request block data only.</li>
      <li>BLOCK_INFO = Request block information and transaction IDs.</li>
      <li>BLOCK_TRANSACTIONS = Request block information and transaction details.</li>
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
    Height of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">hash</td>
    <td> string</td>
    <td>
    Hash of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">data</td>
    <td> string</td>
    <td>
    Block data, only available if the verbosity level is set to BLOCK_DATA.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">block_time</td>
    <td> numeric</td>
    <td>
    Block timestamp.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">header</td>
    <td> object</td>
    <td>
    Block header information.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">header.version</td>
        <td> numeric</td>
        <td>
        Block version.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.prev_block_hash</td>
        <td> string</td>
        <td>
        Hash of the previous block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.state_root</td>
        <td> string</td>
        <td>
        State root of the block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.sortition_seed</td>
        <td> string</td>
        <td>
        Sortition seed of the block.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">header.proposer_address</td>
        <td> string</td>
        <td>
        Address of the proposer of the block.
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
        Hash of the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.round</td>
        <td> numeric</td>
        <td>
        Round of the certificate.
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
        Certificate signature.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">txs</td>
    <td>repeated object</td>
    <td>
    List of transactions in the block.
Transaction information is available when the verbosity level is set to BLOCK_TRANSACTIONS.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">txs[].id</td>
        <td> string</td>
        <td>
        Transaction ID.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].data</td>
        <td> string</td>
        <td>
        Transaction data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].version</td>
        <td> numeric</td>
        <td>
        Transaction version.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> numeric</td>
        <td>
        Lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> numeric</td>
        <td>
        Transaction value in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> numeric</td>
        <td>
        Transaction fee in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> string</td>
        <td>
        (Enum) Type of transaction payload.
        <br>Available values:<ul>
          <li>UNKNOWN = Unknown payload type.</li>
          <li>TRANSFER_PAYLOAD = Transfer payload type.</li>
          <li>BOND_PAYLOAD = Bond payload type.</li>
          <li>SORTITION_PAYLOAD = Sortition payload type.</li>
          <li>UNBOND_PAYLOAD = Unbond payload type.</li>
          <li>WITHDRAW_PAYLOAD = Withdraw payload type.</li>
          </ul>
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].transfer</td>
        <td> object</td>
        <td>
        (OneOf) Transfer payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].transfer.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].transfer.amount</td>
            <td> numeric</td>
            <td>
            Transaction amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> object</td>
        <td>
        (OneOf) Bond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].bond.sender</td>
            <td> string</td>
            <td>
            Sender's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.receiver</td>
            <td> string</td>
            <td>
            Receiver's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].bond.stake</td>
            <td> numeric</td>
            <td>
            Stake amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].sortition</td>
        <td> object</td>
        <td>
        (OneOf) Sortition payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].sortition.address</td>
            <td> string</td>
            <td>
            Address associated with the sortition.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].sortition.proof</td>
            <td> string</td>
            <td>
            Proof for the sortition.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].unbond</td>
        <td> object</td>
        <td>
        (OneOf) Unbond payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].unbond.validator</td>
            <td> string</td>
            <td>
            Address of the validator to unbond from.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].withdraw</td>
        <td> object</td>
        <td>
        (OneOf) Withdraw payload.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">txs[].withdraw.from</td>
            <td> string</td>
            <td>
            Address to withdraw from.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.to</td>
            <td> string</td>
            <td>
            Address to withdraw to.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">txs[].withdraw.amount</td>
            <td> numeric</td>
            <td>
            Withdrawal amount in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].memo</td>
        <td> string</td>
        <td>
        Transaction memo.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].public_key</td>
        <td> string</td>
        <td>
        Public key associated with the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].signature</td>
        <td> string</td>
        <td>
        Transaction signature.
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
    Height of the block.
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
    Hash of the block.
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
    Hash of the block.
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
    Height of the block.
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
    Height of the last block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">last_block_hash</td>
    <td> string</td>
    <td>
    Hash of the last block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_accounts</td>
    <td> numeric</td>
    <td>
    Total number of accounts.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_validators</td>
    <td> numeric</td>
    <td>
    Total number of validators.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_power</td>
    <td> numeric</td>
    <td>
    Total power in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">committee_power</td>
    <td> numeric</td>
    <td>
    Power of the committee.
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
        Hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].data</td>
        <td> string</td>
        <td>
        Validator data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].public_key</td>
        <td> string</td>
        <td>
        Public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].number</td>
        <td> numeric</td>
        <td>
        Validator number.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].stake</td>
        <td> numeric</td>
        <td>
        Validator stake in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_bonding_height</td>
        <td> numeric</td>
        <td>
        Last bonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_sortition_height</td>
        <td> numeric</td>
        <td>
        Last sortition height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].unbonding_height</td>
        <td> numeric</td>
        <td>
        Unbonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].address</td>
        <td> string</td>
        <td>
        Address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].availability_score</td>
        <td> numeric</td>
        <td>
        Availability score of the validator.
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
        Address of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].Active</td>
        <td> boolean</td>
        <td>
        Whether the consensus instance is active.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].height</td>
        <td> numeric</td>
        <td>
        Height of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].round</td>
        <td> numeric</td>
        <td>
        Round of the consensus instance.
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
            <td> string</td>
            <td>
            (Enum) Type of the vote.
            <br>Available values:<ul>
              <li>VOTE_UNKNOWN = Unknown vote type.</li>
              <li>VOTE_PREPARE = Prepare vote type.</li>
              <li>VOTE_PRECOMMIT = Precommit vote type.</li>
              <li>VOTE_CHANGE_PROPOSER = Change proposer vote type.</li>
              </ul>
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].voter</td>
            <td> string</td>
            <td>
            Voter's address.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].block_hash</td>
            <td> string</td>
            <td>
            Hash of the block being voted on.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].round</td>
            <td> numeric</td>
            <td>
            Round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_round</td>
            <td> numeric</td>
            <td>
            Consensus round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_value</td>
            <td> numeric</td>
            <td>
            Consensus value of the vote.
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
    Address of the account.
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
    Account information.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">account.hash</td>
        <td> string</td>
        <td>
        Hash of the account.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.data</td>
        <td> string</td>
        <td>
        Account data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.number</td>
        <td> numeric</td>
        <td>
        Account number.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.balance</td>
        <td> numeric</td>
        <td>
        Account balance in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.address</td>
        <td> string</td>
        <td>
        Address of the account.
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
    Address of the validator.
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
    Validator information.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">validator.hash</td>
        <td> string</td>
        <td>
        Hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.data</td>
        <td> string</td>
        <td>
        Validator data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.public_key</td>
        <td> string</td>
        <td>
        Public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.number</td>
        <td> numeric</td>
        <td>
        Validator number.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> numeric</td>
        <td>
        Validator stake in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> numeric</td>
        <td>
        Last bonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> numeric</td>
        <td>
        Last sortition height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> numeric</td>
        <td>
        Unbonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.address</td>
        <td> string</td>
        <td>
        Address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.availability_score</td>
        <td> numeric</td>
        <td>
        Availability score of the validator.
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
    Validator number.
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
    Validator information.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">validator.hash</td>
        <td> string</td>
        <td>
        Hash of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.data</td>
        <td> string</td>
        <td>
        Validator data.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.public_key</td>
        <td> string</td>
        <td>
        Public key of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.number</td>
        <td> numeric</td>
        <td>
        Validator number.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> numeric</td>
        <td>
        Validator stake in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> numeric</td>
        <td>
        Last bonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> numeric</td>
        <td>
        Last sortition height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> numeric</td>
        <td>
        Unbonding height.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.address</td>
        <td> string</td>
        <td>
        Address of the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.availability_score</td>
        <td> numeric</td>
        <td>
        Availability score of the validator.
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
    Address for which public key is requested.
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
    Public key of the account.
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
    Only returns the peers with connected status
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
    <td class="fw-bold">total_sent_bytes</td>
    <td> numeric</td>
    <td>
    Total bytes sent across the network.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_received_bytes</td>
    <td> numeric</td>
    <td>
    Total bytes received across the network.
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
        Status of the peer.
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
        Agent information of the peer.
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
        Consensus keys used by the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].consensus_address</td>
        <td>repeated string</td>
        <td>
        Consensus address of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].services</td>
        <td> numeric</td>
        <td>
        Services provided by the peer.
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
        Height of the peer in the blockchain.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].received_bundles</td>
        <td> numeric</td>
        <td>
        Count of received bundles.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].invalid_bundles</td>
        <td> numeric</td>
        <td>
        Count of invalid bundles received.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_sent</td>
        <td> numeric</td>
        <td>
        Timestamp of the last sent bundle.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_received</td>
        <td> numeric</td>
        <td>
        Timestamp of the last received bundle.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].sent_bytes</td>
        <td> object</td>
        <td>
        Bytes sent per message type.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].received_bytes</td>
        <td> object</td>
        <td>
        Bytes received per message type.
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
        Direction of connection with the peer.
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
        Total sessions with the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].completed_sessions</td>
        <td> numeric</td>
        <td>
        Completed sessions with the peer.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">sent_bytes</td>
    <td> object</td>
    <td>
    Bytes sent per peer ID.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">received_bytes</td>
    <td> object</td>
    <td>
    Bytes received per peer ID.
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
    Agent information of the node.
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
    Timestamp when the node started.
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
    <td>repeated numeric</td>
    <td>
    List of services provided by the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">services_names</td>
    <td>repeated string</td>
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
    Clock offset
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connection_info</td>
    <td> object</td>
    <td>
    Connection information
    </td>
  </tr>
     <tr>
        <td class="fw-bold">connection_info.connections</td>
        <td> numeric</td>
        <td>
        Total number of the connection.
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
    Name of the new wallet.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    Password for securing the wallet.
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
    Menomic for wallet recovery.
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
    Name of the wallet to restore.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">mnemonic</td>
    <td> string</td>
    <td>
    Menomic for wallet recovery.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    Password for securing the wallet.
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
    Name of the restored wallet.
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
    Name of the wallet to load.
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
    Name of the loaded wallet.
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
    Name of the wallet to unload.
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
    Name of the unloaded wallet.
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
    Name of the wallet.
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
    Name of the wallet.
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
    Name of the wallet used for signing.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">raw_transaction</td>
    <td> string</td>
    <td>
    Raw transaction data to be signed.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">password</td>
    <td> string</td>
    <td>
    Password for unlocking the wallet for signing.
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
    ID of the signed transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">signed_raw_transaction</td>
    <td> string</td>
    <td>
    Signed raw transaction data.
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
    Public key for which the validator address is requested.
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
    Validator address associated with the public key.
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
    The name of the wallet for which the new address is requested.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address_type</td>
    <td> string</td>
    <td>
    (Enum) The type of the new address.
    <br>Available values:<ul>
      <li>ADDRESS_TYPE_TREASURY = </li>
      <li>ADDRESS_TYPE_VALIDATOR = </li>
      <li>ADDRESS_TYPE_BLS_ACCOUNT = </li>
      </ul>
    </td>
  </tr>
  <tr>
    <td class="fw-bold">label</td>
    <td> string</td>
    <td>
    The label for the new address.
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
    The name of the wallet from which the address is created.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">address_info</td>
    <td> object</td>
    <td>
    Information about the new address.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">address_info.address</td>
        <td> string</td>
        <td>
        The string representing the address.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">address_info.public_key</td>
        <td> string</td>
        <td>
        The public key that the address is derived from.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">address_info.label</td>
        <td> string</td>
        <td>
        The label that is associated with the address.
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

<p>GetAddressHistory retrieve transaction history of an address.</p>

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
    Name of the wallet.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">address</td>
    <td> string</td>
    <td>
    Address to get the transaction history of it.
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
    Array of address history and activities.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">history_info[].transaction_id</td>
        <td> string</td>
        <td>
        Hash of transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].time</td>
        <td> numeric</td>
        <td>
        Transaction timestamp.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].payload_type</td>
        <td> string</td>
        <td>
        Type of transaction payload.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].description</td>
        <td> string</td>
        <td>
        Description of transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">history_info[].amount</td>
        <td> numeric</td>
        <td>
        Amount of transaction.
        </td>
      </tr>
         </tbody>
</table>
