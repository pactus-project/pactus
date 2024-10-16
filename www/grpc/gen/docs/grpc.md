---
title: GRPC API Reference
weight: 1
---

Each node in the Pactus network can be configured to use the [gRPC](https://grpc.io/) protocol for communication.
Here you can find the list of all gRPC methods and messages.

All the amounts and values in gRPC endpoints are in NanoPAC units, which are atomic and the smallest unit in the Pactus blockchain.
Each PAC is equivalent to 1,000,000,000 or 10<sup>9</sup> NanoPACs.

<h2>gRPC Services</h2>

<div id="toc-container">
  <ul class="">
  <li> Transaction Service
      <ul>
        <li>
          <a href="#pactus.Transaction.GetTransaction">
          <span class="rpc-badge"></span> GetTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.CalculateFee">
          <span class="rpc-badge"></span> CalculateFee</a>
        </li>
        <li>
          <a href="#pactus.Transaction.BroadcastTransaction">
          <span class="rpc-badge"></span> BroadcastTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.GetRawTransaction">
          <span class="rpc-badge"></span> GetRawTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.GetRawTransferTransaction">
          <span class="rpc-badge"></span> GetRawTransferTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.GetRawBondTransaction">
          <span class="rpc-badge"></span> GetRawBondTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.GetRawUnbondTransaction">
          <span class="rpc-badge"></span> GetRawUnbondTransaction</a>
        </li>
        <li>
          <a href="#pactus.Transaction.GetRawWithdrawTransaction">
          <span class="rpc-badge"></span> GetRawWithdrawTransaction</a>
        </li>
        </ul>
    </li>
    <li> Blockchain Service
      <ul>
        <li>
          <a href="#pactus.Blockchain.GetBlock">
          <span class="rpc-badge"></span> GetBlock</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetBlockHash">
          <span class="rpc-badge"></span> GetBlockHash</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetBlockHeight">
          <span class="rpc-badge"></span> GetBlockHeight</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetBlockchainInfo">
          <span class="rpc-badge"></span> GetBlockchainInfo</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetConsensusInfo">
          <span class="rpc-badge"></span> GetConsensusInfo</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetAccount">
          <span class="rpc-badge"></span> GetAccount</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetValidator">
          <span class="rpc-badge"></span> GetValidator</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetValidatorByNumber">
          <span class="rpc-badge"></span> GetValidatorByNumber</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetValidatorAddresses">
          <span class="rpc-badge"></span> GetValidatorAddresses</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetPublicKey">
          <span class="rpc-badge"></span> GetPublicKey</a>
        </li>
        <li>
          <a href="#pactus.Blockchain.GetTxPoolContent">
          <span class="rpc-badge"></span> GetTxPoolContent</a>
        </li>
        </ul>
    </li>
    <li> Network Service
      <ul>
        <li>
          <a href="#pactus.Network.GetNetworkInfo">
          <span class="rpc-badge"></span> GetNetworkInfo</a>
        </li>
        <li>
          <a href="#pactus.Network.GetNodeInfo">
          <span class="rpc-badge"></span> GetNodeInfo</a>
        </li>
        </ul>
    </li>
    <li> Utils Service
      <ul>
        <li>
          <a href="#pactus.Utils.SignMessageWithPrivateKey">
          <span class="rpc-badge"></span> SignMessageWithPrivateKey</a>
        </li>
        <li>
          <a href="#pactus.Utils.VerifyMessage">
          <span class="rpc-badge"></span> VerifyMessage</a>
        </li>
        </ul>
    </li>
    <li> Wallet Service
      <ul>
        <li>
          <a href="#pactus.Wallet.CreateWallet">
          <span class="rpc-badge"></span> CreateWallet</a>
        </li>
        <li>
          <a href="#pactus.Wallet.RestoreWallet">
          <span class="rpc-badge"></span> RestoreWallet</a>
        </li>
        <li>
          <a href="#pactus.Wallet.LoadWallet">
          <span class="rpc-badge"></span> LoadWallet</a>
        </li>
        <li>
          <a href="#pactus.Wallet.UnloadWallet">
          <span class="rpc-badge"></span> UnloadWallet</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetTotalBalance">
          <span class="rpc-badge"></span> GetTotalBalance</a>
        </li>
        <li>
          <a href="#pactus.Wallet.SignRawTransaction">
          <span class="rpc-badge"></span> SignRawTransaction</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetValidatorAddress">
          <span class="rpc-badge"></span> GetValidatorAddress</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetNewAddress">
          <span class="rpc-badge"></span> GetNewAddress</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetAddressHistory">
          <span class="rpc-badge"></span> GetAddressHistory</a>
        </li>
        <li>
          <a href="#pactus.Wallet.SignMessage">
          <span class="rpc-badge"></span> SignMessage</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetTotalStake">
          <span class="rpc-badge"></span> GetTotalStake</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetAddressInfo">
          <span class="rpc-badge"></span> GetAddressInfo</a>
        </li>
        <li>
          <a href="#pactus.Wallet.SetAddressLabel">
          <span class="rpc-badge"></span> SetAddressLabel</a>
        </li>
        <li>
          <a href="#pactus.Wallet.ListWallet">
          <span class="rpc-badge"></span> ListWallet</a>
        </li>
        <li>
          <a href="#pactus.Wallet.GetWalletInfo">
          <span class="rpc-badge"></span> GetWalletInfo</a>
        </li>
        <li>
          <a href="#pactus.Wallet.ListAddress">
          <span class="rpc-badge"></span> ListAddress</a>
        </li>
        </ul>
    </li>
    </ul>
</div>

<div class="api-doc">

## Transaction Service

<p>Transaction service defines various RPC methods for interacting with
transactions.</p>

### GetTransaction <span id="pactus.Transaction.GetTransaction" class="rpc-badge"></span>

<p>GetTransaction retrieves transaction details based on the provided request
parameters.</p>

<h4>GetTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
    <td> TransactionVerbosity</td>
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
  <h4>GetTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">block_height</td>
    <td> uint32</td>
    <td>
    The height of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">block_time</td>
    <td> uint32</td>
    <td>
    The UNIX timestamp of the block containing the transaction.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">transaction</td>
    <td> TransactionInfo</td>
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
        <td> int32</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.lock_time</td>
        <td> uint32</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.value</td>
        <td> int64</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.fee</td>
        <td> int64</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">transaction.payload_type</td>
        <td> PayloadType</td>
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
        <td> PayloadTransfer</td>
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
            <td> int64</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">transaction.bond</td>
        <td> PayloadBond</td>
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
            <td> int64</td>
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
        <td> PayloadSortition</td>
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
        <td> PayloadUnbond</td>
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
        <td> PayloadWithdraw</td>
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
            <td> int64</td>
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

### CalculateFee <span id="pactus.Transaction.CalculateFee" class="rpc-badge"></span>

<p>CalculateFee calculates the transaction fee based on the specified amount
and payload type.</p>

<h4>CalculateFeeRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">amount</td>
    <td> int64</td>
    <td>
    The amount involved in the transaction, specified in NanoPAC.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> PayloadType</td>
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
    <td> bool</td>
    <td>
    Indicates if the amount should be fixed and include the fee.
    </td>
  </tr>
  </tbody>
</table>
  <h4>CalculateFeeResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">amount</td>
    <td> int64</td>
    <td>
    The calculated amount in NanoPAC.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">fee</td>
    <td> int64</td>
    <td>
    The calculated transaction fee in NanoPAC.
    </td>
  </tr>
     </tbody>
</table>

### BroadcastTransaction <span id="pactus.Transaction.BroadcastTransaction" class="rpc-badge"></span>

<p>BroadcastTransaction broadcasts a signed transaction to the network.</p>

<h4>BroadcastTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>BroadcastTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetRawTransaction <span id="pactus.Transaction.GetRawTransaction" class="rpc-badge"></span>

<p>GetRawTransaction retrieves raw details of transfer, bond, unbond or withdraw transaction.</p>

<h4>GetRawTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> uint32</td>
    <td>
    The lock time for the transaction. If not set, defaults to the last block height.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">memo</td>
    <td> string</td>
    <td>
    A memo string for the transaction.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> int64</td>
    <td>
    The fee for the transaction in NanoPAC.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">transfer</td>
    <td> PayloadTransfer</td>
    <td>
    (OneOf)
    </td>
  </tr>
  <tr>
    <td class="fw-bold">bond</td>
    <td> PayloadBond</td>
    <td>
    (OneOf)
    </td>
  </tr>
  <tr>
    <td class="fw-bold">unbond</td>
    <td> PayloadUnbond</td>
    <td>
    (OneOf)
    </td>
  </tr>
  <tr>
    <td class="fw-bold">withdraw</td>
    <td> PayloadWithdraw</td>
    <td>
    (OneOf)
    </td>
  </tr>
  </tbody>
</table>
  <h4>GetRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetRawTransferTransaction <span id="pactus.Transaction.GetRawTransferTransaction" class="rpc-badge"></span>

<p>Deprecated: GetRawTransferTransaction retrieves raw details of a transfer transaction.
Use GetRawTransaction instead.</p>

<h4>GetRawTransferTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> uint32</td>
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
    <td> int64</td>
    <td>
    The amount to be transferred, specified in NanoPAC. Must be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> int64</td>
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
  <h4>GetRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetRawBondTransaction <span id="pactus.Transaction.GetRawBondTransaction" class="rpc-badge"></span>

<p>Deprecated: GetRawBondTransaction retrieves raw details of a bond transaction.
Use GetRawTransaction instead.</p>

<h4>GetRawBondTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> uint32</td>
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
    <td> int64</td>
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
    <td> int64</td>
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
  <h4>GetRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetRawUnbondTransaction <span id="pactus.Transaction.GetRawUnbondTransaction" class="rpc-badge"></span>

<p>Deprecated: GetRawUnbondTransaction retrieves raw details of an unbond transaction.
Use GetRawTransaction instead.</p>

<h4>GetRawUnbondTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> uint32</td>
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
  <h4>GetRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetRawWithdrawTransaction <span id="pactus.Transaction.GetRawWithdrawTransaction" class="rpc-badge"></span>

<p>Deprecated: GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
Use GetRawTransaction instead.</p>

<h4>GetRawWithdrawTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">lock_time</td>
    <td> uint32</td>
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
    <td> int64</td>
    <td>
    The withdrawal amount in NanoPAC. Must be greater than 0.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">fee</td>
    <td> int64</td>
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
  <h4>GetRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

## Blockchain Service

<p>Blockchain service defines RPC methods for interacting with the blockchain.</p>

### GetBlock <span id="pactus.Blockchain.GetBlock" class="rpc-badge"></span>

<p>GetBlock retrieves information about a block based on the provided request
parameters.</p>

<h4>GetBlockRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> uint32</td>
    <td>
    The height of the block to retrieve.
    </td>
  </tr>
  <tr>
    <td class="fw-bold">verbosity</td>
    <td> BlockVerbosity</td>
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
  <h4>GetBlockResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> uint32</td>
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
    <td> uint32</td>
    <td>
    The timestamp of the block.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">header</td>
    <td> BlockHeaderInfo</td>
    <td>
    Header information of the block.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">header.version</td>
        <td> int32</td>
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
    <td> CertificateInfo</td>
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
        <td> int32</td>
        <td>
        The round of the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.committers</td>
        <td>repeated int32</td>
        <td>
        List of committers in the certificate.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">prev_cert.absentees</td>
        <td>repeated int32</td>
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
    <td>repeated TransactionInfo</td>
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
        <td> int32</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> uint32</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> int64</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> int64</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> PayloadType</td>
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
        <td> PayloadTransfer</td>
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
            <td> int64</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> PayloadBond</td>
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
            <td> int64</td>
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
        <td> PayloadSortition</td>
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
        <td> PayloadUnbond</td>
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
        <td> PayloadWithdraw</td>
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
            <td> int64</td>
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

### GetBlockHash <span id="pactus.Blockchain.GetBlockHash" class="rpc-badge"></span>

<p>GetBlockHash retrieves the hash of a block at the specified height.</p>

<h4>GetBlockHashRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> uint32</td>
    <td>
    The height of the block to retrieve the hash for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>GetBlockHashResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetBlockHeight <span id="pactus.Blockchain.GetBlockHeight" class="rpc-badge"></span>

<p>GetBlockHeight retrieves the height of a block with the specified hash.</p>

<h4>GetBlockHeightRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetBlockHeightResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">height</td>
    <td> uint32</td>
    <td>
    The height of the block.
    </td>
  </tr>
     </tbody>
</table>

### GetBlockchainInfo <span id="pactus.Blockchain.GetBlockchainInfo" class="rpc-badge"></span>

<p>GetBlockchainInfo retrieves general information about the blockchain.</p>

<h4>GetBlockchainInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

Message has no fields.
  <h4>GetBlockchainInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">last_block_height</td>
    <td> uint32</td>
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
    <td> int32</td>
    <td>
    The total number of accounts in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_validators</td>
    <td> int32</td>
    <td>
    The total number of validators in the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_power</td>
    <td> int64</td>
    <td>
    The total power of the blockchain.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">committee_power</td>
    <td> int64</td>
    <td>
    The power of the committee.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">committee_validators</td>
    <td>repeated ValidatorInfo</td>
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
        <td> int32</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].stake</td>
        <td> int64</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_bonding_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].last_sortition_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">committee_validators[].unbonding_height</td>
        <td> uint32</td>
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
        <td> double</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">is_pruned</td>
    <td> bool</td>
    <td>
    If the blocks are subject to pruning.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">pruning_height</td>
    <td> uint32</td>
    <td>
    Lowest-height block stored (only present if pruning is enabled)
    </td>
  </tr>
     <tr>
    <td class="fw-bold">last_block_time</td>
    <td> int64</td>
    <td>
    Timestamp of the last block in Unix format
    </td>
  </tr>
     </tbody>
