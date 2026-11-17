package certificate_test

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/stretchr/testify/require"
)

// A certificate that decodes is encoded back to the same bytes: it has one encoding.
func FuzzCertificateRoundTrip(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzCertificateRoundTrip.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		cert := new(certificate.Certificate)
		if err := cert.Decode(bytes.NewReader(data)); err != nil {
			return
		}

		buf := new(bytes.Buffer)
		require.NoError(t, cert.Encode(buf))
		require.Equal(t, data[:buf.Len()], buf.Bytes())
	})
}
