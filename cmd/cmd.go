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
	"strconv"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet"
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
		Label:   label,
		Mask:    '*',
		Pointer: promptui.PipeCursor,
	}
	password, err := prompt.Run()
	FatalErrorCheck(err)

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
			Pointer:  promptui.PipeCursor,
		}

		_, err := confirmPrompt.Run()
		FatalErrorCheck(err)
	}

	return password
}

// PromptConfirm prompts user to confirm the operation.
func PromptConfirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Pointer:   promptui.PipeCursor,
	}
	result, err := prompt.Run()
	if err != nil {
		if err != promptui.ErrAbort {
			PrintErrorMsg("prompt error: %v", err)
		} else {
			PrintWarnMsg("Aborted.")
		}
		os.Exit(1)
	}

	if len(result) > 0 && strings.ToUpper(result[:1]) == "Y" {
		return true
	}
	return false
}

// PromptInput prompts for an input string.
func PromptInput(label string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Pointer: promptui.PipeCursor,
	}
	result, err := prompt.Run()
	FatalErrorCheck(err)

	return result
}

// PromptInputWithSuggestion prompts the user for an input string with a suggestion.
func PromptInputWithSuggestion(label, suggestion string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: suggestion,
		Pointer: promptui.PipeCursor,
	}
	result, err := prompt.Run()
	FatalErrorCheck(err)

	return result
}

// PromptInputWithRange prompts the user for an input integer within a specified range.
func PromptInputWithRange(label string, def, min, max int) int {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   fmt.Sprintf("%v", def),
		IsVimMode: true,
		Pointer:   promptui.PipeCursor,
		Validate: func(input string) error {
			num, err := strconv.Atoi(input)
			if err != nil {
				return err
			}
			if num < min || num > max {
				return fmt.Errorf("enter a number between %v and %v", min, max)
			}

			return nil
		},
	}
	result, err := prompt.Run()
	FatalErrorCheck(err)

	num, err := strconv.Atoi(result)
	FatalErrorCheck(err)

	return num
}

func FatalErrorCheck(err error) {
	if err != nil {
		if terminalSupported() {
			fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		} else {
			fmt.Printf("%s\n", err.Error())
		}

		os.Exit(1)
	}
}

func PrintErrorMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintSuccessMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintWarnMsg(format string, a ...interface{}) {
	if terminalSupported() {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsg(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgBold(format string, a ...interface{}) {
	if terminalSupported() {
		format = fmt.Sprintf("\033[1m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintLine() {
	fmt.Println()
}

func PrintJSONData(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "   ")
	FatalErrorCheck(err)

	PrintInfoMsg(out.String())
}

func PrintJSONObject(obj interface{}) {
	data, err := json.Marshal(obj)
	FatalErrorCheck(err)

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

		// TODO: remove it before the mainnet launch
		home = path.Join(home, "testnet")
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

func CreateNode(numValidators int, testnet bool, workingDir string,
	mnemonic string, walletPassword string) (
	validatorAddrs []string, rewardAddrs []string, err error) {
	// To make process faster, we update the password after creating the addresses
	network := wallet.NetworkMainNet
	if testnet {
		network = wallet.NetworkTestNet
	}
	walletPath := PactusDefaultWalletPath(workingDir)
	wallet, err := wallet.Create(walletPath, mnemonic, "", network)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < numValidators; i++ {
		addr, err := wallet.DeriveNewAddress(fmt.Sprintf("Validator address %v", i+1))
		if err != nil {
			return nil, nil, err
		}
		validatorAddrs = append(validatorAddrs, addr)
	}

	for i := 0; i < numValidators; i++ {
		addr, err := wallet.DeriveNewAddress(fmt.Sprintf("Reward address %v", i+1))
		if err != nil {
			return nil, nil, err
		}
		rewardAddrs = append(rewardAddrs, addr)
	}

	confPath := PactusConfigPath(workingDir)
	genPath := PactusGenesisPath(workingDir)

	if testnet {
		err = genesis.Testnet().SaveToFile(genPath)
		if err != nil {
			return nil, nil, err
		}

		err = config.SaveTestnetConfig(confPath, numValidators)
		if err != nil {
			return nil, nil, err
		}
	} else {
		panic("not yet!")
	}

	err = wallet.UpdatePassword("", walletPassword)
	if err != nil {
		return nil, nil, err
	}

	err = wallet.Save()
	if err != nil {
		return nil, nil, err
	}

	return validatorAddrs, rewardAddrs, nil
}

func GetKeys(workingDir string, passwordFetcher func(*wallet.Wallet) (string, bool)) (
	*genesis.Genesis, *config.Config, []crypto.Signer, []crypto.Address, *wallet.Wallet, error) {
	gen, err := genesis.LoadFromFile(PactusGenesisPath(workingDir))
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if gen.Params().IsTestnet() {
		crypto.AddressHRP = "tpc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	conf, err := config.LoadFromFile(PactusConfigPath(workingDir))
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = conf.SanityCheck()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	walletPath := PactusDefaultWalletPath(workingDir)
	wallet, err := wallet.Open(walletPath, true)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	addrLabels := wallet.AddressLabels()

	// Create signers
	if len(addrLabels) < conf.Node.NumValidators {
		return nil, nil, nil, nil, nil, fmt.Errorf("not enough addresses in wallet")
	}
	validatorAddrs := make([]string, conf.Node.NumValidators)
	for i := 0; i < conf.Node.NumValidators; i++ {
		validatorAddrs[i] = addrLabels[i].Address
	}
	signers := make([]crypto.Signer, conf.Node.NumValidators)
	password, ok := passwordFetcher(wallet)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("aborted")
	}
	prvKeys, err := wallet.PrivateKeys(password, validatorAddrs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	for i, key := range prvKeys {
		signers[i] = crypto.NewSigner(key)
	}

	// Create reward addresses
	rewardAddrs := make([]crypto.Address, conf.Node.NumValidators)
	if len(conf.Node.RewardAddresses) != 0 {
		for i, addrStr := range conf.Node.RewardAddresses {
			rewardAddrs[i], _ = crypto.AddressFromString(addrStr)
		}
	} else {
		if len(addrLabels) < 2*conf.Node.NumValidators {
			return nil, nil, nil, nil, nil, fmt.Errorf("not enough addresses in wallet")
		}
		for i := 0; i < conf.Node.NumValidators; i++ {
			rewardAddrs[i], _ =
				crypto.AddressFromString(addrLabels[conf.Node.NumValidators+i].Address)
		}
	}

	return gen, conf, signers, rewardAddrs, wallet, nil
}
