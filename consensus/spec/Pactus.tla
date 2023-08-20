-------------------------------- MODULE Pactus --------------------------------
(***************************************************************************)
(* The specification of the Pactus consensus algorithm:                    *)
(* `^\url{https://pactus.org/learn/consensus/protocol/}^'                  *)
(***************************************************************************)
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANT
    \* The maximum number of height.
    \* this is to restrict the allowed behaviours that TLC scans through.
    MaxHeight,
    \* The maximum number of round per height.
    \* this is to restrict the allowed behaviours that TLC scans through.
    MaxRound,
    \* The maximum number of cp-round per height.
    \* this is to restrict the allowed behaviours that TLC scans through.
    MaxCPRound,
    \* The total number of faulty nodes
    NumFaulty,
    \* The index of faulty nodes
    FaultyNodes

VARIABLES
    \* `log` is a set of received messages in the system.
    log,
    \* `states` represents the state of each replica in the consensus protocol.
    states

\* Total number of replicas, which is `3f+1', where `f' is the number of faulty nodes.
Replicas == (3 * NumFaulty) + 1
\* Quorum is 2/3+ of total replicas that is `2f+1'
Quorum == (2 * NumFaulty) + 1
\* OneThird is 1/3+ of total  replicas that is `f+1'
OneThird == NumFaulty + 1

\* A tuple with all variables in the spec (for ease of use in temporal conditions)
vars == <<states, log>>

ASSUME
    /\ NumFaulty >= 1
    /\ FaultyNodes \subseteq 0..Replicas-1

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
    states[index].round % Replicas = index

\* Helper function to check if a node is faulty or not.
IsFaulty(index) == index \in FaultyNodes

\* HasPrepareQuorum checks if there is a quorum of
\* the PREPARE votes in this round.
HasPrepareQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> 0])) >= Quorum

\* HasPrecommitQuorum checks if there is a quorum of
\* the PRECOMMIT votes in this round.
HasPrecommitQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "PRECOMMIT",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> 0])) >= Quorum

CPHasPreVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= Quorum

CPHasPreVotesQuorumForOne(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= Quorum

CPHasPreVotesQuorumForZero(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= Quorum

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

CPHasOneMainVotesZeroInPrvRound(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round - 1,
        cp_val   |-> 0])) > 0

CPHasOneMainVotesOneInPrvRound(index) ==
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
        cp_val   |-> 2])) >= Quorum


CPHasMainVotesQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round])) >= Quorum

CPHasMainVotesQuorumForOne(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 1])) >= Quorum

CPHasMainVotesQuorumForZero(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> states[index].cp_round,
        cp_val   |-> 0])) >= Quorum

GetProposal(height, round) ==
    SubsetOfMsgs([type |-> "PROPOSAL", height |-> height, round |-> round])

HasProposal(index) ==
    Cardinality(GetProposal(states[index].height, states[index].round)) > 0

HasBlockAnnounce(index) ==
    Cardinality(SubsetOfMsgs([
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        cp_round |-> 0,
        cp_val   |-> 0])) >= 1

\* Helper function to check if the block is committed or not.
\* A block is considered committed iff supermajority of non-faulty replicas announce the same block.
IsCommitted(height) ==
    LET subset == SubsetOfMsgs([
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> height,
        cp_round |-> 0,
        cp_val   |-> 0])
    IN /\ Cardinality(subset) >= Quorum
       /\ \A m1, m2 \in subset : m1.round = m2.round

-----------------------------------------------------------------------------
(***************************************************************************)
(* Network functions                                                       *)
(***************************************************************************)

\* `SendMsg` simulates a replica sending a message by appending it to the `log`.
SendMsg(msg) ==
    log' = log \cup msg

\* SendProposal is used to broadcast the PROPOSAL into the network.
SendProposal(index) ==
    SendMsg({[
        type     |-> "PROPOSAL",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0]})

