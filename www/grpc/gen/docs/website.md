---
layout: learn
title: gRPC API Reference
sidebar: gRPC API Reference
---

<h1 id="title">gRPC API Reference</h1>

Each node in the Pactus network can be configured to use the [gRPC](https://grpc.io/) protocol for communication.
Here you can find the list of all gRPC methods and messages.

<h2>Table of Contents</h2>

<div id="toc-container">
  <ul class="">  
    <li> Transaction Service
      <ul>  
        <li>
          <a href="#pactus.Transaction.GetTransaction">
          <span class="badge text-bg-primary">rpc</span> GetTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.CalculateFee">
          <span class="badge text-bg-primary">rpc</span> CalculateFee</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.BroadcastTransaction">
          <span class="badge text-bg-primary">rpc</span> BroadcastTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.GetRawTransferTransaction">
          <span class="badge text-bg-primary">rpc</span> GetRawTransferTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.GetRawBondTransaction">
          <span class="badge text-bg-primary">rpc</span> GetRawBondTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.GetRawUnBondTransaction">
          <span class="badge text-bg-primary">rpc</span> GetRawUnBondTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Transaction.GetRawWithdrawTransaction">
          <span class="badge text-bg-primary">rpc</span> GetRawWithdrawTransaction</a>
        </li> 
      </ul>
    </li>    
    <li> Blockchain Service
      <ul>  
        <li>
          <a href="#pactus.Blockchain.GetBlock">
          <span class="badge text-bg-primary">rpc</span> GetBlock</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetBlockHash">
          <span class="badge text-bg-primary">rpc</span> GetBlockHash</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetBlockHeight">
          <span class="badge text-bg-primary">rpc</span> GetBlockHeight</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetBlockchainInfo">
          <span class="badge text-bg-primary">rpc</span> GetBlockchainInfo</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetConsensusInfo">
          <span class="badge text-bg-primary">rpc</span> GetConsensusInfo</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetAccount">
          <span class="badge text-bg-primary">rpc</span> GetAccount</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetValidator">
          <span class="badge text-bg-primary">rpc</span> GetValidator</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetValidatorByNumber">
          <span class="badge text-bg-primary">rpc</span> GetValidatorByNumber</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetValidatorAddresses">
          <span class="badge text-bg-primary">rpc</span> GetValidatorAddresses</a>
        </li> 
        <li>
          <a href="#pactus.Blockchain.GetPublicKey">
          <span class="badge text-bg-primary">rpc</span> GetPublicKey</a>
        </li> 
      </ul>
    </li>    
    <li> Network Service
      <ul>  
        <li>
          <a href="#pactus.Network.GetNetworkInfo">
          <span class="badge text-bg-primary">rpc</span> GetNetworkInfo</a>
        </li> 
        <li>
          <a href="#pactus.Network.GetNodeInfo">
          <span class="badge text-bg-primary">rpc</span> GetNodeInfo</a>
        </li> 
      </ul>
    </li>    
    <li> Wallet Service
      <ul>  
        <li>
          <a href="#pactus.Wallet.CreateWallet">
          <span class="badge text-bg-primary">rpc</span> CreateWallet</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.LoadWallet">
          <span class="badge text-bg-primary">rpc</span> LoadWallet</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.UnloadWallet">
          <span class="badge text-bg-primary">rpc</span> UnloadWallet</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.LockWallet">
          <span class="badge text-bg-primary">rpc</span> LockWallet</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.UnlockWallet">
          <span class="badge text-bg-primary">rpc</span> UnlockWallet</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.SignRawTransaction">
          <span class="badge text-bg-primary">rpc</span> SignRawTransaction</a>
        </li> 
        <li>
          <a href="#pactus.Wallet.GetValidatorAddress">
          <span class="badge text-bg-primary">rpc</span> GetValidatorAddress</a>
        </li> 
      </ul>
    </li>   

    <li>Messages and Enums
      <ul>  
        <li>
          <a href="#pactus.BroadcastTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> BroadcastTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.BroadcastTransactionResponse">
            <span class="badge text-bg-secondary">msg</span> BroadcastTransactionResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.CalculateFeeRequest">
            <span class="badge text-bg-secondary">msg</span> CalculateFeeRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.CalculateFeeResponse">
            <span class="badge text-bg-secondary">msg</span> CalculateFeeResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetRawBondTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> GetRawBondTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetRawTransactionResponse">
            <span class="badge text-bg-secondary">msg</span> GetRawTransactionResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetRawTransferTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> GetRawTransferTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetRawUnBondTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> GetRawUnBondTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetRawWithdrawTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> GetRawWithdrawTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> GetTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetTransactionResponse">
            <span class="badge text-bg-secondary">msg</span> GetTransactionResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.PayloadBond">
            <span class="badge text-bg-secondary">msg</span> PayloadBond
          </a>
        </li> 
        <li>
          <a href="#pactus.PayloadSortition">
            <span class="badge text-bg-secondary">msg</span> PayloadSortition
          </a>
        </li> 
        <li>
          <a href="#pactus.PayloadTransfer">
            <span class="badge text-bg-secondary">msg</span> PayloadTransfer
          </a>
        </li> 
        <li>
          <a href="#pactus.PayloadUnbond">
            <span class="badge text-bg-secondary">msg</span> PayloadUnbond
          </a>
        </li> 
        <li>
          <a href="#pactus.PayloadWithdraw">
            <span class="badge text-bg-secondary">msg</span> PayloadWithdraw
          </a>
        </li> 
        <li>
          <a href="#pactus.TransactionInfo">
            <span class="badge text-bg-secondary">msg</span> TransactionInfo
          </a>
        </li>   
        <li>
          <a href="#pactus.AccountInfo">
            <span class="badge text-bg-secondary">msg</span> AccountInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.BlockHeaderInfo">
            <span class="badge text-bg-secondary">msg</span> BlockHeaderInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.CertificateInfo">
            <span class="badge text-bg-secondary">msg</span> CertificateInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.ConsensusInfo">
            <span class="badge text-bg-secondary">msg</span> ConsensusInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.GetAccountRequest">
            <span class="badge text-bg-secondary">msg</span> GetAccountRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetAccountResponse">
            <span class="badge text-bg-secondary">msg</span> GetAccountResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockHashRequest">
            <span class="badge text-bg-secondary">msg</span> GetBlockHashRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockHashResponse">
            <span class="badge text-bg-secondary">msg</span> GetBlockHashResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockHeightRequest">
            <span class="badge text-bg-secondary">msg</span> GetBlockHeightRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockHeightResponse">
            <span class="badge text-bg-secondary">msg</span> GetBlockHeightResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockRequest">
            <span class="badge text-bg-secondary">msg</span> GetBlockRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockResponse">
            <span class="badge text-bg-secondary">msg</span> GetBlockResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockchainInfoRequest">
            <span class="badge text-bg-secondary">msg</span> GetBlockchainInfoRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetBlockchainInfoResponse">
            <span class="badge text-bg-secondary">msg</span> GetBlockchainInfoResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetConsensusInfoRequest">
            <span class="badge text-bg-secondary">msg</span> GetConsensusInfoRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetConsensusInfoResponse">
            <span class="badge text-bg-secondary">msg</span> GetConsensusInfoResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetPublicKeyRequest">
            <span class="badge text-bg-secondary">msg</span> GetPublicKeyRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetPublicKeyResponse">
            <span class="badge text-bg-secondary">msg</span> GetPublicKeyResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorAddressesRequest">
            <span class="badge text-bg-secondary">msg</span> GetValidatorAddressesRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorAddressesResponse">
            <span class="badge text-bg-secondary">msg</span> GetValidatorAddressesResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorByNumberRequest">
            <span class="badge text-bg-secondary">msg</span> GetValidatorByNumberRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorRequest">
            <span class="badge text-bg-secondary">msg</span> GetValidatorRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorResponse">
            <span class="badge text-bg-secondary">msg</span> GetValidatorResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.ValidatorInfo">
            <span class="badge text-bg-secondary">msg</span> ValidatorInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.VoteInfo">
            <span class="badge text-bg-secondary">msg</span> VoteInfo
          </a>
        </li>   
        <li>
          <a href="#pactus.GetNetworkInfoRequest">
            <span class="badge text-bg-secondary">msg</span> GetNetworkInfoRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetNetworkInfoResponse">
            <span class="badge text-bg-secondary">msg</span> GetNetworkInfoResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetNetworkInfoResponse.ReceivedBytesEntry">
            <span class="badge text-bg-secondary">msg</span> GetNetworkInfoResponse.ReceivedBytesEntry
          </a>
        </li> 
        <li>
          <a href="#pactus.GetNetworkInfoResponse.SentBytesEntry">
            <span class="badge text-bg-secondary">msg</span> GetNetworkInfoResponse.SentBytesEntry
          </a>
        </li> 
        <li>
          <a href="#pactus.GetNodeInfoRequest">
            <span class="badge text-bg-secondary">msg</span> GetNodeInfoRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetNodeInfoResponse">
            <span class="badge text-bg-secondary">msg</span> GetNodeInfoResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.PeerInfo">
            <span class="badge text-bg-secondary">msg</span> PeerInfo
          </a>
        </li> 
        <li>
          <a href="#pactus.PeerInfo.ReceivedBytesEntry">
            <span class="badge text-bg-secondary">msg</span> PeerInfo.ReceivedBytesEntry
          </a>
        </li> 
        <li>
          <a href="#pactus.PeerInfo.SentBytesEntry">
            <span class="badge text-bg-secondary">msg</span> PeerInfo.SentBytesEntry
          </a>
        </li>   
        <li>
          <a href="#pactus.CreateWalletRequest">
            <span class="badge text-bg-secondary">msg</span> CreateWalletRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.CreateWalletResponse">
            <span class="badge text-bg-secondary">msg</span> CreateWalletResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorAddressRequest">
            <span class="badge text-bg-secondary">msg</span> GetValidatorAddressRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.GetValidatorAddressResponse">
            <span class="badge text-bg-secondary">msg</span> GetValidatorAddressResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.LoadWalletRequest">
            <span class="badge text-bg-secondary">msg</span> LoadWalletRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.LoadWalletResponse">
            <span class="badge text-bg-secondary">msg</span> LoadWalletResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.LockWalletRequest">
            <span class="badge text-bg-secondary">msg</span> LockWalletRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.LockWalletResponse">
            <span class="badge text-bg-secondary">msg</span> LockWalletResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.SignRawTransactionRequest">
            <span class="badge text-bg-secondary">msg</span> SignRawTransactionRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.SignRawTransactionResponse">
            <span class="badge text-bg-secondary">msg</span> SignRawTransactionResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.UnloadWalletRequest">
            <span class="badge text-bg-secondary">msg</span> UnloadWalletRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.UnloadWalletResponse">
            <span class="badge text-bg-secondary">msg</span> UnloadWalletResponse
          </a>
        </li> 
        <li>
          <a href="#pactus.UnlockWalletRequest">
            <span class="badge text-bg-secondary">msg</span> UnlockWalletRequest
          </a>
        </li> 
        <li>
          <a href="#pactus.UnlockWalletResponse">
            <span class="badge text-bg-secondary">msg</span> UnlockWalletResponse
          </a>
        </li>  
         
        <li>
          <a href="#pactus.PayloadType">
            <span class="badge text-bg-info">enum</span> PayloadType
          </a>
        </li> 
        <li>
          <a href="#pactus.TransactionVerbosity">
            <span class="badge text-bg-info">enum</span> TransactionVerbosity
          </a>
        </li>   
        <li>
          <a href="#pactus.BlockVerbosity">
            <span class="badge text-bg-info">enum</span> BlockVerbosity
          </a>
        </li> 
        <li>
          <a href="#pactus.VoteType">
            <span class="badge text-bg-info">enum</span> VoteType
          </a>
        </li>      

        <li>
          <a href="#scalar-value-types">Scalar Value Types</a>
        </li>
      </ul>
    </li>
  </ul>
</div>
  
<h2>Transaction Service <span class="badge text-bg-warning fs-6 align-top">transaction.proto</span></h2>
<p>Transaction service defines various RPC methods for interacting with</p><p>transactions.</p>  
<h3 id="pactus.Transaction.GetTransaction">GetTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetTransactionRequest">GetTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetTransactionResponse">GetTransactionResponse</a></div>
<p>GetTransaction retrieves transaction details based on the provided request</p><p>parameters.</p> 
<h3 id="pactus.Transaction.CalculateFee">CalculateFee <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.CalculateFeeRequest">CalculateFeeRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.CalculateFeeResponse">CalculateFeeResponse</a></div>
<p>CalculateFee calculates the transaction fee based on the specified amount</p><p>and payload type.</p> 
<h3 id="pactus.Transaction.BroadcastTransaction">BroadcastTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.BroadcastTransactionRequest">BroadcastTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.BroadcastTransactionResponse">BroadcastTransactionResponse</a></div>
<p>BroadcastTransaction broadcasts a signed transaction to the network.</p> 
<h3 id="pactus.Transaction.GetRawTransferTransaction">GetRawTransferTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetRawTransferTransactionRequest">GetRawTransferTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetRawTransactionResponse">GetRawTransactionResponse</a></div>
<p>GetRawTransferTransaction retrieves raw details of a transfer transaction.</p> 
<h3 id="pactus.Transaction.GetRawBondTransaction">GetRawBondTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetRawBondTransactionRequest">GetRawBondTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetRawTransactionResponse">GetRawTransactionResponse</a></div>
<p>GetRawBondTransaction retrieves raw details of a bond transaction.</p> 
<h3 id="pactus.Transaction.GetRawUnBondTransaction">GetRawUnBondTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetRawUnBondTransactionRequest">GetRawUnBondTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetRawTransactionResponse">GetRawTransactionResponse</a></div>
<p>GetRawUnBondTransaction retrieves raw details of an unbond transaction.</p> 
<h3 id="pactus.Transaction.GetRawWithdrawTransaction">GetRawWithdrawTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetRawWithdrawTransactionRequest">GetRawWithdrawTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetRawTransactionResponse">GetRawTransactionResponse</a></div>
<p>GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.</p>     
<h2>Blockchain Service <span class="badge text-bg-warning fs-6 align-top">blockchain.proto</span></h2>
<p>Blockchain service defines RPC methods for interacting with the blockchain.</p>  
<h3 id="pactus.Blockchain.GetBlock">GetBlock <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetBlockRequest">GetBlockRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetBlockResponse">GetBlockResponse</a></div>
<p>GetBlock retrieves information about a block based on the provided request</p><p>parameters.</p> 
<h3 id="pactus.Blockchain.GetBlockHash">GetBlockHash <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetBlockHashRequest">GetBlockHashRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetBlockHashResponse">GetBlockHashResponse</a></div>
<p>GetBlockHash retrieves the hash of a block at the specified height.</p> 
<h3 id="pactus.Blockchain.GetBlockHeight">GetBlockHeight <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetBlockHeightRequest">GetBlockHeightRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetBlockHeightResponse">GetBlockHeightResponse</a></div>
<p>GetBlockHeight retrieves the height of a block with the specified hash.</p> 
<h3 id="pactus.Blockchain.GetBlockchainInfo">GetBlockchainInfo <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetBlockchainInfoRequest">GetBlockchainInfoRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetBlockchainInfoResponse">GetBlockchainInfoResponse</a></div>
<p>GetBlockchainInfo retrieves general information about the blockchain.</p> 
<h3 id="pactus.Blockchain.GetConsensusInfo">GetConsensusInfo <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetConsensusInfoRequest">GetConsensusInfoRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetConsensusInfoResponse">GetConsensusInfoResponse</a></div>
<p>GetConsensusInfo retrieves information about the consensus instances.</p> 
<h3 id="pactus.Blockchain.GetAccount">GetAccount <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetAccountRequest">GetAccountRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetAccountResponse">GetAccountResponse</a></div>
<p>GetAccount retrieves information about an account based on the provided</p><p>address.</p> 
<h3 id="pactus.Blockchain.GetValidator">GetValidator <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetValidatorRequest">GetValidatorRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetValidatorResponse">GetValidatorResponse</a></div>
<p>GetValidator retrieves information about a validator based on the provided</p><p>address.</p> 
<h3 id="pactus.Blockchain.GetValidatorByNumber">GetValidatorByNumber <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetValidatorByNumberRequest">GetValidatorByNumberRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetValidatorResponse">GetValidatorResponse</a></div>
<p>GetValidatorByNumber retrieves information about a validator based on the</p><p>provided number.</p> 
<h3 id="pactus.Blockchain.GetValidatorAddresses">GetValidatorAddresses <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetValidatorAddressesRequest">GetValidatorAddressesRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetValidatorAddressesResponse">GetValidatorAddressesResponse</a></div>
<p>GetValidatorAddresses retrieves a list of all validator addresses.</p> 
<h3 id="pactus.Blockchain.GetPublicKey">GetPublicKey <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetPublicKeyRequest">GetPublicKeyRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetPublicKeyResponse">GetPublicKeyResponse</a></div>
<p>GetPublicKey retrieves the public key of an account based on the provided</p><p>address.</p>     
<h2>Network Service <span class="badge text-bg-warning fs-6 align-top">network.proto</span></h2>
<p>Network service provides RPCs for retrieving information about the network.</p>  
<h3 id="pactus.Network.GetNetworkInfo">GetNetworkInfo <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetNetworkInfoRequest">GetNetworkInfoRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetNetworkInfoResponse">GetNetworkInfoResponse</a></div>
<p>GetNetworkInfo retrieves information about the overall network.</p> 
<h3 id="pactus.Network.GetNodeInfo">GetNodeInfo <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetNodeInfoRequest">GetNodeInfoRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetNodeInfoResponse">GetNodeInfoResponse</a></div>
<p>GetNodeInfo retrieves information about a specific node in the network.</p>     
<h2>Wallet Service <span class="badge text-bg-warning fs-6 align-top">wallet.proto</span></h2>
<p>Define the Wallet service with various RPC methods for wallet management.</p>  
<h3 id="pactus.Wallet.CreateWallet">CreateWallet <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.CreateWalletRequest">CreateWalletRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.CreateWalletResponse">CreateWalletResponse</a></div>
<p>CreateWallet creates a new wallet with the specified parameters.</p> 
<h3 id="pactus.Wallet.LoadWallet">LoadWallet <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.LoadWalletRequest">LoadWalletRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.LoadWalletResponse">LoadWalletResponse</a></div>
<p>LoadWallet loads an existing wallet with the given name.</p> 
<h3 id="pactus.Wallet.UnloadWallet">UnloadWallet <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.UnloadWalletRequest">UnloadWalletRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.UnloadWalletResponse">UnloadWalletResponse</a></div>
<p>UnloadWallet unloads a currently loaded wallet with the specified name.</p> 
<h3 id="pactus.Wallet.LockWallet">LockWallet <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.LockWalletRequest">LockWalletRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.LockWalletResponse">LockWalletResponse</a></div>
<p>LockWallet locks a currently loaded wallet with the provided password and</p><p>timeout.</p> 
<h3 id="pactus.Wallet.UnlockWallet">UnlockWallet <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.UnlockWalletRequest">UnlockWalletRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.UnlockWalletResponse">UnlockWalletResponse</a></div>
<p>UnlockWallet unlocks a locked wallet with the provided password and</p><p>timeout.</p> 
<h3 id="pactus.Wallet.SignRawTransaction">SignRawTransaction <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.SignRawTransactionRequest">SignRawTransactionRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.SignRawTransactionResponse">SignRawTransactionResponse</a></div>
<p>SignRawTransaction signs a raw transaction for a specified wallet.</p> 
<h3 id="pactus.Wallet.GetValidatorAddress">GetValidatorAddress <span class="badge text-bg-primary fs-6 align-top">rpc</span></h3>
<div class="request pt-3">Request message: <a href="#pactus.GetValidatorAddressRequest">GetValidatorAddressRequest</a></div>
<div class="response pb-3">Response message: <a href="#pactus.GetValidatorAddressResponse">GetValidatorAddressResponse</a></div>
<p>GetValidatorAddress retrieves the validator address associated with a</p><p>public key.</p>   
<h2>Messages and Enums</h2> 
<h3 id="pactus.BroadcastTransactionRequest">
BroadcastTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for broadcasting a signed transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">signed_raw_transaction</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Signed raw transaction data. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.BroadcastTransactionResponse">
BroadcastTransactionResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the ID of the broadcasted transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Transaction ID. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.CalculateFeeRequest">
CalculateFeeRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for calculating transaction fee.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">amount</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction amount. </td>
    </tr>
    <tr>
      <td class="fw-bold">payload_type</td>
      <td>
        <a href="#pactus.PayloadType">PayloadType</a>
      </td>
      <td>Type of transaction payload. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.CalculateFeeResponse">
CalculateFeeResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the calculated transaction fee.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">fee</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Calculated transaction fee. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetRawBondTransactionRequest">
GetRawBondTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving raw details of a bond transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">lock_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Lock time for the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">sender</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Sender's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">receiver</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Receiver's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">stake</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Stake amount. </td>
    </tr>
    <tr>
      <td class="fw-bold">public_key</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Public key of the validator. </td>
    </tr>
    <tr>
      <td class="fw-bold">fee</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction fee. </td>
    </tr>
    <tr>
      <td class="fw-bold">memo</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Transaction memo. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetRawTransactionResponse">
GetRawTransactionResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing raw transaction data.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">raw_transaction</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Raw transaction data. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetRawTransferTransactionRequest">
GetRawTransferTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving raw details of a transfer transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">lock_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Lock time for the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">sender</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Sender's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">receiver</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Receiver's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">amount</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction amount. </td>
    </tr>
    <tr>
      <td class="fw-bold">fee</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction fee. </td>
    </tr>
    <tr>
      <td class="fw-bold">memo</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Transaction memo. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetRawUnBondTransactionRequest">
GetRawUnBondTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving raw details of an unbond transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">lock_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Lock time for the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">validator_address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the validator to unbond from. </td>
    </tr>
    <tr>
      <td class="fw-bold">memo</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Transaction memo. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetRawWithdrawTransactionRequest">
GetRawWithdrawTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving raw details of a withdraw transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">lock_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Lock time for the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">validator_address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the validator to withdraw from. </td>
    </tr>
    <tr>
      <td class="fw-bold">account_address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the account to withdraw to. </td>
    </tr>
    <tr>
      <td class="fw-bold">fee</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction fee. </td>
    </tr>
    <tr>
      <td class="fw-bold">amount</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Withdrawal amount. </td>
    </tr>
    <tr>
      <td class="fw-bold">memo</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Transaction memo. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetTransactionRequest">
GetTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving transaction details.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Transaction ID. </td>
    </tr>
    <tr>
      <td class="fw-bold">verbosity</td>
      <td>
        <a href="#pactus.TransactionVerbosity">TransactionVerbosity</a>
      </td>
      <td>Verbosity level for transaction details. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetTransactionResponse">
GetTransactionResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing details of a transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">block_height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the block containing the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">block_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Time of the block containing the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">transaction</td>
      <td>
        <a href="#pactus.TransactionInfo">TransactionInfo</a>
      </td>
      <td>Information about the transaction. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PayloadBond">
PayloadBond
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Payload for a bond transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">sender</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Sender's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">receiver</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Receiver's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">stake</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Stake amount. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PayloadSortition">
PayloadSortition
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Payload for a sortition transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address associated with the sortition. </td>
    </tr>
    <tr>
      <td class="fw-bold">proof</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Proof for the sortition. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PayloadTransfer">
PayloadTransfer
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Payload for a transfer transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">sender</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Sender's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">receiver</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Receiver's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">amount</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction amount. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PayloadUnbond">
PayloadUnbond
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Payload for an unbond transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">validator</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the validator to unbond from. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PayloadWithdraw">
PayloadWithdraw
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Payload for a withdraw transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">from</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address to withdraw from. </td>
    </tr>
    <tr>
      <td class="fw-bold">to</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address to withdraw to. </td>
    </tr>
    <tr>
      <td class="fw-bold">amount</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Withdrawal amount. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.TransactionInfo">
TransactionInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Information about a transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Transaction ID. </td>
    </tr>
    <tr>
      <td class="fw-bold">data</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Transaction data. </td>
    </tr>
    <tr>
      <td class="fw-bold">version</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Transaction version. </td>
    </tr>
    <tr>
      <td class="fw-bold">lock_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Lock time for the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">value</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction value. </td>
    </tr>
    <tr>
      <td class="fw-bold">fee</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Transaction fee. </td>
    </tr>
    <tr>
      <td class="fw-bold">payload_type</td>
      <td>
        <a href="#pactus.PayloadType">PayloadType</a>
      </td>
      <td>Type of transaction payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">transfer</td>
      <td>
        <a href="#pactus.PayloadTransfer">PayloadTransfer</a>
      </td>
      <td>Transfer payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">bond</td>
      <td>
        <a href="#pactus.PayloadBond">PayloadBond</a>
      </td>
      <td>Bond payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">sortition</td>
      <td>
        <a href="#pactus.PayloadSortition">PayloadSortition</a>
      </td>
      <td>Sortition payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">unbond</td>
      <td>
        <a href="#pactus.PayloadUnbond">PayloadUnbond</a>
      </td>
      <td>Unbond payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">withdraw</td>
      <td>
        <a href="#pactus.PayloadWithdraw">PayloadWithdraw</a>
      </td>
      <td>Withdraw payload. </td>
    </tr>
    <tr>
      <td class="fw-bold">memo</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Transaction memo. </td>
    </tr>
    <tr>
      <td class="fw-bold">public_key</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Public key associated with the transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">signature</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Transaction signature. </td>
    </tr>
  </tbody>
</table>    
<h3 id="pactus.AccountInfo">
AccountInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about an account.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the account. </td>
    </tr>
    <tr>
      <td class="fw-bold">data</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Account data. </td>
    </tr>
    <tr>
      <td class="fw-bold">number</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Account number. </td>
    </tr>
    <tr>
      <td class="fw-bold">balance</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Account balance. </td>
    </tr>
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the account. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.BlockHeaderInfo">
BlockHeaderInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about the header of a block.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">version</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Block version. </td>
    </tr>
    <tr>
      <td class="fw-bold">prev_block_hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the previous block. </td>
    </tr>
    <tr>
      <td class="fw-bold">state_root</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>State root of the block. </td>
    </tr>
    <tr>
      <td class="fw-bold">sortition_seed</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Sortition seed of the block. </td>
    </tr>
    <tr>
      <td class="fw-bold">proposer_address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the proposer of the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.CertificateInfo">
CertificateInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about a certificate.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the certificate. </td>
    </tr>
    <tr>
      <td class="fw-bold">round</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Round of the certificate. </td>
    </tr>
    <tr>
      <td class="fw-bold">committers</td>
      <td>repeated
        <a href="#int32">int32</a>
      </td>
      <td>List of committers in the certificate. </td>
    </tr>
    <tr>
      <td class="fw-bold">absentees</td>
      <td>repeated
        <a href="#int32">int32</a>
      </td>
      <td>List of absentees in the certificate. </td>
    </tr>
    <tr>
      <td class="fw-bold">signature</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Certificate signature. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.ConsensusInfo">
ConsensusInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about consensus.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the consensus instance. </td>
    </tr>
    <tr>
      <td class="fw-bold">Active</td>
      <td>
        <a href="#bool">bool</a>
      </td>
      <td>Whether the consensus instance is active. </td>
    </tr>
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the consensus instance. </td>
    </tr>
    <tr>
      <td class="fw-bold">round</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Round of the consensus instance. </td>
    </tr>
    <tr>
      <td class="fw-bold">votes</td>
      <td>repeated
        <a href="#pactus.VoteInfo">VoteInfo</a>
      </td>
      <td>List of votes in the consensus instance. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetAccountRequest">
GetAccountRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request account information based on an address.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the account. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetAccountResponse">
GetAccountResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with account information.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">account</td>
      <td>
        <a href="#pactus.AccountInfo">AccountInfo</a>
      </td>
      <td>Account information. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockHashRequest">
GetBlockHashRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request block hash based on height.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockHashResponse">
GetBlockHashResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with the block hash.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockHeightRequest">
GetBlockHeightRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request block height based on hash.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockHeightResponse">
GetBlockHeightResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with the block height.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockRequest">
GetBlockRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request block information based on height and verbosity.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the block. </td>
    </tr>
    <tr>
      <td class="fw-bold">verbosity</td>
      <td>
        <a href="#pactus.BlockVerbosity">BlockVerbosity</a>
      </td>
      <td>Verbosity level for block information. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockResponse">
GetBlockResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with block information.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the block. </td>
    </tr>
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the block. </td>
    </tr>
    <tr>
      <td class="fw-bold">data</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Block data. </td>
    </tr>
    <tr>
      <td class="fw-bold">block_time</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Block timestamp. </td>
    </tr>
    <tr>
      <td class="fw-bold">header</td>
      <td>
        <a href="#pactus.BlockHeaderInfo">BlockHeaderInfo</a>
      </td>
      <td>Block header information. </td>
    </tr>
    <tr>
      <td class="fw-bold">prev_cert</td>
      <td>
        <a href="#pactus.CertificateInfo">CertificateInfo</a>
      </td>
      <td>Certificate information of the previous block. </td>
    </tr>
    <tr>
      <td class="fw-bold">txs</td>
      <td>repeated
        <a href="#pactus.TransactionInfo">TransactionInfo</a>
      </td>
      <td>List of transactions in the block. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetBlockchainInfoRequest">
GetBlockchainInfoRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request general information about the blockchain.</p>
 Message has no fields.  
<h3 id="pactus.GetBlockchainInfoResponse">
GetBlockchainInfoResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with general blockchain information.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">last_block_height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the last block. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_block_hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the last block. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_accounts</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Total number of accounts. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_validators</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Total number of validators. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_power</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Total power in the blockchain. </td>
    </tr>
    <tr>
      <td class="fw-bold">committee_power</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Power of the committee. </td>
    </tr>
    <tr>
      <td class="fw-bold">committee_validators</td>
      <td>repeated
        <a href="#pactus.ValidatorInfo">ValidatorInfo</a>
      </td>
      <td>List of committee validators. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetConsensusInfoRequest">
GetConsensusInfoRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request consensus information.</p>
 Message has no fields.  
<h3 id="pactus.GetConsensusInfoResponse">
GetConsensusInfoResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with consensus information.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">instances</td>
      <td>repeated
        <a href="#pactus.ConsensusInfo">ConsensusInfo</a>
      </td>
      <td>List of consensus instances. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetPublicKeyRequest">
GetPublicKeyRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request public key based on an address.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address for which public key is requested. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetPublicKeyResponse">
GetPublicKeyResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with the public key.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">public_key</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Public key of the account. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorAddressesRequest">
GetValidatorAddressesRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request validator addresses.</p>
 Message has no fields.  
<h3 id="pactus.GetValidatorAddressesResponse">
GetValidatorAddressesResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with a list of validator addresses.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">addresses</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>List of validator addresses. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorByNumberRequest">
GetValidatorByNumberRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request validator information based on a validator number.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">number</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Validator number. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorRequest">
GetValidatorRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message to request validator information based on an address.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the validator. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorResponse">
GetValidatorResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing the response with validator information.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">validator</td>
      <td>
        <a href="#pactus.ValidatorInfo">ValidatorInfo</a>
      </td>
      <td>Validator information. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.ValidatorInfo">
ValidatorInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about a validator.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the validator. </td>
    </tr>
    <tr>
      <td class="fw-bold">data</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Validator data. </td>
    </tr>
    <tr>
      <td class="fw-bold">public_key</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Public key of the validator. </td>
    </tr>
    <tr>
      <td class="fw-bold">number</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Validator number. </td>
    </tr>
    <tr>
      <td class="fw-bold">stake</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Validator stake. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_bonding_height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Last bonding height. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_sortition_height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Last sortition height. </td>
    </tr>
    <tr>
      <td class="fw-bold">unbonding_height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Unbonding height. </td>
    </tr>
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Address of the validator. </td>
    </tr>
    <tr>
      <td class="fw-bold">availability_score</td>
      <td>
        <a href="#double">double</a>
      </td>
      <td>Availability score of the validator. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.VoteInfo">
VoteInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Message containing information about a vote.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">type</td>
      <td>
        <a href="#pactus.VoteType">VoteType</a>
      </td>
      <td>Type of the vote. </td>
    </tr>
    <tr>
      <td class="fw-bold">voter</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Voter's address. </td>
    </tr>
    <tr>
      <td class="fw-bold">block_hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the block being voted on. </td>
    </tr>
    <tr>
      <td class="fw-bold">round</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Round of the vote. </td>
    </tr>
    <tr>
      <td class="fw-bold">cp_round</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Consensus round of the vote. </td>
    </tr>
    <tr>
      <td class="fw-bold">cp_value</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Consensus value of the vote. </td>
    </tr>
  </tbody>
</table>    
<h3 id="pactus.GetNetworkInfoRequest">
GetNetworkInfoRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving overall network information.</p>
 Message has no fields.  
<h3 id="pactus.GetNetworkInfoResponse">
GetNetworkInfoResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing information about the overall network.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">network_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the network. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_sent_bytes</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Total bytes sent across the network. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_received_bytes</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Total bytes received across the network. </td>
    </tr>
    <tr>
      <td class="fw-bold">connected_peers_count</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Number of connected peers. </td>
    </tr>
    <tr>
      <td class="fw-bold">connected_peers</td>
      <td>repeated
        <a href="#pactus.PeerInfo">PeerInfo</a>
      </td>
      <td>List of connected peers. </td>
    </tr>
    <tr>
      <td class="fw-bold">sent_bytes</td>
      <td>repeated
        <a href="#pactus.GetNetworkInfoResponse.SentBytesEntry">GetNetworkInfoResponse.SentBytesEntry</a>
      </td>
      <td>Bytes sent per peer ID. </td>
    </tr>
    <tr>
      <td class="fw-bold">received_bytes</td>
      <td>repeated
        <a href="#pactus.GetNetworkInfoResponse.ReceivedBytesEntry">GetNetworkInfoResponse.ReceivedBytesEntry</a>
      </td>
      <td>Bytes received per peer ID. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetNetworkInfoResponse.ReceivedBytesEntry">
GetNetworkInfoResponse.ReceivedBytesEntry
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p></p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">key</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td> </td>
    </tr>
    <tr>
      <td class="fw-bold">value</td>
      <td>
        <a href="#uint64">uint64</a>
      </td>
      <td> </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetNetworkInfoResponse.SentBytesEntry">
GetNetworkInfoResponse.SentBytesEntry
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p></p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">key</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td> </td>
    </tr>
    <tr>
      <td class="fw-bold">value</td>
      <td>
        <a href="#uint64">uint64</a>
      </td>
      <td> </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetNodeInfoRequest">
GetNodeInfoRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for retrieving information about a specific node in the</p><p>network.</p>
 Message has no fields.  
<h3 id="pactus.GetNodeInfoResponse">
GetNodeInfoResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing information about a specific node in the network.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">moniker</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Moniker of the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">agent</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Agent information of the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">peer_id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Peer ID of the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">started_at</td>
      <td>
        <a href="#uint64">uint64</a>
      </td>
      <td>Timestamp when the node started. </td>
    </tr>
    <tr>
      <td class="fw-bold">reachability</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Reachability status of the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">services</td>
      <td>repeated
        <a href="#int32">int32</a>
      </td>
      <td>List of services provided by the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">services_names</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>Names of services provided by the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">addrs</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>List of addresses associated with the node. </td>
    </tr>
    <tr>
      <td class="fw-bold">protocols</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>List of protocols supported by the node. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PeerInfo">
PeerInfo
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Information about a peer in the network.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">status</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Status of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">moniker</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Moniker of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">agent</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Agent information of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">peer_id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Peer ID of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">consensus_keys</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>Consensus keys used by the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">consensus_address</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>Consensus address of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">services</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Services provided by the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_block_hash</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Hash of the last block the peer knows. </td>
    </tr>
    <tr>
      <td class="fw-bold">height</td>
      <td>
        <a href="#uint32">uint32</a>
      </td>
      <td>Height of the peer in the blockchain. </td>
    </tr>
    <tr>
      <td class="fw-bold">received_messages</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Count of received messages. </td>
    </tr>
    <tr>
      <td class="fw-bold">invalid_messages</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Count of invalid messages received. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_sent</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Timestamp of the last sent message. </td>
    </tr>
    <tr>
      <td class="fw-bold">last_received</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td>Timestamp of the last received message. </td>
    </tr>
    <tr>
      <td class="fw-bold">sent_bytes</td>
      <td>repeated
        <a href="#pactus.PeerInfo.SentBytesEntry">PeerInfo.SentBytesEntry</a>
      </td>
      <td>Bytes sent per message type. </td>
    </tr>
    <tr>
      <td class="fw-bold">received_bytes</td>
      <td>repeated
        <a href="#pactus.PeerInfo.ReceivedBytesEntry">PeerInfo.ReceivedBytesEntry</a>
      </td>
      <td>Bytes received per message type. </td>
    </tr>
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Network address of the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">direction</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Direction of connection with the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">protocols</td>
      <td>repeated
        <a href="#string">string</a>
      </td>
      <td>List of protocols supported by the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">total_sessions</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Total sessions with the peer. </td>
    </tr>
    <tr>
      <td class="fw-bold">completed_sessions</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Completed sessions with the peer. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PeerInfo.ReceivedBytesEntry">
PeerInfo.ReceivedBytesEntry
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p></p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">key</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td> </td>
    </tr>
    <tr>
      <td class="fw-bold">value</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td> </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.PeerInfo.SentBytesEntry">
PeerInfo.SentBytesEntry
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p></p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">key</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td> </td>
    </tr>
    <tr>
      <td class="fw-bold">value</td>
      <td>
        <a href="#int64">int64</a>
      </td>
      <td> </td>
    </tr>
  </tbody>
</table>    
<h3 id="pactus.CreateWalletRequest">
CreateWalletRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for creating a new wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the new wallet. </td>
    </tr>
    <tr>
      <td class="fw-bold">mnemonic</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Mnemonic for wallet recovery. </td>
    </tr>
    <tr>
      <td class="fw-bold">language</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Language for the mnemonic. </td>
    </tr>
    <tr>
      <td class="fw-bold">password</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Password for securing the wallet. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.CreateWalletResponse">
CreateWalletResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the name of the created wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the created wallet. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorAddressRequest">
GetValidatorAddressRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for obtaining the validator address associated with a public</p><p>key.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">public_key</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Public key for which the validator address is requested. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.GetValidatorAddressResponse">
GetValidatorAddressResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the validator address corresponding to a public</p><p>key.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">address</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Validator address associated with the public key. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.LoadWalletRequest">
LoadWalletRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for loading an existing wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the wallet to load. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.LoadWalletResponse">
LoadWalletResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the name of the loaded wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the loaded wallet. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.LockWalletRequest">
LockWalletRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for locking a currently loaded wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the wallet to lock. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.LockWalletResponse">
LockWalletResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the name of the locked wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the locked wallet. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.SignRawTransactionRequest">
SignRawTransactionRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for signing a raw transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the wallet used for signing. </td>
    </tr>
    <tr>
      <td class="fw-bold">raw_transaction</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Raw transaction data to be signed. </td>
    </tr>
    <tr>
      <td class="fw-bold">password</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Password for unlocking the wallet for signing. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.SignRawTransactionResponse">
SignRawTransactionResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the transaction ID and signed raw transaction.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">transaction_id</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>ID of the signed transaction. </td>
    </tr>
    <tr>
      <td class="fw-bold">signed_raw_transaction</td>
      <td>
        <a href="#bytes">bytes</a>
      </td>
      <td>Signed raw transaction data. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.UnloadWalletRequest">
UnloadWalletRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for unloading a currently loaded wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the wallet to unload. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.UnloadWalletResponse">
UnloadWalletResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the name of the unloaded wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the unloaded wallet. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.UnlockWalletRequest">
UnlockWalletRequest
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Request message for unlocking a wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the wallet to unlock. </td>
    </tr>
    <tr>
      <td class="fw-bold">password</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Password for unlocking the wallet. </td>
    </tr>
    <tr>
      <td class="fw-bold">timeout</td>
      <td>
        <a href="#int32">int32</a>
      </td>
      <td>Timeout duration for the unlocked state. </td>
    </tr>
  </tbody>
</table>  
<h3 id="pactus.UnlockWalletResponse">
UnlockWalletResponse
<span class="badge text-bg-secondary fs-6 align-top">msg</span>
</h3>
  <p>Response message containing the name of the unlocked wallet.</p>

<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Field</td><td>Type</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
    <tr>
      <td class="fw-bold">wallet_name</td>
      <td>
        <a href="#string">string</a>
      </td>
      <td>Name of the unlocked wallet. </td>
    </tr>
  </tbody>
</table>   
 
<h3 id="pactus.PayloadType">
PayloadType
<span class="badge text-bg-info fs-6 align-top">enum</span>
</h3>
<p>Enumeration for different types of transaction payloads.</p>
<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Name</td><td>Number</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
    
      <tr>
        <td class="fw-bold">UNKNOWN</td>
        <td>0</td>
        <td>Unknown payload type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">TRANSFER_PAYLOAD</td>
        <td>1</td>
        <td>Transfer payload type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">BOND_PAYLOAD</td>
        <td>2</td>
        <td>Bond payload type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">SORTITION_PAYLOAD</td>
        <td>3</td>
        <td>Sortition payload type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">UNBOND_PAYLOAD</td>
        <td>4</td>
        <td>Unbond payload type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">WITHDRAW_PAYLOAD</td>
        <td>5</td>
        <td>Withdraw payload type.</td>
      </tr>
    
  </tbody>
</table> 
<h3 id="pactus.TransactionVerbosity">
TransactionVerbosity
<span class="badge text-bg-info fs-6 align-top">enum</span>
</h3>
<p>Enumeration for verbosity level when requesting transaction details.</p>
<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Name</td><td>Number</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
    
      <tr>
        <td class="fw-bold">TRANSACTION_DATA</td>
        <td>0</td>
        <td>Request only transaction data.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">TRANSACTION_INFO</td>
        <td>1</td>
        <td>Request detailed transaction information.</td>
      </tr>
    
  </tbody>
</table>   
<h3 id="pactus.BlockVerbosity">
BlockVerbosity
<span class="badge text-bg-info fs-6 align-top">enum</span>
</h3>
<p>Enumeration for verbosity level when requesting block information.</p>
<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Name</td><td>Number</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
    
      <tr>
        <td class="fw-bold">BLOCK_DATA</td>
        <td>0</td>
        <td>Request block data only.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">BLOCK_INFO</td>
        <td>1</td>
        <td>Request block information only.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">BLOCK_TRANSACTIONS</td>
        <td>2</td>
        <td>Request block transactions only.</td>
      </tr>
    
  </tbody>
</table> 
<h3 id="pactus.VoteType">
VoteType
<span class="badge text-bg-info fs-6 align-top">enum</span>
</h3>
<p>Enumeration for types of votes.</p>
<table class="table table-bordered table-sm">
  <thead>
    <tr><td>Name</td><td>Number</td><td>Description</td></tr>
  </thead>
  <tbody class="table-group-divider">
    
      <tr>
        <td class="fw-bold">VOTE_UNKNOWN</td>
        <td>0</td>
        <td>Unknown vote type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">VOTE_PREPARE</td>
        <td>1</td>
        <td>Prepare vote type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">VOTE_PRECOMMIT</td>
        <td>2</td>
        <td>Precommit vote type.</td>
      </tr>
    
      <tr>
        <td class="fw-bold">VOTE_CHANGE_PROPOSER</td>
        <td>3</td>
        <td>Change proposer vote type.</td>
      </tr>
    
  </tbody>
</table>      

<h3 id="scalar-value-types">Scalar Value Types</h3>
<table class="table table-bordered table-sm">
  <thead>
    <tr><td>.proto Type</td><td>C++</td><td>Java</td><td>Python</td><td>Go</td><td>C#</td><td>PHP</td></tr>
  </thead>
  <tbody class="table-group-divider"> 
      <tr id="double">
        <td class="fw-bold">double</td>
        <td>double</td>
        <td>double</td>
        <td>float</td>
        <td>float64</td>
        <td>double</td>
        <td>float</td>
      </tr> 
      <tr id="float">
        <td class="fw-bold">float</td>
        <td>float</td>
        <td>float</td>
        <td>float</td>
        <td>float32</td>
        <td>float</td>
        <td>float</td>
      </tr> 
      <tr id="int32">
        <td class="fw-bold">int32</td>
        <td>int32</td>
        <td>int</td>
        <td>int</td>
        <td>int32</td>
        <td>int</td>
        <td>integer</td>
      </tr> 
      <tr id="int64">
        <td class="fw-bold">int64</td>
        <td>int64</td>
        <td>long</td>
        <td>int/long</td>
        <td>int64</td>
        <td>long</td>
        <td>integer/string</td>
      </tr> 
      <tr id="uint32">
        <td class="fw-bold">uint32</td>
        <td>uint32</td>
        <td>int</td>
        <td>int/long</td>
        <td>uint32</td>
        <td>uint</td>
        <td>integer</td>
      </tr> 
      <tr id="uint64">
        <td class="fw-bold">uint64</td>
        <td>uint64</td>
        <td>long</td>
        <td>int/long</td>
        <td>uint64</td>
        <td>ulong</td>
        <td>integer/string</td>
      </tr> 
      <tr id="sint32">
        <td class="fw-bold">sint32</td>
        <td>int32</td>
        <td>int</td>
        <td>int</td>
        <td>int32</td>
        <td>int</td>
        <td>integer</td>
      </tr> 
      <tr id="sint64">
        <td class="fw-bold">sint64</td>
        <td>int64</td>
        <td>long</td>
        <td>int/long</td>
        <td>int64</td>
        <td>long</td>
        <td>integer/string</td>
      </tr> 
      <tr id="fixed32">
        <td class="fw-bold">fixed32</td>
        <td>uint32</td>
        <td>int</td>
        <td>int</td>
        <td>uint32</td>
        <td>uint</td>
        <td>integer</td>
      </tr> 
      <tr id="fixed64">
        <td class="fw-bold">fixed64</td>
        <td>uint64</td>
        <td>long</td>
        <td>int/long</td>
        <td>uint64</td>
        <td>ulong</td>
        <td>integer/string</td>
      </tr> 
      <tr id="sfixed32">
        <td class="fw-bold">sfixed32</td>
        <td>int32</td>
        <td>int</td>
        <td>int</td>
        <td>int32</td>
        <td>int</td>
        <td>integer</td>
      </tr> 
      <tr id="sfixed64">
        <td class="fw-bold">sfixed64</td>
        <td>int64</td>
        <td>long</td>
        <td>int/long</td>
        <td>int64</td>
        <td>long</td>
        <td>integer/string</td>
      </tr> 
      <tr id="bool">
        <td class="fw-bold">bool</td>
        <td>bool</td>
        <td>boolean</td>
        <td>boolean</td>
        <td>bool</td>
        <td>bool</td>
        <td>boolean</td>
      </tr> 
      <tr id="string">
        <td class="fw-bold">string</td>
        <td>string</td>
        <td>String</td>
        <td>str/unicode</td>
        <td>string</td>
        <td>string</td>
        <td>string</td>
      </tr> 
      <tr id="bytes">
        <td class="fw-bold">bytes</td>
        <td>string</td>
        <td>ByteString</td>
        <td>str</td>
        <td>[]byte</td>
        <td>ByteString</td>
        <td>string</td>
      </tr> 
  </tbody>
</table>
