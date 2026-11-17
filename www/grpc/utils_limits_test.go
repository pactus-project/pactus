package grpc_test

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// An aggregation request holds at most 1024 public keys or signatures.
func TestAggregationLimits(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	td := setup(t, testConfig())
	client := td.utilClient(t)

	const maxItems = 1024
	pub, prv := ts.RandBLSKeyPair()
	sig := prv.Sign([]byte("message"))

	pubs := make([]string, maxItems+1)
	sigs := make([]string, maxItems+1)
	for i := range pubs {
		pubs[i] = pub.String()
		sigs[i] = sig.String()
	}

	_, err := client.PublicKeyAggregation(t.Context(),
		&pactus.PublicKeyAggregationRequest{PublicKeys: pubs[:maxItems]})
	require.NoError(t, err)

	_, err = client.PublicKeyAggregation(t.Context(),
		&pactus.PublicKeyAggregationRequest{PublicKeys: pubs})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	_, err = client.SignatureAggregation(t.Context(),
		&pactus.SignatureAggregationRequest{Signatures: sigs[:maxItems]})
	require.NoError(t, err)

	_, err = client.SignatureAggregation(t.Context(),
		&pactus.SignatureAggregationRequest{Signatures: sigs})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
