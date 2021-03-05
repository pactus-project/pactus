package capnp

import "github.com/zarbchain/zarb-go/crypto"

func (zs zarbServer) GetAccount(args ZarbServer_getAccount) error {
	s, _ := args.Params.Address()
	addr, err := crypto.AddressFromString(string(s))
	if err != nil {
		zs.logger.Error("Error on retriving account", "err", err)
		return err
	}
	acc, err := zs.store.Account(addr)
	if err != nil {
		zs.logger.Error("Error on retriving account", "address", addr, "err", err)
		return err
	}

	d, _ := acc.Encode()
	res, _ := args.Results.NewResult()
	return res.SetData(d)
}
