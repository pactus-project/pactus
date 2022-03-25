package capnp

func (zs *zarbServer) GetBlockchainInfo(args ZarbServer_getBlockchainInfo) error {
	res, _ := args.Results.NewResult()
	res.SetLastBlockHeight(zs.state.LastBlockHeight())
	return res.SetLastBlockHash(zs.state.LastBlockHash().RawBytes())
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

		if err := p.SetMoniker(peer.Moniker); err != nil {
			return err
		}
		if err := p.SetAgent(peer.Agent); err != nil {
			return err
		}
		if err := p.SetPeerID(string(peer.PeerID)); err != nil {
			return err
		}
		if err := p.SetPublicKey(peer.PublicKey.String()); err != nil {
			return err
		}
		p.SetStatus(int32(peer.Status))
		p.SetFlags(int32(peer.Flags))
		p.SetHeight(int32(peer.Height))
		p.SetReceivedMessages(int32(peer.ReceivedBundles))
		p.SetInvalidMessages(int32(peer.InvalidBundles))
		p.SetReceivedBytes(int32(peer.ReceivedBytes))
	}

	return nil
}
