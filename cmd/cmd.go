//nolint:forbidigo // enable printing function for cmd package
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
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/addresspath"
)

const (
	DefaultHomeDirName    = "pactus"
	DefaultWalletsDirName = "wallets"
	DefaultWalletName     = "default_wallet"
)

var terminalSupported = false

func init() {
	terminalSupported = CheckTerminalSupported()
}

// CheckTerminalSupported returns true if the current terminal supports
// line editing features.
func CheckTerminalSupported() bool {
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
		if terminalSupported {
			fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		} else {
			fmt.Printf("%s\n", err.Error())
		}

		os.Exit(1)
	}
}

func PrintErrorMsgf(format string, a ...interface{}) {
	format = "[ERROR] " + format
	if terminalSupported {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintSuccessMsgf(format string, a ...interface{}) {
	if terminalSupported {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintWarnMsgf(format string, a ...interface{}) {
	if terminalSupported {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgBoldf(format string, a ...interface{}) {
	if terminalSupported {
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

func PactusDefaultHomeDir() string {
	home := ""
	usr, err := user.Current()
	if err != nil {
		PrintWarnMsgf("unable to get current user: %v", err)
	} else {
		home = filepath.Join(usr.HomeDir, home, DefaultHomeDirName)
	}

	return home
}

func PactusWalletDir(home string) string {
	return filepath.Join(home, "wallets")
}

func PactusGenesisPath(home string) string {
	return filepath.Join(home, "genesis.json")
}

func PactusConfigPath(home string) string {
	return filepath.Join(home, "config.toml")
}

func PactusDefaultWalletPath(home string) string {
	return filepath.Join(PactusWalletDir(home), DefaultWalletName)
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

// TODO: write test for me.
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
		genDoc := genesis.MainnetGenesis()
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, nil, err
		}
		conf := config.DefaultConfigMainnet()
		if err := conf.Save(confPath); err != nil {
			return nil, nil, err
		}
	case genesis.Testnet:
		genDoc := genesis.TestnetGenesis()
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, nil, err
		}
		conf := config.DefaultConfigTestnet()
		if err := conf.Save(confPath); err != nil {
			return nil, nil, err
		}

	case genesis.Localnet:
		genDoc := makeLocalGenesis(*walletInstance)
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, nil, err
		}

		conf := config.DefaultConfigLocalnet()
		if err := conf.Save(confPath); err != nil {
			return nil, nil, err
		}
	}

	if err := walletInstance.UpdatePassword("", walletPassword); err != nil {
		return nil, nil, err
	}

	if err := walletInstance.Save(); err != nil {
		return nil, nil, err
	}

	return validatorAddrs, rewardAddrs, nil
}

// TODO: write test for me.
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

	walletsDir := PactusWalletDir(workingDir)
	confPath := PactusConfigPath(workingDir)

	conf, err := MakeConfig(gen, confPath, walletsDir)
	if err != nil {
		return nil, nil, err
	}

	err = conf.BasicCheck()
	if err != nil {
		return nil, nil, err
	}

	defaultWalletPath := PactusDefaultWalletPath(workingDir)
	walletInstance, err := wallet.Open(defaultWalletPath, true)
	if err != nil {
		return nil, nil, err
	}

	valAddrsInfo := walletInstance.AllValidatorAddresses()
	if len(valAddrsInfo) == 0 {
		return nil, nil, fmt.Errorf("no validator addresses found in the wallet")
	}

	if len(valAddrsInfo) > 32 {
		PrintWarnMsgf("wallet has more than 32 validator addresses, only the first 32 will be used")
		valAddrsInfo = valAddrsInfo[:32]
	}

	if len(conf.Node.RewardAddresses) > 0 &&
		len(conf.Node.RewardAddresses) != len(valAddrsInfo) {
		return nil, nil, fmt.Errorf("reward addresses should be %v", len(valAddrsInfo))
	}

	valAddrs := make([]string, len(valAddrsInfo))
	for i := 0; i < len(valAddrs); i++ {
		valAddr, _ := crypto.AddressFromString(valAddrsInfo[i].Address)
		if !valAddr.IsValidatorAddress() {
			return nil, nil, fmt.Errorf("invalid validator address: %s", valAddrsInfo[i].Address)
		}
		valAddrs[i] = valAddr.String()
	}

	valKeys := make([]*bls.ValidatorKey, len(valAddrsInfo))
	password, ok := passwordFetcher(walletInstance)
	if !ok {
		return nil, nil, fmt.Errorf("aborted")
	}
	prvKeys, err := walletInstance.PrivateKeys(password, valAddrs)
	if err != nil {
		return nil, nil, err
	}
	for i, prv := range prvKeys {
		valKeys[i] = bls.NewValidatorKey(prv.(*bls.PrivateKey))
	}

	// Create reward addresses
	rewardAddrs := make([]crypto.Address, 0, len(valAddrsInfo))
	if len(conf.Node.RewardAddresses) != 0 {
		for _, addrStr := range conf.Node.RewardAddresses {
			addr, _ := crypto.AddressFromString(addrStr)
			rewardAddrs = append(rewardAddrs, addr)
		}
	} else {
		for i := 0; i < len(valAddrsInfo); i++ {
			valAddrPath, _ := addresspath.FromString(valAddrsInfo[i].Path)
			accAddrPath := addresspath.NewPath(
				valAddrPath.Purpose(),
				valAddrPath.CoinType(),
				uint32(crypto.AddressTypeBLSAccount)+hdkeychain.HardenedKeyStart,
				valAddrPath.AddressIndex())

			addrInfo := walletInstance.AddressFromPath(accAddrPath.String())
			if addrInfo == nil {
				return nil, nil, fmt.Errorf("unable to find reward address for: %s [%s]",
					valAddrsInfo[i].Address, accAddrPath)
			}

			addr, _ := crypto.AddressFromString(addrInfo.Address)
			rewardAddrs = append(rewardAddrs, addr)
		}
	}

	// Check if reward addresses are account address
	for _, addr := range rewardAddrs {
		if !addr.IsAccountAddress() {
			return nil, nil, fmt.Errorf("reward address is not an account address: %s", addr)
		}
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

// TODO: write test for me.
func MakeConfig(genDoc *genesis.Genesis, confPath, walletsDir string) (*config.Config, error) {
	var defConf *config.Config
	chainType := genDoc.ChainType()

	switch chainType {
	case genesis.Mainnet:
		defConf = config.DefaultConfigMainnet()
	case genesis.Testnet:
		defConf = config.DefaultConfigTestnet()
	case genesis.Localnet:
		defConf = config.DefaultConfigLocalnet()
	}

	conf, err := config.LoadFromFile(confPath, true, defConf)
	if err != nil {
		PrintWarnMsgf("Unable to load the config: %s", err)
		PrintInfoMsgf("Attempting to update or restore the config file...")

		conf, err = RecoverConfig(confPath, defConf, chainType)
		if err != nil {
			return nil, err
		}
	}

	// Now we can update the private filed, if any
	genParams := genDoc.Params()

	conf.Store.TxCacheSize = genParams.TransactionToLiveInterval
	conf.Store.SortitionCacheSize = genParams.SortitionInterval
	conf.Store.AccountCacheSize = 1024
	conf.Store.PublicKeyCacheSize = 1024

	conf.GRPC.DefaultWalletName = DefaultWalletName
	conf.GRPC.WalletsDir = walletsDir

	return conf, nil
}

func RecoverConfig(confPath string, defConf *config.Config, chainType genesis.ChainType) (*config.Config, error) {
	// Try to attempt to load config in non-strict mode
	conf, err := config.LoadFromFile(confPath, false, defConf)

	// Create a backup of the config
	if util.PathExists(confPath) {
		confBackupPath := fmt.Sprintf("%v_bak_%s", confPath, time.Now().Format("2006-01-02T15-04-05"))
		renameErr := os.Rename(confPath, confBackupPath)
		if renameErr != nil {
			return nil, renameErr
		}
	}

	if err == nil {
		err := conf.Save(confPath)
		if err != nil {
			return nil, err
		}
		PrintSuccessMsgf("Config updated.")
	} else {
		switch chainType {
		case genesis.Mainnet:
			err = config.SaveMainnetConfig(confPath)
			if err != nil {
				return nil, err
			}

		case genesis.Testnet,
			genesis.Localnet:
			err = defConf.Save(confPath)
			if err != nil {
				return nil, err
			}
		}

		PrintSuccessMsgf("Config restored to the default values")
		conf, _ = config.LoadFromFile(confPath, true, defConf) // This time it should be OK
	}

	return conf, err
}