</table>

### GetConsensusInfo <span id="pactus.Blockchain.GetConsensusInfo" class="rpc-badge"></span>

<p>GetConsensusInfo retrieves information about the consensus instances.</p>

<h4>GetConsensusInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

Message has no fields.
  <h4>GetConsensusInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">proposal</td>
    <td> ProposalInfo</td>
    <td>
    The proposal of the consensus info.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">proposal.height</td>
        <td> uint32</td>
        <td>
        The height of the proposal.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">proposal.round</td>
        <td> int32</td>
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
    <td>repeated ConsensusInfo</td>
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
        <td> bool</td>
        <td>
        Indicates whether the consensus instance is active and part of the
committee.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].height</td>
        <td> uint32</td>
        <td>
        The height of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].round</td>
        <td> int32</td>
        <td>
        The round of the consensus instance.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">instances[].votes</td>
        <td>repeated VoteInfo</td>
        <td>
        List of votes in the consensus instance.
        </td>
      </tr>
         <tr>
            <td class="fw-bold">instances[].votes[].type</td>
            <td> VoteType</td>
            <td>
            (Enum)The type of the vote.
            <br>Available values:<ul>
              <li>VOTE_UNKNOWN = 0 (Unknown vote type.)</li>
              <li>VOTE_PREPARE = 1 (Prepare vote type.)</li>
              <li>VOTE_PRECOMMIT = 2 (Precommit vote type.)</li>
              <li>VOTE_CHANGE_PROPOSER = 3 (Change proposer vote type.)</li>
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
            <td> int32</td>
            <td>
            The consensus round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_round</td>
            <td> int32</td>
            <td>
            The change-proposer round of the vote.
            </td>
          </tr>
          <tr>
            <td class="fw-bold">instances[].votes[].cp_value</td>
            <td> int32</td>
            <td>
            The change-proposer value of the vote.
            </td>
          </tr>
          </tbody>
