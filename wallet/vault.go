package wallet

type vault struct {
	Addresses []address `json:"addresses"`
	Seed      seed      `json:"seed"`
	Keystore  keystore  `json:"keystore"`
}

type seed struct {
	Function string    `json:"function"`
	Seed     encrypted `json:"seed"`
}

type address struct {
	Label   string `json:"label"`
	From    string `json:"from"`
	KeyInfo string `json:"key_info"`
}

type keystore struct {
	Prv []string `json:"prv"`
}
