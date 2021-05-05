package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
)

func (zs zarbServer) GetAccount(args ZarbServer_getAccount) error {
	s, _ := args.Params.Address()
	addr, err := crypto.AddressFromString(string(s))
	if err != nil {
		return fmt.Errorf("invalid address: %s", err)
	}
	acc := zs.state.Account(addr)
	if acc == nil {
		return fmt.Errorf("account not found")
	}

	d, _ := acc.Encode()
	res, _ := args.Results.NewResult()
	return res.SetData(d)
}