</table>

### GetAccount <span id="pactus.Blockchain.GetAccount" class="rpc-badge"></span>

<p>GetAccount retrieves information about an account based on the provided
address.</p>

<h4>GetAccountRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetAccountResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">account</td>
    <td> AccountInfo</td>
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
        <td> int32</td>
        <td>
        The unique number assigned to the account.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">account.balance</td>
        <td> int64</td>
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

### GetValidator <span id="pactus.Blockchain.GetValidator" class="rpc-badge"></span>

<p>GetValidator retrieves information about a validator based on the provided
address.</p>

<h4>GetValidatorRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetValidatorResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">validator</td>
    <td> ValidatorInfo</td>
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
        <td> int32</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> int64</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> uint32</td>
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
        <td> double</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         </tbody>
</table>

### GetValidatorByNumber <span id="pactus.Blockchain.GetValidatorByNumber" class="rpc-badge"></span>

<p>GetValidatorByNumber retrieves information about a validator based on the
provided number.</p>

<h4>GetValidatorByNumberRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">number</td>
    <td> int32</td>
    <td>
    The unique number of the validator to retrieve information for.
    </td>
  </tr>
  </tbody>
</table>
  <h4>GetValidatorResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">validator</td>
    <td> ValidatorInfo</td>
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
        <td> int32</td>
        <td>
        The unique number assigned to the validator.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.stake</td>
        <td> int64</td>
        <td>
        The stake of the validator in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_bonding_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last bonded.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.last_sortition_height</td>
        <td> uint32</td>
        <td>
        The height at which the validator last participated in sortition.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">validator.unbonding_height</td>
        <td> uint32</td>
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
        <td> double</td>
        <td>
        The availability score of the validator.
        </td>
      </tr>
         </tbody>
