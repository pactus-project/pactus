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
