package main

import (
	"encoding/hex"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPassword = "password"

// createTestWallet creates an encrypted wallet with a BLS account and a validator address.
func createTestWallet(t *testing.T) (path, account, validator string) {
	t.Helper()

	path = util.TempFilePath()
	mnemonic, err := wallet.GenerateMnemonic(128)
	require.NoError(t, err)

	wlt, err := wallet.Create(t.Context(), path, mnemonic, testPassword, genesis.Mainnet)
	require.NoError(t, err)
	defer wlt.Close()

	accInfo, err := wlt.NewAddress(crypto.AddressTypeBLSAccount, "account")
	require.NoError(t, err)
	valInfo, err := wlt.NewAddress(crypto.AddressTypeValidator, "validator")
	require.NoError(t, err)

	return path, accInfo.Address, valInfo.Address
}

// runWallet runs a command of the wallet in offline mode, and returns what it prints.
func runWallet(t *testing.T, path string, args ...string) string {
	t.Helper()

	stdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer

	output := make(chan string)
	go func() {
		data, _ := io.ReadAll(reader)
		output <- string(data)
	}()

	rootCmd := newRootCmd()
	rootCmd.SetArgs(append([]string{"--path", path, "--offline"}, args...))
	execErr := rootCmd.Execute()

	os.Stdout = stdout
	require.NoError(t, writer.Close())
	printed := <-output
	require.NoError(t, execErr, printed)

	return printed
}

// withWallet opens the wallet, and closes it before the next command of the test opens it.
func withWallet(t *testing.T, path string, check func(wlt *wallet.Wallet)) {
	t.Helper()

	wlt, err := wallet.Open(t.Context(), path, wallet.WithOfflineProvider())
	require.NoError(t, err)
	defer wlt.Close()

	check(wlt)
}

// A command that fails prints its error.
func TestExecutePrintsTheError(t *testing.T) {
	stdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{"unknown-command"})
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	execute(rootCmd)

	os.Stdout = stdout
	require.NoError(t, writer.Close())
	printed, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Contains(t, string(printed), "unknown command")
}

// The entry point of the program runs the command of its arguments.
func TestMainRunsTheCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	args := os.Args
	stdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Args = []string{"pactus-wallet", "--path", path, "--offline", "info"}
	os.Stdout = writer
	main()
	os.Stdout = stdout
	os.Args = args

	require.NoError(t, writer.Close())
	printed, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Contains(t, string(printed), "Network: Mainnet")
}

func TestInfoCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	output := runWallet(t, path, "info")
	assert.Contains(t, output, "Encrypted: true")
	assert.Contains(t, output, "Neutered: false")
	assert.Contains(t, output, "Network: Mainnet")
}

func TestFeeCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	output := runWallet(t, path, "fee", "0.25")
	assert.Contains(t, output, "Default fee is set to 0.25 PAC")

	withWallet(t, path, func(wlt *wallet.Wallet) {
		assert.Equal(t, amount.Amount(250_000_000), wlt.Info().DefaultFee)
	})
}

func TestSeedCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	output := runWallet(t, path, "seed", "--password", testPassword)

	withWallet(t, path, func(wlt *wallet.Wallet) {
		mnemonic, err := wlt.Mnemonic(testPassword)
		require.NoError(t, err)
		assert.Contains(t, output, mnemonic)
	})
}

func TestAddressCommands(t *testing.T) {
	path, account, validator := createTestWallet(t)

	t.Run("all", func(t *testing.T) {
		output := runWallet(t, path, "address", "all")
		assert.Contains(t, output, account)
		assert.Contains(t, output, validator)
		assert.Contains(t, output, "account")
	})

	t.Run("all with the balances and the stakes", func(t *testing.T) {
		output := runWallet(t, path, "address", "all", "--balance", "--stake")
		assert.Contains(t, output, account)
		assert.Contains(t, output, "0 PAC")
	})

	t.Run("new", func(t *testing.T) {
		output := runWallet(t, path, "address", "new", "--type", "bls", "--label", "second")

		withWallet(t, path, func(wlt *wallet.Wallet) {
			infos := wlt.ListAddresses(wallet.OnlyAccountAddresses())
			require.Len(t, infos, 2)
			assert.Equal(t, "second", infos[1].Label)
			assert.Contains(t, output, infos[1].Address)
		})
	})

	t.Run("balance", func(t *testing.T) {
		output := runWallet(t, path, "address", "balance", account)
		assert.Contains(t, output, "0 PAC")
	})

	t.Run("public key", func(t *testing.T) {
		output := runWallet(t, path, "address", "pub", account)

		unknown := runWallet(t, path, "address", "pub", "pc1zzqkzzu4vyddss052as6c37qrdcfptegquw826x")

		withWallet(t, path, func(wlt *wallet.Wallet) {
			info, err := wlt.AddressInfo(account)
			require.NoError(t, err)
			assert.Contains(t, output, info.PublicKey)
			assert.NotContains(t, unknown, info.PublicKey)
		})
	})

	t.Run("private key", func(t *testing.T) {
		output := runWallet(t, path, "address", "priv", account, "--password", testPassword)

		withWallet(t, path, func(wlt *wallet.Wallet) {
			prv, err := wlt.PrivateKey(testPassword, account)
			require.NoError(t, err)
			assert.Contains(t, output, prv.String())
		})
	})
}

func TestTransactionListCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	// A new wallet has no transaction.
	runWallet(t, path, "transaction", "list")
	runWallet(t, path, "transaction", "list", "--direction", "incoming", "--count", "5", "--skip", "1")
}

func TestNeuterCommand(t *testing.T) {
	path, _, _ := createTestWallet(t)

	output := runWallet(t, path, "neuter")
	assert.Contains(t, output, path+".neutered")
	assert.True(t, util.PathExists(path+".neutered"))

	// The wallet itself is not changed.
	withWallet(t, path, func(wlt *wallet.Wallet) {
		assert.False(t, wlt.Info().Neutered)
	})
}

// The send commands of an offline wallet sign the transaction and print it, without broadcasting it.
func TestSendCommandsOffline(t *testing.T) {
	path, account, validator := createTestWallet(t)
	receiver := "pc1zzqkzzu4vyddss052as6c37qrdcfptegquw826x"
	signedData := regexp.MustCompile(`Signed transaction data: ([0-9a-f]+)`)

	signedTx := func(t *testing.T, output string) *tx.Tx {
		t.Helper()

		match := signedData.FindStringSubmatch(output)
		require.Len(t, match, 2, output)
		data, err := hex.DecodeString(match[1])
		require.NoError(t, err)
		trx, err := tx.FromBytes(data)
		require.NoError(t, err)
		require.NoError(t, trx.BasicCheck())

		return trx
	}

	common := []string{"--fee", "0.01", "--lock-time", "100", "--no-confirm", "--password", testPassword}

	t.Run("transfer", func(t *testing.T) {
		args := append([]string{"send", "transfer", account, receiver, "1.5", "--memo", "memo"}, common...)
		trx := signedTx(t, runWallet(t, path, args...))

		assert.True(t, trx.IsTransferTx())
		assert.Equal(t, account, trx.Payload().Signer().String())
		assert.Equal(t, amount.Amount(1_500_000_000), trx.Payload().Value())
		assert.Equal(t, amount.Amount(10_000_000), trx.Fee())
		assert.Equal(t, "memo", trx.Memo())
		assert.EqualValues(t, 100, trx.LockTime())
	})

	t.Run("bond", func(t *testing.T) {
		var publicKey string
		withWallet(t, path, func(wlt *wallet.Wallet) {
			info, err := wlt.AddressInfo(validator)
			require.NoError(t, err)
			publicKey = info.PublicKey
		})

		args := append([]string{"send", "bond", account, validator, "2", "--pub", publicKey}, common...)
		trx := signedTx(t, runWallet(t, path, args...))

		assert.True(t, trx.IsBondTx())
		assert.Equal(t, amount.Amount(2_000_000_000), trx.Payload().Value())
	})

	t.Run("unbond", func(t *testing.T) {
		args := append([]string{"send", "unbond", validator}, common...)
		trx := signedTx(t, runWallet(t, path, args...))

		assert.True(t, trx.IsUnbondTx())
		assert.Equal(t, validator, trx.Payload().Signer().String())
	})

	t.Run("withdraw", func(t *testing.T) {
		args := append([]string{"send", "withdraw", validator, account, "1"}, common...)
		trx := signedTx(t, runWallet(t, path, args...))

		assert.True(t, trx.IsWithdrawTx())
		assert.Equal(t, amount.Amount(1_000_000_000), trx.Payload().Value())
	})

	// Nothing is broadcast, so nothing is recorded.
	withWallet(t, path, func(wlt *wallet.Wallet) {
		assert.Empty(t, wlt.ListTransactions())
	})
}
