package cmd

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A node that fails to start is stopped: what it opened, like the store, is released.
func TestStartNodeFailureStopsTheNode(t *testing.T) {
	dir := util.TempDirPath()
	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"
	wlt, _, err := CreateNode(t.Context(), 4, genesis.Localnet, dir, mnemonic, "password")
	require.NoError(t, err)
	wlt.Close()
	t.Chdir(dir)

	// The gRPC server can't listen on an address that is in use: the start fails.
	busy, err := new(net.ListenConfig).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = busy.Close() }()

	modifier := func(cfg *config.Config) *config.Config {
		cfg.Network.ListenAddrStrings = []string{"/ip4/127.0.0.1/tcp/0"}
		cfg.Network.EnableMdns = false
		cfg.GRPC.Enable = true
		cfg.GRPC.Listen = busy.Addr().String()

		return cfg
	}
	before := openFiles(t, dir)
	_, err = StartNode(t.Context(), dir, func() (string, bool) { return "password", true }, modifier)
	require.Error(t, err)

	// The files of the node, like the ones of its store, are closed.
	assert.Empty(t, openFiles(t, dir), "files opened before the start: %v", before)
}

// openFiles returns the files under dir that this process has open.
func openFiles(t *testing.T, dir string) []string {
	t.Helper()

	realDir, err := filepath.EvalSymlinks(dir)
	require.NoError(t, err)

	targets := []string{}
	switch runtime.GOOS {
	case "linux":
		const fdDir = "/proc/self/fd"
		entries, err := os.ReadDir(fdDir)
		require.NoError(t, err)
		for _, entry := range entries {
			target, err := os.Readlink(filepath.Join(fdDir, entry.Name()))
			if err == nil {
				targets = append(targets, target)
			}
		}

	case "darwin":
		out, err := exec.CommandContext(t.Context(), "lsof", "-Fn", "-p", strconv.Itoa(os.Getpid())).Output()
		if err != nil {
			t.Skip("lsof is not available")
		}
		for _, line := range strings.Split(string(out), "\n") {
			if name, ok := strings.CutPrefix(line, "n"); ok {
				targets = append(targets, name)
			}
		}

	default:
		t.Skip("the open files are not listed on " + runtime.GOOS)
	}

	files := []string{}
	for _, target := range targets {
		if target == realDir || target == dir || strings.HasSuffix(target, ".log") {
			// The working directory, and the file of the logger, which is global to the process.
			continue
		}
		if strings.HasPrefix(target, realDir) || strings.HasPrefix(target, dir) {
			files = append(files, target)
		}
	}

	return files
}
