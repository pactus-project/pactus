package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestDownloadRequestType(t *testing.T) {
	p := &DownloadRequestPayload{}
	assert.Equal(t, p.Type(), PayloadTypeDownloadRequest)
}

func TestDownloadRequestPayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		p := NewDownloadRequestPayload(1, "", 100, 200)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid from", func(t *testing.T) {
		p := NewDownloadRequestPayload(1, util.RandomPeerID(), -100, 200)
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid range", func(t *testing.T) {
		p := NewDownloadRequestPayload(1, util.RandomPeerID(), 200, 100)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewDownloadRequestPayload(1, util.RandomPeerID(), 100, 200)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
