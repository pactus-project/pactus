-------------------------------- MODULE Liveness --------------------------------
EXTENDS Zarb

CONSTANT MaxHeight


-----------------------------------------------------------------------------
(***************************************************************************)
(* Liveness: A node eventually moves to a new height.                      *)
(***************************************************************************)
Success == <>(
    \E index \in 0..NumValidators-1:
        /\ states[index].name = "new-height"
        /\ states[index].height = MaxHeight
)

LiveSpec ==
    /\ Spec /\ SF_log(Next) /\ SF_states(Next)


=============================================================================
