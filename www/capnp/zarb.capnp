using Go = import "/go.capnp";
@0x84b56bd0975dfd33;
$Go.package("capnp");
$Go.import("capnp");


struct Header {
  version             @0 :Int32;
  time                @1 :Int64;
  lastBlockHash       @2 :Data;
  stateHash           @3 :Data;
  txsHash             @4 :Data;
  lastReceiptsHash    @5 :Data;
  lastCommitHash      @6 :Data;
  committersHash      @7 :Data;
  proposerAddress     @8 :Data;
}

struct Txs {
  hashes             @0 :List(Data);
}

struct Commit {
  blockHash           @0 :Data;
  round               @1 :UInt32;
	signed              @2 :List(Int32);
  missed              @3 :List(Int32);
  signature           @4 :Data;
}

struct Block {
  header              @0 :Header;
  lastCommit          @1 :Commit;
  txs                 @2 :Txs;
}

struct BlockchainResult {

}

struct BlockResult {
  hash                @0 :Data;
  data                @1 :Data;
  block               @2 :Block;
}

struct Receipt {
  hash                @0 :Data;
  data                @1 :Data;
}

struct TransactionResult {
  id                  @0 :Data;
  data                @1 :Data;
  transaction         @2 :Data; # TODO: define tx struct
  receipt             @3 :Receipt;
}

struct AccountResult {
  data                @0 :Data;
}

struct ValidatorResult {
  data                @0 :Data;
}

struct SendTransactionResult {
  status                @0 :Int32;
  id                    @1 :Data;
}

interface ZarbServer {
  getBlockchainInfo    @0 ()                                       -> (result: BlockchainResult);
	getBlock             @1 (height: UInt64, verbosity: Int32)       -> (result :BlockResult);
	getTransaction       @2 (id: Data, verbosity: Int32)           -> (result :TransactionResult);
	getBlockHeight       @3 (hash: Data)                             -> (result :UInt64);
	getAccount           @4 (address: Data, verbosity: Int32)        -> (result :AccountResult);
	getValidator         @5 (address: Data, verbosity: Int32)        -> (result :ValidatorResult);
  sendRawTransaction   @6 (rawTx: Data)                             -> (result:SendTransactionResult);
}

