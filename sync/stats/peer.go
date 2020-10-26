package stats

type Peer struct {
	ReceivedMsg int
}

func NewPeer() *Peer {
	return &Peer{}
}
