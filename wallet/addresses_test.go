package wallet_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPrivateKey(t *testing.T) {
	td := setup(t)

	t.Run("Unknown address", func(t *testing.T) {
		addr := td.RandAccAddress().String()
		td.mockStorage.EXPECT().AddressInfo(addr).Return(nil, storage.ErrNotFound)

		_, err := td.wallet.PrivateKey(td.password, addr)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
}

func TestAddressCount(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().AddressCount().Return(3)
	assert.Equal(t, 3, td.wallet.AddressCount())
}

func TestHasAddress(t *testing.T) {
	td := setup(t)

	knownAddr := td.RandAccAddress().String()
	td.mockStorage.EXPECT().HasAddress(knownAddr).Return(true)
	td.mockStorage.EXPECT().HasAddress(gomock.Any()).Return(false)

	t.Run("Known address", func(t *testing.T) {
		assert.True(t, td.wallet.HasAddress(knownAddr))
	})

	t.Run("Unknown address", func(t *testing.T) {
		unknownAddr := td.RandAccAddress().String()
		assert.False(t, td.wallet.HasAddress(unknownAddr))
	})
}

func TestListAddresses(t *testing.T) {
	td := setup(t)

	valInfo, _ := td.testVault.NewValidatorAddress("addr-1")
	accInfo, _ := td.testVault.NewBLSAccountAddress("addr-2")
	edInfo, _ := td.testVault.NewEd25519AccountAddress("addr-3", td.password)

	_, prv1 := td.RandBLSKeyPair()
	impAcc, impVal, _ := td.testVault.ImportBLSPrivateKey(td.password, prv1)
	_, prv2 := td.RandEd25519KeyPair()
	impEd, _ := td.testVault.ImportEd25519PrivateKey(td.password, prv2)

	existing := []types.AddressInfo{
		*impAcc,
		*impVal,
		*impEd,
		*valInfo,
		*accInfo,
		*edInfo,
	}
	td.mockStorage.EXPECT().AllAddresses().Return(existing).AnyTimes()

	t.Run("List should be sorted", func(t *testing.T) {
		listed := td.wallet.ListAddresses()

		assert.Equal(t, "m/44'/21888'/3'/0'", listed[0].Path)
		assert.Equal(t, "m/12381'/21888'/1'/0", listed[1].Path)
		assert.Equal(t, "m/12381'/21888'/2'/0", listed[2].Path)
		assert.Equal(t, "m/65535'/21888'/1'/0'", listed[3].Path)
		assert.Equal(t, "m/65535'/21888'/2'/0'", listed[4].Path)
		assert.Equal(t, "m/65535'/21888'/3'/1'", listed[5].Path)
	})

	t.Run("Only account addresses", func(t *testing.T) {
		infos := td.wallet.ListAddresses(wallet.OnlyAccountAddresses())

		for _, i := range infos {
			addr, _ := crypto.AddressFromString(i.Address)
			assert.True(t, addr.IsAccountAddress())
		}
	})

	t.Run("Only validator addresses", func(t *testing.T) {
		infos := td.wallet.ListAddresses(wallet.OnlyValidatorAddresses())

		for _, i := range infos {
			addr, _ := crypto.AddressFromString(i.Address)
			assert.True(t, addr.IsValidatorAddress())
		}
	})
}

func TestNewValidatorAddress(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
	td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)
	label := td.RandString(16)
	addressInfo, err := td.wallet.NewValidatorAddress(label)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Equal(t, label, addressInfo.Label)
	assert.Equal(t, "m/12381'/21888'/1'/0", addressInfo.Path)

	pub, _ := bls.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.ValidatorAddress().String(), addressInfo.Address)
}

func TestNewBLSAccountAddress(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
	td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)
	label := td.RandString(16)
	addressInfo, err := td.wallet.NewBLSAccountAddress(label)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Equal(t, label, addressInfo.Label)
	assert.Equal(t, "m/12381'/21888'/2'/0", addressInfo.Path)

	pub, _ := bls.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.AccountAddress().String(), addressInfo.Address)
}

