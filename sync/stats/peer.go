package stats

type Peer struct {
	receivedMsg int
}

func NewPeer() *Peer {
	return &Peer{}
}
