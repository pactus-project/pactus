package crypto

var (
	// AddressHRP is the Human Readable Part (HRP) for address.
	AddressHRP = "pc"
	// PublicKeyHRP is the Human Readable Part (HRP) for public key.
	PublicKeyHRP = "public"
	// PrivateKeyHRP is the Human Readable Part (HRP) for private key.
	PrivateKeyHRP = "secret"
	// XPublicKeyHRP is the Human Readable Part (HRP) for extended public key.
	XPublicKeyHRP = "xpublic"
	// XPrivateKeyHRP is the Human Readable Part (HRP) for extended private key.
	XPrivateKeyHRP = "xsecret"
)

// ToTestnetHRP makes HRPs testnet specified.
func ToTestnetHRP() {
	AddressHRP = "tpc"
	PublicKeyHRP = "tpublic"
	PrivateKeyHRP = "tsecret"
	XPublicKeyHRP = "txpublic"
	XPrivateKeyHRP = "txsecret"
}
