package shell

import (
	"errors"
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestBuildCompletionArgs_Empty(t *testing.T) {
	args, err := buildCompletionArgs("")
	require.NoError(t, err)

	expected := []string{"__complete", ""}
	require.Equal(t, expected, args)
}

func TestBuildCompletionArgs_CurrentArg(t *testing.T) {
	args, err := buildCompletionArgs("a b")
	require.NoError(t, err)

	expected := []string{"__complete", "a", "b"}
	require.Equal(t, expected, args)
}

func TestBuildCompletionArgs_MultiwordString(t *testing.T) {
	args, err := buildCompletionArgs(`a "b c"`)
	require.NoError(t, err)

	expected := []string{"__complete", "a", "b c"}
	require.Equal(t, expected, args)
}

func TestBuildCompletionArgs_NextArg(t *testing.T) {
	args, err := buildCompletionArgs("a b ")
	require.NoError(t, err)

	expected := []string{"__complete", "a", "b", ""}
	require.Equal(t, expected, args)
}

func TestReadCommandOutput_Stdout(t *testing.T) {
	cmd := &cobra.Command{
		Use: "command",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print("out")
		},
	}

	out, err := readCommandOutput(cmd, []string{})
	require.NoError(t, err)
	require.Equal(t, "out", out)
}

func TestReadCommandOutput_Stderr(t *testing.T) {
	cmd := &cobra.Command{
		Use: "command",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.PrintErr("out")
		},
	}

	out, err := readCommandOutput(cmd, []string{})
	require.NoError(t, err)
	require.Empty(t, out)
}

func TestReadCommandOutput_Err(t *testing.T) {
	cmd := &cobra.Command{
		Use: "command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("err")
		},
	}

	_, err := readCommandOutput(cmd, []string{})
	require.Error(t, err)
}

func TestParseSuggestions_WithDescription(t *testing.T) {
	out := `command-with-description	description
:4
Completion ended with directive: ShellCompDirectiveNoFileComp`
	expected := []prompt.Suggest{{Text: "command-with-description", Description: "description"}}
	require.Equal(t, expected, parseSuggestions(out))
}

func TestParseSuggestions_WithoutDescription(t *testing.T) {
	out := `command-without-description
:4
Completion ended with directive: ShellCompDirectiveNoFileComp`
	expected := []prompt.Suggest{{Text: "command-without-description"}}
	require.Equal(t, expected, parseSuggestions(out))
}

func TestParseSuggestions_HideShorthandFlags(t *testing.T) {
	out := `--flag	A flag.
-f	A flag.
:4
Completion ended with directive: ShellCompDirectiveNoFileComp`
	expected := []prompt.Suggest{{Text: "--flag", Description: "A flag."}}
	require.Equal(t, expected, parseSuggestions(out))
}

func TestParseSuggestions_Sort(t *testing.T) {
	out := `b
a
:4
Completion ended with directive: ShellCompDirectiveNoFileComp`
	expected := []prompt.Suggest{{Text: "a"}, {Text: "b"}}
	require.Equal(t, expected, parseSuggestions(out))
}

func TestEscapeSpecialCharacters_Spaces(t *testing.T) {
	require.Equal(t, `"string with spaces"`, escapeSpecialCharacters("string with spaces"))
}

func TestEscapeSpecialCharacters_All(t *testing.T) {
	require.Equal(t, "\\\\\\\"\\$\\`\\!", escapeSpecialCharacters("\\\"$`!"))
}

func TestEditCommandTree_RemoveShell(t *testing.T) {
	root := &cobra.Command{}
	sh := &cobra.Command{Use: "lexer"}
	root.AddCommand(sh)

	s := &lexer{root: root}
	s.editCommandTree(sh)
	require.False(t, hasSubcommand(root, "lexer"))
}

func TestEditCommandTree_AddExit(t *testing.T) {
	root := &cobra.Command{}

	s := &lexer{root: root}
	s.editCommandTree(nil)
	require.True(t, hasSubcommand(root, "exit"))
}

func hasSubcommand(cmd *cobra.Command, name string) bool {
	for _, subcommand := range cmd.Commands() {
		if subcommand.Name() == name {
			return true
		}
	}

	return false
}

func TestNew_CommandProperties(t *testing.T) {
	// Test that the New function creates a command with correct properties
	root := &cobra.Command{}

	interactiveCmd := New(root, nil)

	require.Equal(t, "interactive", interactiveCmd.Use)
	require.Equal(t, "Start pactus-shell in interactive mode.", interactiveCmd.Short)
	require.NotNil(t, interactiveCmd.Run)
}

func TestNew_CommandName(t *testing.T) {
	// Test that the command is named "interactive" not "shell"
	root := &cobra.Command{}

	interactiveCmd := New(root, nil)

	require.Equal(t, "interactive", interactiveCmd.Name())
	require.NotEqual(t, "shell", interactiveCmd.Name())
}

func TestNew_WithOptions(t *testing.T) {
	// Test that New function accepts prompt options
	root := &cobra.Command{}

	// Test with some dummy options
	interactiveCmd := New(root, nil,
		prompt.OptionPrefix("test> "),
		prompt.OptionShowCompletionAtStart(),
	)

	require.Equal(t, "interactive", interactiveCmd.Use)
	require.NotNil(t, interactiveCmd.Run)
}
func TestNew_InteractiveCommandCreation(t *testing.T) {
	// Test that the New function creates a command with the exact changed properties
	root := &cobra.Command{}

	interactiveCmd := New(root, nil)

	// Verify the specific changed lines
	require.Equal(t, "interactive", interactiveCmd.Use)
	require.Equal(t, "Start pactus-shell in interactive mode.", interactiveCmd.Short)
	require.NotNil(t, interactiveCmd.Run)

	// Test that it's not the old values
	require.NotEqual(t, "shell", interactiveCmd.Use)
	require.NotEqual(t, "Start an interactive shell.", interactiveCmd.Short)
}

func TestNew_WithAllOptions(t *testing.T) {
	root := &cobra.Command{}
	refresh := func() *cobra.Command { return root }

	interactiveCmd := New(root, refresh,
		prompt.OptionPrefix("test> "),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
	)

	// Verify the changed properties are correct
	require.Equal(t, "interactive", interactiveCmd.Use)
	require.Equal(t, "Start pactus-shell in interactive mode.", interactiveCmd.Short)
	require.NotNil(t, interactiveCmd.Run)
}

func TestNew_CommandNameNotShell(t *testing.T) {
	// Explicitly test that the command is NOT named "shell" anymore
	root := &cobra.Command{}

	interactiveCmd := New(root, nil)

	// This ensures we're testing the actual change
	require.NotEqual(t, "shell", interactiveCmd.Name())
	require.Equal(t, "interactive", interactiveCmd.Name())
}

func TestChangedLinesExecution(t *testing.T) {
	root := &cobra.Command{}

	cmd := New(root, nil)

	// Verify the exact changes
	require.Equal(t, "interactive", cmd.Use)
	require.Equal(t, "Start pactus-shell in interactive mode.", cmd.Short)

	require.NotEqual(t, "shell", cmd.Use)
	require.NotEqual(t, "Start an interactive shell.", cmd.Short)
}
