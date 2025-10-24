package cmd

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeConfig(t *testing.T) {
	t.Run("No genesis file, Should return error", func(t *testing.T) {
		workingDir := util.TempDirPath()

		_, _, err := MakeConfig(workingDir)
		assert.Error(t, err)
	})

	t.Run("No Config file, Should recover it", func(t *testing.T) {
		workingDir := util.TempDirPath()
		genPath := PactusGenesisPath(workingDir)
		gen := genesis.MainnetGenesis()
		err := gen.SaveToFile(genPath)
		require.NoError(t, err)

		_, _, err = MakeConfig(workingDir)
		assert.NoError(t, err)
	})

	t.Run("Invalid Config file, Should recover it", func(t *testing.T) {
		workingDir := util.TempDirPath()
		genPath := PactusGenesisPath(workingDir)
		confPath := PactusConfigPath(workingDir)

		gen := genesis.MainnetGenesis()
		err := gen.SaveToFile(genPath)
		require.NoError(t, err)

		err = util.WriteFile(confPath, []byte("invalid-config"))
		require.NoError(t, err)

		_, _, err = MakeConfig(workingDir)
		assert.NoError(t, err)
	})

	t.Run("Everything is good", func(t *testing.T) {
		workingDir := util.TempDirPath()
		genPath := PactusGenesisPath(workingDir)
		confPath := PactusConfigPath(workingDir)

		gen := genesis.MainnetGenesis()
		err := gen.SaveToFile(genPath)
		require.NoError(t, err)

		err = config.SaveMainnetConfig(confPath)
		require.NoError(t, err)

		_, _, err = MakeConfig(workingDir)
		assert.NoError(t, err)
	})
}

func TestPathsUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	tests := []struct {
		home                      string
		expectedWalletDir         string
		expectedDefaultWalletPath string
		expectedGenesisPath       string
		expectedConfigPath        string
	}{
		{
			"/home/pactus",
			"/home/pactus/wallets",
			"/home/pactus/wallets/default_wallet",
			"/home/pactus/genesis.json",
			"/home/pactus/config.toml",
		},
		{
			"/home/pactus/",
			"/home/pactus/wallets",
			"/home/pactus/wallets/default_wallet",
			"/home/pactus/genesis.json",
			"/home/pactus/config.toml",
		},
	}

	for _, tt := range tests {
		walletDir := PactusWalletDir(tt.home)
		defaultWalletPath := PactusDefaultWalletPath(tt.home)
		genesisPath := PactusGenesisPath(tt.home)
		configPath := PactusConfigPath(tt.home)

		assert.Equal(t, tt.expectedWalletDir, walletDir)
		assert.Equal(t, tt.expectedDefaultWalletPath, defaultWalletPath)
		assert.Equal(t, tt.expectedGenesisPath, genesisPath)
		assert.Equal(t, tt.expectedConfigPath, configPath)
	}
}

func TestPathsWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		return
	}
	tests := []struct {
		home                      string
		expectedWalletDir         string
		expectedDefaultWalletPath string
		expectedGenesisPath       string
		expectedConfigPath        string
	}{
		{
			"c:\\pactus",
			"c:\\pactus\\wallets",
			"c:\\pactus\\wallets\\default_wallet",
			"c:\\pactus\\genesis.json",
			"c:\\pactus\\config.toml",
		},
		{
			"c:\\home\\",
			"c:\\home\\wallets",
			"c:\\home\\wallets\\default_wallet",
			"c:\\home\\genesis.json",
			"c:\\home\\config.toml",
		},
	}

	for _, tt := range tests {
		walletDir := PactusWalletDir(tt.home)
		defaultWalletPath := PactusDefaultWalletPath(tt.home)
		genesisPath := PactusGenesisPath(tt.home)
		configPath := PactusConfigPath(tt.home)

		assert.Equal(t, tt.expectedWalletDir, walletDir)
		assert.Equal(t, tt.expectedDefaultWalletPath, defaultWalletPath)
		assert.Equal(t, tt.expectedGenesisPath, genesisPath)
		assert.Equal(t, tt.expectedConfigPath, configPath)
	}
}

