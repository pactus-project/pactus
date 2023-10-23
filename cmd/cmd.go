package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

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
		if !errors.Is(promptui.ErrAbort, err) {
			PrintErrorMsgf("prompt error: %v", err)
		} else {
			PrintWarnMsgf("Aborted.")
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

func PrintErrorMsgf(format string, a ...interface{}) {
	if terminalSupported() {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintSuccessMsgf(format string, a ...interface{}) {
	if terminalSupported() {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintWarnMsgf(format string, a ...interface{}) {
	if terminalSupported() {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgBoldf(format string, a ...interface{}) {
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

	PrintInfoMsgf(out.String())
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
			home = filepath.Join(usr.HomeDir, "pactus")
		}
	}
	return home
}

func PactusDefaultWalletPath(home string) string {
	return filepath.Join(home, "wallets", "default_wallet")
}

func PactusGenesisPath(home string) string {
	return filepath.Join(home, "genesis.json")
}

func PactusConfigPath(home string) string {
	return filepath.Join(home, "config.toml")
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

func CreateNode(numValidators int, chain genesis.ChainType, workingDir string,
	mnemonic string, walletPassword string,
) ([]string, []string, error) {
	// To make process faster, we update the password after creating the addresses
	walletPath := PactusDefaultWalletPath(workingDir)
	walletInstance, err := wallet.Create(walletPath, mnemonic, "", chain)
	if err != nil {
		return nil, nil, err
	}

	validatorAddrs := []string{}
	for i := 0; i < numValidators; i++ {
		addr, err := walletInstance.NewValidatorAddress(fmt.Sprintf("Validator address %v", i+1))
		if err != nil {
			return nil, nil, err
		}
		validatorAddrs = append(validatorAddrs, addr)
	}

	rewardAddrs := []string{}
	for i := 0; i < numValidators; i++ {
		addr, err := walletInstance.NewBLSAccountAddress(fmt.Sprintf("Reward address %v", i+1))
		if err != nil {
			return nil, nil, err
		}
		rewardAddrs = append(rewardAddrs, addr)
	}

	confPath := PactusConfigPath(workingDir)
	genPath := PactusGenesisPath(workingDir)

	switch chain {
	case genesis.Mainnet:
		panic("not yet!")
	case genesis.Testnet:
		err = genesis.TestnetGenesis().SaveToFile(genPath)
		if err != nil {
			return nil, nil, err
		}

		err = config.SaveTestnetConfig(confPath, numValidators)
		if err != nil {
			return nil, nil, err
		}
	case genesis.Localnet:
		err = makeLocalGenesis(*walletInstance).SaveToFile(genPath)
		if err != nil {
			return nil, nil, err
		}

		err := config.SaveLocalnetConfig(confPath, numValidators)
		if err != nil {
			return nil, nil, err
		}
	}

	err = walletInstance.UpdatePassword("", walletPassword)
	if err != nil {
		return nil, nil, err
	}

	err = walletInstance.Save()
	if err != nil {
		return nil, nil, err
	}

	return validatorAddrs, rewardAddrs, nil
}

func StartNode(workingDir string, passwordFetcher func(*wallet.Wallet) (string, bool)) (
	*node.Node, *wallet.Wallet, error,
) {
	gen, err := genesis.LoadFromFile(PactusGenesisPath(workingDir))
	if err != nil {
		return nil, nil, err
	}

	if !gen.ChainType().IsMainnet() {
		crypto.AddressHRP = "tpc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	confPath := PactusConfigPath(workingDir)
	conf, err := config.LoadFromFile(confPath, true)
	if err != nil {
		PrintWarnMsgf("Unable to load the config: %s", err)
		PrintInfoMsgf("Attempting to restore the config to the default values...")

		// First, try to open the old config file in non-strict mode
		confBack, err := config.LoadFromFile(confPath, false)
		if err != nil {
			return nil, nil, err
		}

		// Let's create a backup of the config
		confBackupPath := fmt.Sprintf("%v_bak_%s", confPath, time.Now().Format("2006_01_02"))
		err = os.Rename(confPath, confBackupPath)
		if err != nil {
			return nil, nil, err
		}

		// Now, attempt to restore the config file with the number of validators from the old config.
		switch gen.ChainType() {
		case genesis.Mainnet:
			panic("not yet implemented!")

		case genesis.Testnet:
			err = config.SaveTestnetConfig(confPath, confBack.Node.NumValidators)
			if err != nil {
				return nil, nil, err
			}

		case genesis.Localnet:
			err = config.SaveLocalnetConfig(confPath, confBack.Node.NumValidators)
			if err != nil {
				return nil, nil, err
			}

		default:
			return nil, nil, fmt.Errorf("invalid chain type")
		}

		PrintSuccessMsgf("Config restored to the default values")
		conf, _ = config.LoadFromFile(confPath, true) // This time it should be OK
	}

	err = conf.BasicCheck()
	if err != nil {
		return nil, nil, err
	}

	walletPath := PactusDefaultWalletPath(workingDir)
	walletInstance, err := wallet.Open(walletPath, true)
	if err != nil {
		return nil, nil, err
	}
	addrLabels := walletInstance.AddressInfos()

	if len(addrLabels) < conf.Node.NumValidators {
		return nil, nil, fmt.Errorf("not enough addresses in wallet")
	}
	validatorAddrs := make([]string, conf.Node.NumValidators)
	for i := 0; i < conf.Node.NumValidators; i++ {
		valAddr, _ := crypto.AddressFromString(addrLabels[i].Address)
		if !valAddr.IsValidatorAddress() {
			return nil, nil, fmt.Errorf("invalid validator address: %s", addrLabels[i].Address)
		}
		validatorAddrs[i] = valAddr.String()
	}
	valKeys := make([]*bls.ValidatorKey, conf.Node.NumValidators)
	password, ok := passwordFetcher(walletInstance)
	if !ok {
		return nil, nil, fmt.Errorf("aborted")
	}
	prvKeys, err := walletInstance.PrivateKeys(password, validatorAddrs)
	if err != nil {
		return nil, nil, err
	}
	for i, prv := range prvKeys {
		valKeys[i] = bls.NewValidatorKey(prv.(*bls.PrivateKey))
	}

	// Create reward addresses
	rewardAddrs := make([]crypto.Address, 0, conf.Node.NumValidators)
	if len(conf.Node.RewardAddresses) != 0 {
		for _, addrStr := range conf.Node.RewardAddresses {
			addr, _ := crypto.AddressFromString(addrStr)
			rewardAddrs = append(rewardAddrs, addr)
		}
	} else {
		for i := conf.Node.NumValidators; i < len(addrLabels); i++ {
			addr, _ := crypto.AddressFromString(addrLabels[i].Address)
			if addr.IsAccountAddress() {
				rewardAddrs = append(rewardAddrs, addr)
				if len(rewardAddrs) == conf.Node.NumValidators {
					break
				}
			}
		}
	}
	if len(rewardAddrs) != conf.Node.NumValidators {
		return nil, nil, fmt.Errorf("not enough addresses in wallet")
	}

	nodeInstance, err := node.NewNode(gen, conf, valKeys, rewardAddrs)
	if err != nil {
		return nil, nil, err
	}

	err = nodeInstance.Start()
	if err != nil {
		return nil, nil, err
	}

	return nodeInstance, walletInstance, nil
}

// makeLocalGenesis makes genesis file for the local network.
func makeLocalGenesis(w wallet.Wallet) *genesis.Genesis {
	// Treasury account
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{
		crypto.TreasuryAddress: acc,
	}

	vals := make([]*validator.Validator, 4)
	for i := 0; i < 4; i++ {
		info := w.AddressInfo(w.AddressInfos()[i].Address)
		pub, _ := bls.PublicKeyFromString(info.PublicKey)
		vals[i] = validator.NewValidator(pub, int32(i))
	}

	// create genesis
	params := param.DefaultParams()
	params.BlockVersion = 0
	gen := genesis.MakeGenesis(util.RoundNow(60), accs, vals, params)
	return gen
}
