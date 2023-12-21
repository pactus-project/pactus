-------------------------------- MODULE Pactus --------------------------------
(***************************************************************************)
(* The specification of the Pactus consensus algorithm:                    *)
(* `^\url{https://pactus.org/learn/consensus/protocol/}^'                  *)
(***************************************************************************)
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANT
    \* The maximum number of height.
    \* This limits the range of behaviors evaluated by TLC
    MaxHeight,
    \* The maximum number of round per height.
    \* This limits the range of behaviors evaluated by TLC
    MaxRound,
    \* The maximum number of cp-round per height.
    \* This limits the range of behaviors evaluated by TLC
    MaxCPRound,
    \* The total number of nodes in the network, denoted as `n` in the protocol.
    NumNodes,
    \* The total number of faulty nodes, denoted as `f` in the protocol.
    NumFaulty,
    \* The indices of faulty nodes.
    FaultyNodes

VARIABLES
    \* `log` is a set of messages received by the system.
    log,
    \* `states` represents the state of each replica in the consensus protocol.
    states

\* ThreeFPlusOne is equal to `3f+1', where `f' is the number of faulty nodes.
ThreeFPlusOne == (3 * NumFaulty) + 1
\* TwoFPlusOne is equal to `2f+1', where `f' is the number of faulty nodes.
TwoFPlusOne == (2 * NumFaulty) + 1
\* OneFPlusOne is equal to `f+1', where `f' is the number of faulty nodes.
OneFPlusOne == (1 * NumFaulty) + 1


\* A tuple containing all variables in the spec (for ease of use in temporal conditions).
vars == <<states, log>>

ASSUME
    \* Ensure that the number of nodes is sufficient to tolerate the specified number of faults.
    /\ NumNodes >= ThreeFPlusOne
    \* Ensure that `FaultyNodes` is a valid subset of node indices.
    /\ FaultyNodes \subseteq 0..NumNodes-1

-----------------------------------------------------------------------------
(***************************************************************************)
(* Helper functions                                                        *)
(***************************************************************************)

\* Fetch a subset of messages in the network based on the params filter.
SubsetOfMsgs(params) ==
   {msg \in log: \A field \in DOMAIN params: msg[field] = params[field]}

\* IsProposer checks if the replica is the proposer for this round.
\* To simplify, we assume the proposer always starts with the first replica,
\* and moves to the next by the change-proposer phase.
IsProposer(index) ==
    states[index].round % NumNodes = index

\* Helper function to check if a node is faulty or not.
IsFaulty(index) == index \in FaultyNodes

\* HasPrepareAbsoluteQuorum checks whether the node with the given index
\* has received all the PREPARE votes in this round.
HasPrepareAbsoluteQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round])) >= ThreeFPlusOne

\* HasPrepareQuorum checks whether the node with the given index
\* has received 2f+1 the PREPARE votes in this round.
HasPrepareQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round])) >= TwoFPlusOne

\* HasPrecommitQuorum checks whether the node with the given index
\* has received 2f+1 the PRECOMMIT votes in this round.
HasPrecommitQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PRECOMMIT",
        height   |-> states[index].height,
        round    |-> states[index].round])) >= TwoFPlusOne

CPHasPreVotesMinorityQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> 0,
        cp_val   |-> 1])) >= OneFPlusOne

CPHasPreVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= TwoFPlusOne

CPHasPreVotesQuorumForOne(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= TwoFPlusOne

CPHasPreVotesQuorumForZero(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= TwoFPlusOne

CPHasPreVotesForZeroAndOne(index) ==
    /\ Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= 1
    /\ Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= 1

CPHasAMainVotesZeroInPrvRound(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 0])) > 0

CPHasAMainVotesOneInPrvRound(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 1])) > 0

CPAllMainVotesAbstainInPrvRound(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 2])) >= TwoFPlusOne

CPOneFPlusOneMainVotesAbstainInPrvRound(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 2])) >= OneFPlusOne

CPHasMainVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= TwoFPlusOne

CPHasMainVotesQuorumForOne(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= TwoFPlusOne

CPHasMainVotesQuorumForZero(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= TwoFPlusOne

CPHasDecideVotesForZero(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:DECIDE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_val   |-> 0])) > 0

CPHasDecideVotesForOne(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:DECIDE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_val   |-> 1])) > 0

GetProposal(height, round) ==
    SubsetOfMsgs([type |-> "PROPOSAL", height |-> height, round |-> round])

HasProposal(index) ==
    Cardinality(GetProposal(states[index].height, states[index].round)) > 0

HasPrepared(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index])) = 1

HasBlockAnnounce(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> states[index].height,
        round    |-> states[index].round])) >= 1

