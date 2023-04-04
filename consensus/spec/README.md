# Consensus specification

This folder contains the consensus specification for the Pactus blockchain,
which is based on the TLA+ formal language.
The specification defines the consensus algorithm used by the blockchain.

More info can be found [here](https://pactus.org/learn/consensus/specification/)

## Model checking

To run the model checker, you will need to download and install the [TLA+ Toolbox](https://lamport.azurewebsites.net/tla/toolbox.html),
which includes the TLC model checker. Follow the steps below to run the TLC model checker:

- Add the `Pactus_Liveness.tla` spec to your TLA+ Toolbox project.
- Create a new model and specify the temporal formula as `LiveSpec`.
- Specify the invariants formula as `TypeOK`.
- Specify the properties formula as `Success`.
- Since we are checking for liveness, add the `Constraint` formula as the "State Constraint".
- Define the required constants:
    - `NumFaulty`: the number of faulty nodes (e.g. 1)
    - `MaxHeight`: the maximum height of the system (e.g. 2)
    - `MaxRound`: the maximum round of the consensus algorithm (e.g. 2)
- Run the TLC checker to check the correctness of the specification.
