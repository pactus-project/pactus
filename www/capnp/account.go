package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
)

func (zs zarbServer) GetAccount(args ZarbServer_getAccount) error {
	s, _ := args.Params.Address()
	addr, err := crypto.AddressFromString(s)
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
