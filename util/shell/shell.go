package shell

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

type lexer struct {
	root    *cobra.Command
	refresh func() *cobra.Command
	cache   map[string][]prompt.Suggest
	stdin   *term.State
}

// New creates a Cobra CLI command named "shell" which runs an interactive shell prompt for the root command.
func New(root *cobra.Command, refresh func() *cobra.Command, opts ...prompt.Option) *cobra.Command {
	sh := &lexer{
		root:    root,
		refresh: refresh,
		cache:   make(map[string][]prompt.Suggest),
	}

	prefix := fmt.Sprintf("> %s ", root.Name())
	opts = append(opts, prompt.OptionPrefix(prefix), prompt.OptionShowCompletionAtStart())

	return &cobra.Command{
		Use:   "shell",
		Short: "Start an interactive shell.",
		Run: func(cmd *cobra.Command, _ []string) {
			sh.saveStdin()

			sh.editCommandTree(cmd)
			prompt.New(sh.executor, sh.completer, opts...).Run()

			sh.restoreStdin()
		},
	}
}

func (s *lexer) editCommandTree(shell *cobra.Command) {
	s.root.RemoveCommand(shell)

	// Hide the "completion" command
	if cmd, _, err := s.root.Find([]string{"completion"}); err == nil {
		// TODO: Remove this command
		cmd.Hidden = true
	}

	s.root.AddCommand(&cobra.Command{
		Use:   "exit",
		Short: "Exit the interactive shell.",
		Run: func(*cobra.Command, []string) {
			// TODO: Exit cleanly without help from the os package
			os.Exit(0)
		},
	})

	initDefaultHelpFlag(s.root)
}

func initDefaultHelpFlag(cmd *cobra.Command) {
	cmd.InitDefaultHelpFlag()

	for _, subcommand := range cmd.Commands() {
		initDefaultHelpFlag(subcommand)
	}
}

func (s *lexer) saveStdin() {
	state, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	s.stdin = state
}

func (s *lexer) executor(line string) {
	// Allow command to read from stdin
	s.restoreStdin()

	args, _ := shlex.Split(line)
	_ = execute(s.root, args)

	if s.refresh != nil {
		s.root = s.refresh()
		s.editCommandTree(s.root)
	} else {
		if cmd, _, err := s.root.Find(args); err == nil {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				flag.Changed = false
			})
		}
	}

	s.cache = make(map[string][]prompt.Suggest)
}

func (s *lexer) restoreStdin() {
	if s.stdin != nil {
		_ = term.Restore(int(os.Stdin.Fd()), s.stdin)
	}
}

func (s *lexer) completer(d prompt.Document) []prompt.Suggest {
	args, err := buildCompletionArgs(d.CurrentLine())
	if err != nil {
		return nil
	}

	if !isFlag(args[len(args)-1]) {
		// Clear partial strings to generate all possible completions
		args[len(args)-1] = ""
	}
	key := strings.Join(args, " ")

	suggestions, ok := s.cache[key]
	if !ok {
		out, err := readCommandOutput(s.root, args)
		if err != nil {
			return nil
		}
		suggestions = parseSuggestions(out)
		s.cache[key] = suggestions
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func buildCompletionArgs(input string) ([]string, error) {
	args, err := shlex.Split(input)

	args = append([]string{"__complete"}, args...)
	if input == "" || input[len(input)-1] == ' ' {
		args = append(args, "")
	}

	return args, err
}

func readCommandOutput(cmd *cobra.Command, args []string) (string, error) {
	buf := new(bytes.Buffer)

	stdout := cmd.OutOrStdout()
	stderr := os.Stderr

	cmd.SetOut(buf)
	_, os.Stderr, _ = os.Pipe()

	err := execute(cmd, args)

	cmd.SetOut(stdout)
	os.Stderr = stderr

	return buf.String(), err
}

func execute(cmd *cobra.Command, args []string) error {
	if c, _, err := cmd.Find(args); err == nil {
		// Reset flag values between runs due to a limitation in Cobra
		c.Flags().VisitAll(func(flag *pflag.Flag) {
			if val, ok := flag.Value.(pflag.SliceValue); ok {
				_ = val.Replace([]string{})
			} else {
				_ = flag.Value.Set(flag.DefValue)
			}

			_ = c.Flags().SetAnnotation(flag.Name, cobra.BashCompOneRequiredFlag, []string{"false"})
		})

		c.InitDefaultHelpFlag()
	}

	cmd.SetArgs(args)

	return cmd.Execute()
}

func parseSuggestions(out string) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, 0)

	x := strings.Split(out, "\n")
	if len(x) < 2 {
		return nil
	}

	for _, line := range x[:len(x)-2] {
		l := strings.SplitN(line, "\t", 2)

		if isShorthandFlag(l[0]) {
			continue
		}

		suggestion := prompt.Suggest{Text: escapeSpecialCharacters(l[0])}
		if len(l) > 1 {
			suggestion.Description = l[1]
		}

		suggestions = append(suggestions, suggestion)
	}

	sort.Slice(suggestions, func(i, j int) bool {
		it := suggestions[i].Text
		jt := suggestions[j].Text

		if isFlag(it) && isFlag(jt) {
			return it < jt
		}

		if isFlag(it) {
			return false
		}

		if isFlag(jt) {
			return true
		}

		return it < jt
	})

	return suggestions
}

func escapeSpecialCharacters(val string) string {
	for _, c := range []string{`\`, `"`, "$", "`", "!"} {
		val = strings.ReplaceAll(val, c, `\`+c)
	}

	if strings.ContainsAny(val, " #&*;<>?[]|~") {
		val = fmt.Sprintf(`"%q"`, val)
	}

	return val
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func isShorthandFlag(arg string) bool {
	return isFlag(arg) && !strings.HasPrefix(arg, "--")
}
