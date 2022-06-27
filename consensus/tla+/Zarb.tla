-------------------------------- MODULE Zarb --------------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANT
    \* The total number of faulty nodes
    NumFaulty,
    MaxRound

NumValidators == (3 * NumFaulty) + 1
QuorumCnt == (2 * NumFaulty) + 1
OneThird == NumFaulty + 1

ASSUME
    /\ NumFaulty >= 1



VARIABLES
    log,
    states


vars == <<states, log>>


-----------------------------------------------------------------------------
(***************************************************************************)
(* Helper functions                                                        *)
(***************************************************************************)
\* Fetch a subset of messages in the network based on the params filter.
SubsetOfMsgs(params) ==
    {msg \in log: \A field \in DOMAIN params: msg[field] = params[field]}


\* In Zarb isProposer is chosen based on the time a validator was joined the network
\* here we assume the validators joined sequentially
IsProposer(index) ==
    (states[index].round + states[index].proposerIndex) % NumValidators = index

HasPrepareQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "PREPARE",
        height |-> states[index].height,
        round  |-> states[index].round])) >= QuorumCnt

HasPrecommitQuorum(index) ==
    Cardinality(SubsetOfMsgs([
        type   |-> "PRECOMMIT",
        height |-> states[index].height,
        round  |-> states[index].round])) >= QuorumCnt

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

SendMsg(msg) ==
    log' = log \cup {msg}

GetProposal(height, round) ==
    SubsetOfMsgs([type |-> "PROPOSAL", height |-> height, round |-> round])

HasProposal(height, round) ==
    Cardinality(GetProposal(height, round)) > 0

IsCommitted(height) ==
    Cardinality(SubsetOfMsgs([type |-> "BLOCK-ANNOUNCE", height |-> height])) > 0

\* SendProposal is used to broadcase proposal into the network
SendProposal(index) ==
    SendMsg([
        type    |-> "PROPOSAL",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index
        ])

\*
SendPrepareVote(index) ==
    SendMsg([
        type    |-> "PREPARE",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index
        ])

\*
SendPrecommitVote(index) ==
    SendMsg([
        type    |-> "PRECOMMIT",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index
        ])



\*
SendChangeProposerRequest(index) ==
    SendMsg([
        type    |-> "CHANGE-PROPOSER",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> index
        ])


\*
AnnounceBlock(index)  ==
    log' = {msg \in log: (msg.type = "BLOCK-ANNOUNCE") \/ msg.height > states[index].height } \cup {[
        type    |-> "BLOCK-ANNOUNCE",
        height  |-> states[index].height,
        round   |-> states[index].round,
        index   |-> -1
        ]}


NewHeight(index) ==
    /\ states[index].name = "new-height"
    /\ states' = [states EXCEPT
        ![index].name = "propose",
        ![index].height = states[index].height + 1,
        ![index].round = 0]
    /\ UNCHANGED <<log>>


Propose(index) ==
    /\ states[index].name = "propose"
    /\ IF IsProposer(index)
       THEN SendProposal(index)
       ELSE log' = log
    /\ states' = [states EXCEPT ![index].name = "prepare"]


Prepare(index) ==
    /\ states[index].name = "prepare"
    /\ IF /\ HasProposal(states[index].height, states[index].round)
          /\ ~HasOneThirdOfChangeProposer(index)
          \/ states[index].round < MaxRound
       THEN /\ SendPrepareVote(index)
            /\ IF HasPrepareQuorum(index)
               THEN states' = [states EXCEPT ![index].name = "precommit"]
               ELSE states' = states
       ELSE /\ SendChangeProposerRequest(index)
            /\ states' = [states EXCEPT ![index].name = "change-proposer"]


Precommit(index) ==
    /\ states[index].name = "precommit"
    /\ SendPrecommitVote(index)
    /\ IF HasPrecommitQuorum(index) /\ ~HasOneThirdOfChangeProposer(index)
       THEN states' = [states EXCEPT ![index].name = "commit"]
       ELSE states' = states


Commit(index) ==
    /\ states[index].name = "commit"
    /\ AnnounceBlock(index)
    /\ states' = [states EXCEPT
        ![index].name = "new-height",
        ![index].proposerIndex = (states[index].round + 1) % NumValidators]

ChangeProposer(index) ==
    /\ states[index].name = "change-proposer"
    /\ IF HasChangeProposerQuorum(index)
       THEN states' = [states EXCEPT
            ![index].name = "propose",
            ![index].round = states[index].round + 1]
       ELSE states' = states
    /\ UNCHANGED <<log>>


Sync(index) ==
    LET
        blocks == SubsetOfMsgs([type |-> "BLOCK-ANNOUNCE", height |-> states[index].height])
    IN
        /\ Cardinality(blocks) > 0
        /\ states' = [states EXCEPT
            ![index].name = "propose",
            ![index].height = states[index].height + 1,
            ![index].round = 0,
            ![index].proposerIndex = ((CHOOSE b \in blocks: TRUE).round + 1) % NumValidators]
        /\ log' = log


Init ==
    /\ log = {}
    /\ states = [index \in 0..NumValidators-1 |-> [
        name            |-> "new-height",
        height          |-> 0,
        round           |-> 0,
        proposerIndex   |-> 0
       ]]

Next ==
    \E index \in 0..NumValidators-1:
        \/ Sync(index)
        \/ NewHeight(index)
        \/ Propose(index)
        \/ Prepare(index)
        \/ Precommit(index)
        \/ Commit(index)
        \/ ChangeProposer(index)

\* The specification must start with the initial state and transition according
\* to Next.
Spec ==
    Init /\ [][Next]_vars


TypeOK ==
    /\ \A index \in 0..NumValidators-1:
        /\ states[index].name \in {"new-height", "propose", "prepare",
            "precommit", "commit", "change-proposer"}
        /\ states[index].name = "propose" =>
            \/ IsCommitted(states[index].height)
            \/ Cardinality(SubsetOfMsgs([index |-> index, height |-> states[index].height, round |-> states[index].round])) = 0
        /\ states[index].name = "precommit" =>
            \/ IsCommitted(states[index].height)
            \/ HasPrepareQuorum(index)
        /\ states[index].name = "commit" =>
            \/ IsCommitted(states[index].height)
            \/ HasPrecommitQuorum(index)
        /\ \A round \in 0..states[index].round:
            /\ Cardinality(GetProposal(states[index].height, round)) <= 1 \* not more than two proposals per round
            /\ round > 0 => HasChangeProposerQuorum(index)




=============================================================================
