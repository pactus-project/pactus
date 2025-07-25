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

	"github.com/k0kubun/go-ansi"
	"github.com/manifoldco/promptui"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/schollz/progressbar/v3"
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
		if !errors.Is(err, promptui.ErrAbort) {
			PrintErrorMsgf("prompt error: %v", err)
		} else {
			PrintWarnMsgf("Aborted.")
		}
		os.Exit(1)
	}

	if result != "" && strings.ToUpper(result[:1]) == "Y" {
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

// PromptSelect prompts create choice menu for select by user.
func PromptSelect(label string, items []string) int {
	prompt := promptui.Select{
		Label:   label,
		Items:   items,
		Pointer: promptui.PipeCursor,
	}

	choice, _, err := prompt.Run()
	FatalErrorCheck(err)

	return choice
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

func PrintErrorMsgf(format string, args ...any) {
	format = "[ERROR] " + format
	if terminalSupported {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", args...)
}

func PrintSuccessMsgf(format string, a ...any) {
	if terminalSupported {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintWarnMsgf(format string, a ...any) {
	if terminalSupported {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgBoldf(format string, a ...any) {
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

func PrintJSONObject(obj any) {
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

func CreateNode(numValidators int, chain genesis.ChainType, workingDir string,
	mnemonic string, walletPassword string,
) ([]string, string, error) {
	// To make process faster, we update the password after creating the addresses
	walletPath := PactusDefaultWalletPath(workingDir)
	wlt, err := wallet.Create(walletPath, mnemonic, "", chain)
	if err != nil {
		return nil, "", err
	}

	validatorAddrs := []string{}
	for i := 0; i < numValidators; i++ {
		addressInfo, err := wlt.NewValidatorAddress(fmt.Sprintf("Validator address %v", i+1))
		if err != nil {
			return nil, "", err
		}
		validatorAddrs = append(validatorAddrs, addressInfo.Address)
	}

	addressInfo, err := wlt.NewEd25519AccountAddress(
		"Reward address", "")
	if err != nil {
		return nil, "", err
	}
	rewardAddr := addressInfo.Address

	confPath := PactusConfigPath(workingDir)
	genPath := PactusGenesisPath(workingDir)

	switch chain {
	case genesis.Mainnet:
		genDoc := genesis.MainnetGenesis()
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, "", err
		}
		err := config.SaveMainnetConfig(confPath)
		if err != nil {
			return nil, "", err
		}
	case genesis.Testnet:
		genDoc := genesis.TestnetGenesis()
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, "", err
		}
		conf := config.DefaultConfigTestnet()
		if err := conf.Save(confPath); err != nil {
			return nil, "", err
		}

	case genesis.Localnet:
		if numValidators < 4 {
			return nil, "", errors.New("LocalNeed needs at least 4 validators")
		}
		genDoc := makeLocalGenesis(*wlt)
		if err := genDoc.SaveToFile(genPath); err != nil {
			return nil, "", err
		}

		conf := config.DefaultConfigLocalnet()
		if err := conf.Save(confPath); err != nil {
			return nil, "", err
		}
	}

	if err := wlt.UpdatePassword("", walletPassword); err != nil {
		return nil, "", err
	}

	if err := wlt.Save(); err != nil {
		return nil, "", err
	}

	return validatorAddrs, rewardAddr, nil
}

// StartNode starts the node from the given working directory.
// The passwordFetcher will be used to fetch the password for the default_wallet if it is encrypted.
// It returns an error if the genesis doc or default_wallet can't be found inside the working directory.
// TODO: write test for me.
func StartNode(workingDir string, passwordFetcher func(*wallet.Wallet) (string, bool),
	configModifier func(cfg *config.Config) *config.Config,
) (*node.Node, *wallet.Wallet, error) {
	conf, gen, err := MakeConfig(workingDir)
	if err != nil {
		return nil, nil, err
	}

	if configModifier != nil {
		conf = configModifier(conf)
	}

	defaultWalletPath := PactusDefaultWalletPath(workingDir)
	wlt, err := wallet.Open(defaultWalletPath, true,
		wallet.WithCustomServers([]string{conf.GRPC.Listen}))
	if err != nil {
		return nil, nil, err
	}

	valAddrsInfo := wlt.AllValidatorAddresses()
	if len(valAddrsInfo) == 0 {
		return nil, nil, errors.New("no validator addresses found in the wallet")
	}

	if len(valAddrsInfo) > 32 {
		PrintWarnMsgf("wallet has more than 32 validator addresses, only the first 32 will be used")
		valAddrsInfo = valAddrsInfo[:32]
	}

	rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, conf.Node.RewardAddresses)
	if err != nil {
		return nil, nil, err
	}

	valKeys, err := MakeValidatorKey(wlt, valAddrsInfo, passwordFetcher)
	if err != nil {
		return nil, nil, err
	}

	node, err := node.NewNode(gen, conf, valKeys, rewardAddrs)
	if err != nil {
		return nil, nil, err
	}

	err = node.Start()
	if err != nil {
		return nil, nil, err
	}

	return node, wlt, nil
}

// makeLocalGenesis makes genesis file for the local network.
func makeLocalGenesis(wlt wallet.Wallet) *genesis.Genesis {
	// Treasury account
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{
		crypto.TreasuryAddress: acc,
	}

	genValNum := 4
	vals := make([]*validator.Validator, genValNum)
	for i := 0; i < genValNum; i++ {
		info := wlt.AddressInfo(wlt.AllValidatorAddresses()[i].Address)
		pub, _ := bls.PublicKeyFromString(info.PublicKey)
		vals[i] = validator.NewValidator(pub, int32(i))
	}

	// create genesis
	params := genesis.DefaultGenesisParams()
	params.BlockVersion = 0
	gen := genesis.MakeGenesis(util.RoundNow(60), accs, vals, params)

	return gen
}

// MakeConfig attempts to load the configuration file and
// returns an instance of the configuration along with the genesis document.
// The genesis document is required to determine the chain type, which influences the configuration settings.
// The function sets various private configurations, such as the "wallets directory" and chain-specific HRP values.
// If the configuration file cannot be loaded, it tries to recover or restore the configuration.
func MakeConfig(workingDir string) (*config.Config, *genesis.Genesis, error) {
	gen, err := genesis.LoadFromFile(PactusGenesisPath(workingDir))
	if err != nil {
		return nil, nil, err
	}

	if !gen.ChainType().IsMainnet() {
		crypto.ToTestnetHRP()
	}

	walletsDir := PactusWalletDir(workingDir)
	confPath := PactusConfigPath(workingDir)

	var defConf *config.Config
	chainType := gen.ChainType()

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
			return nil, nil, err
		}
	}

	// Now we can update the private filed, if any
	genParams := gen.Params()

	conf.Store.TxCacheWindow = genParams.TransactionToLiveInterval
	conf.Store.SeedCacheWindow = genParams.SortitionInterval
	conf.Store.AccountCacheSize = 1024
	conf.Store.PublicKeyCacheSize = 1024

	conf.GRPC.DefaultWalletName = DefaultWalletName
	conf.GRPC.WalletsDir = walletsDir

	conf.WalletManager.ChainType = chainType
	conf.WalletManager.WalletsDir = walletsDir

	if err := conf.BasicCheck(); err != nil {
		return nil, nil, err
	}

	return conf, gen, nil
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

// MakeRewardAddresses generates a list of reward addresses based on wallet and configuration.
// If no reward addresses are provided in the config,
// the function attempts to use Ed25519 or BLS addresses from the wallet.
func MakeRewardAddresses(wlt *wallet.Wallet, valAddrsInfo []vault.AddressInfo,
	confRewardAddrs []string,
) ([]crypto.Address, error) {
	rewardAddrs := make([]crypto.Address, 0, len(valAddrsInfo))

	switch {
	// Case 1: No reward addresses in the config file.
	case len(confRewardAddrs) == 0:
		// Try to use the first Ed25519 address from the wallet as the reward address.
		firstEd25519AddrPath := addresspath.NewPath(
			vault.PurposeBIP44Hardened,
			wlt.CoinType()+addresspath.HardenedKeyStart,
			uint32(crypto.AddressTypeEd25519Account)+addresspath.HardenedKeyStart,
			uint32(0)+addresspath.HardenedKeyStart)

		addrInfo := wlt.AddressFromPath(firstEd25519AddrPath.String())
		if addrInfo == nil {
			// If no Ed25519 address is found, try the first BLS address instead.
			firstBLSAddrPath := addresspath.NewPath(
				vault.PurposeBLS12381Hardened,
				wlt.CoinType()+addresspath.HardenedKeyStart,
				uint32(crypto.AddressTypeBLSAccount)+addresspath.HardenedKeyStart,
				uint32(0))

			addrInfo = wlt.AddressFromPath(firstBLSAddrPath.String())

			if addrInfo == nil {
				return nil, errors.New("unable to find a reward address in the wallet")
			}
		}

		addr, _ := crypto.AddressFromString(addrInfo.Address)
		for i := 0; i < len(valAddrsInfo); i++ {
			rewardAddrs = append(rewardAddrs, addr)
		}

	// Case 2: One reward address is specified in the config file.
	case len(confRewardAddrs) == 1:
		// Use this single address for all validators.
		addr, _ := crypto.AddressFromString(confRewardAddrs[0])
		for i := 0; i < len(valAddrsInfo); i++ {
			rewardAddrs = append(rewardAddrs, addr)
		}

	// Case 3: Each validator has a corresponding reward address in the config file.
	case len(confRewardAddrs) == len(valAddrsInfo):
		for i := 0; i < len(valAddrsInfo); i++ {
			addr, _ := crypto.AddressFromString(confRewardAddrs[i])
			rewardAddrs = append(rewardAddrs, addr)
		}

	default:
		return nil, fmt.Errorf("expected %v reward addresses, but got %v",
			len(valAddrsInfo), len(confRewardAddrs))
	}

	return rewardAddrs, nil
}

func MakeValidatorKey(walletInstance *wallet.Wallet, valAddrsInfo []vault.AddressInfo,
	passwordFetcher func(*wallet.Wallet) (string, bool),
) ([]*bls.ValidatorKey, error) {
	valAddrs := make([]string, len(valAddrsInfo))
	for i := 0; i < len(valAddrs); i++ {
		valAddr, _ := crypto.AddressFromString(valAddrsInfo[i].Address)
		if !valAddr.IsValidatorAddress() {
			return nil, fmt.Errorf("invalid validator address: %s", valAddrsInfo[i].Address)
		}
		valAddrs[i] = valAddr.String()
	}

	valKeys := make([]*bls.ValidatorKey, len(valAddrsInfo))
	password, ok := passwordFetcher(walletInstance)
	if !ok {
		return nil, errors.New("aborted")
	}
	prvKeys, err := walletInstance.PrivateKeys(password, valAddrs)
	if err != nil {
		return nil, err
	}
	for i, prv := range prvKeys {
		valKeys[i] = bls.NewValidatorKey(prv.(*bls.PrivateKey))
	}

	return valKeys, nil
}

func TerminalProgressBar(totalSize int64, barWidth int) *progressbar.ProgressBar {
	if barWidth < 15 {
		barWidth = 15
	}

	opts := []progressbar.Option{
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(barWidth),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	}

	return progressbar.NewOptions64(totalSize, opts...)
}
