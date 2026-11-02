package crypto

import "sync"

type SignatureType byte

var (
	SignatureTypeBLS       SignatureType = 1
	SignatureTypeEd25519   SignatureType = 3
	SignatureTypeSecp256k1 SignatureType = 4
)

// hrpMu guards the HRP variables below: they are read from arbitrary
// goroutines (via AddressHRP(), etc.) and rewritten wholesale by
// ToTestnetHRP/ToMainnetHRP and SetAddressHRP.
var hrpMu sync.RWMutex

var (
	addressHRP     = "pc"
	publicKeyHRP   = "public"
	privateKeyHRP  = "secret"
	xPublicKeyHRP  = "xpublic"
	xPrivateKeyHRP = "xsecret"
)

// AddressHRP returns the Human Readable Part (HRP) for address.
func AddressHRP() string {
	hrpMu.RLock()
	defer hrpMu.RUnlock()

	return addressHRP
}

// PublicKeyHRP returns the Human Readable Part (HRP) for public key.
func PublicKeyHRP() string {
	hrpMu.RLock()
	defer hrpMu.RUnlock()

	return publicKeyHRP
}

// PrivateKeyHRP returns the Human Readable Part (HRP) for private key.
func PrivateKeyHRP() string {
	hrpMu.RLock()
	defer hrpMu.RUnlock()

	return privateKeyHRP
}

// XPublicKeyHRP returns the Human Readable Part (HRP) for extended public key.
func XPublicKeyHRP() string {
	hrpMu.RLock()
	defer hrpMu.RUnlock()

	return xPublicKeyHRP
}

// XPrivateKeyHRP returns the Human Readable Part (HRP) for extended private key.
func XPrivateKeyHRP() string {
	hrpMu.RLock()
	defer hrpMu.RUnlock()

	return xPrivateKeyHRP
}

// SetAddressHRP overrides the Human Readable Part (HRP) for address, leaving
// the other HRPs untouched. It is intended for tests that need a network's
// address prefix without switching the whole process to that network.
func SetAddressHRP(hrp string) {
	hrpMu.Lock()
	defer hrpMu.Unlock()

	addressHRP = hrp
}

// ToTestnetHRP makes HRPs testnet specified.
func ToTestnetHRP() {
	hrpMu.Lock()
	defer hrpMu.Unlock()

	addressHRP = "tpc"
	publicKeyHRP = "tpublic"
	privateKeyHRP = "tsecret"
	xPublicKeyHRP = "txpublic"
	xPrivateKeyHRP = "txsecret"
}

// ToMainnetHRP makes HRPs mainnet specified.
func ToMainnetHRP() {
	hrpMu.Lock()
	defer hrpMu.Unlock()

	addressHRP = "pc"
	publicKeyHRP = "public"
	privateKeyHRP = "secret"
	xPublicKeyHRP = "xpublic"
	xPrivateKeyHRP = "xsecret"
}
