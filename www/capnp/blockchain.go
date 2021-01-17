package capnp

func (f factory) GetBlockchainInfo(args ZarbServer_getBlockchainInfo) error {
	height := f.state.LastBlockHeight()
	res, _ := args.Results.NewResult()
	res.SetHeight(int64(height))
	return nil
}