\* Helper function to check if the block is committed or not.
\* A block is considered committed iff supermajority of non-faulty replicas announce the same block.
IsCommitted ==
    LET subset == SubsetOfMsgs([
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> MaxHeight])
    IN /\ Cardinality(subset) >= TwoFPlusOne
       /\ \A m1, m2 \in subset : m1.round = m2.round

-----------------------------------------------------------------------------
(***************************************************************************)
(* Network functions                                                       *)
(***************************************************************************)

\* `SendMsg` simulates a replica sending a message by appending it to the `log`.
SendMsg(msg) ==
    IF msg.cp_round < MaxCPRound
    THEN log' = log \cup {msg}
    ELSE log' = log

\* SendProposal is used to broadcast the PROPOSAL into the network.
SendProposal(index) ==
    SendMsg([
        type     |-> "PROPOSAL",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

\* SendPrepareVote is used to broadcast PREPARE votes into the network.
SendPrepareVote(index) ==
    SendMsg([
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

\* SendPrecommitVote is used to broadcast PRECOMMIT votes into the network.
SendPrecommitVote(index) ==
    SendMsg([
        type     |-> "PRECOMMIT",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

\* SendCPPreVote is used to broadcast CP:PRE-VOTE votes into the network.
SendCPPreVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val])

\* SendCPMainVote is used to broadcast CP:MAIN-VOTE votes into the network.
SendCPMainVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val])

\* SendCPDeciedVote is used to broadcast CP:DECIDE votes into the network.
SendCPDeciedVote(index, cp_val) ==
    SendMsg([
        type     |-> "CP:DECIDE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        index    |-> -1,  \* reduce the model size
        cp_val   |-> cp_val])

\* AnnounceBlock is used to broadcast BLOCK-ANNOUNCE messages into the network.
AnnounceBlock(index)  ==
    SendMsg([
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0])

-----------------------------------------------------------------------------
(***************************************************************************)
(* States functions                                                        *)
(***************************************************************************)

\* NewHeight state
NewHeight(index) ==
    IF states[index].height >= MaxHeight
    THEN UNCHANGED <<states, log>>
    ELSE
        /\ ~IsFaulty(index)
        /\ states[index].name = "new-height"
        /\ states' = [states EXCEPT
            ![index].name = "propose",
            ![index].height = states[index].height + 1,
            ![index].round = 0]
        /\ log' = log


\* Propose state
Propose(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "propose"
    /\ IF IsProposer(index)
       THEN SendProposal(index)
       ELSE log' = log
    /\ states' = [states EXCEPT
        ![index].name = "prepare",
        ![index].cp_round = 0]


\* Prepare state
Prepare(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "prepare"
    /\ HasProposal(index)
    /\ SendPrepareVote(index)
    /\ states' = states

\* Precommit state
Precommit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "precommit"
    /\ IF HasPrecommitQuorum(index)
       THEN /\ states' = [states EXCEPT ![index].name = "commit"]
            /\ log' = log
       ELSE /\ HasProposal(index)
            /\ SendPrecommitVote(index)
            /\ states' = states

\* Commit state
Commit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "commit"
    /\ AnnounceBlock(index)
    /\ states' = [states EXCEPT
        ![index].name = "new-height"]

\* Timeout: A non-faulty Replica try to change the proposer if its timer expires.
Timeout(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "prepare"
    /\ IF states[index].round >= MaxRound
       THEN
            /\ HasPrepareQuorum(index)
            /\ states' = [states EXCEPT ![index].name = "cp:pre-vote"]
            /\ log' = log
       ELSE
            /\ states' = [states EXCEPT ![index].name = "cp:pre-vote"]
            /\ log' = log



CPPreVote(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:pre-vote"
        /\ IF states[index].cp_round = 0
           THEN
                IF HasPrepareQuorum(index)
                THEN /\ SendCPPreVote(index, 0)
                     /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
                ELSE IF HasPrepared(index)
                     THEN /\ CPHasPreVotesMinorityQuorum(index)
                          /\ SendCPPreVote(index, 1)
                          /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
                     ELSE /\ SendCPPreVote(index, 1)
                          /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]
            ELSE
                /\
                    \/
                        /\ CPHasAMainVotesOneInPrvRound(index)
                        /\ SendCPPreVote(index, 1)
                    \/
                        /\ CPHasAMainVotesZeroInPrvRound(index)
                        /\ SendCPPreVote(index, 0)
                    \/
                        /\ CPAllMainVotesAbstainInPrvRound(index)
                        /\ SendCPPreVote(index, 0) \* biased to zero
                /\ states' = [states EXCEPT ![index].name = "cp:main-vote"]


CPMainVote(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:main-vote"
    /\ CPHasPreVotesQuorum(index)
    /\
        \/
               \* all votes for 1
            /\ CPHasPreVotesQuorumForOne(index)
            /\ SendCPMainVote(index, 1)
            /\ states' = [states EXCEPT ![index].name = "cp:decide"]
        \/
               \* all votes for 0
            /\ CPHasPreVotesQuorumForZero(index)
            /\ SendCPMainVote(index, 0)
            /\ states' = [states EXCEPT ![index].name = "cp:decide"]
        \/
               \* Abstain vote
            /\ CPHasPreVotesForZeroAndOne(index)
            /\ SendCPMainVote(index, 2)
            /\ states' = [states EXCEPT ![index].name = "cp:decide"]

CPDecide(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:decide"
    /\ CPHasMainVotesQuorum(index)
    /\
        IF CPHasMainVotesQuorumForZero(index)
        THEN
            /\ SendCPDeciedVote(index, 0)
            /\ states' = states
        ELSE IF CPHasMainVotesQuorumForOne(index)
        THEN
            /\ SendCPDeciedVote(index, 1)
            /\ states' = states
        ELSE
            /\ states' = [states EXCEPT ![index].name = "cp:pre-vote",
                                        ![index].cp_round = states[index].cp_round + 1]
            /\ log' = log


CPStrongTerminate(index) ==
    /\ ~IsFaulty(index)
    /\
        \/ states[index].name = "cp:pre-vote"
        \/ states[index].name = "cp:main-vote"
        \/ states[index].name = "cp:decide"
    /\
        IF CPHasDecideVotesForOne(index)
        THEN /\ states' = [states EXCEPT ![index].name = "propose",
                                         ![index].round = states[index].round + 1]
             /\ log' = log
        ELSE IF CPHasDecideVotesForZero(index)
        THEN
             /\ states' = [states EXCEPT ![index].name = "precommit"]
             /\ log' = log
        ELSE IF /\ states[index].cp_round = MaxCPRound
                /\ CPOneFPlusOneMainVotesAbstainInPrvRound(index)
        THEN
            /\ states' = [states EXCEPT ![index].name = "precommit"]
            /\ log' = log
        ELSE
            /\ states' = states
            /\ log' = log

StrongCommit(index) ==
    /\ ~IsFaulty(index)
    /\
        \/ states[index].name = "prepare"
        \/ states[index].name = "precommit"
        \/ states[index].name = "cp:pre-vote"
        \/ states[index].name = "cp:main-vote"
        \/ states[index].name = "cp:decide"
    /\ HasPrepareAbsoluteQuorum(index)
    /\ states' = [states EXCEPT ![index].name = "commit"]
    /\ log' = log

-----------------------------------------------------------------------------

Init ==
    /\ log = {}
    /\ states = [index \in 0..NumNodes-1 |-> [
        name       |-> "new-height",
        height     |-> 0,
        round      |-> 0,
        cp_round   |-> 0]]

Next ==
    \E index \in 0..NumNodes-1:
        \/ NewHeight(index)
        \/ Propose(index)
        \/ Prepare(index)
        \/ Precommit(index)
        \/ Timeout(index)
        \/ Commit(index)
        \/ StrongCommit(index)
        \/ CPPreVote(index)
        \/ CPMainVote(index)
        \/ CPDecide(index)
        \/ CPStrongTerminate(index)

Spec ==
    Init /\ [][Next]_vars /\ WF_vars(Next)


(***************************************************************************)
(* Success: All non-faulty nodes eventually commit at MaxHeight.           *)
(***************************************************************************)
Success == <>(IsCommitted)

(***************************************************************************)
(* TypeOK is the type-correctness invariant.                               *)
(***************************************************************************)
TypeOK ==
    /\ \A index \in 0..NumNodes-1:
        /\ states[index].name \in {"new-height", "propose", "prepare",
            "precommit", "commit", "cp:pre-vote", "cp:main-vote", "cp:decide"}
        /\ states[index].height <= MaxHeight
        /\ states[index].round <= MaxRound
        /\ states[index].cp_round <= MaxCPRound
        /\ states[index].name = "new-height" /\ states[index].height > 0 =>
            /\ HasBlockAnnounce(index)
        /\ states[index].name = "precommit" =>
            /\ HasPrepareQuorum(index)
            /\ HasProposal(index)
        /\ states[index].name = "commit" =>
            /\ HasPrepareQuorum(index)
            /\ HasProposal(index)
        /\ \A round \in 0..states[index].round:
            \* Not more than one proposal per round
            /\ Cardinality(GetProposal(states[index].height, round)) <= 1

=============================================================================