\* SendPrepareVote is used to broadcast PREPARE votes into the network.
SendPrepareVote(index) ==
    SendMsg({[
        type     |-> "PREPARE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0]})

\* SendPrecommitVote is used to broadcast PRECOMMIT votes into the network.
SendPrecommitVote(index) ==
    SendMsg({[
        type     |-> "PRECOMMIT",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0]})

\* SendCPPreVote is used to broadcast CP:PRE-VOTE votes into the network.
SendCPPreVote(index, cp_val) ==
    SendMsg({[
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val]})

\* SendCPMainVote is used to broadcast CP:MAIN-VOTE votes into the network.
SendCPMainVote(index, cp_val) ==
    SendMsg({[
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round,
        cp_val   |-> cp_val]})

SendCPVotesForNextRound(index, cp_val) ==
    SendMsg({
    [
        type     |-> "CP:PRE-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round + 1,
        cp_val   |-> cp_val],
    [
        type     |-> "CP:MAIN-VOTE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> states[index].cp_round + 1,
        cp_val   |-> cp_val]})

\* AnnounceBlock is used to broadcast BLOCK-ANNOUNCE messages into the network.
AnnounceBlock(index)  ==
    SendMsg({[
        type     |-> "BLOCK-ANNOUNCE",
        height   |-> states[index].height,
        round    |-> states[index].round,
        index    |-> index,
        cp_round |-> 0,
        cp_val   |-> 0]})

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
        /\ states[index].height < MaxHeight
        /\ states' = [states EXCEPT
            ![index].name = "propose",
            ![index].height = states[index].height + 1,
            ![index].round = 0]
        /\ UNCHANGED <<log>>


\* Propose state
Propose(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "propose"
    /\ IF IsProposer(index)
       THEN SendProposal(index)
       ELSE UNCHANGED <<log>>
    /\ states' = [states EXCEPT
        ![index].name = "prepare",
        ![index].timeout = FALSE,
        ![index].cp_round = 0]


\* Prepare state
Prepare(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "prepare"
    /\ IF HasPrepareQuorum(index)
       THEN /\ states' = [states EXCEPT ![index].name = "precommit"]
            /\ UNCHANGED <<log>>
       ELSE /\ HasProposal(index)
            /\ SendPrepareVote(index)
            /\ UNCHANGED <<states>>

\* Precommit state
Precommit(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "precommit"
    /\ IF HasPrecommitQuorum(index)
       THEN /\ states' = [states EXCEPT ![index].name = "commit"]
            /\ UNCHANGED <<log>>
       ELSE /\ HasProposal(index)
            /\ SendPrecommitVote(index)
            /\ UNCHANGED <<states>>

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
    /\ states[index].round < MaxRound
    /\ states[index].timeout = FALSE
    /\
        \/
            /\ states[index].name = "prepare"
            /\ SendCPPreVote(index, 1)

        \/
            /\ states[index].name = "precommit"
            /\ SendCPPreVote(index, 0)
    /\ states' = [states EXCEPT
            ![index].name = "cp:main-vote",
            ![index].timeout = TRUE]



CPPreVote(index) ==
    /\ ~IsFaulty(index)
    /\ states[index].name = "cp:pre-vote"
    /\
        \/
            /\ CPHasOneMainVotesOneInPrvRound(index)
            /\ SendCPPreVote(index, 1)
        \/
            /\ CPHasOneMainVotesZeroInPrvRound(index)
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
    /\
        \/
            /\ states[index].cp_decided = 1
            /\ states' = [states EXCEPT ![index].name = "propose",
                                        ![index].round = states[index].round + 1]
        \/
            /\ states[index].cp_decided = 0
            /\ states' = [states EXCEPT ![index].name = "prepare"]
        \/
            /\ states[index].cp_decided = -1
            /\ CPHasMainVotesQuorum(index)
            /\
                IF  /\ CPHasMainVotesQuorumForOne(index)
                    /\ states[index].cp_round /= MaxCPRound - 1
                THEN states' = [states EXCEPT ![index].name = "cp:pre-vote",
                                              ![index].cp_decided = 1,
                                              ![index].cp_round = states[index].cp_round + 1]
                ELSE IF \/ CPHasMainVotesQuorumForZero(index)
                        \/ states[index].cp_round = MaxCPRound - 1
                    THEN states' = [states EXCEPT ![index].name = "cp:pre-vote",
                                                  ![index].cp_decided = 0,
                                                  ![index].cp_round = states[index].cp_round + 1]
                    ELSE states' = [states EXCEPT ![index].name = "cp:pre-vote",
                                                  ![index].cp_round = states[index].cp_round + 1]

    /\ log' = log


Sync(index) ==
    /\ ~IsFaulty(index)
    /\
        \/ states[index].name = "cp:pre-vote"
        \/ states[index].name = "cp:main-vote"
        \/ states[index].name = "cp:decide"
    /\ HasBlockAnnounce(index)
    /\ states' = [states EXCEPT ![index].name = "prepare"]
    /\ log' = log

-----------------------------------------------------------------------------

Init ==
    /\ log = {}
    /\ states = [index \in 0..Replicas-1 |-> [
        name       |-> "new-height",
        height     |-> 0,
        round      |-> 0,
        timeout    |-> FALSE,
        cp_round   |-> 0,
        cp_decided |-> -1]]

Next ==
    \E index \in 0..Replicas-1:
        \/ NewHeight(index)
        \/ Propose(index)
        \/ Prepare(index)
        \/ Precommit(index)
        \/ Timeout(index)
        \/ Commit(index)
        \/ Sync(index)
        \/ CPPreVote(index)
        \/ CPMainVote(index)
        \/ CPDecide(index)

Spec ==
    Init /\ [][Next]_vars /\ WF_vars(Next)


(***************************************************************************)
(* Success: All non-faulty nodes eventually commit at MaxHeight.           *)
(***************************************************************************)
Success == <>(IsCommitted(MaxHeight))

(***************************************************************************)
(* TypeOK is the type-correctness invariant.                               *)
(***************************************************************************)
TypeOK ==
    /\ \A index \in 0..Replicas-1:
        /\ states[index].name \in {"new-height", "propose", "prepare",
            "precommit", "commit", "cp:pre-vote", "cp:main-vote", "cp:decide"}
        /\ states[index].height <= MaxHeight
        /\ states[index].round <= MaxRound
        /\ states[index].cp_round <= MaxCPRound + 1
        /\ states[index].name = "new-height" /\ states[index].height > 1 =>
            /\ IsCommitted(states[index].height - 1)
        /\ states[index].name = "precommit" =>
            /\ HasPrepareQuorum(index)
            /\ HasProposal(index)
        /\ states[index].name = "commit" =>
            /\ HasPrepareQuorum(index)
            /\ HasPrecommitQuorum(index)
            /\ HasProposal(index)
        /\ \A round \in 0..states[index].round:
            \* Not more than one proposal per round
            /\ Cardinality(GetProposal(states[index].height, round)) <= 1

=============================================================================