func TestMakeRewardAddresses(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	setupWallet := func() *wallet.Wallet {
		walletPath := util.TempFilePath()
		mnemonic, _ := wallet.GenerateMnemonic(128)
		wlt, err := wallet.Create(walletPath, mnemonic, "", genesis.Mainnet)
		assert.NoError(t, err)

		_, _ = wlt.NewValidatorAddress("Validator 1")
		_, _ = wlt.NewValidatorAddress("Validator 2")
		_, _ = wlt.NewValidatorAddress("Validator 3")

		return wlt
	}

	t.Run("No reward addresses in wallet", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{}
		_, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.ErrorContains(t, err, "unable to find a reward address in the wallet")
	})

	t.Run("Wallet with one Ed25519 address", func(t *testing.T) {
		wlt := setupWallet()

		addr1Info, _ := wlt.NewEd25519AccountAddress("", "")
		_, _ = wlt.NewEd25519AccountAddress("", "")
		_, _ = wlt.NewBLSAccountAddress("")

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[1].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[2].String(), addr1Info.Address)
	})

	t.Run("Wallet with one BLS address", func(t *testing.T) {
		wlt := setupWallet()

		addr1Info, _ := wlt.NewBLSAccountAddress("")
		_, _ = wlt.NewBLSAccountAddress("")

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[1].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[2].String(), addr1Info.Address)
	})

	t.Run("One reward address in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{
			ts.RandAccAddress().String(),
		}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[1].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[2].String(), confRewardAddresses[0])
	})

	t.Run("Three reward addresses in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
		}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[1].String(), confRewardAddresses[1])
		assert.Equal(t, rewardAddrs[2].String(), confRewardAddresses[2])
	})

	t.Run("Insufficient reward addresses in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.AllValidatorAddresses()
		confRewardAddresses := []string{
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
		}
		_, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.ErrorContains(t, err, "expected 3 reward addresses, but got 2")
	})
}

func TestCreateNode(t *testing.T) {
	tests := []struct {
		name           string
		numValidators  int
		chain          genesis.ChainType
		workingDir     string
		mnemonic       string
		withErr        bool
		validatorAddrs []string
		rewardAddrs    string
	}{
		{
			name:           "Create node for Mainnet",
			numValidators:  1,
			chain:          genesis.Mainnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: []string{"pc1pqpu5tkuctj6ecxjs85f9apm802hhc65amwhuyw"},
			rewardAddrs:    "pc1rkg0nhswqj85wnz9sm0g9kfkxj68lfx9lhftl8n",
			withErr:        false,
		},
		{
			name:           "Create node for Testnet",
			numValidators:  1,
			chain:          genesis.Testnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: []string{"tpc1p54ex6jvqkz6qyld5wgm77qm7walgy664hxz2pc"},
			rewardAddrs:    "tpc1rps3xncfvepre5w754xtxxqmrmhwuackjvaft5y",
			withErr:        false,
		},

		{
			name:          "Create node for Localnet",
			numValidators: 4,
			chain:         genesis.Localnet,
			workingDir:    util.TempDirPath(),
			mnemonic:      "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: []string{
				"tpc1p54ex6jvqkz6qyld5wgm77qm7walgy664hxz2pc",
				"tpc1pdf5e0q4d6eaww3uq5pmw5aayqpaqplra0pj8z2",
				"tpc1pe5px2dddn6g4zgnu3wpwgrqpdjrufvda57a4wm",
				"tpc1p8yyhysp380j9q9gxa6vlhstgkd94238kunttpr",
			},
			rewardAddrs: "tpc1rps3xncfvepre5w754xtxxqmrmhwuackjvaft5y",
			withErr:     false,
		},
		{
			name:           "Localnet with one validator",
			numValidators:  1,
			chain:          genesis.Localnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: nil,
			rewardAddrs:    "",
			withErr:        true,
		},
		{
			name:           "Invalid mnemonic",
			numValidators:  4,
			chain:          genesis.Mainnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "",
			validatorAddrs: nil,
			rewardAddrs:    "",
			withErr:        true,
		},
	}

	for _, tt := range tests {
		wlt, rewardAddrs, err := CreateNode(tt.numValidators, tt.chain, tt.workingDir, tt.mnemonic, "")

		if tt.withErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.validatorAddrs, wlt.AllValidatorAddresses())
			assert.Equal(t, tt.rewardAddrs, rewardAddrs)
		}
	}
}