</table>

### GetValidatorAddresses <span id="pactus.Blockchain.GetValidatorAddresses" class="rpc-badge"></span>

<p>GetValidatorAddresses retrieves a list of all validator addresses.</p>

<h4>GetValidatorAddressesRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

Message has no fields.
  <h4>GetValidatorAddressesResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetPublicKey <span id="pactus.Blockchain.GetPublicKey" class="rpc-badge"></span>

<p>GetPublicKey retrieves the public key of an account based on the provided
address.</p>

<h4>GetPublicKeyRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetPublicKeyResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetTxPoolContent <span id="pactus.Blockchain.GetTxPoolContent" class="rpc-badge"></span>

<p>GetTxPoolContent retrieves current transactions in the transaction pool.</p>

<h4>GetTxPoolContentRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">payload_type</td>
    <td> PayloadType</td>
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
  <h4>GetTxPoolContentResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">txs</td>
    <td>repeated TransactionInfo</td>
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
        <td> int32</td>
        <td>
        The version of the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].lock_time</td>
        <td> uint32</td>
        <td>
        The lock time for the transaction.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].value</td>
        <td> int64</td>
        <td>
        The value of the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].fee</td>
        <td> int64</td>
        <td>
        The fee for the transaction in NanoPAC.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">txs[].payload_type</td>
        <td> PayloadType</td>
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
        <td> PayloadTransfer</td>
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
            <td> int64</td>
            <td>
            The amount to be transferred in NanoPAC.
            </td>
          </tr>
          <tr>
        <td class="fw-bold">txs[].bond</td>
        <td> PayloadBond</td>
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
            <td> int64</td>
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
        <td> PayloadSortition</td>
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
        <td> PayloadUnbond</td>
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
        <td> PayloadWithdraw</td>
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
            <td> int64</td>
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

