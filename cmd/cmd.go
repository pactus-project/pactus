package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/peterh/liner"
)

var ZARB = `
 ███████╗  █████╗  ██████╗  ██████╗
 ╚══███╔╝ ██╔══██╗ ██╔══██╗ ██╔══██╗
   ███╔╝  ███████║ ██████╔╝ ██████╔╝
  ███╔╝   ██╔══██║ ██╔══██╗ ██╔══██╗
 ███████╗ ██║  ██║ ██║  ██║ ██████╔╝
 ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═════╝
`

type terminalPrompter struct {
	*liner.State
	warned     bool
	supported  bool
	normalMode liner.ModeApplier
	rawMode    liner.ModeApplier
}

// Stdin holds the stdin line reader (also using stdout for printing prompts).
// Only this reader may be used for input because it keeps an internal buffer.
// var
var Stdin = newTerminalPrompter()

// newTerminalPrompter creates a liner based user input prompter working off the
// standard input and output streams.
func newTerminalPrompter() *terminalPrompter {
	p := new(terminalPrompter)
	// Get the original mode before calling NewLiner.
	// This is usually regular "cooked" mode where characters echo.
	normalMode, _ := liner.TerminalMode()
	// Turn on liner. It switches to raw mode.
	p.State = liner.NewLiner()
	rawMode, err := liner.TerminalMode()
	if err != nil || !liner.TerminalSupported() {
		p.supported = false
	} else {
		p.supported = true
		p.normalMode = normalMode
		p.rawMode = rawMode
		// Switch back to normal mode while we're not prompting.
		applyMode(normalMode)
	}
	p.SetCtrlCAborts(true)
	p.SetTabCompletionStyle(liner.TabPrints)
	p.SetMultiLineMode(false)
	return p
}

// PromptPassword displays the given prompt to the user and requests some textual
// data to be entered, but one which must not be echoed out into the terminal.
// The method returns the input provided by the user.
func (p *terminalPrompter) PromptPassword(prompt string) (string, error) {
	if p.supported {
		applyMode(p.rawMode)
		defer applyMode(p.normalMode)
		return p.State.PasswordPrompt(prompt)
	}
	if !p.warned {
		PrintWarnMsg("!! Unsupported terminal, password will be echoed.")
		p.warned = true
	}
	// Just as in Prompt, handle printing the prompt here instead of relying on liner.
	fmt.Print(prompt)
	pass, err := p.State.Prompt("")
	fmt.Println()
	return pass, err
}

// PromptInput displays the given prompt to the user and requests some textual
// data to be entered, returning the input of the user.
func (p *terminalPrompter) PromptInput(prompt string) (string, error) {
	if p.supported {
		applyMode(p.rawMode)
		defer applyMode(p.normalMode)
	} else {
		// liner tries to be smart about printing the prompt
		// and doesn't print anything if input is redirected.
		// Un-smart it by printing the prompt always.
		fmt.Print(prompt)
		prompt = ""
		defer fmt.Println()
	}
	return p.State.Prompt(prompt)
}

// PromptConfirm displays the given prompt to the user and requests a boolean
// choice to be made, returning that choice.
func (p *terminalPrompter) PromptConfirm(prompt string) (bool, error) {
	input, err := p.PromptInput(prompt + " [y/N] ")
	if len(input) > 0 && strings.ToUpper(input[:1]) == "Y" {
		return true, nil
	}
	return false, err
}

// PromptPassword prompts the user for a password. Set confirmation to true
// to require the user to confirm the password.
func PromptPassword(prompt string, confirmation bool) string {
	password, err := Stdin.PromptPassword(prompt)
	if err != nil {
		PrintErrorMsg("Failed to read password: %v", err)
		os.Exit(1)
	}

	if confirmation {
		confirm, err := Stdin.PromptPassword("Repeat password: ")
		if err != nil {
			PrintErrorMsg("Failed to read password confirmation: %v", err)
			os.Exit(1)
		}
		if password != confirm {
			PrintErrorMsg("passwords do not match")
			os.Exit(1)
		}
	}

	return password
}

// PromptConfirm prompts user to confirm the operation
func PromptConfirm(prompt string) bool {
	input, err := Stdin.PromptConfirm(prompt)
	if err != nil {
		PrintErrorMsg("Failed to read input: %v", err)
		os.Exit(1)
	}
	return input
}

// Promptlabel prompts for an input string
func PromptInput(prompt string) string {
	input, err := Stdin.PromptInput(prompt)
	if err != nil {
		PrintErrorMsg("Failed to read input: %v", err)
		os.Exit(1)
	}
	return input
}

// Promptlabel prompts for an input string with a suggestion
func PromptInputWithSuggestion(prompt, suggestion string) string {
	input, err := Stdin.PromptWithSuggestion(prompt, suggestion, 0)
	if err != nil {
		PrintErrorMsg("Failed to read input: %v", err)
		os.Exit(1)
	}
	return input
}

func PrintDangerMsg(format string, a ...interface{}) {
	if liner.TerminalSupported() {
		format = fmt.Sprintf("\033[31m%s\033[0m\n", format)
	}
	fmt.Printf(format, a...)
}

func PrintErrorMsg(format string, a ...interface{}) {
	if liner.TerminalSupported() {
		format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format) //Print error msg with red color
	}
	fmt.Printf(format, a...)
}

func PrintSuccessMsg(format string, a ...interface{}) {
	if liner.TerminalSupported() {
		format = fmt.Sprintf("\033[32m%s\033[0m\n", format) //Print successful msg with green color
	}
	fmt.Printf(format, a...)
}

func PrintWarnMsg(format string, a ...interface{}) {
	if liner.TerminalSupported() {
		format = fmt.Sprintf("\033[33m%s\033[0m\n", format) //Print warning msg with yellow color
	}
	fmt.Printf(format, a...)
}

func PrintInfoMsg(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func PrintLine() {
	fmt.Println()
}

func PrintJSONData(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "   ")
	if err != nil {
		PrintErrorMsg("json.Indent error: %v", err)
		return
	}
	PrintInfoMsg(out.String())
}

func PrintJSONObject(obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		PrintErrorMsg("json.Marshal error: %v", err)
		return
	}
	PrintJSONData(data)
}

func ZarbHomeDir() string {
	home := ""
	usr, err := user.Current()
	if err == nil {
		// Running as root, probably inside docker
		if usr.HomeDir == "/root" {
			home = "/zarb/"
		} else {
			home = path.Join(usr.HomeDir, "zarb")
		}
	}
	return home
}

func ZarbDefaultWalletPath(home string) string {
	return path.Join(home, "wallets"+string(os.PathSeparator)+"default_wallet")
}

func ZarbGenesisPath(home string) string {
	return path.Join(home, "genesis.json")
}

func ZarbConfigPath(home string) string {
	return path.Join(home, "config.toml")
}

// TrapSignal traps SIGINT and SIGTERM and terminates the server correctly.
func TrapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}

func applyMode(m liner.ModeApplier) {
	if err := m.ApplyMode(); err != nil {
		panic(err)
	}
}
