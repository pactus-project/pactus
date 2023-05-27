-------------------------------- MODULE Pactus --------------------------------
(***************************************************************************)
(* The specification of the Pactus consensus algorithm based on            *)
(* Practical Byzantine Fault Tolerant.                                     *)
(* For more information check here:                                        *)
(* `^\url{https://pactus.org/learn/consensus/protocol/}^'                  *)
(***************************************************************************)
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANT
    \* The total number of faulty nodes
    NumFaulty,
    \* The maximum number of round per height.
    \* this is to restrict the allowed behaviours that TLC scans through.
    MaxRound

ASSUME
    /\ NumFaulty >= 1

VARIABLES
    log,
    states

\* Total number of replicas that is `3f+1' where `f' is number of faulty nodes.
Replicas == (3 * NumFaulty) + 1
\* 2/3 of total replicas that is `2f+1'
QuorumCnt == (2 * NumFaulty) + 1
\* 1/3 of total replicas that is `f+1'
OneThird == NumFaulty + 1

\* A tuple with all variables in the spec (for ease of use in temporal conditions)
vars == <<states, log>>

-----------------------------------------------------------------------------
(***************************************************************************)
(* Helper functions                                                        *)
(***************************************************************************)
\* Fetch a subset of messages in the network based on the params filter.
SubsetOfMsgs(params) ==
    {msg \in log: \A field \in DOMAIN params: msg[field] = params[field]}


\* IsProposer checks if the replica is the proposer for this round
\* To simplify, we assume the proposer always starts with the first replica
\* and moves to the next by the change-proposer phase.
IsProposer(index) ==
    states[index].round % Replicas = index

\* HasPrepareQuorum checks if there is a quorum of the PREPARE votes in each round.
HasPrepareQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "PREPARE",
        height |-> states[index].height,
        round  |-> states[index].round])) >= QuorumCnt

\* HasPrecommitQuorum checks if there is a quorum of the PRECOMMIT votes in each round.
HasPrecommitQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "PRECOMMIT",
        height |-> states[index].height,
        round  |-> states[index].round])) >= QuorumCnt

\* HasChangeProposerQuorum checks if there is a quorum of the CHANGE-PROPOSER votes in each round.
HasChangeProposerQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "CHANGE-PROPOSER",
        height |-> states[index].height,
        round  |-> states[index].round])) >= QuorumCnt

HasOneThirdOfChangeProposer(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "CHANGE-PROPOSER",
        height |-> states[index].height,
        round  |-> states[index].round])) >= OneThird

GetProposal(height, round) ==
    SubsetOfMsgs([type |-> "PROPOSAL", height |-> height, round |-> round])

HasProposal(height, round) ==
    Cardinality(GetProposal(height, round)) > 0

IsCommitted(height) ==
    Cardinality(SubsetOfMsgs([type |-> "BLOCK-ANNOUNCE", height |-> height])) > 0

HasVoted(index, type) ==
    Cardinality(SubsetOfMsgs([
        type |-> type,
        height |-> states[index].height,
        round |-> states[index].round,
        index |-> index])) > 0

-----------------------------------------------------------------------------
(***************************************************************************)
(* Network functions                                                       *)
(***************************************************************************)

\* SendMsg broadcasts the message iff the current height is not committed yet.
SendMsg(msg) ==
    IF ~IsCommitted(msg.height)
    THEN log' = log \cup {msg}
    ELSE log' = log


\* SendProposal is used to broadcast the PROPOSAL into the network.
SendProposal(index) ==
    SendMsg([
        type    |-> "PROPOSAL",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index])

\* SendPrepareVote is used to broadcast PREPARE votes into the network.
SendPrepareVote(index) ==
    SendMsg([
        type    |-> "PREPARE",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index])

\* SendPrecommitVote is used to broadcast PRECOMMIT votes into the network.
SendPrecommitVote(index) ==
    SendMsg([
        type    |-> "PRECOMMIT",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index])

\* SendChangeProposerRequest is used to broadcast CHANGE-PROPOSER votes into the network.
SendChangeProposerRequest(index) ==
    SendMsg([
        type    |-> "CHANGE-PROPOSER",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index])

\* AnnounceBlock announces the block for the current height and clears the logs.
AnnounceBlock(index)  ==
    log' = {msg \in log: (msg.type = "BLOCK-ANNOUNCE") \/ msg.height > states[index].height } \cup {[
        type    |-> "BLOCK-ANNOUNCE",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> -1]}

IsFaulty(index) == index >= 3 * NumFaulty
-----------------------------------------------------------------------------
(***************************************************************************)
(* States functions                                                        *)
(***************************************************************************)

\* NewHeight state
NewHeight(index) ==
    /\ states[index].name = "new-height"
    /\ ~IsFaulty(index)
    /\ states' = [states EXCEPT
        ![index].name = "propose",
        ![index].height = states[index].height + 1,
        ![index].round = 0]
    /\ UNCHANGED <<log>>


