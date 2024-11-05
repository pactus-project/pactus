package crypto

type AddressType byte

const (
	AddressTypeTreasury       AddressType = 0
	AddressTypeValidator      AddressType = 1
	AddressTypeBLSAccount     AddressType = 2
	AddressTypeEd25519Account AddressType = 3
)

func (t AddressType) String() string {
	switch t {
	case AddressTypeTreasury:
		return "treasury"
	case AddressTypeBLSAccount:
		return "bls_account"
	case AddressTypeEd25519Account:
		return "ed25519_account"
	case AddressTypeValidator:
		return "validator"

	default:
		return "unknown-address-type"
	}
}
