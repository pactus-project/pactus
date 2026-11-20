//go:build gtk

package main

import (
	"path/filepath"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoteBasicAuth(t *testing.T) {
	t.Run("no user, no authentication", func(t *testing.T) {
		auth, err := remoteBasicAuth("", "")
		require.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("a user needs a password file", func(t *testing.T) {
		_, err := remoteBasicAuth("user", "")
		require.Error(t, err)
	})

	t.Run("the password is read from the file, without the end of line", func(t *testing.T) {
		path := filepath.Join(util.TempDirPath(), "password")
		require.NoError(t, util.WriteFile(path, []byte("secret\n")))

		auth, err := remoteBasicAuth("user", path)
		require.NoError(t, err)
		assert.Equal(t, "user", auth.Username)
		assert.Equal(t, "secret", auth.Password)
	})
}
