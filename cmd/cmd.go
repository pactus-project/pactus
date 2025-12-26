package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/signal"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
)

const (
	DefaultHomeDirName    = "pactus"
	DefaultWalletsDirName = "wallets"
	DefaultWalletName     = "default_wallet"
)

func PactusDefaultHomeDir() string {
	home := ""
	usr, err := user.Current()
	if err != nil {
		terminal.PrintWarnMsgf("unable to get current user: %v", err)
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

func PactusLockFilePath(home string) string {
	return filepath.Join(home, ".pactus.lock")
}

func PactusDefaultWalletPath(home string) string {
	return filepath.Join(PactusWalletDir(home), DefaultWalletName)
}

func PactusDaemonName() string {
	if runtime.GOOS == "windows" {
		return "pactus-daemon.exe"
	}

	return "./pactus-daemon"
}

func CreateNode(numValidators int, chain genesis.ChainType, workingDir string,
	mnemonic string, walletPassword string,
) (*wallet.Wallet, string, error) {
	// To make process faster, we update the password after creating the addresses
	walletPath := PactusDefaultWalletPath(workingDir)
	wlt, err := wallet.Create(context.Background(),
		walletPath, mnemonic, "", chain, wallet.WithOfflineMode())
	if err != nil {
		return nil, "", err
	}

	for i := 0; i < numValidators; i++ {
		_, _ = wlt.NewAddress(crypto.AddressTypeValidator, fmt.Sprintf("Validator address %v", i+1))
	}
	rewardAddrInfo, _ := wlt.NewAddress(crypto.AddressTypeEd25519Account, "Reward address",
		wallet.WithPassword(walletPassword))

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
			return nil, "", errors.New("localnet needs at least 4 validators")
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

	return wlt, rewardAddrInfo.Address, nil
}

// StartNode starts the node from the given working directory.
// The passwordFetcher will be used to fetch the password for the default_wallet if it is encrypted.
// It returns an error if the genesis doc or default_wallet can't be found inside the working directory.
// TODO: write test for me.
func StartNode(workingDir string, passwordFetcher func() (string, bool),
	configModifier func(cfg *config.Config) *config.Config,
) (*node.Node, error) {
	conf, gen, err := MakeConfig(workingDir)
	if err != nil {
		return nil, err
	}

	if configModifier != nil {
		conf = configModifier(conf)
	}

	defaultWalletPath := PactusDefaultWalletPath(workingDir)
	wlt, err := wallet.Open(context.Background(), defaultWalletPath, wallet.WithOfflineMode())
	if err != nil {
		return nil, err
	}

	valList := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
	if len(valList) == 0 {
		return nil, errors.New("no validator addresses found in the wallet")
	}

	if len(valList) > 32 {
		terminal.PrintWarnMsgf("wallet has more than 32 validator addresses, only the first 32 will be used")
		valList = valList[:32]
	}

	rewardAddrs, err := MakeRewardAddresses(wlt, valList, conf.Node.RewardAddresses)
	if err != nil {
		return nil, err
	}

	valKeys, err := MakeValidatorKey(wlt, valList, passwordFetcher)
	if err != nil {
		return nil, err
	}

	node, err := node.NewNode(gen, conf, valKeys, rewardAddrs)
	if err != nil {
		return nil, err
	}

	err = node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
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
	addrs := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
	for i := 0; i < genValNum; i++ {
		info, _ := wlt.AddressInfo(addrs[i].Address)
		pub, _ := bls.PublicKeyFromString(info.PublicKey)
		vals[i] = validator.NewValidator(pub, int32(i))
	}

	// create genesis
	params := genesis.DefaultGenesisParams()
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
		terminal.PrintWarnMsgf("Unable to load the config: %s", err)
		terminal.PrintInfoMsgf("Attempting to update or restore the config file...")

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
	conf.WalletManager.DefaultWalletName = DefaultWalletName

	if conf.GRPC.Enable {
		conf.WalletManager.GRPCAddress = conf.GRPC.Listen
	}

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
		terminal.PrintSuccessMsgf("Config updated.")
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

		terminal.PrintSuccessMsgf("Config restored to the default values")
		conf, _ = config.LoadFromFile(confPath, true, defConf) // This time it should be OK
	}

	return conf, err
}

// MakeRewardAddresses generates a list of reward addresses based on wallet and configuration.
// If no reward addresses are provided in the config,
// the function attempts to use Ed25519 or BLS addresses from the wallet.
func MakeRewardAddresses(wlt *wallet.Wallet, valList []types.AddressInfo,
	confRewardAddrs []string,
) ([]crypto.Address, error) {
	rewardAddrs := make([]crypto.Address, 0, len(valList))

	switch {
	// Case 1: No reward addresses in the config file.
	case len(confRewardAddrs) == 0:
		var addrInfo *types.AddressInfo
		// Try to use the first Ed25519 address from the wallet as the reward address.
		ed25519Addrs := wlt.ListAddresses(wallet.WithAddressType(crypto.AddressTypeEd25519Account))
		if len(ed25519Addrs) == 0 {
			// If no Ed25519 address is found, try the first BLS address instead.
			blsAddrs := wlt.ListAddresses(wallet.WithAddressType(crypto.AddressTypeBLSAccount))
			if len(blsAddrs) == 0 {
				return nil, errors.New("unable to find a reward address in the wallet")
			}

			addrInfo = &blsAddrs[0]
		} else {
			addrInfo = &ed25519Addrs[0]
		}

		addr, _ := crypto.AddressFromString(addrInfo.Address)
		for i := 0; i < len(valList); i++ {
			rewardAddrs = append(rewardAddrs, addr)
		}

	// Case 2: One reward address is specified in the config file.
	case len(confRewardAddrs) == 1:
		// Use this single address for all validators.
		addr, _ := crypto.AddressFromString(confRewardAddrs[0])
		for i := 0; i < len(valList); i++ {
			rewardAddrs = append(rewardAddrs, addr)
		}

	// Case 3: Each validator has a corresponding reward address in the config file.
	case len(confRewardAddrs) == len(valList):
		for _, addrStr := range confRewardAddrs {
			addr, _ := crypto.AddressFromString(addrStr)
			rewardAddrs = append(rewardAddrs, addr)
		}

	default:
		return nil, fmt.Errorf("expected %v reward addresses, but got %v",
			len(valList), len(confRewardAddrs))
	}

	return rewardAddrs, nil
}

func MakeValidatorKey(walletInstance *wallet.Wallet, valAddrsInfo []types.AddressInfo,
	passwordFetcher func() (string, bool),
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

	walletPassword := ""
	if walletInstance.IsEncrypted() {
		password, ok := passwordFetcher()
		if !ok {
			return nil, errors.New("aborted")
		}

		walletPassword = password
	}

	prvKeys, err := walletInstance.PrivateKeys(walletPassword, valAddrs)
	if err != nil {
		return nil, err
	}
	for i, prv := range prvKeys {
		valKeys[i] = bls.NewValidatorKey(prv.(*bls.PrivateKey))
	}

	return valKeys, nil
}

func RecoverWalletAddresses(wlt *wallet.Wallet, password string) {
	ctx, cancel := context.WithCancel(context.Background())
	signal.HandleSignals(func(os.Signal) {
		cancel()
	}, syscall.SIGINT, syscall.SIGTERM)

	terminal.PrintInfoMsgf("🔄 Recovering wallet addresses...")
	terminal.PrintInfoMsgf("   Press 'Ctrl+C' to abort if needed")
	terminal.PrintLine()

	index := 0
	err := wlt.RecoveryAddresses(ctx, password, func(addr string) {
		terminal.PrintInfoMsgf("%d. %s", index+1, addr)
		index++
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			terminal.PrintWarnMsgf("Address recovery aborted")
		} else {
			terminal.PrintErrorMsgf("Address recovery failed: %v", err)
		}
	}
}
