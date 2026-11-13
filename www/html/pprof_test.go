package html_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/grpc/fake"
	"github.com/pactus-project/pactus/www/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// The profiles are behind the authentication, like the other pages.
func TestPprofNeedsTheCredential(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	gRPCServer := fake.NewFakeGRPCServer(t, ts, &grpc.Config{Enable: true, Listen: "127.0.0.1:0"})
	require.NoError(t, gRPCServer.Server.StartServer())

	free, err := new(net.ListenConfig).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := free.Addr().String()
	require.NoError(t, free.Close())

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	require.NoError(t, err)
	conf := &html.Config{Enable: true, Listen: address, EnablePprof: true}
	server := html.NewServer(t.Context(), conf, "user:"+string(hash))
	require.NoError(t, server.StartServer(gRPCServer.Server.Address()))
	t.Cleanup(server.StopServer)

	get := func(user, password string) int {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
			"http://"+address+"/debug/pprof/", http.NoBody)
		require.NoError(t, err)
		if user != "" {
			req.SetBasicAuth(user, password)
		}

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = res.Body.Close() }()

		return res.StatusCode
	}

	assert.Equal(t, http.StatusUnauthorized, get("", ""))
	assert.Equal(t, http.StatusUnauthorized, get("user", "wrong"))
	assert.Equal(t, http.StatusOK, get("user", "password"))
}
