package crypto

type AddressType byte

const (
	AddressTypeTreasury         AddressType = 0
	AddressTypeValidator        AddressType = 1
	AddressTypeBLSAccount       AddressType = 2
	AddressTypeEd25519Account   AddressType = 3
	AddressTypeSecp256k1Account AddressType = 4
)

func AddressTypeFromString(str string) (AddressType, error) {
	switch str {
	case AddressTypeValidator.String():
		return AddressTypeValidator, nil
	case AddressTypeBLSAccount.String():
		return AddressTypeBLSAccount, nil
	case AddressTypeEd25519Account.String():
		return AddressTypeEd25519Account, nil
	case AddressTypeSecp256k1Account.String():
		return AddressTypeSecp256k1Account, nil
	default:
		return AddressType(255), ErrInvalidAddressType
	}
}

func (t AddressType) String() string {
	switch t {
	case AddressTypeTreasury:
		return "treasury"
	case AddressTypeValidator:
		return "validator"
	case AddressTypeBLSAccount:
		return "bls_account"
	case AddressTypeEd25519Account:
		return "ed25519_account"
	case AddressTypeSecp256k1Account:
		return "secp256k1_account"

	default:
		return "unknown-address-type"
	}
}