### GetNetworkInfo <span id="pactus.Network.GetNetworkInfo" class="rpc-badge"></span>

<p>GetNetworkInfo retrieves information about the overall network.</p>

<h4>GetNetworkInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">only_connected</td>
    <td> bool</td>
    <td>
    If true, only returns peers with connected status.
    </td>
  </tr>
  </tbody>
</table>
  <h4>GetNetworkInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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
    <td> int64</td>
    <td>
    Total bytes sent across the network.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">total_received_bytes</td>
    <td> int64</td>
    <td>
    Total bytes received across the network.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connected_peers_count</td>
    <td> uint32</td>
    <td>
    Number of connected peers.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connected_peers</td>
    <td>repeated PeerInfo</td>
    <td>
    List of connected peers.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">connected_peers[].status</td>
        <td> int32</td>
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
        <td class="fw-bold">connected_peers[].consensus_addresses</td>
        <td>repeated string</td>
        <td>
        Consensus addresses of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].services</td>
        <td> uint32</td>
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
        <td> uint32</td>
        <td>
        Blockchain height of the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].received_bundles</td>
        <td> int32</td>
        <td>
        Number of received bundles.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].invalid_bundles</td>
        <td> int32</td>
        <td>
        Number of invalid bundles received.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_sent</td>
        <td> int64</td>
        <td>
        Timestamp of the last sent bundle.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].last_received</td>
        <td> int64</td>
        <td>
        Timestamp of the last received bundle.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].sent_bytes</td>
        <td> map&lt;int32, int64&gt;</td>
        <td>
        Bytes sent per message type.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].received_bytes</td>
        <td> map&lt;int32, int64&gt;</td>
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
        <td> int32</td>
        <td>
        Total download sessions with the peer.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connected_peers[].completed_sessions</td>
        <td> int32</td>
        <td>
        Completed download sessions with the peer.
        </td>
      </tr>
         <tr>
    <td class="fw-bold">sent_bytes</td>
    <td> map&lt;int32, int64&gt;</td>
    <td>
    Bytes sent per peer ID.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">received_bytes</td>
    <td> map&lt;int32, int64&gt;</td>
    <td>
    Bytes received per peer ID.
    </td>
  </tr>
     </tbody>
