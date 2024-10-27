package cmd

import (
	"bytes"
	"io"
	"os"
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

// captureOutput is a helper function to capture the printed output of a function.
func captureOutput(fun func()) string {
	// Redirect stdout to a buffer
	oldStdout := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer

	// Capture the printed output
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		outC <- buf.String()
	}()

	// Execute the function
	fun()

	// Reset stdout
	_ = writer.Close()
	os.Stdout = oldStdout
	out := <-outC

	return out
}

func TestPrintNotSupported(t *testing.T) {
	terminalSupported = false
	output := captureOutput(func() {
		PrintJSONObject([]int{1, 2, 3})
		PrintLine()
		PrintInfoMsgBoldf("This is PrintInfoMsgBoldf: %s", "msg")
		PrintInfoMsgf("This is PrintInfoMsgf: %s", "msg")
		PrintSuccessMsgf("This is PrintSuccessMsgf: %s", "msg")
		PrintWarnMsgf("This is PrintWarnMsgf: %s", "msg")
		PrintErrorMsgf("This is PrintErrorMsgf: %s", "msg")
	})

	expected := "[\n   1,\n   2,\n   3\n]\n" +
		"\n" +
		"This is PrintInfoMsgBoldf: msg\n" +
		"This is PrintInfoMsgf: msg\n" +
		"This is PrintSuccessMsgf: msg\n" +
		"This is PrintWarnMsgf: msg\n" +
		"[ERROR] This is PrintErrorMsgf: msg\n"

	assert.Equal(t, expected, output)
}

