package vault

import (
	"context"
	"errors"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tPassword = "super_secret_password"

type testData struct {
	*testsuite.TestSuite

	vault     *Vault
	mnemonic  string
	testAddrs []*types.AddressInfo
}

// setup returns an instances of vault fo testing.
func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	mnemonic, _ := GenerateMnemonic(128)
	vault, err := CreateVaultFromMnemonic(mnemonic, 21888)
	assert.NoError(t, err)

	// Create some test address
	addr1, err := vault.NewBLSAccountAddress("bls-account-address")
	assert.NoError(t, err)
	addr2, err := vault.NewEd25519AccountAddress("ed25519-account-address", "")
	assert.NoError(t, err)
	addr3, err := vault.NewValidatorAddress("validator-address")
	assert.NoError(t, err)

	_, importedBLSPrv := ts.RandBLSKeyPair()
	addr4, addr5, err := vault.ImportBLSPrivateKey("", importedBLSPrv)
	assert.NoError(t, err)

	_, importedEd25519Prv := ts.RandEd25519KeyPair()
	addr6, err := vault.ImportEd25519PrivateKey("", importedEd25519Prv)
	assert.NoError(t, err)

	testAddrs := []*types.AddressInfo{addr1, addr2, addr3, addr4, addr5, addr6}

	assert.False(t, vault.IsEncrypted())

	// Set encryption options to minimal values for faster test execution.
	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(8),
		encrypter.OptionParallelism(1),
	}

	err = vault.UpdatePassword("", tPassword, opts...)
	assert.NoError(t, err)
	assert.True(t, vault.IsEncrypted())

	return &testData{
		TestSuite: ts,
		vault:     vault,
		mnemonic:  mnemonic,
		testAddrs: testAddrs,
	}
}

func TestCreateVaultFromMnemonic(t *testing.T) {
	td := setup(t)

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := CreateVaultFromMnemonic("invalid mnemonic phrase seed", 21888)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		recovered, err := CreateVaultFromMnemonic(td.mnemonic, 21888)
		assert.NoError(t, err)

		vaultMnemonic, err := recovered.Mnemonic("")
		assert.NoError(t, err)
		assert.Equal(t, vaultMnemonic, td.mnemonic)

		assert.Zero(t, recovered.Purposes.PurposeBLS.NextAccountIndex)
		assert.Zero(t, recovered.Purposes.PurposeBLS.NextValidatorIndex)
		assert.Zero(t, recovered.Purposes.PurposeBIP44.NextEd25519Index)

		// Recover addresses
		_, err = recovered.NewBLSAccountAddress("bls-account-address")
		assert.NoError(t, err)
		_, err = recovered.NewEd25519AccountAddress("ed25519-account-address", "")
		assert.NoError(t, err)
		_, err = recovered.NewValidatorAddress("validator-address")
		assert.NoError(t, err)

		assert.Equal(t, recovered.Purposes, td.vault.Purposes)
	})
}

func TestGetPrivateKeys(t *testing.T) {
	td := setup(t)

	t.Run("Unknown purpose", func(t *testing.T) {
		path, _ := addresspath.FromString("m/0")
		_, err := td.vault.PrivateKeys(tPassword, []addresspath.Path{path})
		assert.ErrorIs(t, err, ErrUnsupportedPurpose)
	})

	t.Run("No password", func(t *testing.T) {
		path, _ := addresspath.FromString("m/44/21888/3/0")
		_, err := td.vault.PrivateKeys("", []addresspath.Path{path})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Invalid password", func(t *testing.T) {
		path, _ := addresspath.FromString("m/44/21888/3/0")
		_, err := td.vault.PrivateKeys("wrong_password", []addresspath.Path{path})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Check all the private keys", func(t *testing.T) {
		for _, info := range td.testAddrs {
			path, _ := addresspath.FromString(info.Path)
			prv, err := td.vault.PrivateKeys(tPassword, []addresspath.Path{path})
			assert.NoError(t, err)

			switch path.AddressType() {
			case crypto.AddressTypeBLSAccount,
				crypto.AddressTypeValidator:
				pub, _ := bls.PublicKeyFromString(info.PublicKey)
				require.True(t, prv[0].PublicKey().EqualsTo(pub))
			case crypto.AddressTypeEd25519Account:
				pub, _ := ed25519.PublicKeyFromString(info.PublicKey)
				require.True(t, prv[0].PublicKey().EqualsTo(pub))
			case crypto.AddressTypeTreasury:
				assert.Fail(t, "not supported")
			}
		}
	})
}

