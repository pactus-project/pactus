package wallet

type vault struct {
	Seed        seed          `json:"seed"`
	Addresses   []address     `json:"addresses"`
	ImportedKey []importedKey `json:"imported_keys"`
}

type seed struct {
	Function string `json:"function"`
	Seed     string `json:"seed"`
}

type address struct {
	Label   string `json:"label"`
	From    string `json:"from"`
	KeyInfo string `json:"key_info"`
}

type importedKey struct {
	Prv string `json:"label"`
}
