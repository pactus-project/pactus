package certificate_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertificateStrictDecoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert := ts.GenerateTestCertificate(ts.RandHeight())
	valid, err := cert.MarshalCBOR()
	require.NoError(t, err)

	t.Run("a validator number beyond int32 is refused", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		require.NoError(t, encoding.WriteElements(buf, types.Height(10), types.Round(0)))
		require.NoError(t, encoding.WriteVarInt(buf, 1))
		require.NoError(t, encoding.WriteVarInt(buf, uint64(math.MaxInt32)+6)) // 2^31 + 5
		require.NoError(t, encoding.WriteVarInt(buf, 0))
		require.NoError(t, cert.Signature().Encode(buf))

		decoded := new(certificate.Certificate)
		err := decoded.Decode(buf)
		require.ErrorIs(t, err, certificate.ErrInvalidValidatorNumber)
	})

	t.Run("a failed decoding leaves the certificate untouched", func(t *testing.T) {
		original := ts.GenerateTestCertificate(ts.RandHeight())
		before := original.Hash()

		buf := bytes.NewBuffer(nil)
		require.NoError(t, encoding.WriteElements(buf, types.Height(10), types.Round(3)))
		// The rest of the data is missing.

		require.Error(t, original.Decode(buf))
		assert.Equal(t, before, original.Hash())
	})

	t.Run("bytes after the certificate are refused", func(t *testing.T) {
		data := make([]byte, 0, 1)
		require.NoError(t, cbor.Unmarshal(valid, &data))
		data = append(data, 0)
		bs, err := cbor.Marshal(data)
		require.NoError(t, err)

		require.ErrorIs(t, new(certificate.Certificate).UnmarshalCBOR(bs), certificate.ErrTrailingBytes)
		require.NoError(t, new(certificate.Certificate).UnmarshalCBOR(valid))
	})

	t.Run("a certificate without signature doesn't panic", func(t *testing.T) {
		unsigned := certificate.NewCertificate(1, 0)

		require.NotPanics(t, func() {
			_, err := unsigned.MarshalCBOR()
			require.ErrorIs(t, err, certificate.ErrNoSignature)
			assert.Equal(t, hash.UndefHash, unsigned.Hash())
			assert.Nil(t, unsigned.Clone().Signature())
		})
	})
}
