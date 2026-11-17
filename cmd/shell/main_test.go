package main

import (
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	rootCommandOnce sync.Once
	rootCommand     *cobra.Command
)

// testRootCommand creates the root command once: the flags of the services are registered
// in a global binder, and creating the command again would define them twice.
func testRootCommand() *cobra.Command {
	rootCommandOnce.Do(func() {
		rootCommand = createRootCommand()
	})

	return rootCommand
}

func TestCreateRootCommand(t *testing.T) {
	// This executes the ENTIRE main() logic including all changed lines
	rootCmd := testRootCommand()

	// Verify the root command properties
	require.Equal(t, "interactive", rootCmd.Use)
	require.Equal(t, "Pactus Shell", rootCmd.Short)
	require.Contains(t, rootCmd.Long, "pactus-shell is a command line tool")
	require.True(t, rootCmd.SilenceUsage)

	// Verify interactive command was added
	var interactiveCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "interactive" {
			interactiveCmd = cmd

			break
		}
	}

	require.NotNil(t, interactiveCmd, "interactive command should be added")
	require.Equal(t, "interactive", interactiveCmd.Use)
	require.Equal(t, "Start pactus-shell in interactive mode", interactiveCmd.Short)

	// Verify flags were set up
	serverAddrFlag := interactiveCmd.Flags().Lookup("server-addr")
	require.NotNil(t, serverAddrFlag)
	require.Equal(t, defaultServerAddr, serverAddrFlag.DefValue)

	usernameFlag := interactiveCmd.Flags().Lookup("auth-username")
	require.NotNil(t, usernameFlag)

	passwordFlag := interactiveCmd.Flags().Lookup("auth-password")
	require.NotNil(t, passwordFlag)

	// Verify PreRun and PersistentPreRun are set
	require.NotNil(t, interactiveCmd.PreRun)
	require.NotNil(t, interactiveCmd.PersistentPreRun)

	// Execute PreRun to hit
	interactiveCmd.PreRun(nil, nil)

	// Execute PersistentPreRun
	testCmd := &cobra.Command{}
	testCmd.SetContext(t.Context())
	interactiveCmd.PersistentPreRun(testCmd, nil)

	// Verify other commands were added
	var clearCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "clear" {
			clearCmd = cmd

			break
		}
	}
	require.NotNil(t, clearCmd)

	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	require.NoError(t, err)

	rootCmd.SetArgs([]string{"interactive", "--help"})
	err = rootCmd.Execute()
	require.NoError(t, err)

	// Test that main function would work
	require.NotNil(t, rootCmd)
	require.Equal(t, "interactive", rootCmd.Use)
}

// TestLivePrefix tests the livePrefix function.
func TestLivePrefix(t *testing.T) {
	_prefix = "test@localhost > "

	prefix, ok := livePrefix()
	require.True(t, ok)
	require.Equal(t, "test@localhost > ", prefix)
}

// TestClearScreenCommand tests the clearScreen function.
func TestClearScreenCommand(t *testing.T) {
	clearCmd := clearScreen()

	require.NotNil(t, clearCmd)
	require.Equal(t, "clear", clearCmd.Use)
	require.Equal(t, "clear screen", clearCmd.Short)
	require.NotNil(t, clearCmd.Run)

	// Test that the Run function works without error
	clearCmd.Run(nil, nil)
}

// TestSetAuthContext tests the setAuthContext function.
func TestSetAuthContext(t *testing.T) {
	rootCmd := &cobra.Command{}
	rootCmd.SetContext(t.Context())

	// Test with empty credentials (should not modify context)
	originalCtx := rootCmd.Context()
	setAuthContext(rootCmd, "", "")
	require.Equal(t, originalCtx, rootCmd.Context())

	// Test with valid credentials (should set auth context)
	setAuthContext(rootCmd, "testuser", "testpass")
	require.NotNil(t, rootCmd.Context())
	require.NotEqual(t, originalCtx, rootCmd.Context())
}

// TestConstants tests the defined constants.
func TestConstants(t *testing.T) {
	require.Equal(t, "localhost:50051", defaultServerAddr)
	require.Equal(t, "prettyjson", defaultResponseFormat)
}

// TestClsFunction tests the cls function doesn't panic.
func TestClsFunction(t *testing.T) {
	require.NotPanics(t, func() {
		cls()
	}, "calling cls() should not panic")
}

func TestSetFlagDefault(t *testing.T) {
	var username, password string

	rootCmd := &cobra.Command{Use: "root"}
	service := &cobra.Command{Use: "service"}
	service.PersistentFlags().StringVar(&username, "auth-username", "", "")
	service.PersistentFlags().StringVar(&password, "auth-password", "", "")
	method := &cobra.Command{Use: "method", Run: func(*cobra.Command, []string) {}}
	service.AddCommand(method)
	rootCmd.AddCommand(service)

	setFlagDefault(rootCmd, "auth-username", "user")
	setFlagDefault(rootCmd, "auth-password", "secret")

	// The interactive shell resets the flags of a command to their defaults before each run.
	// The credentials of the session must survive it.
	for range 3 {
		for _, name := range []string{"auth-username", "auth-password"} {
			flag := method.InheritedFlags().Lookup(name)
			require.NotNil(t, flag)
			require.NoError(t, flag.Value.Set(flag.DefValue))
		}

		assert.Equal(t, "user", username)
		assert.Equal(t, "secret", password)
	}

	// A flag that a command doesn't have is ignored.
	setFlagDefault(rootCmd, "unknown-flag", "value")
}

// The credentials given to the interactive shell become the defaults of the commands,
// which the shell resets before each run.
func TestInteractiveKeepsTheCredentials(t *testing.T) {
	rootCmd := testRootCommand()

	var interactiveCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "interactive" {
			interactiveCmd = cmd
		}
	}
	require.NotNil(t, interactiveCmd)

	require.NoError(t, interactiveCmd.Flags().Set("auth-username", "user"))
	require.NoError(t, interactiveCmd.Flags().Set("auth-password", "secret"))
	interactiveCmd.PreRun(interactiveCmd, nil)

	found := 0
	var check func(cmd *cobra.Command)
	check = func(cmd *cobra.Command) {
		if cmd != interactiveCmd {
			if flag := cmd.PersistentFlags().Lookup("auth-username"); flag != nil {
				assert.Equal(t, "user", flag.DefValue, cmd.CommandPath())
				assert.Equal(t, "secret", cmd.PersistentFlags().Lookup("auth-password").DefValue)
				found++
			}
		}
		for _, sub := range cmd.Commands() {
			check(sub)
		}
	}
	check(rootCmd)
	assert.Positive(t, found, "the services have the credential flags")
}