\* Propose state
Propose(index) ==
    /\ states[index].name = "propose"
    /\ ~IsFaulty(index)
    /\ IF IsProposer(index)
       THEN SendProposal(index)
       ELSE log' = log
    /\ states' = [states EXCEPT
        ![index].name = "prepare"]


\* Prepare state
Prepare(index) ==
    /\ states[index].name = "prepare"
    /\ ~IsFaulty(index)
    /\ ~HasOneThirdOfChangeProposer(index) \* This check is optional
    /\ HasProposal(states[index].height, states[index].round)
    /\ SendPrepareVote(index)
    /\ IF  /\ HasPrepareQuorum(index)
       THEN states' = [states EXCEPT ![index].name = "precommit"]
       ELSE states' = states


\* Precommit state
Precommit(index) ==
    /\ states[index].name = "precommit"
    /\ ~IsFaulty(index)
    /\ ~HasOneThirdOfChangeProposer(index) \* This check is optional
    /\ SendPrecommitVote(index)
    /\ IF  /\ HasPrecommitQuorum(index)
           /\ HasVoted(index, "PRECOMMIT")
       THEN states' = [states EXCEPT ![index].name = "commit"]
       ELSE states' = states


Timeout(index) ==
    /\
        \/ states[index].name = "prepare"
        \/ states[index].name = "precommit"
    /\ ~IsFaulty(index)
    /\ (states[index].round < MaxRound)
    /\ SendChangeProposerRequest(index)
    /\ states' = [states EXCEPT
        ![index].name = "change-proposer"]


Byzantine(index) ==
    /\ IsFaulty(index)
    /\ SendChangeProposerRequest(index)
    /\ states' = [states EXCEPT
        ![index].name = "change-proposer"]

\* Commit state
Commit(index) ==
    /\ states[index].name = "commit"
    /\ ~IsFaulty(index)
    /\ AnnounceBlock(index) \* this clear the logs
    /\ states' = [states EXCEPT
        ![index].name = "new-height"]


\* ChangeProposer state
ChangeProposer(index) ==
    /\ states[index].name = "change-proposer"
    /\ ~IsFaulty(index)
    /\ IF HasChangeProposerQuorum(index)
       THEN states' = [states EXCEPT
            ![index].name = "propose",
            ![index].round = states[index].round + 1]
       ELSE states' = states
    /\ UNCHANGED <<log>>

\* Sync checks the log for the committed blocks at the current height.
\* If such a block exists, it commits and moves to the next height.
Sync(index) ==
    LET
        blocks == SubsetOfMsgs([type |-> "BLOCK-ANNOUNCE", height |-> states[index].height])
    IN
        /\ Cardinality(blocks) > 0
        /\ states' = [states EXCEPT
            ![index].name = "propose",
            ![index].height = states[index].height + 1,
            ![index].round = 0]
        /\ log' = log

-----------------------------------------------------------------------------

Init ==
    /\ log = {}
    /\ states = [index \in 0..Replicas-1 |-> [
        name            |-> "new-height",
        height          |-> 0,
        round           |-> 0]]

Next ==
    \E index \in 0..Replicas-1:
       \/ Sync(index)
       \/ NewHeight(index)
       \/ Propose(index)
       \/ Prepare(index)
       \/ Precommit(index)
       \/ Timeout(index)
       \/ Commit(index)
       \/ ChangeProposer(index)
       \/ Byzantine(index)

Spec ==
    Init /\ [][Next]_vars

(***************************************************************************)
(* TypeOK is the type-correctness invariant.                               *)
(***************************************************************************)
TypeOK ==
    /\ \A index \in 0..Replicas-1:
        /\ states[index].name \in {"new-height", "propose", "prepare",
            "precommit", "commit", "change-proposer"}
        /\ ~IsCommitted(states[index].height) =>
            /\ states[index].name = "new-height" /\ states[index].height > 1 =>
                IsCommitted(states[index].height - 1)
            /\ states[index].name = "propose" /\ states[index].round > 0 =>
                /\ Cardinality(SubsetOfMsgs([
                        type   |-> "CHANGE-PROPOSER",
                        height |-> states[index].height,
                        round  |-> states[index].round - 1])) >= QuorumCnt
            /\ states[index].name = "precommit" =>
                /\ HasPrepareQuorum(index)
                /\ HasProposal(states[index].height, states[index].round)
            /\ states[index].name = "commit" =>
                /\ HasPrepareQuorum(index)
                /\ HasPrecommitQuorum(index)
                /\ HasProposal(states[index].height, states[index].round)
            /\ \A round \in 0..states[index].round:
                \* Not more than one proposal per round
                /\ Cardinality(GetProposal(states[index].height, round)) <= 1

=============================================================================
