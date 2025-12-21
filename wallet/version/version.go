package version

const (
	Version1 = 1 // Initial version
	Version2 = 2 // Supporting Ed25519
	Version3 = 3 // Supporting AEC-256-CBC encryption method
	Version4 = 4 // Set Default Fee for the Wallet
	Version5 = 5 // Define entries and decouple it from vault

	VersionLatest = Version5
)
