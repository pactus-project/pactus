package cmd

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/stretchr/testify/assert"
)

// captureOutput is a helper function to capture the printed output of a function.
func captureOutput(f func()) string {
	// Redirect stdout to a buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Capture the printed output
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Execute the function
	f()

	// Reset stdout
	w.Close()
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

	for _, test := range tests {
		walletDir := PactusWalletDir(test.home)
		defaultWalletPath := PactusDefaultWalletPath(test.home)
		genesisPath := PactusGenesisPath(test.home)
		configPath := PactusConfigPath(test.home)

		assert.Equal(t, test.expectedWalletDir, walletDir)
		assert.Equal(t, test.expectedDefaultWalletPath, defaultWalletPath)
		assert.Equal(t, test.expectedGenesisPath, genesisPath)
		assert.Equal(t, test.expectedConfigPath, configPath)
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

	for _, test := range tests {
		walletDir := PactusWalletDir(test.home)
		defaultWalletPath := PactusDefaultWalletPath(test.home)
		genesisPath := PactusGenesisPath(test.home)
		configPath := PactusConfigPath(test.home)

		assert.Equal(t, test.expectedWalletDir, walletDir)
		assert.Equal(t, test.expectedDefaultWalletPath, defaultWalletPath)
		assert.Equal(t, test.expectedGenesisPath, genesisPath)
		assert.Equal(t, test.expectedConfigPath, configPath)
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
	rewardAddr1, _ := walletInstance.NewBLSAccountAddress("")
	rewardAddr2, _ := walletInstance.NewBLSAccountAddress("")

	_, err = MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.ErrorContains(t, err, "unable to find reward address for")

	// Test 3 - Get reward addresses from wallet
	rewardAddr3, _ := walletInstance.NewBLSAccountAddress("")

	rewardAddrs, err := MakeRewardAddresses(walletInstance, valAddrsInfo, confRewardAddresses)
	assert.NoError(t, err)
	assert.Equal(t, rewardAddrs[0].String(), rewardAddr1)
	assert.Equal(t, rewardAddrs[1].String(), rewardAddr2)
	assert.Equal(t, rewardAddrs[2].String(), rewardAddr3)

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
