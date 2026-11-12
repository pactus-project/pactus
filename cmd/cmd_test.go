package cmd

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
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
		require.Error(t, err)
	})

	t.Run("No Config file, Should recover it", func(t *testing.T) {
		workingDir := util.TempDirPath()
		genPath := PactusGenesisPath(workingDir)
		gen := genesis.MainnetGenesis()
		err := gen.SaveToFile(genPath)
		require.NoError(t, err)

		_, _, err = MakeConfig(workingDir)
		require.NoError(t, err)
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
		require.NoError(t, err)
	})

	t.Run("Broken Config file, Should not keep a part of it", func(t *testing.T) {
		for _, gen := range []*genesis.Genesis{genesis.MainnetGenesis(), genesis.TestnetGenesis()} {
			t.Cleanup(crypto.ToMainnetHRP)

			workingDir := util.TempDirPath()
			confPath := PactusConfigPath(workingDir)
			require.NoError(t, gen.SaveToFile(PactusGenesisPath(workingDir)))

			// The file is valid up to the last line.
			// The values before it must not end up in the restored config.
			broken := "[grpc]\n" +
				"enable = true\n" +
				"listen = \"0.0.0.0:50051\"\n" +
				"enable_wallet = true\n" +
				"[network]\n" +
				"max_connections = 'oops'\n"
			require.NoError(t, util.WriteFile(confPath, []byte(broken)))

			conf, _, err := MakeConfig(workingDir)
			require.NoError(t, err)

			defConf := config.DefaultConfigForChain(gen.ChainType())
			assert.Equal(t, defConf.GRPC.Listen, conf.GRPC.Listen)
			assert.Equal(t, defConf.GRPC.EnableWallet, conf.GRPC.EnableWallet)

			saved, err := util.ReadFile(confPath)
			require.NoError(t, err)
			assert.NotContains(t, string(saved), "0.0.0.0:50051")
		}
	})

	t.Run("Invalid custom genesis file, Should return error", func(t *testing.T) {
		t.Cleanup(crypto.ToMainnetHRP)

		workingDir := util.TempDirPath()
		genPath := PactusGenesisPath(workingDir)
		data := `{
  "genesis_time": "2024-01-24T20:24:00Z",
  "params": {"block_version": 1, "block_interval_in_second": 10, "committee_size": 4},
  "accounts": [{"address": "not-an-address", "balance": 1}],
  "validators": []
}`
		require.NoError(t, util.WriteFile(genPath, []byte(data)))

		_, _, err := MakeConfig(workingDir)
		require.ErrorAs(t, err, &genesis.InvalidGenesisError{})
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
		require.NoError(t, err)
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

	password := "password"
	setupWallet := func() *wallet.Wallet {
		walletPath := util.TempFilePath()
		mnemonic, _ := wallet.GenerateMnemonic(128)
		wlt, err := wallet.Create(t.Context(), walletPath, mnemonic, password, genesis.Mainnet)
		require.NoError(t, err)

		_, _ = wlt.NewValidatorAddress("Validator 1")
		_, _ = wlt.NewValidatorAddress("Validator 2")
		_, _ = wlt.NewValidatorAddress("Validator 3")

		return wlt
	}

	t.Run("No reward addresses in wallet", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
		confRewardAddresses := []string{}
		_, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		assert.ErrorContains(t, err, "unable to find a reward address in the wallet")
	})

	t.Run("Wallet with one Ed25519 address", func(t *testing.T) {
		wlt := setupWallet()

		addr1Info, _ := wlt.NewEd25519AccountAddress("", password)
		_, _ = wlt.NewEd25519AccountAddress("", password)
		_, _ = wlt.NewBLSAccountAddress("")

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
		confRewardAddresses := []string{}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		require.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[1].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[2].String(), addr1Info.Address)
	})

	t.Run("Wallet with one BLS address", func(t *testing.T) {
		wlt := setupWallet()

		addr1Info, _ := wlt.NewBLSAccountAddress("")
		_, _ = wlt.NewBLSAccountAddress("")

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
		confRewardAddresses := []string{}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		require.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[1].String(), addr1Info.Address)
		assert.Equal(t, rewardAddrs[2].String(), addr1Info.Address)
	})

	t.Run("One reward address in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
		confRewardAddresses := []string{
			ts.RandAccAddress().String(),
		}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		require.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[1].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[2].String(), confRewardAddresses[0])
	})

	t.Run("Three reward addresses in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
		confRewardAddresses := []string{
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
		}
		rewardAddrs, err := MakeRewardAddresses(wlt, valAddrsInfo, confRewardAddresses)
		require.NoError(t, err)

		assert.Equal(t, rewardAddrs[0].String(), confRewardAddresses[0])
		assert.Equal(t, rewardAddrs[1].String(), confRewardAddresses[1])
		assert.Equal(t, rewardAddrs[2].String(), confRewardAddresses[2])
	})

	t.Run("Insufficient reward addresses in config", func(t *testing.T) {
		wlt := setupWallet()

		valAddrsInfo := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
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

	password := "password"
	for _, tt := range tests {
		wlt, rewardAddrs, err := CreateNode(t.Context(),
			tt.numValidators, tt.chain, tt.workingDir, tt.mnemonic, password)

		if tt.withErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)

			valInfos := wlt.ListAddresses(wallet.OnlyValidatorAddresses())
			for i, addr := range tt.validatorAddrs {
				assert.Equal(t, valInfos[i].Address, addr)
			}
			assert.Equal(t, tt.rewardAddrs, rewardAddrs)
		}
	}
}