</table>

### GetNodeInfo <span id="pactus.Network.GetNodeInfo" class="rpc-badge"></span>

<p>GetNodeInfo retrieves information about a specific node in the network.</p>

<h4>GetNodeInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

Message has no fields.
  <h4>GetNodeInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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
    <td> uint64</td>
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
    <td> int32</td>
    <td>
    A bitfield indicating the services provided by the node.
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
    <td> double</td>
    <td>
    Clock offset of the node.
    </td>
  </tr>
     <tr>
    <td class="fw-bold">connection_info</td>
    <td> ConnectionInfo</td>
    <td>
    Information about the node's connections.
    </td>
  </tr>
     <tr>
        <td class="fw-bold">connection_info.connections</td>
        <td> uint64</td>
        <td>
        Total number of connections.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connection_info.inbound_connections</td>
        <td> uint64</td>
        <td>
        Number of inbound connections.
        </td>
      </tr>
         <tr>
        <td class="fw-bold">connection_info.outbound_connections</td>
        <td> uint64</td>
        <td>
        Number of outbound connections.
        </td>
      </tr>
         </tbody>
</table>

## Utils Service

<p>Utils service defines RPC methods for utility functions such as message
signing and verification.</p>

### SignMessageWithPrivateKey <span id="pactus.Utils.SignMessageWithPrivateKey" class="rpc-badge"></span>

<p>SignMessageWithPrivateKey sign message with provided private key.</p>

<h4>SignMessageWithPrivateKeyRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>SignMessageWithPrivateKeyResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### VerifyMessage <span id="pactus.Utils.VerifyMessage" class="rpc-badge"></span>

<p>VerifyMessage verify signature with public key and message</p>

<h4>VerifyMessageRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>VerifyMessageResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">is_valid</td>
    <td> bool</td>
    <td>
    Indicates if the signature is valid (true) or not (false).
    </td>
  </tr>
     </tbody>
</table>

## Wallet Service

<p>Define the Wallet service with various RPC methods for wallet management.</p>

### CreateWallet <span id="pactus.Wallet.CreateWallet" class="rpc-badge"></span>

<p>CreateWallet creates a new wallet with the specified parameters.</p>

<h4>CreateWalletRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>CreateWalletResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### RestoreWallet <span id="pactus.Wallet.RestoreWallet" class="rpc-badge"></span>

<p>RestoreWallet restores an existing wallet with the given mnemonic.</p>

<h4>RestoreWalletRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>RestoreWalletResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### LoadWallet <span id="pactus.Wallet.LoadWallet" class="rpc-badge"></span>

