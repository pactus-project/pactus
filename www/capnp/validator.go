package capnp

import "github.com/zarbchain/zarb-go/crypto"

func (f factory) GetValidator(b ZarbServer_getValidator) error {
	s, _ := b.Params.Address()
	addr, err := crypto.AddressFromString(string(s))
	if err != nil {
		f.logger.Error("Error on retriving validator", "err", err)
		return err
	}
	val, err := f.store.Validator(addr)
	if err != nil {
		f.logger.Error("Error on retriving validator", "address", addr, "err", err)
		return err
	}

	d, _ := val.Encode()
	res, _ := b.Results.NewResult()
	return res.SetData(d)
}
