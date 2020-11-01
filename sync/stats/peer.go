package stats

type Peer struct {
	ReceivedMsg int
	InvalidMsg  int
}

func NewPeer() *Peer {
	return &Peer{}
}
