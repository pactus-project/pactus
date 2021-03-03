package capnp

import "github.com/zarbchain/zarb-go/crypto"

func (zs zarbServer) GetValidator(b ZarbServer_getValidator) error {
	s, _ := b.Params.Address()
	addr, err := crypto.AddressFromString(string(s))
	if err != nil {
		zs.logger.Error("Error on retriving validator", "err", err)
		return err
	}
	val, err := zs.store.Validator(addr)
	if err != nil {
		zs.logger.Error("Error on retriving validator", "address", addr, "err", err)
		return err
	}

	d, _ := val.Encode()
	res, _ := b.Results.NewResult()
	return res.SetData(d)
}
