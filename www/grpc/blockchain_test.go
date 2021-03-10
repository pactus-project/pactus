package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetNetworkInfo(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return Nodes PeerID as peerid", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
		// assert.Error(t, err)
		assert.Nil(t, err)
		assert.Equal(t, tMockSync.ID.String(), res.PeerId)
		assert.Equal(t, 2, len(res.Peers))
	})

	newPeer := tMockSync.AddPeer("newPeer", 12)

	t.Run("Should return newly added Peer", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
		// assert.Error(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(res.Peers))
		for _, p := range res.Peers {
			if p.Moniker == "newPeer" {
				assert.Equal(t, newPeer.PeerID().String(), p.PeerId)
				assert.Equal(t, int32(12), p.Height)
				return
			}
		}
		t.Error("new Peer Not Found")
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}
