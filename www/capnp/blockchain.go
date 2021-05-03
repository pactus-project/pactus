package capnp

func (zs *zarbServer) GetBlockchainInfo(args ZarbServer_getBlockchainInfo) error {
	height := zs.state.LastBlockHeight()
	res, _ := args.Results.NewResult()
	res.SetHeight(int64(height))
	return nil
}

func (zs *zarbServer) GetNetworkInfo(args ZarbServer_getNetworkInfo) error {
	res, _ := args.Results.NewResult()

	err := res.SetPeerID(zs.sync.SelfID().String())
	if err != nil {
		return err
	}
	pl, err := res.NewPeers(int32(len(zs.sync.Peers())))
	if err != nil {
		return err
	}
	for i, peer := range zs.sync.Peers() {
		p := pl.At(i)

		if err := p.SetMoniker(peer.Moniker()); err != nil {
			return err
		}
		if err := p.SetNodeVersion(peer.NodeVersion()); err != nil {
			return err
		}
		if err := p.SetPeerID(string(peer.PeerID())); err != nil {
			return err
		}
		if err := p.SetPublicKey(peer.PublicKey().String()); err != nil {
			return err
		}
		p.SetInitialBlockDownload(peer.InitialBlockDownload())
		p.SetHeight(int32(peer.Height()))
		p.SetReceivedMessages(int32(peer.ReceivedMessages()))
		p.SetInvalidMessages(int32(peer.InvalidMessages()))
		p.SetReceivedBytes(int32(peer.ReceivedBytes()))
	}

	return nil
}