func TestNewE225519AccountAddress(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
	td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)
	label := td.RandString(16)
	addressInfo, err := td.wallet.NewEd25519AccountAddress(label, td.password)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Equal(t, label, addressInfo.Label)
	assert.Equal(t, "m/44'/21888'/3'/0'", addressInfo.Path)

	pub, _ := ed25519.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.AccountAddress().String(), addressInfo.Address)
}

func TestImportBLSPrivateKey(t *testing.T) {
	td := setup(t)

	pub, prv := td.RandBLSKeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(false)

		err := td.wallet.ImportBLSPrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(false)
		td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil).Times(2)
		td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)

		err := td.wallet.ImportBLSPrivateKey(td.password, prv)
		assert.NoError(t, err)

		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(true)
		td.mockStorage.EXPECT().HasAddress(pub.ValidatorAddress().String()).Return(true)

		assert.True(t, td.wallet.HasAddress(pub.AccountAddress().String()))
		assert.True(t, td.wallet.HasAddress(pub.ValidatorAddress().String()))
	})

	t.Run("Reimporting private key", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(true)

		err := td.wallet.ImportBLSPrivateKey(td.password, prv)
		assert.ErrorIs(t, err, wallet.ErrAddressExists)
	})
}

func TestImportEd25519PrivateKey(t *testing.T) {
	td := setup(t)

	pub, prv := td.RandEd25519KeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(false)

		err := td.wallet.ImportEd25519PrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(false)
		td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
		td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)

		err := td.wallet.ImportEd25519PrivateKey(td.password, prv)
		assert.NoError(t, err)

		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(true)
		assert.True(t, td.wallet.HasAddress(pub.AccountAddress().String()))
	})

	t.Run("Reimporting private key", func(t *testing.T) {
		td.mockStorage.EXPECT().HasAddress(pub.AccountAddress().String()).Return(true)

		err := td.wallet.ImportEd25519PrivateKey(td.password, prv)
		assert.ErrorIs(t, err, wallet.ErrAddressExists)
	})
}

func TestSetAddressLabel(t *testing.T) {
	td := setup(t)

	td.mockStorage.EXPECT().InsertAddress(gomock.Any()).Return(nil)
	td.mockStorage.EXPECT().UpdateVault(td.testVault).Return(nil)
	testAddr, err := td.wallet.NewBLSAccountAddress("test")
	require.NoError(t, err)

	t.Run("Set label for unknown address", func(t *testing.T) {
		invAddr := td.RandAccAddress().String()
		td.mockStorage.EXPECT().AddressInfo(invAddr).Return(nil, storage.ErrNotFound).AnyTimes()

		err := td.wallet.SetAddressLabel(invAddr, "i have label")
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Empty(t, td.wallet.AddressLabel(invAddr))
	})

	t.Run("Update label", func(t *testing.T) {
		updatedInfo := *testAddr
		updatedInfo.Label = "I have a label"

		td.mockStorage.EXPECT().AddressInfo(testAddr.Address).Return(testAddr, nil)
		td.mockStorage.EXPECT().UpdateAddress(&updatedInfo).Return(nil)
		td.mockStorage.EXPECT().AddressInfo(testAddr.Address).Return(&updatedInfo, nil)

		err := td.wallet.SetAddressLabel(testAddr.Address, "I have a label")
		assert.NoError(t, err)
		assert.Equal(t, "I have a label", td.wallet.AddressLabel(testAddr.Address))
	})

	t.Run("Remove label", func(t *testing.T) {
		noLabelInfo := *testAddr
		noLabelInfo.Label = ""

		td.mockStorage.EXPECT().AddressInfo(testAddr.Address).Return(testAddr, nil)
		td.mockStorage.EXPECT().UpdateAddress(&noLabelInfo).Return(nil)
		td.mockStorage.EXPECT().AddressInfo(testAddr.Address).Return(&noLabelInfo, nil)

		err := td.wallet.SetAddressLabel(testAddr.Address, "")
		assert.NoError(t, err)
		assert.Empty(t, td.wallet.AddressLabel(testAddr.Address))
	})
}
