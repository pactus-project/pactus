package capnp

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

func (zs pactusServer) GetAccount(args PactusServer_getAccount) error {
	capAddr, _ := args.Params.Address()
	addr, err := crypto.AddressFromString(capAddr)
	if err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}
	acc := zs.state.AccountByAddress(addr)
	if acc == nil {
		return fmt.Errorf("account not found")
	}

	d, _ := acc.Bytes()
	res, _ := args.Results.NewResult()
	return res.SetData(d)
}
