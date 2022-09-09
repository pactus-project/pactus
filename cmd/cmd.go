package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
)

var Pactus = ``

// terminalSupported returns true if the current terminal supports
// line editing features.
func terminalSupported() bool {
	bad := map[string]bool{"": true, "dumb": true, "cons25": true}
	return !bad[strings.ToLower(os.Getenv("TERM"))]
}

// PromptPassword prompts the user for a password. Set confirmation to true
// to require the user to confirm the password.
func PromptPassword(label string, confirmation bool) string {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		PrintErrorMsg("Failed to read password: %v", err)
		os.Exit(1)
	}

	if confirmation {
		validate := func(input string) error {
			if input != password {
				return errors.New("passwords do not match")
			}
			return nil
		}

		confirmPrompt := promptui.Prompt{
			Label:    "Confirm password",
			Validate: validate,
			Mask:     '*',
		}

		_, err := confirmPrompt.Run()
		if err != nil {
			PrintErrorMsg("prompt error: %v", err)
			os.Exit(1)
		}
	}

	return password
}

// PromptConfirm prompts user to confirm the operation.
func PromptConfirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		if err != promptui.ErrAbort {
			PrintErrorMsg("prompt error: %v", err)
		}
		os.Exit(1)
	}

	if len(result) > 0 && strings.ToUpper(result[:1]) == "Y" {
		return true
	}
	return false
}

// Promptlabel prompts for an input string.
func PromptInput(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		PrintErrorMsg("prompt error: %v", err)
		os.Exit(1)
	}
	return result
}

// Promptlabel prompts for an input string with a suggestion.
func PromptInputWithSuggestion(label, suggestion string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: suggestion,
	}
	result, err := prompt.Run()
	if err != nil {
		PrintErrorMsg("prompt error: %v", err)
		os.Exit(1)
	}
	return result
}

func PrintDangerMsg(format string, a ...interface{}) {
	if terminalSupported() {
		format = fmt.Sprintf("\033[31m%s\033[0m\n", format)
	}
	fmt.Printf(format, a...)
}

func PrintErrorMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format)
	}
	fmt.Printf(format, a...)
}

func PrintSuccessMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m\n", format)
	}
	fmt.Printf(format, a...)
}

func PrintWarnMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m\n", format)
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

func PactusHomeDir() string {
	home := ""
	usr, err := user.Current()
	if err == nil {
		// Running as root, probably inside docker
		if usr.HomeDir == "/root" {
			home = "/pactus/"
		} else {
			home = path.Join(usr.HomeDir, "pactus")
		}
	}
	return home
}

func PactusDefaultWalletPath(home string) string {
	return path.Join(home, "wallets"+string(os.PathSeparator)+"default_wallet")
}

func PactusGenesisPath(home string) string {
	return path.Join(home, "genesis.json")
}

func PactusConfigPath(home string) string {
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
