package autostart

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgsFromFlags(t *testing.T) {
	flagSet := flag.NewFlagSet("pactus-gui", flag.ContinueOnError)
	flagSet.String("working-dir", "/default", "")
	flagSet.String("password", "", "")
	flagSet.Bool("testnet", false, "")
	flagSet.String("grpc-addr", "", "")

	err := flagSet.Parse([]string{"-working-dir", "/home/user/pactus", "-password", "secret", "-testnet"})
	require.NoError(t, err)

	args := ArgsFromFlags(flagSet, "password")
	assert.Equal(t, []string{"-testnet=true", "-working-dir=/home/user/pactus"}, args)
}

func TestArgsFromFlagsNoneSet(t *testing.T) {
	flagSet := flag.NewFlagSet("pactus-gui", flag.ContinueOnError)
	flagSet.String("working-dir", "/default", "")

	require.NoError(t, flagSet.Parse(nil))
	assert.Empty(t, ArgsFromFlags(flagSet))
}
