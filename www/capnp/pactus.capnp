using Go = import "go.capnp";
@0x84b56bd0975dfd33;
$Go.package("capnp");
$Go.import("capnp");


struct Header {
  version             @0 :UInt8;
  time                @1 :UInt32;
  prevBlockHash       @2 :Data;
  stateRoot           @3 :Data;
  sortitionSeed       @4 :Data;
  proposerAddress     @5 :Text;
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

struct BlockchainInfoResult {
  lastBlockHeight     @0 :UInt32;
  lastBlockHash       @1 :Data;
  committee           @2 :Committee;
}

struct BlockResult {
  height              @0 :UInt32;
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

struct Committee {
  totalPower          @0 :Int64;
  committeePower      @1 :Int64;
  validators          @2 :List(ValidatorResult);
}

struct Peer {
  status              @0 :Int32;
  moniker             @1 :Text;
  agent               @2 :Text;
  peerID              @3 :Text;
  publicKey           @4 :Text;
  lastSeen            @5 :Int32;
  flags               @6 :Int32;
  height              @7 :UInt32;
  receivedMessages    @8 :Int32;
  invalidMessages     @9 :Int32;
  receivedBytes       @10 :Int32;
}

struct NetworkInfoResult {
  peerID              @0 :Text;
  peers               @1 :List(Peer);
}

struct Vote {
  type                @0 :Int8;
  voter               @1 :Text;
  blockHash           @2 :Data;
  round               @3 :Int16;
}

struct ConsensusInfoResult {
  height              @0 :UInt32;
  round               @1 :Int16;
  votes               @2 :List(Vote);
}

struct SendTransactionResult {
  status              @0 :Int32;
  id                  @1 :Data;
}

interface PactusServer {
  getBlock            @0 (height: UInt32, verbosity: Int32) -> (result: BlockResult);
  getBlockHash        @1 (height: UInt32)                   -> (result: Data);
  getBlockHeight      @2 (hash: Data)                       -> (result: UInt32);
  getTransaction      @3 (id: Data, verbosity: Int32)       -> (result: TransactionResult);
  getAccount          @4 (address: Text)                    -> (result: AccountResult);
  getValidator        @5 (address: Text)                    -> (result: ValidatorResult);
  getBlockchainInfo   @6 ()                                 -> (result: BlockchainInfoResult);
  getNetworkInfo      @7 ()                                 -> (result: NetworkInfoResult);
  getConsensusInfo    @8 ()                                 -> (result: ConsensusInfoResult);
  sendRawTransaction  @9 (rawTx: Data)                      -> (result: SendTransactionResult);
}