func TestImportBLSPrivateKey(t *testing.T) {
	td := setup(t)

	_, prv := td.RandBLSKeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		_, _, err := td.vault.ImportBLSPrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		accInfo, valInfo, err := td.vault.ImportBLSPrivateKey(tPassword, prv)
		assert.NoError(t, err)

		assert.Equal(t, prv.PublicKeyNative().String(), accInfo.PublicKey)
		assert.Equal(t, prv.PublicKeyNative().String(), valInfo.PublicKey)

		assert.Equal(t, "m/65535'/21888'/1'/2'", accInfo.Path)
		assert.Equal(t, "m/65535'/21888'/2'/2'", valInfo.Path)
	})
}

func TestImportEd25519PrivateKey(t *testing.T) {
	td := setup(t)

	_, prv := td.RandEd25519KeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		_, err := td.vault.ImportEd25519PrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		info, err := td.vault.ImportEd25519PrivateKey(tPassword, prv)
		assert.NoError(t, err)

		assert.Equal(t, prv.PublicKeyNative().String(), info.PublicKey)
		assert.Equal(t, "m/65535'/21888'/3'/2'", info.Path)
	})
}

func TestGetMnemonic(t *testing.T) {
	td := setup(t)

	t.Run("Invalid password", func(t *testing.T) {
		_, err := td.vault.Mnemonic("invalid-password")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("No password", func(t *testing.T) {
		_, err := td.vault.Mnemonic("")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		m, err := td.vault.Mnemonic(tPassword)
		assert.NoError(t, err)
		assert.Equal(t, m, td.mnemonic)
	})
}

func TestUpdatePassword(t *testing.T) {
	td := setup(t)

	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(1),
		encrypter.OptionParallelism(1),
	}

	newPassword := "new-password"

	t.Run("Rejects empty current password", func(t *testing.T) {
		err := td.vault.UpdatePassword("", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Rejects incorrect current password", func(t *testing.T) {
		err := td.vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Updates password with valid current password", func(t *testing.T) {
		assert.NoError(t, td.vault.UpdatePassword(tPassword, newPassword, opts...))
		assert.True(t, td.vault.IsEncrypted())
	})

	t.Run("Old password is no longer valid after update", func(t *testing.T) {
		err := td.vault.UpdatePassword(tPassword, newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Clears vault password when new password is empty", func(t *testing.T) {
		assert.NoError(t, td.vault.UpdatePassword(newPassword, ""))
		assert.False(t, td.vault.IsEncrypted())
	})
}

func TestNeuter(t *testing.T) {
	td := setup(t)

	td.vault.Neuter()

	assert.True(t, td.vault.IsNeutered())

	_, err := td.vault.Mnemonic(tPassword)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = td.vault.PrivateKeys(tPassword, []addresspath.Path{})
	assert.ErrorIs(t, err, ErrNeutered)

	_, _, err = td.vault.ImportBLSPrivateKey("any", nil)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = td.vault.ImportEd25519PrivateKey("any", nil)
	assert.ErrorIs(t, err, ErrNeutered)

	err = td.vault.UpdatePassword("any", "any")
	assert.ErrorIs(t, err, ErrNeutered)
}

// TestAddressRecovery tests the address recovery functionality according to PIP-41 specification.
// This test verifies that the RecoverAddresses function correctly identifies and recovers
// previously used addresses when restoring a wallet from a mnemonic phrase.
//
// The first 8 BLS account addresses for the test mnemonic are:
// pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf (index 0)
// pc1z4xuja689hg2434yhr32clhn97x6afw58a5n9ns (index 1)
// pc1zaj6dzh6zg8zsgzy2rrtvyyeg0l4d32p8e6xn5h (index 2)
// pc1ztmex7taes23h6z4jf0awwmps0zpzmecuzcsev0 (index 3)
// pc1zkry0kt7fxufqjql6zus54a397w4ukqqg0l2sz4 (index 4)
// pc1zqar4tm23a3k0cyy3n86fq59psajah3wgm3hc4x (index 5)
// pc1zpmxu83gp7y84ekn89rfkyf099sj6f9jlmututf (index 6)
// pc1zydjhrq06ngg6nwqs8n8jkyw6u58qlqc5cqqxht (index 7)
//
// The first 8 Ed25519 account addresses for the test mnemonic are:
// pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3 (index 0)
// pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94 (index 1)
// pc1ruumtknmwr6ns32rkezfph38tawwx7gesmykk4g (index 2)
// pc1r4waddcacrxw2vg4ge8vtlnk9mnccnuv0374xuv (index 3)
// pc1re5an4nasvgpmxmuptxxd8hqy6adncqy4qyhj8w (index 4)
// pc1rul34wczhq44s5chtxvlgmrgf6dp0xx47zzg9ud (index 5)
// pc1r77rvd98gld8vfgzfa89har678dlpm9pkxex4zf (index 6)
// pc1rmzpqfhs4ekrevmwwj2gsz6m4kjym3eg99x7zk5 (index 7)
//
// The test uses a mock hasActivity function to simulate blockchain activity checks.
func TestAddressRecovery(t *testing.T) {
	//nolint:dupword // has duplicated words
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus"

	t.Run("recover addresses from a fresh wallet without any active addresses", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		// Mock hasActivity to return false for all addresses (no active addresses)
		hasActivity := func(_ string) (bool, error) {
			return false, nil
		}

		recovered, err := vault.RecoverAddresses(context.Background(), "", hasActivity)
		assert.NoError(t, err)

		assert.Empty(t, recovered)
	})

	t.Run("recover addresses with one gap at the beginning", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		// Mock hasActivity to return true only for the first call (address at index 0)
		hasActivity := func(addr string) (bool, error) {
			return addr == "pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94" ||
				addr == "pc1z4xuja689hg2434yhr32clhn97x6afw58a5n9ns", nil
		}

		recovered, err := vault.RecoverAddresses(context.Background(), "", hasActivity)
		assert.NoError(t, err)

		// Should have 4 addresses
		assert.Len(t, recovered, 4)
		assert.Equal(t, "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", recovered[0].Address)
		assert.Equal(t, "pc1z4xuja689hg2434yhr32clhn97x6afw58a5n9ns", recovered[1].Address)
		assert.Equal(t, "pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3", recovered[2].Address)
		assert.Equal(t, "pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94", recovered[3].Address)
	})

	t.Run("recover addresses with gaps in the middle of the address list", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		hasActivity := func(addr string) (bool, error) {
			return addr == "pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3" ||
				addr == "pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94" ||
				addr == "pc1r4waddcacrxw2vg4ge8vtlnk9mnccnuv0374xuv" ||
				addr == "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf" ||
				addr == "pc1ztmex7taes23h6z4jf0awwmps0zpzmecuzcsev0", nil
		}

		recovered, err := vault.RecoverAddresses(context.Background(), "", hasActivity)
		assert.NoError(t, err)

		assert.Len(t, recovered, 8)

		assert.Equal(t, "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", recovered[0].Address)
		assert.Equal(t, "pc1z4xuja689hg2434yhr32clhn97x6afw58a5n9ns", recovered[1].Address)
		assert.Equal(t, "pc1zaj6dzh6zg8zsgzy2rrtvyyeg0l4d32p8e6xn5h", recovered[2].Address)
		assert.Equal(t, "pc1ztmex7taes23h6z4jf0awwmps0zpzmecuzcsev0", recovered[3].Address)
		assert.Equal(t, "pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3", recovered[4].Address)
		assert.Equal(t, "pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94", recovered[5].Address)
		assert.Equal(t, "pc1ruumtknmwr6ns32rkezfph38tawwx7gesmykk4g", recovered[6].Address)
		assert.Equal(t, "pc1r4waddcacrxw2vg4ge8vtlnk9mnccnuv0374xuv", recovered[7].Address)
	})

	t.Run("prevent recovering existing address", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		_, _ = vault.NewEd25519AccountAddress("existing address", "")

		hasActivity := func(addr string) (bool, error) {
			return addr == "pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3", nil
		}

		recovered, err := vault.RecoverAddresses(context.Background(), "", hasActivity)
		assert.NoError(t, err)

		assert.Len(t, recovered, 0)
	})

	t.Run("error handling", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		// Mock hasActivity to return an error
		hasActivity := func(_ string) (bool, error) {
			return false, errors.New("blockchain connection error")
		}

		_, err = vault.RecoverAddresses(context.Background(), "", hasActivity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "blockchain connection error")
	})

	t.Run("cancel recovery with context cancel signal", func(t *testing.T) {
		vault, err := CreateVaultFromMnemonic(testMnemonic, 21888) // Mainnet
		assert.NoError(t, err)

		// Create a cancellable context
		ctx, cancel := context.WithCancel(context.Background())

		// Counter to track how many times hasActivity is called
		callCount := 0

		// Mock hasActivity to cancel context after a few calls
		hasActivity := func(_ string) (bool, error) {
			callCount++
			// Cancel the context after 3 calls to simulate interruption during recovery
			if callCount >= 3 {
				cancel()
			}

			return false, nil
		}

		_, err = vault.RecoverAddresses(ctx, "", hasActivity)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}
