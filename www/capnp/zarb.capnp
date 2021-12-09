using Go = import "go.capnp";
@0x84b56bd0975dfd33;
$Go.package("capnp");
$Go.import("capnp");


struct Header {
  version             @0 :Int32;
  time                @1 :Int64;
  prevBlockHash       @2 :Data;
  stateHash           @3 :Data;
  txsHash             @4 :Data;
  prevCertificateHash @5 :Data;
  sortitionSeed       @6 :Data;
  proposerAddress     @7 :Data;
}

struct Txs {
  hashes              @0 :List(Data);
}

struct Certificate {
  blockHash           @0 :Data;
  round               @1 :UInt32;
  committers          @2 :List(Int32);
  absentees            @3 :List(Int32);
  signature           @4 :Data;
}

struct Block {
  header              @0 :Header;
  lastCertificate     @1 :Certificate;
  txs                 @2 :Txs;
}

struct BlockchainResult {
  height             @0 :Int64;
}

struct BlockResult {
  hash                @0 :Data;
  data                @1 :Data;
  block               @2 :Block;
}


struct TransactionResult {
  id                  @0 :Data;
  data                @1 :Data;
  transaction         @2 :Data; # TODO: define tx struct
}

struct AccountResult {
  data                @0 :Data;
}

struct ValidatorResult {
  data                @0 :Data;
}

struct Peer {
  status                @0 :Int32;
  moniker               @1 :Text;
  nodeVersion           @2 :Text;
  peerID                @3 :Text;
  publicKey             @4 :Text;
  initialBlockDownload  @5 :Bool;
  height                @6 :Int32;
  receivedMessages      @7 :Int32;
  invalidMessages       @8 :Int32;
  receivedBytes         @9 :Int32;
}

struct NetworkResult {
  peerID              @0 :Text;
  peers               @1 :List(Peer);
}

struct SendTransactionResult {
  status              @0 :Int32;
  id                  @1 :Data;
}

interface ZarbServer {
  getBlock             @0 (height: UInt64, verbosity: Int32)       -> (result :BlockResult);
  getBlockHeight       @1 (hash: Data)                             -> (result :UInt64);
  getTransaction       @2 (id: Data, verbosity: Int32)             -> (result :TransactionResult);
  getAccount           @3 (address: Data, verbosity: Int32)        -> (result :AccountResult);
  getValidator         @4 (address: Data, verbosity: Int32)        -> (result :ValidatorResult);
  getBlockchainInfo    @5 ()                                       -> (result :BlockchainResult);
  getNetworkInfo       @6 ()                                       -> (result :NetworkResult);
  sendRawTransaction   @7 (rawTx: Data)                            -> (result :SendTransactionResult);
}