func TestPrintSupported(t *testing.T) {
	terminalSupported = true
	output := captureOutput(func() {
		PrintJSONObject([]int{1, 2, 3})
		PrintLine()
		PrintInfoMsgBoldf("This is PrintInfoMsgBoldf: %s", "msg")
		PrintInfoMsgf("This is PrintInfoMsgf: %s", "msg")
		PrintSuccessMsgf("This is PrintSuccessMsgf: %s", "msg")
		PrintWarnMsgf("This is PrintWarnMsgf: %s", "msg")
		PrintErrorMsgf("This is PrintErrorMsgf: %s", "msg")
	})

	expected := "[\n   1,\n   2,\n   3\n]\n" +
		"\n" +
		"\x1b[1mThis is PrintInfoMsgBoldf: msg\x1b[0m\n" +
		"This is PrintInfoMsgf: msg\n" +
		"\x1b[32mThis is PrintSuccessMsgf: msg\x1b[0m\n" +
		"\x1b[33mThis is PrintWarnMsgf: msg\x1b[0m\n" +
		"\x1b[31m[ERROR] This is PrintErrorMsgf: msg\x1b[0m\n"

	assert.Equal(t, expected, output)
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

	walletPath := util.TempFilePath()
	mnemonic, _ := wallet.GenerateMnemonic(128)
	walletInstance, err := wallet.Create(walletPath, mnemonic, "", genesis.Mainnet)
	assert.NoError(t, err)

	_, _ = walletInstance.NewValidatorAddress("")
	_, _ = walletInstance.NewValidatorAddress("")
	_, _ = walletInstance.NewValidatorAddress("")

	// Test 1 - Wallet without reward addresses
	valAddrsInfo := walletInstance.AllValidatorAddresses()
	confRewardAddresses := []string{}
	_, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.ErrorContains(t, err, "unable to find reward address for")

	// Test 2 - Not enough reward addresses in wallet
	rewardAddr1Info, _ := walletInstance.NewBLSAccountAddress("")
	rewardAddr2Info, _ := walletInstance.NewBLSAccountAddress("")

	_, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.ErrorContains(t, err, "unable to find reward address for")

	// Test 3 - Get reward addresses from wallet
	rewardAddr3Info, _ := walletInstance.NewBLSAccountAddress("")

	rewardAddrs, err := MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.NoError(t, err)
	assert.Equal(t, rewardAddrs[0].String(), rewardAddr1Info.Address)
	assert.Equal(t, rewardAddrs[1].String(), rewardAddr2Info.Address)
	assert.Equal(t, rewardAddrs[2].String(), rewardAddr3Info.Address)

	// Test 4 - Not enough reward addresses in config
	confRewardAddr1 := ts.RandAccAddress().String()
	confRewardAddr2 := ts.RandAccAddress().String()
	confRewardAddresses = []string{confRewardAddr1, confRewardAddr2}

	_, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.ErrorContains(t, err, "reward addresses should be 3")

	// Test 5 - Get reward addresses from config
	confRewardAddr3 := ts.RandAccAddress().String()
	confRewardAddresses = []string{confRewardAddr1, confRewardAddr2, confRewardAddr3}

	rewardAddrs, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.NoError(t, err)
	assert.Equal(t, rewardAddrs[0].String(), confRewardAddr1)
	assert.Equal(t, rewardAddrs[1].String(), confRewardAddr2)
	assert.Equal(t, rewardAddrs[2].String(), confRewardAddr3)

	// Test 6 - Set one reward addresses in config
	confRewardAddr := ts.RandAccAddress().String()
	confRewardAddresses = []string{confRewardAddr}

	rewardAddrs, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.NoError(t, err)
	assert.Equal(t, rewardAddrs[0].String(), confRewardAddr)
	assert.Equal(t, rewardAddrs[1].String(), confRewardAddr)
	assert.Equal(t, rewardAddrs[2].String(), confRewardAddr)

	// Test 7 - Set validator address as reward addresses in config
	confRewardAddr = ts.RandValAddress().String()
	confRewardAddresses = []string{confRewardAddr}

	_, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.ErrorContains(t, err, "reward address is not an account address")
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
		rewardAddrs    []string
	}{
		{
			name:           "Create node for Mainnet",
			numValidators:  1,
			chain:          genesis.Mainnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: []string{"pc1pqpu5tkuctj6ecxjs85f9apm802hhc65amwhuyw"},
			rewardAddrs:    []string{"pc1zmpnme0xrgzhml77e3k70ey9hwwwsfed6l04pqc"},
			withErr:        false,
		},
		{
			name:           "Create node for Testnet",
			numValidators:  1,
			chain:          genesis.Testnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: []string{"tpc1p54ex6jvqkz6qyld5wgm77qm7walgy664hxz2pc"},
			rewardAddrs:    []string{"tpc1zlkjrgfkrh7f9enpt730tp5vgx7tgtqzplhfksa"},
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
			rewardAddrs: []string{
				"tpc1zlkjrgfkrh7f9enpt730tp5vgx7tgtqzplhfksa",
				"tpc1ztzwc9x98j88wctmzm5t09z592lqw0sqc3rn6lu",
				"tpc1zslef8hjkwqxdcekcqxra6djgjr5gryrj8l3fyf",
				"tpc1zru3xxmgz5dqqkv0mesqq3t3luepzg3e6jeqkeu",
			},
			withErr: false,
		},
		{
			name:           "Localnet with one validator",
			numValidators:  1,
			chain:          genesis.Localnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "legal winner thank year wave sausage worth useful legal winner thank yellow",
			validatorAddrs: nil,
			rewardAddrs:    nil,
			withErr:        true,
		},
		{
			name:           "Invalid mnemonic",
			numValidators:  4,
			chain:          genesis.Mainnet,
			workingDir:     util.TempDirPath(),
			mnemonic:       "",
			validatorAddrs: nil,
			rewardAddrs:    nil,
			withErr:        true,
		},
	}

	for _, tt := range tests {
		validatorAddrs, rewardAddrs, err := CreateNode(
			tt.numValidators, tt.chain, tt.workingDir, tt.mnemonic, "")

		if tt.withErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.validatorAddrs, validatorAddrs)
			assert.Equal(t, tt.rewardAddrs, rewardAddrs)
		}
	}
}
