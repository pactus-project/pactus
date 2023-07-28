# Consensus specification

This folder contains the consensus specification for the Pactus blockchain,
which is based on the TLA+ formal language.
The specification defines the consensus algorithm used by the blockchain.

More info can be found [here](https://pactus.org/learn/consensus/specification/)

## Model checking

To run the model checker, you will need to download and install the [TLA+ Toolbox](https://lamport.azurewebsites.net/tla/toolbox.html),
which includes the TLC model checker. Follow the steps below to run the TLC model checker:

- Add the `Pactus.tla` spec to your TLA+ Toolbox project.
- Create a new model and specify a temporal formula as `Spec`.
- Specify an invariants formula as `TypeOK`.
- Specify a properties formula as `Success`.
- Define the required constants:
    - `NumFaulty`: the number of faulty nodes (e.g. 1)
    - `FaultyNodes`: the index of faulty nodes (e.g. {3})
    - `MaxHeight`: the maximum height of the system (e.g. 1)
    - `MaxRound`: the maximum block-creation round of the consensus algorithm (e.g. 1)
    - `MaxCPRound`: the maximum change-proposer round of the consensus algorithm (e.g. 1)
- Run the TLC checker to check the correctness of the specification.
