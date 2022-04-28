package capnp

import (
	"fmt"

	"github.com/zarbchain/zarb-go/types/crypto"
)

func (zs *zarbServer) GetValidator(b ZarbServer_getValidator) error {
	s, _ := b.Params.Address()
	addr, err := crypto.AddressFromString(s)
	if err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}
	val := zs.state.ValidatorByAddress(addr)
	if val == nil {
		return fmt.Errorf("validator not found")
	}

	d, _ := val.Bytes()
	res, _ := b.Results.NewResult()
	return res.SetData(d)
}