<p>LoadWallet loads an existing wallet with the given name.</p>

<h4>LoadWalletRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>LoadWalletResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### UnloadWallet <span id="pactus.Wallet.UnloadWallet" class="rpc-badge"></span>

<p>UnloadWallet unloads a currently loaded wallet with the specified name.</p>

<h4>UnloadWalletRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>UnloadWalletResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetTotalBalance <span id="pactus.Wallet.GetTotalBalance" class="rpc-badge"></span>

<p>GetTotalBalance returns the total available balance of the wallet.</p>

<h4>GetTotalBalanceRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetTotalBalanceResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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
    <td> int64</td>
    <td>
    The total balance of the wallet in NanoPAC.
    </td>
  </tr>
     </tbody>
</table>

### SignRawTransaction <span id="pactus.Wallet.SignRawTransaction" class="rpc-badge"></span>

<p>SignRawTransaction signs a raw transaction for a specified wallet.</p>

<h4>SignRawTransactionRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>SignRawTransactionResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetValidatorAddress <span id="pactus.Wallet.GetValidatorAddress" class="rpc-badge"></span>

<p>GetValidatorAddress retrieves the validator address associated with a
public key.</p>

<h4>GetValidatorAddressRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetValidatorAddressResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetNewAddress <span id="pactus.Wallet.GetNewAddress" class="rpc-badge"></span>

<p>GetNewAddress generates a new address for the specified wallet.</p>

<h4>GetNewAddressRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
    <td> AddressType</td>
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
  <h4>GetNewAddressResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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
    <td> AddressInfo</td>
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

### GetAddressHistory <span id="pactus.Wallet.GetAddressHistory" class="rpc-badge"></span>

<p>GetAddressHistory retrieves the transaction history of an address.</p>

<h4>GetAddressHistoryRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetAddressHistoryResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">history_info</td>
    <td>repeated HistoryInfo</td>
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
        <td> uint32</td>
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
        <td> int64</td>
        <td>
        The amount involved in the transaction.
        </td>
      </tr>
         </tbody>
</table>

### SignMessage <span id="pactus.Wallet.SignMessage" class="rpc-badge"></span>

<p>SignMessage signs an arbitrary message.</p>

<h4>SignMessageRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>SignMessageResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetTotalStake <span id="pactus.Wallet.GetTotalStake" class="rpc-badge"></span>

<p>GetTotalStake return total stake of wallet.</p>

<h4>GetTotalStakeRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetTotalStakeResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">total_stake</td>
    <td> int64</td>
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

### GetAddressInfo <span id="pactus.Wallet.GetAddressInfo" class="rpc-badge"></span>

<p>GetAddressInfo return address information.</p>

<h4>GetAddressInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetAddressInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### SetAddressLabel <span id="pactus.Wallet.SetAddressLabel" class="rpc-badge"></span>

<p>SetAddressLabel set label for given address.</p>

<h4>SetLabelRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>SetLabelResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  </tbody>
</table>

### ListWallet <span id="pactus.Wallet.ListWallet" class="rpc-badge"></span>

<p>ListWallet return list wallet name.</p>

<h4>ListWalletRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

Message has no fields.
  <h4>ListWalletResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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

### GetWalletInfo <span id="pactus.Wallet.GetWalletInfo" class="rpc-badge"></span>

<p>GetWalletInfo return wallet information.</p>

<h4>GetWalletInfoRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>GetWalletInfoResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

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
    <td> int64</td>
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
    <td> bool</td>
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
    <td class="fw-bold">crc</td>
    <td> uint32</td>
    <td>
    
    </td>
  </tr>
     <tr>
    <td class="fw-bold">created_at</td>
    <td> int64</td>
    <td>
    
    </td>
  </tr>
     </tbody>
</table>

### ListAddress <span id="pactus.Wallet.ListAddress" class="rpc-badge"></span>

<p>ListAddress return list address in wallet.</p>

