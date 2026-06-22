package main

import (
	"testing"

	pactuscmd "github.com/pactus-project/pactus/cmd"
	"github.com/spf13/cobra"
)

func TestResolveSnapshotURL(t *testing.T) {
	tests := []struct {
		name        string
		setSnapshot bool
		snapshotURL string
		setServer   bool
		serverAddr  string
		expectedURL string
	}{
		{
			name:        "uses default when no flags are set",
			expectedURL: pactuscmd.DefaultSnapshotURL,
		},
		{
			name:        "uses server-addr when only server-addr is set",
			setServer:   true,
			serverAddr:  "https://server.example.org",
			expectedURL: "https://server.example.org",
		},
		{
			name:        "uses snapshot-url when only snapshot-url is set",
			setSnapshot: true,
			snapshotURL: "https://snapshot.example.org",
			expectedURL: "https://snapshot.example.org",
		},
		{
			name:        "snapshot-url wins when both are set",
			setSnapshot: true,
			snapshotURL: "https://snapshot.example.org",
			setServer:   true,
			serverAddr:  "https://server.example.org",
			expectedURL: "https://snapshot.example.org",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newImportTestCommand()

			if tc.setSnapshot {
				if err := c.Flags().Set("snapshot-url", tc.snapshotURL); err != nil {
					t.Fatalf("set snapshot-url: %v", err)
				}
			}

			if tc.setServer {
				if err := c.Flags().Set("server-addr", tc.serverAddr); err != nil {
					t.Fatalf("set server-addr: %v", err)
				}
			}

			snapshotURL, err := c.Flags().GetString("snapshot-url")
			if err != nil {
				t.Fatalf("get snapshot-url: %v", err)
			}

			serverAddr, err := c.Flags().GetString("server-addr")
			if err != nil {
				t.Fatalf("get server-addr: %v", err)
			}

			got := resolveSnapshotURL(c, snapshotURL, serverAddr)
			if got != tc.expectedURL {
				t.Fatalf("expected %q, got %q", tc.expectedURL, got)
			}
		})
	}
}

func newImportTestCommand() *cobra.Command {
	c := &cobra.Command{Use: "import"}
	c.Flags().String("snapshot-url", pactuscmd.DefaultSnapshotURL, "")
	c.Flags().String("server-addr", pactuscmd.DefaultSnapshotURL, "")

	return c
}
