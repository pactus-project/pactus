# Changelog

## 0.9.3
- Adding sortition_seed to block header
- Sortition runs on sortition_seed instead of block_hash
- Updating round-robin mechanism for choosing proposer
- Block certificate includes committers and absentees

## 0.9.2
- Refactoring commit base on zip-0001 proposal
- Introducing Download topic for downloading blockchain
- Transaction ID is hash of transaction without signature and public-key
- Using blake2b for hashing
- Using ARC cache for syncing module
- Reap transaction from txpool before proposing block
- Executing transaction
- Updating capnp and http server
- Making sandbox more isolate
- Updating state merkle tree
- Add number to account and validator structure (Used for Merklizing state)
- Add chain params to genesis
- Stamp validation check for transaction
- Refactoring transactions and transaction receipts
- Assigning version 1001 for testnet blocks
- Generating keys based on BIP-0039 (mnemonic or seed phrase)
- Add new command for recovering key by seed
- Add new command for sending raw transaction
- Add new command for making `Bond` transaction
- Add new argument for signing transaction to `key sign` command
- Improve consensus mechanism
- Improve syncing process

## 0.9.1
- Ensure messages belongs to same network

## 0.9.0
- Refactoring block structure
- Refactoring Certificate structure
- Aggregating validators' signatures
- Validating Certificate and committers
- Calculating root hash of committers
- Report UndefHash as a sanity error
- Try to load last state info upon starting the node
- Saving the last state info when a new block is committed
- Updating store interface
- Adding more tests for consensus, state and txPool

## 0.8.0
- Adding syncer package for syncing blockchain
- Adding message package that includes network messages
- Decoupling network from state and consensus package
- Kademlia DHT added to network package
- Keeping statistical report for nodes and peers
- Reporting number of invalid messages from peers
- Adding tests for consensus and sync package
- Move config files into packages
- Adding `TestConfig` method for testing purpose
- Detecting duplicated vote, tests are included
- Adding helper methods for testing purpose

## 0.7.0

 First version