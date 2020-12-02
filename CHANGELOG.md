# Changelog

## 0.9.1
- Adding wizard argument to start command
- Ensure messages belongs to same network

## 0.9.0
- Refactoring block structure
- Refactoring Commit structure
- Aggregating validators' signatures
- Validating Commit and committers
- Merklizing committers
- Repost UndefHash as a sanity error
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