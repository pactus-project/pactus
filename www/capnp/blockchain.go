package capnp

import "github.com/fxamacker/cbor/v2"

func (f factory) GetBlockchainInfo(args ZarbServer_getBlockchainInfo) error {
	height := f.state.LastBlockHeight()
	res, _ := args.Results.NewResult()
	res.SetHeight(int64(height))
	return nil
}

func (f factory) GetNetworkInfo(args ZarbServer_getNetworkInfo) error {
	res, _ := args.Results.NewResult()

	err := res.SetPeerID(f.sync.PeerID().String())
	if err != nil {
		return err
	}
	pl, err := res.NewPeers(int32(len(f.sync.Peers())))
	if err != nil {
		return err
	}
	for i, peer := range f.sync.Peers() {
		p := pl.At(i)

		if err := p.SetMoniker(peer.Moniker()); err != nil {
			return err
		}
		bs, _ := cbor.Marshal(peer.NodeVersion())
		if err := p.SetNodeVersion(bs); err != nil {
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
		p.SetReceivedMsg(int32(peer.ReceivedMsg()))
		p.SetInvalidMsg(int32(peer.InvalidMsg()))
		p.SetReceivedBytes(int32(peer.ReceivedBytes()))
	}

	return nil
}
