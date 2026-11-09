package protocol

import (
	"fmt"
	"strconv"
)

type Version uint8

const (
	ProtocolVersionUnknown Version = 0
	ProtocolVersion1       Version = 1 // Initial version
	ProtocolVersion2       Version = 2 // Split Reward Fork (PIP-43)
	ProtocolVersion3       Version = 3 // Validator Delegation (PIP-49)
	ProtocolVersion4       Version = 4 // Block Reward Halving (PIP-55)
	ProtocolVersion5       Version = 5 // Native State Anchoring (PIP-50)

	// ProtocolVersion6 enforces the hardened rules (see HasHardenedRules).
	// It is prepared but it is NOT active: no block can have this version until
	// ProtocolVersionLatest is raised to it in a release.
	ProtocolVersion6 Version = 6 // Security hardening

	ProtocolVersionLatest = ProtocolVersion5
)

// HasHardenedRules returns true if the protocol version enforces the hardened rules:
//   - A delegated validator can only withdraw its whole stake, with a bounded fee.
//   - The justifications of the change-proposer votes are bound to the height and round of the vote.
//   - The pre-vote `yes` of the first change-proposer round has no block hash.
//   - The subsidy transaction pays exactly the expected recipients.
//   - Transactions that are not free pay a fee, and transfers have a positive amount.
func (v Version) HasHardenedRules() bool {
	return v >= ProtocolVersion6
}

func ParseVersion(s string) (Version, error) {
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}

	return Version(v), nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d", v)
}