<h4>ListAddressRequest <span class="badge text-bg-info fs-6 align-top">Request</span></h4>

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
  <h4>ListAddressResponse <span class="badge text-bg-warning fs-6 align-top">Response</span></h4>

<table class="table table-bordered table-responsive table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
  <tr>
    <td class="fw-bold">data</td>
    <td>repeated AddressInfo</td>
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

## Scalar Value Types

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>.proto Type</td><td>Go</td><td>C++</td><td>Rust</td><td>Java</td><td>Python</td><td>C#</td></tr>
  </thead>
  <tbody class="table-group-divider">
      <tr id="double">
        <td class="fw-bold">double</td>
        <td>float64</td>
        <td>double</td>
        <td>f64</td>
        <td>double</td>
        <td>float</td>
        <td>double</td>
      </tr>
      <tr id="float">
        <td class="fw-bold">float</td>
        <td>float32</td>
        <td>float</td>
        <td>f32</td>
        <td>float</td>
        <td>float</td>
        <td>float</td>
      </tr>
      <tr id="int32">
        <td class="fw-bold">int32</td>
        <td>int32</td>
        <td>int32</td>
        <td>i32</td>
        <td>int</td>
        <td>int</td>
        <td>int</td>
      </tr>
      <tr id="int64">
        <td class="fw-bold">int64</td>
        <td>int64</td>
        <td>int64</td>
        <td>i64</td>
        <td>long</td>
        <td>int/long</td>
        <td>long</td>
      </tr>
      <tr id="uint32">
        <td class="fw-bold">uint32</td>
        <td>uint32</td>
        <td>uint32</td>
        <td>u32</td>
        <td>int</td>
        <td>int/long</td>
        <td>uint</td>
      </tr>
      <tr id="uint64">
        <td class="fw-bold">uint64</td>
        <td>uint64</td>
        <td>uint64</td>
        <td>u64</td>
        <td>long</td>
        <td>int/long</td>
        <td>ulong</td>
      </tr>
      <tr id="sint32">
        <td class="fw-bold">sint32</td>
        <td>int32</td>
        <td>int32</td>
        <td>i32</td>
        <td>int</td>
        <td>int</td>
        <td>int</td>
      </tr>
      <tr id="sint64">
        <td class="fw-bold">sint64</td>
        <td>int64</td>
        <td>int64</td>
        <td>i64</td>
        <td>long</td>
        <td>int/long</td>
        <td>long</td>
      </tr>
      <tr id="fixed32">
        <td class="fw-bold">fixed32</td>
        <td>uint32</td>
        <td>uint32</td>
        <td>u64</td>
        <td>int</td>
        <td>int</td>
        <td>uint</td>
      </tr>
      <tr id="fixed64">
        <td class="fw-bold">fixed64</td>
        <td>uint64</td>
        <td>uint64</td>
        <td>u64</td>
        <td>long</td>
        <td>int/long</td>
        <td>ulong</td>
      </tr>
      <tr id="sfixed32">
        <td class="fw-bold">sfixed32</td>
        <td>int32</td>
        <td>int32</td>
        <td>i32</td>
        <td>int</td>
        <td>int</td>
        <td>int</td>
      </tr>
      <tr id="sfixed64">
        <td class="fw-bold">sfixed64</td>
        <td>int64</td>
        <td>int64</td>
        <td>i64</td>
        <td>long</td>
        <td>int/long</td>
        <td>long</td>
      </tr>
      <tr id="bool">
        <td class="fw-bold">bool</td>
        <td>bool</td>
        <td>bool</td>
        <td>bool</td>
        <td>boolean</td>
        <td>boolean</td>
        <td>bool</td>
      </tr>
      <tr id="string">
        <td class="fw-bold">string</td>
        <td>string</td>
        <td>string</td>
        <td>String</td>
        <td>String</td>
        <td>str/unicode</td>
        <td>string</td>
      </tr>
      <tr id="bytes">
        <td class="fw-bold">bytes</td>
        <td>[]byte</td>
        <td>string</td>
        <td>Vec<u8></td>
        <td>ByteString</td>
        <td>str</td>
        <td>ByteString</td>
      </tr>
  </tbody>
</table>
