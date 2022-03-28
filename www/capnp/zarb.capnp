using Go = import "go.capnp";
@0x84b56bd0975dfd33;
$Go.package("capnp");
$Go.import("capnp");


struct Header {
  version             @0 :UInt8;
  time                @1 :Int32;
  prevBlockHash       @2 :Data;
  stateRoot           @3 :Data;
  sortitionSeed       @4 :Data;
  proposerAddress     @5 :Data;
}

struct Certificate {
  round               @0 :Int16;
  committers          @1 :List(Int32);
  absentees           @2 :List(Int32);
  signature           @3 :Data;
}

struct Block {
  header              @0 :Header;
  prevCert            @1 :Certificate;
  txs                 @2 :List(Data);
}

struct BlockchainResult {
  lastBlockHeight     @0 :Int32;
  lastBlockHash       @1 :Data;
}

struct BlockResult {
  height              @0 :Data;
  hash                @1 :Data;
  block               @2 :Block;
  data                @3 :Data;
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
  agent                 @2 :Text;
  peerID                @3 :Text;
  publicKey             @4 :Text;
  flags                 @5 :Int32;
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
  getBlock             @0 (hash: Data, verbosity: Int32)          -> (result :BlockResult);
  getBlockHash         @1 (height: Int32)                         -> (result :Data);
  getTransaction       @2 (id: Data, verbosity: Int32)            -> (result :TransactionResult);
  getAccount           @3 (address: Data, verbosity: Int32)       -> (result :AccountResult);
  getValidator         @4 (address: Data, verbosity: Int32)       -> (result :ValidatorResult);
  getBlockchainInfo    @5 ()                                      -> (result :BlockchainResult);
  getNetworkInfo       @6 ()                                      -> (result :NetworkResult);
  sendRawTransaction   @7 (rawTx: Data)                           -> (result :SendTransactionResult);
}

