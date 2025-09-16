-------------------------------- MODULE Pactus --------------------------------
(***************************************************************************)
(* The specification of the Pactus consensus algorithm:                    *)
(* `^\url{https://docs.pactus.org/protocol/consensus/protocol/}^'          *)
(***************************************************************************)
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANT
    \* The maximum number of rounds, limiting the range of behaviors evaluated by TLC.
    MaxRound,
    \* The maximum number of change-proposer (CP) rounds, limiting the range of behaviors evaluated by TLC.
    MaxCPRound,
    \* The total number of nodes in the network, denoted as `N` in the protocol.
    N,
    \* The maximum number of faulty nodes in the network, denoted as `F` in the protocol.
    F,
    \* The indices of faulty nodes.
    FaultyNodes

VARIABLES
    \* The set of messages received by the network.
    network,
    \* The set of messages delivered to each replica.
    logs,
    \* The state of each replica in the consensus protocol.
    states

\* Helper expressions for common values.
ThreeFPlusOne == (3 * F) + 1
TwoFPlusOne   == (2 * F) + 1
OneFPlusOne   == (1 * F) + 1

\* A tuple containing all variables in the spec for ease of use in temporal conditions.
vars == <<network, logs, states>>

ASSUME
    \* Ensure the number of nodes is sufficient to tolerate the specified number of faults.
    /\ N >= ThreeFPlusOne
    \* Ensure that `FaultyNodes` is a valid subset of node indices.
    /\ FaultyNodes \subseteq 0..N-1
    \* Ensure that the number of faulty nodes does not exceed the maximum allowed.
    /\ Cardinality(FaultyNodes) <= F
    \* Ensure that `MaxRound` is greater than 0.
    /\ MaxRound > 0

-----------------------------------------------------------------------------
(***************************************************************************)
(* Helper functions                                                        *)
(***************************************************************************)

\* Check if the replica is the proposer for this round.
\* The proposer starts with the first replica and moves to the next in the change-proposer phase.
IsProposer(index) ==
    states[index].round % N = index

\* Check if a node is faulty.
IsFaulty(index) == index \in FaultyNodes

\* Returns a subset of `bag` where each element matches all criteria specified in `params`.
SubsetOfMsgs(bag, params) ==
   {i \in bag: \A field \in DOMAIN params: i[field] = params[field]}

\* Check if the node has received at least `3f+1` PRECOMMIT votes for a proposal in the current round.
HasPreCommitAbsolute(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "PRECOMMIT",
        round    |-> states[index].round])) >= ThreeFPlusOne

\* Check if the node has received at least `2f+1` PRECOMMIT votes for a proposal in the current round.
HasPreCommitQuorum(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "PRECOMMIT",
        round    |-> states[index].round])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:PRE-VOTE votes in the current CP round.
CPHasPreVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:PRE-VOTE votes with value 1 (yes) in the current CP round.
CPHasPreVotesQuorumForYes(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:PRE-VOTE votes with value 0 (no) in the current CP round.
CPHasPreVotesQuorumForNo(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= TwoFPlusOne

\* Check if the node has received at least `f+1` CP:PRE-VOTE votes with value 1 (yes) in the current CP round.
CPHasPreVotesMinorityForYes(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= OneFPlusOne

\* Check if the node has received both yes and no CP:PRE-VOTE votes in the current CP round.
CPHasPreVotesForYesAndNo(index) ==
    /\ Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= 1
    /\ Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= 1

\* Check if the node has received at least one CP:MAIN-VOTE with value 0 (no) in the previous CP round.
CPHasOneMainVotesNoInPrvRound(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 0])) > 0

\* Check if the node has received at least one CP:MAIN-VOTE with value 1 (yes) in the previous CP round.
CPHasOneMainVotesYesInPrvRound(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 1])) > 0

\* Check if the node has received at least `2f+1` CP:MAIN-VOTE votes with value 2 (abstain) in the previous CP round.
CPAllMainVotesAbstainInPrvRound(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 2])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:MAIN-VOTE votes in the current CP round.
CPHasMainVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:MAIN-VOTE votes with value 1 (yes) in the current CP round.
CPHasMainVotesQuorumForYes(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= TwoFPlusOne

\* Check if the node has received at least `2f+1` CP:MAIN-VOTE votes with value 2 (abstain) in the current CP round.
CPHasMainVotesQuorumForAbstain(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 2])) >= TwoFPlusOne

\* Check if the node has received at least one CP:DECIDED vote with value 1 (yes) in the current round.
CPHasDecideVotesForYes(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "CP:DECIDED",
        round    |-> states[index].round,
        cp_val   |-> 1])) > 0

\* Check if the node has received a proposal in the current round.
HasProposal(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "PROPOSAL",
        round    |-> states[index].round])) > 0

\* Check if the node has sent its own PRECOMMIT vote in the current round.
HasPrecommited(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "PRECOMMIT",
        round    |-> states[index].round,
        index    |-> index])) = 1

\* Check if the node has received an announcement message in the current round.
HasAnnouncement(index) ==
    Cardinality(SubsetOfMsgs(logs[index], [
        type     |-> "ANNOUNCEMENT",
        round    |-> states[index].round])) > 0

\* Check if the proposal is committed.
\* A proposal is considered committed if a super-majority of non-faulty replicas announce the same proposal.
IsCommitted ==
    LET subset == SubsetOfMsgs(network, [type |-> "ANNOUNCEMENT"])
    IN /\ Cardinality(subset) >= TwoFPlusOne
       /\ \A m1, m2 \in subset : m1.round = m2.round

-----------------------------------------------------------------------------
(***************************************************************************)
(* Network functions                                                       *)
(***************************************************************************)

\* Simulate a replica sending a message by appending it to the `network`.
\* The message is delivered to the sender's log immediately.
SendMsg(msg) ==
    /\ network' = network \cup {msg}
    /\ logs' = [logs EXCEPT ![msg.index] = logs[msg.index] \cup {msg}]

\* Deliver a message to the specified replica's log.
DeliverMsg(index) ==
    LET undeliveredMsgs == network \ logs[index]
    IN IF Cardinality(undeliveredMsgs) = 0 THEN
        UNCHANGED <<vars>>
    ELSE
        LET msg == CHOOSE x \in undeliveredMsgs: TRUE
        IN
            /\ logs' = [logs EXCEPT ![index] = logs[index] \cup {msg}]
            /\ UNCHANGED <<states, network>>

\* Broadcast a PROPOSAL message into the network.
SendProposal(index) ==
    SendMsg([
        type     |-> "PROPOSAL",
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

\* Broadcast PRECOMMIT votes into the network.
SendPreCommitVote(index) ==
    SendMsg([
        type     |-> "PRECOMMIT",
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

\* Broadcast CP:PRE-VOTE votes into the network.
SendCPPreVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:PRE-VOTE",
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val])

\* Broadcast CP:MAIN-VOTE votes into the network.
SendCPMainVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:MAIN-VOTE",
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val])

\* Broadcast CP:DECIDED votes into the network.
SendCPDecideVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:DECIDED",
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        index    |-> index,
        cp_val   |-> cp_val])

\* Broadcast ANNOUNCEMENT messages into the network.
Announce(index)  ==
    SendMsg([
        type     |-> "ANNOUNCEMENT",
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

-----------------------------------------------------------------------------
(***************************************************************************)
(* State transition functions                                              *)
(***************************************************************************)

\* Transition to the propose state.
Propose(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "propose"
    /\
        IF IsProposer(index) THEN
            SendProposal(index)
        ELSE
            UNCHANGED <<logs, network>>
    /\ states' = [states EXCEPT ![index].name = "precommit"]

\* Transition to the precommit state.
PreCommit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "precommit"
    /\ HasProposal(index)
    /\ SendPreCommitVote(index)
    /\ states' = states

\* AbsoluteCommit checks if 3F+1 replicas voted for the proposal.
AbsoluteCommit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name # "commit" \* to prevent shuttering
    /\ HasPreCommitAbsolute(index)
    /\ states' = [states EXCEPT ![index].name = "commit"]
    /\ UNCHANGED <<network, logs>>

\* QuorumCommit checks if 2F+1 replicas voted for the proposal after the change-proposer phase.
QuorumCommit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "precommit"
    /\ states[index].decided = TRUE
    /\ HasPreCommitQuorum(index)
    /\ states' = [states EXCEPT ![index].name = "commit"]
    /\ UNCHANGED <<network, logs>>

\* Transition to the commit state.
Commit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "commit"
    /\ Announce(index)
    /\ UNCHANGED <<states>>

\* Transition for timeout: a non-faulty replica changes the proposer if its timer expires.
Timeout(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "precommit"
    /\ states[index].decided = FALSE
    /\
        \* To limit the the behaviors.
        \/ states[index].round < MaxRound
        \/ HasPreCommitQuorum(index)
    /\ states' = [states EXCEPT ![index].name = "cp:pre-vote"]
    /\ UNCHANGED <<network, logs>>

\* Transition to the CP pre-vote state.
CPPreVote(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:pre-vote"
    /\
        IF states[index].cp_round = 0 THEN
            IF ~HasPrecommited(index) THEN
                /\ SendCPPreVote(index, 1)
                /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
            ELSE IF Cardinality(
                            SubsetOfMsgs(logs[index], [type |-> "PRECOMMIT",   round |-> states[index].round]) \cup
                            SubsetOfMsgs(logs[index], [type |-> "CP:PRE-VOTE", round |-> states[index].round, cp_round |-> states[index].cp_round])
                        ) >= TwoFPlusOne THEN
                    \* Check if there is quorum of PRECOMMIT votes
                    IF HasPreCommitQuorum(index) THEN
                        /\ SendCPPreVote(index, 0)
                        /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
                    ELSE
                        /\ SendCPPreVote(index, 1)
                        /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
            ELSE
                /\ UNCHANGED <<vars>>
        ELSE
            /\
                \/
                    /\ CPHasOneMainVotesNoInPrvRound(index)
                    /\ SendCPPreVote(index, 0)
                \/
                    /\ CPHasOneMainVotesYesInPrvRound(index)
                    /\ SendCPPreVote(index, 1)
                \/
                    /\ CPAllMainVotesAbstainInPrvRound(index)
                    /\ SendCPPreVote(index, 0) \* biased to zero when all votes abstain
            /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]

\* Transition to the CP main-vote state.
CPMainVote(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:main-vote"
    /\ CPHasPreVotesQuorum(index)
    /\
        \/
            \* all votes for 0
            /\ CPHasPreVotesQuorumForNo(index)
            /\ states' = [states EXCEPT ![index].name = "precommit",
                                        ![index].decided = TRUE]
            /\ UNCHANGED <<network, logs>>
        \/
            \* all votes for 1
            /\ CPHasPreVotesQuorumForYes(index)
            /\ SendCPMainVote(index, 1)
            /\ states' = [states EXCEPT ![index].name = "cp:decide"]
        \/
            \* Abstain vote
            /\ CPHasPreVotesForYesAndNo(index)
            /\ SendCPMainVote(index, 2) \* Abstain
            /\ states' = [states EXCEPT ![index].name = "cp:decide"]

\* Transition to the CP decide state.
CPDecide(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:decide"
    /\  CPHasMainVotesQuorum(index)
    /\
        \/
            /\ CPHasMainVotesQuorumForYes(index)
            /\ SendCPDecideVote(index, 1)
            /\ states' = [states EXCEPT ![index].name = "propose",
                                        ![index].round = states[index].round + 1]
        \/
            /\ states[index].cp_round < MaxCPRound
            /\ CPHasMainVotesQuorumForAbstain(index)
            /\ states' = [states EXCEPT ![index].name = "cp:pre-vote",
                                        ![index].cp_round = states[index].cp_round + 1]
            /\ UNCHANGED <<network, logs>>

\* Transition for strong termination of Change-Proposer phase.
CPStrongTerminate(index) ==
    /\ ~IsFaulty(index)
    /\
        \/ states[index].name = "cp:pre-vote"
        \/ states[index].name = "cp:main-vote"
        \/ states[index].name = "cp:decide"
    /\
        \* To limit the the behaviors.
        IF /\ states[index].cp_round = MaxCPRound
           /\ HasPreCommitQuorum(index) THEN
               /\ states' = [states EXCEPT ![index].name = "precommit",
                                           ![index].decided = TRUE]

        ELSE IF CPHasDecideVotesForYes(index) THEN
            /\ states' = [states EXCEPT ![index].name = "propose",
                                        ![index].round = states[index].round + 1,
                                        ![index].cp_round = 0]
        ELSE
            /\ states' = states

    /\ UNCHANGED <<network, logs>>

-----------------------------------------------------------------------------

\* Initial state
Init ==
    /\ network = {}
    /\ logs = [index \in 0..N-1 |-> {}]
    /\ states = [index \in 0..N-1 |-> [
        name       |-> "propose",
        decided    |-> FALSE,
        round      |-> 0,
        cp_round   |-> 0]]

\* State transition relation
Next ==
    \E index \in 0..N-1:
        \/ Propose(index)
        \/ PreCommit(index)
        \/ Timeout(index)
        \/ Commit(index)
        \/ AbsoluteCommit(index)
        \/ QuorumCommit(index)
        \/ CPPreVote(index)
        \/ CPMainVote(index)
        \/ CPDecide(index)
        \/ CPStrongTerminate(index)
        \/ DeliverMsg(index)

\* Specification
Spec ==
    Init /\ [][Next]_vars /\ WF_vars(Next)

(***************************************************************************)
(* Success: All non-faulty nodes eventually commit.                        *)
(***************************************************************************)
Success == <>(IsCommitted)

(***************************************************************************)
(* TypeOK is the type-correctness invariant.                               *)
(***************************************************************************)
TypeOK ==
    /\ \A index \in 0..N-1:
        /\ states[index].round <= MaxRound
        /\ states[index].cp_round <= MaxCPRound
        /\ states[index].round >= 0
        /\ states[index].cp_round >= 0
        /\ states[index].name \in {"propose", "precommit", "commit", "cp:pre-vote", "cp:main-vote", "cp:decide"}
        /\ states[index].decided \in {TRUE, FALSE}
        /\ states[index].name = "propose" /\ states[index].round > 0 =>
            /\ Cardinality(SubsetOfMsgs(network, [
                type     |-> "CP:DECIDED",
                round    |-> states[index].round-1,
                cp_val   |-> 1])) > 0
            /\ Cardinality(SubsetOfMsgs(network, [
                type     |-> "ANNOUNCEMENT",
                round    |-> states[index].round-1])) = 0
        /\ states[index].name = "commit" =>
            /\ Cardinality(SubsetOfMsgs(network, [
                type     |-> "PRECOMMIT",
                round    |-> states[index].round])) >= TwoFPlusOne
            /\ Cardinality(SubsetOfMsgs(network, [
                type     |-> "PROPOSAL",
                round    |-> states[index].round])) = 1
            /\  LET subset == SubsetOfMsgs(network, [type |-> "ANNOUNCEMENT"])
                IN /\ \A m1, m2 \in subset : m1.round = m2.round
    /\ \A msg \in network:
        /\ msg.round <= MaxRound
        /\ msg.cp_round <= MaxCPRound
        /\ msg.round >= 0
        /\ msg.cp_round >= 0
        /\ msg.type \in {"PROPOSAL", "PRECOMMIT", "CP:PRE-VOTE", "CP:MAIN-VOTE", "CP:DECIDED", "ANNOUNCEMENT"}
        /\ msg.index \in 0..N-1



=============================================================================
