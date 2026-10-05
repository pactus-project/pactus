package grpc

import (
	"errors"
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAnchor reports "not found" only for store.ErrNotFound.
// Any other store error is an internal error, not a missing anchor.
func TestGetAnchorStoreErrors(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	mockState := state.NewMockState(ts.MockController())
	server := newBlockchainServer(&Server{state: mockState})
	addr := ts.RandAccAddress()
	req := &pactus.GetAnchorRequest{Address: addr.String()}

	mockState.EXPECT().AccountByAddress(addr).Return(nil, errors.New("disk failure"))
	_, err := server.GetAnchor(t.Context(), req)
	require.Equal(t, codes.Internal, status.Code(err))

	mockState.EXPECT().AccountByAddress(addr).Return(nil, store.ErrNotFound)
	res, err := server.GetAnchor(t.Context(), req)
	require.NoError(t, err)
	require.False(t, res.Found)
	require.Equal(t, addr.String(), res.Address)
	require.Nil(t, res.Anchor)
}
