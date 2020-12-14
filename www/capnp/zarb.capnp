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

struct Committer {
  address             @0 :Data;
	status              @1 :Int32;
}

struct Commit {
  round               @0 :UInt32;
	signature           @1 :Data;
	committers          @2 :List(Committer);
}

struct Block {
  header              @0 :Header;
  lastCommit          @1 :Commit;
  txs                 @2 :Txs;
}

struct BlockInfo {
  hash                @0 :Data;
  height              @1 :UInt32;
  data                @2 :Data;
  block               @3 :Block;
}


struct Tx {
  stamp               @0 :Data;
	sender              @1 :Data;
	receiver            @2 :Data;
	amount              @3 :UInt64;
	fee                 @4 :UInt64;
	data                @5 :Data;
	memo                @6 :Text;
}

struct TxInfo {
  hash                @0 :Data;
  height              @1 :UInt32;
  data                @2 :Data;
  tx                  @3 :Tx;
}


interface ZarbServer {
	blockAt @0 (height: UInt32) -> (blockInfo :BlockInfo);
	block @1 (hash: Data) -> (blockInfo :BlockInfo);
	tx @2 (hash: Data) -> (txInfo :BlockInfo);

}

