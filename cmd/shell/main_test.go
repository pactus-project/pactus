package main

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)
func TestCreateRootCommand(t *testing.T) {
	// This executes the ENTIRE main() logic including all changed lines
	rootCmd := createRootCommand()

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
	testCmd.SetContext(context.Background())
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

// TestLivePrefix tests the livePrefix function
func TestLivePrefix(t *testing.T) {
	_prefix = "test@localhost > "

	prefix, ok := livePrefix()
	require.True(t, ok)
	require.Equal(t, "test@localhost > ", prefix)
}

// TestClearScreenCommand tests the clearScreen function
func TestClearScreenCommand(t *testing.T) {
	clearCmd := clearScreen()

	require.NotNil(t, clearCmd)
	require.Equal(t, "clear", clearCmd.Use)
	require.Equal(t, "clear screen", clearCmd.Short)
	require.NotNil(t, clearCmd.Run)

	// Test that the Run function works without error
	clearCmd.Run(nil, nil)
}

// TestSetAuthContext tests the setAuthContext function
func TestSetAuthContext(t *testing.T) {
	rootCmd := &cobra.Command{}
	rootCmd.SetContext(context.Background())

	// Test with empty credentials (should not modify context)
	originalCtx := rootCmd.Context()
	setAuthContext(rootCmd, "", "")
	require.Equal(t, originalCtx, rootCmd.Context())

	// Test with valid credentials (should set auth context)
	setAuthContext(rootCmd, "testuser", "testpass")
	require.NotNil(t, rootCmd.Context())
	require.NotEqual(t, originalCtx, rootCmd.Context())
}

// TestConstants tests the defined constants
func TestConstants(t *testing.T) {
	require.Equal(t, "localhost:50051", defaultServerAddr)
	require.Equal(t, "prettyjson", defaultResponseFormat)
}

// TestClsFunction tests the cls function doesn't panic
func TestClsFunction(t *testing.T) {
	// Test that cls() doesn't panic
	cls()
}