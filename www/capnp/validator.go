package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
)

func (zs *zarbServer) GetValidator(b ZarbServer_getValidator) error {
	s, _ := b.Params.Address()
	addr, err := crypto.AddressFromString(string(s))
	if err != nil {
		return fmt.Errorf("Invalid address: %s", err)
	}
	val := zs.state.Validator(addr)
	if val == nil {
		return fmt.Errorf("Validator not found")
	}

	d, _ := val.Encode()
	res, _ := b.Results.NewResult()
	return res.SetData(d)
}
